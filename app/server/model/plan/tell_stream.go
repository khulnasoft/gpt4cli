package plan

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"gpt4cli-server/db"
	"gpt4cli-server/hooks"
	"gpt4cli-server/model"
	"gpt4cli-server/types"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/sashabaranov/go-openai"
)

const MaxAutoContinueIterations = 100
const MaxSendRate = 30 * time.Millisecond
const MaxTellStreamRetries = 4

func (state *activeTellStreamState) listenStream(stream *openai.ChatCompletionStream) {
	defer stream.Close()

	clients := state.clients
	auth := state.auth
	req := state.req
	plan := state.plan
	planId := plan.Id
	branch := state.branch
	currentOrgId := state.currentOrgId
	currentUserId := state.currentUserId
	convo := state.convo
	summaries := state.summaries
	summarizedToMessageId := state.summarizedToMessageId
	iteration := state.iteration
	missingFileResponse := state.missingFileResponse
	replyId := state.replyId
	replyParser := state.replyParser
	settings := state.settings

	active := GetActivePlan(planId, branch)

	if active == nil {
		log.Printf("listenStream - Active plan not found for plan ID %s on branch %s\n", planId, branch)
		return
	}

	replyFiles := []string{}
	chunksReceived := 0
	maybeRedundantBacktickContent := ""
	fileOpen := false

	// Create a timer that will trigger if no chunk is received within the specified duration
	timer := time.NewTimer(model.OPENAI_STREAM_CHUNK_TIMEOUT)
	defer timer.Stop()

	streamFinished := false

	execHookOnStop := func(sendStreamErr bool) {
		_, apiErr := hooks.ExecHook(hooks.DidSendModelRequest, hooks.HookParams{
			Auth: auth,
			Plan: plan,
			DidSendModelRequestParams: &hooks.DidSendModelRequestParams{
				InputTokens:   state.totalRequestTokens,
				OutputTokens:  active.NumTokens,
				ModelName:     state.settings.ModelPack.Planner.BaseModelConfig.ModelName,
				ModelProvider: state.settings.ModelPack.Planner.BaseModelConfig.Provider,
				ModelPackName: state.settings.ModelPack.Name,
				ModelRole:     shared.ModelRolePlanner,
				Purpose:       "Generated plan reply",
			},
		})

		if apiErr != nil {
			log.Printf("Error executing did send model request hook after cancel or error: %v\n", apiErr)

			if sendStreamErr {
				activePlan := GetActivePlan(planId, branch)

				if activePlan == nil {
					log.Printf(" Active plan not found for plan ID %s on branch %s\n", planId, branch)
					return
				}

				activePlan.StreamDoneCh <- apiErr
			}
		}
	}

mainLoop:
	for {
		select {
		case <-active.Ctx.Done():
			// The main modelContext was canceled (not the timer)
			log.Println("\nTell: stream canceled")
			execHookOnStop(false)
			return
		case <-timer.C:
			// Timer triggered because no new chunk was received in time
			log.Println("\nTell: stream timeout due to inactivity")
			if streamFinished {
				log.Println("Tell stream finished—timed out waiting for usage chunk")
				execHookOnStop(true)
				return
			} else {
				state.onError(fmt.Errorf("stream timeout due to inactivity | This usually means the model is not responding."), true, "", "")
				continue mainLoop
			}

		default:
			response, err := stream.Recv()

			if err == nil {
				// Successfully received a chunk, reset the timer
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(model.OPENAI_STREAM_CHUNK_TIMEOUT)
			} else {
				if err.Error() == "context canceled" {
					log.Println("Tell: stream context canceled")
					execHookOnStop(false)
					return
				}

				log.Printf("Tell: error receiving stream chunk: %v\n", err)
				execHookOnStop(true)
				return
			}

			if len(response.Choices) == 0 {
				if response.Usage != nil {

					log.Println("Tell stream usage:")
					spew.Dump(response.Usage)

					_, apiErr := hooks.ExecHook(hooks.DidSendModelRequest, hooks.HookParams{
						Auth: auth,
						Plan: plan,
						DidSendModelRequestParams: &hooks.DidSendModelRequestParams{
							InputTokens:   response.Usage.PromptTokens,
							OutputTokens:  response.Usage.CompletionTokens,
							ModelName:     state.settings.ModelPack.Planner.BaseModelConfig.ModelName,
							ModelProvider: state.settings.ModelPack.Planner.BaseModelConfig.Provider,
							ModelPackName: state.settings.ModelPack.Name,
							ModelRole:     shared.ModelRolePlanner,
							Purpose:       "Generated plan reply",
						},
					})

					if apiErr != nil {
						log.Printf("Tell stream: error executing did send model request hook: %v\n", err)

						// ensure the active plan is still available
						activePlan := GetActivePlan(planId, branch)

						if activePlan == nil {
							log.Printf(" Active plan not found for plan ID %s on branch %s\n", planId, branch)
							return
						}

						activePlan.StreamDoneCh <- apiErr
					}
					return
				}

				state.onError(fmt.Errorf("stream finished with no choices | This usually means the model failed to generate a valid response."), true, "", "")
				continue mainLoop
			}

			// if stream finished and it's not a usage chunk, keeep listening for usage chunk
			if streamFinished {
				// log.Println("Tell stream finished—no usage chunk-will keep listening")
				// continue
			}

			if len(response.Choices) > 1 {
				state.onError(fmt.Errorf("stream finished with more than one choice | This usually means the model failed to generate a valid response."), true, "", "")
				continue mainLoop
			}

			choice := response.Choices[0]

			if choice.FinishReason != "" {
				log.Println("Model stream finished")
				log.Println("Finish reason: ", choice.FinishReason)

				if choice.FinishReason == "error" {
					state.onError(fmt.Errorf("model stopped with error status | This usually means the model is not responding."), true, "", "")
					continue mainLoop
				}

				time.Sleep(30 * time.Millisecond)
				active.FlushStreamBuffer()
				time.Sleep(100 * time.Millisecond)

				active.Stream(shared.StreamMessage{
					Type: shared.StreamMessageDescribing,
				})
				active.FlushStreamBuffer()

				err := db.SetPlanStatus(planId, branch, shared.PlanStatusDescribing, "")
				if err != nil {
					state.onError(fmt.Errorf("failed to set plan status to describing: %v", err), true, "", "")
					continue mainLoop
				}

				var generatedDescription *db.ConvoMessageDescription
				var shouldContinue bool
				var subtaskFinished bool

				autoLoadContextFiles := state.checkAutoLoadContext()
				hasNewSubtasks := state.checkNewSubtasks()
				moveFiles := state.checkMoveFileOps()
				removeFiles := state.checkRemoveFileOps()
				resetFiles := state.checkResetFileOps()
				var allSubtasksFinished bool

				if req.BuildMode == shared.BuildModeAuto {
					getBuildState := func() *activeBuildStreamState {
						return &activeBuildStreamState{
							tellState:     state,
							clients:       clients,
							auth:          auth,
							currentOrgId:  currentOrgId,
							currentUserId: currentUserId,
							plan:          plan,
							branch:        branch,
							settings:      settings,
							modelContext:  state.modelContext,
						}
					}

					if len(moveFiles) > 0 {
						log.Println("Detected move files, queuing builds")
						for i, moveFile := range moveFiles {
							getBuildState().queueBuilds([]*types.ActiveBuild{
								{
									ReplyId:         replyId,
									Idx:             i,
									Path:            moveFile.Source,
									IsMoveOp:        true,
									MoveDestination: moveFile.Destination,
								},
							})
						}
					}

					if len(removeFiles) > 0 {
						log.Println("Detected remove files, queuing builds")
						for i, removeFile := range removeFiles {
							getBuildState().queueBuilds([]*types.ActiveBuild{
								{
									ReplyId:    replyId,
									Idx:        i,
									Path:       removeFile,
									IsRemoveOp: true,
								},
							})
						}
					}

					if len(resetFiles) > 0 {
						log.Println("Detected reset files, queuing builds")

						for i, resetFile := range resetFiles {
							getBuildState().queueBuilds([]*types.ActiveBuild{
								{
									ReplyId:   replyId,
									Idx:       i,
									Path:      resetFile,
									IsResetOp: true,
								},
							})
						}
					}
				}

				var errCh = make(chan error, 2)

				go func() {
					if len(replyFiles) > 0 {
						log.Println("Generating plan description")

						res, err := state.genPlanDescription()
						if err != nil {
							errCh <- fmt.Errorf("failed to generate plan description: %v", err)
							return
						}

						generatedDescription = res
						generatedDescription.OrgId = currentOrgId
						generatedDescription.SummarizedToMessageId = summarizedToMessageId
						generatedDescription.MadePlan = true
						generatedDescription.Files = replyFiles

						log.Println("Generated plan description.")
					}
					errCh <- nil
				}()

				if req.IsChatOnly || len(autoLoadContextFiles) > 0 || hasNewSubtasks || !req.AutoContinue {

					// if we're auto-loading context files, we always want to continue for at least another iteration with the loaded context (even if it's chat only)
					if len(autoLoadContextFiles) > 0 && !hasNewSubtasks {
						log.Printf("Auto loading context files, so continuing to planning phase")
						shouldContinue = true
					} else if req.IsChatOnly {
						log.Printf("Chat only, won't continue")
						shouldContinue = false
					} else if hasNewSubtasks {
						log.Printf("Has new subtasks, can continue")
						shouldContinue = req.AutoContinue
					}

					errCh <- nil
				} else {
					go func() {
						log.Println("Getting exec status")
						subtaskFinished, shouldContinue, err = state.execStatusShouldContinue(active.CurrentReplyContent, active.Ctx)
						if err != nil {
							errCh <- fmt.Errorf("failed to get exec status: %v", err)
							return
						}

						log.Printf("Should continue: %v\n", shouldContinue)

						errCh <- nil
					}()
				}

				for i := 0; i < 2; i++ {
					err := <-errCh
					if err != nil {
						state.onError(err, true, "", "")
						continue mainLoop
					}
				}

				log.Println("Locking repo to store assistant reply and description")

				repoLockId, err := db.LockRepo(
					db.LockRepoParams{
						OrgId:    currentOrgId,
						UserId:   currentUserId,
						PlanId:   planId,
						Branch:   branch,
						Scope:    db.LockScopeWrite,
						Ctx:      active.Ctx,
						CancelFn: active.CancelFn,
					},
				)

				if err != nil {
					log.Printf("Error locking repo: %v\n", err)
					active.StreamDoneCh <- &shared.ApiError{
						Type:   shared.ApiErrorTypeOther,
						Status: http.StatusInternalServerError,
						Msg:    "Error locking repo",
					}
					continue mainLoop
				}

				log.Println("Locked repo for assistant reply and description")

				err = func() error {
					defer func() {
						if err != nil {
							log.Printf("Error storing reply and description: %v\n", err)
							err = db.GitClearUncommittedChanges(auth.OrgId, planId)
							if err != nil {
								log.Printf("Error clearing uncommitted changes: %v\n", err)
							}
						}

						log.Println("Unlocking repo for assistant reply and description")

						err = db.DeleteRepoLock(repoLockId)
						if err != nil {
							log.Printf("Error unlocking repo: %v\n", err)
							active.StreamDoneCh <- &shared.ApiError{
								Type:   shared.ApiErrorTypeOther,
								Status: http.StatusInternalServerError,
								Msg:    "Error unlocking repo",
							}
						}
					}()

					assistantMsg, convoCommitMsg, err := state.storeAssistantReply() // updates state.convo
					convo = state.convo

					if err != nil {
						state.onError(fmt.Errorf("failed to store assistant message: %v", err), true, "", "")
						return err
					}

					log.Println("getting description for assistant message: ", assistantMsg.Id)

					var description *db.ConvoMessageDescription
					if len(replyFiles) == 0 {
						description = &db.ConvoMessageDescription{
							OrgId:                 currentOrgId,
							PlanId:                planId,
							ConvoMessageId:        assistantMsg.Id,
							SummarizedToMessageId: summarizedToMessageId,
							BuildPathsInvalidated: map[string]bool{},
							MadePlan:              false,
						}
					} else {
						description = generatedDescription
						description.ConvoMessageId = assistantMsg.Id
					}

					log.Println("Storing description")
					err = db.StoreDescription(description)

					if err != nil {
						state.onError(fmt.Errorf("failed to store description: %v", err), false, assistantMsg.Id, convoCommitMsg)
						return err
					}
					log.Println("Description stored")

					if hasNewSubtasks || subtaskFinished {
						if subtaskFinished && state.currentSubtask != nil {
							log.Println("Subtask finished")
							log.Println("Current subtask:")
							log.Println(state.currentSubtask.Title)
							state.currentSubtask.IsFinished = true

							log.Println("Updated state. Current subtask:")
							log.Println(state.currentSubtask)
						}

						log.Println("Storing plan subtasks")
						err = db.StorePlanSubtasks(currentOrgId, planId, state.subtasks)
						if err != nil {
							log.Printf("Error storing plan subtasks: %v\n", err)
							state.onError(fmt.Errorf("failed to store plan subtasks: %v", err), false, assistantMsg.Id, convoCommitMsg)
							return err
						}

						state.currentSubtask = nil
						allSubtasksFinished = true
						for _, subtask := range state.subtasks {
							if !subtask.IsFinished {
								state.currentSubtask = subtask
								allSubtasksFinished = false
								break
							}
						}

						log.Println("Set new current subtask. Current subtask:")
						log.Println(state.currentSubtask)
						log.Println("All subtasks finished:", allSubtasksFinished)

						log.Println("Update state of subtasks")
						spew.Dump(state.subtasks)
					}

					// spew.Dump(description)

					log.Println("Comitting reply message, description, and subtasks")

					err = db.GitAddAndCommit(currentOrgId, planId, branch, convoCommitMsg)
					if err != nil {
						state.onError(fmt.Errorf("failed to commit: %v", err), false, assistantMsg.Id, convoCommitMsg)
						return err
					}
					log.Println("Assistant reply, description, and subtasks committed")

					return nil
				}()

				if err != nil {
					continue mainLoop
				}

				// summarize convo needs to come *after* the reply is stored in order to correctly summarize the latest message
				log.Println("summarize convo")
				envVar := settings.ModelPack.PlanSummary.BaseModelConfig.ApiKeyEnvVar
				client := clients[envVar]

				// summarize in the background
				go func() {
					err := summarizeConvo(client, settings.ModelPack.PlanSummary, summarizeConvoParams{
						auth:                  auth,
						plan:                  plan,
						branch:                branch,
						convo:                 convo,
						summaries:             summaries,
						userPrompt:            state.userPrompt,
						currentOrgId:          currentOrgId,
						currentReply:          active.CurrentReplyContent,
						currentReplyNumTokens: active.NumTokens,
						modelPackName:         settings.ModelPack.Name,
					}, active.SummaryCtx)

					if err != nil {
						log.Printf("Error summarizing convo: %v\n", err)
						active.StreamDoneCh <- &shared.ApiError{
							Type:   shared.ApiErrorTypeOther,
							Status: http.StatusInternalServerError,
							Msg:    fmt.Sprintf("Error summarizing convo: %v", err),
						}
					}
				}()

				log.Println("Sending active.CurrentReplyDoneCh <- true")

				active.CurrentReplyDoneCh <- true

				log.Println("Resetting active.CurrentReplyDoneCh")

				UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
					ap.CurrentStreamingReplyId = ""
					ap.CurrentReplyDoneCh = nil
				})

				log.Printf("len(autoLoadContextFiles): %d\n", len(autoLoadContextFiles))
				if len(autoLoadContextFiles) > 0 {
					log.Println("Sending stream message to load context files")

					active.Stream(shared.StreamMessage{
						Type:             shared.StreamMessageLoadContext,
						LoadContextFiles: autoLoadContextFiles,
					})
					active.FlushStreamBuffer()

					// Force a small delay to ensure message is processed
					time.Sleep(100 * time.Millisecond)

					log.Println("Waiting for client to auto load context (30s timeout)")

					select {
					case <-active.Ctx.Done():
						log.Println("Context cancelled while waiting for auto load context")
						execHookOnStop(true)
						return
					case <-time.After(30 * time.Second):
						log.Println("Timeout waiting for auto load context")
						state.onError(fmt.Errorf("timeout waiting for auto load context response"), true, "", "")
						continue mainLoop
					case <-active.AutoLoadContextCh:
					}
				}

				// if we're auto-loading context files, we always want to continue for at least another iteration with the loaded context
				log.Printf("req.AutoContinue: %v\n", req.AutoContinue)
				log.Printf("shouldContinue: %v\n", shouldContinue)
				log.Printf("iteration: %d\n", iteration)
				log.Printf("MaxAutoContinueIterations: %d\n", MaxAutoContinueIterations)
				log.Printf("hasNewSubtasks: %v\n", hasNewSubtasks)
				log.Printf("len(state.subtasks): %d\n", len(state.subtasks))
				log.Printf("allSubtasksFinished: %v\n", allSubtasksFinished)

				if (len(autoLoadContextFiles) > 0 && !hasNewSubtasks) ||
					(req.AutoContinue && shouldContinue && iteration < MaxAutoContinueIterations &&
						!(len(state.subtasks) > 0 && allSubtasksFinished)) {
					log.Println("Auto continue plan")
					// continue plan
					execTellPlan(clients, plan, branch, auth, req, iteration+1, "", false, 0)
				} else {
					var buildFinished bool
					UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
						buildFinished = ap.BuildFinished()
						ap.RepliesFinished = true
					})

					log.Printf("Won't continue plan. Build finished: %v\n", buildFinished)

					time.Sleep(50 * time.Millisecond)

					// note: tell-based verification is currently disabled, so 'verifyOrFinish' will pass through to 'finish' logic
					// ShouldVerifyDiff will always be false unless we re-enable tell-based verification
					if buildFinished {
						log.Println("Reply is finished and build is finished, calling verifyOrFinish")
						active := GetActivePlan(planId, branch)

						if active == nil {
							log.Printf("Active plan not found for planId: %s, branch: %s\n", planId, branch)
							continue mainLoop
						}

						active.Finish()
					} else {
						log.Println("Plan is still building")
						log.Println("Updating status to building")
						err := db.SetPlanStatus(planId, branch, shared.PlanStatusBuilding, "")
						if err != nil {
							log.Printf("Error setting plan status to building: %v\n", err)
							active.StreamDoneCh <- &shared.ApiError{
								Type:   shared.ApiErrorTypeOther,
								Status: http.StatusInternalServerError,
								Msg:    "Error setting plan status to building",
							}
							continue mainLoop
						}

						log.Println("Sending RepliesFinished stream message")
						active.Stream(shared.StreamMessage{
							Type: shared.StreamMessageRepliesFinished,
						})

					}
				}

				// Reset the timer for the usage chunk
				if !timer.Stop() {
					<-timer.C
				}
				timer.Reset(model.OPENAI_USAGE_CHUNK_TIMEOUT)
				streamFinished = true
				continue
			}

			chunksReceived++
			delta := choice.Delta
			content := delta.Content

			// log.Printf("content: %s\n", content)

			if missingFileResponse != "" {
				if maybeRedundantBacktickContent != "" {
					if strings.Contains(content, "\n") {
						maybeRedundantBacktickContent = ""
					} else {
						maybeRedundantBacktickContent += content
					}
					continue // skip processing this chunk
				} else if chunksReceived < 3 && strings.Contains(content, "```") {
					// received closing triple backticks in first 3 chunks after missing file response
					// means this is a redundant start of a new file block, so just ignore it

					maybeRedundantBacktickContent += content
					continue // skip processing this chunk
				}
			}

			// log.Printf("Adding chunk to parser: %s\n", content)
			// log.Printf("fileOpen: %v\n", fileOpen)

			replyParser.AddChunk(content, true)
			parserRes := replyParser.Read()

			if !fileOpen && parserRes.CurrentFilePath != "" {
				log.Printf("File open: %s\n", parserRes.CurrentFilePath)
				fileOpen = true
			}

			if fileOpen && strings.HasSuffix(content, "```") {
				log.Println("FinishAndRead because of closing backticks")
				parserRes = replyParser.FinishAndRead()
				fileOpen = false
			}

			if fileOpen && parserRes.CurrentFilePath == "" {
				// log.Println("File open but current file path is empty, closing file")
				fileOpen = false
			}

			files := parserRes.Files
			fileContents := parserRes.FileContents
			state.replyNumTokens = parserRes.TotalTokens
			currentFile := parserRes.CurrentFilePath
			fileDescriptions := parserRes.FileDescriptions

			// log.Printf("currentFile: %s\n", currentFile)
			// log.Println("files:")
			// spew.Dump(files)

			// Handle file that is present in project paths but not in context
			// Prompt user for what to do on the client side, stop the stream, and wait for user response before proceeding
			if currentFile != "" &&
				!req.IsChatOnly &&
				active.ContextsByPath[currentFile] == nil &&
				req.ProjectPaths[currentFile] && !active.AllowOverwritePaths[currentFile] {
				log.Printf("Attempting to overwrite a file that isn't in context: %s\n", currentFile)

				// attempting to overwrite a file that isn't in context
				// we will stop the stream and ask the user what to do
				err := db.SetPlanStatus(planId, branch, shared.PlanStatusMissingFile, "")

				if err != nil {
					log.Printf("Error setting plan %s status to prompting: %v\n", planId, err)
					active.StreamDoneCh <- &shared.ApiError{
						Type:   shared.ApiErrorTypeOther,
						Status: http.StatusInternalServerError,
						Msg:    "Error setting plan status to prompting",
					}
					continue mainLoop
				}

				var previousReplyContent string
				var trimmedContent string

				UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
					ap.MissingFilePath = currentFile
					previousReplyContent = ap.CurrentReplyContent
					trimmedContent = replyParser.GetReplyForMissingFile()
					ap.CurrentReplyContent = trimmedContent
				})

				log.Println("Previous reply content:")
				log.Println(previousReplyContent)

				// log.Println("Trimmed content:")
				// log.Println(trimmedContent)

				chunkToStream := getCroppedChunk(previousReplyContent+content, trimmedContent, content)

				// log.Printf("chunkToStream: %s\n", chunkToStream)

				if chunkToStream != "" {
					log.Printf("Streaming remaining chunk before missing file prompt: %s\n", chunkToStream)
					active.Stream(shared.StreamMessage{
						Type:       shared.StreamMessageReply,
						ReplyChunk: chunkToStream,
					})
				}

				log.Printf("Prompting user for missing file: %s\n", currentFile)

				active.Stream(shared.StreamMessage{
					Type:                   shared.StreamMessagePromptMissingFile,
					MissingFilePath:        currentFile,
					MissingFileAutoContext: active.AutoContext,
				})

				log.Printf("Stopping stream for missing file: %s\n", currentFile)
				// log.Printf("Chunk content: %s\n", content)
				// log.Printf("Current reply content: %s\n", active.CurrentReplyContent)

				// stop stream for now
				active.CancelModelStreamFn()

				log.Printf("Stopped stream for missing file: %s\n", currentFile)

				// wait for user response to come in
				var userChoice shared.RespondMissingFileChoice
				select {
				case <-active.Ctx.Done():
					log.Println("Context cancelled while waiting for missing file response")
					execHookOnStop(true)
					return
				case <-time.After(30 * time.Minute): // long timeout here since we're waiting for user input
					log.Println("Timeout waiting for missing file choice")
					state.onError(fmt.Errorf("timeout waiting for missing file choice"), true, "", "")
					continue mainLoop
				case userChoice = <-active.MissingFileResponseCh:
				}

				log.Printf("User choice for missing file: %s\n", userChoice)

				active.ResetModelCtx()

				UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
					ap.MissingFilePath = ""
					ap.CurrentReplyContent = replyParser.GetReplyForMissingFile()
				})

				log.Println("Continuing stream")

				// continue plan
				execTellPlan(
					clients,
					plan,
					branch,
					auth,
					req,
					iteration, // keep the same iteration
					userChoice,
					false,
					0,
				)
				return
			}

			UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
				ap.CurrentReplyContent += content
				ap.NumTokens++
			})

			// log.Printf("Sending stream msg: %s", content)
			active.Stream(shared.StreamMessage{
				Type:       shared.StreamMessageReply,
				ReplyChunk: content,
			})

			// log.Println("Content:", content)
			// log.Println("Current reply content:", active.CurrentReplyContent)
			// log.Println("Current file:", currentFile)
			// log.Println("files:")
			// spew.Dump(files)
			// log.Println("replyFiles:")
			// spew.Dump(replyFiles)

			if !req.IsChatOnly && len(files) > len(replyFiles) {
				log.Printf("%d new files\n", len(files)-len(replyFiles))

				for i, file := range files {
					if i < len(replyFiles) {
						continue
					}

					log.Printf("Detected file: %s\n", file)

					if req.BuildMode == shared.BuildModeAuto {
						log.Printf("Queuing build for %s\n", file)
						// log.Println("Content:")
						// log.Println(fileContents[i])

						buildState := &activeBuildStreamState{
							tellState:     state,
							clients:       clients,
							auth:          auth,
							currentOrgId:  currentOrgId,
							currentUserId: currentUserId,
							plan:          plan,
							branch:        branch,
							settings:      settings,
							modelContext:  state.modelContext,
						}

						fileContentTokens, err := shared.GetNumTokens(fileContents[i])

						if err != nil {
							log.Printf("Error getting num tokens for file %s: %v\n", file, err)
							state.onError(fmt.Errorf("error getting num tokens for file %s: %v", file, err), true, "", "")
							continue mainLoop
						}

						buildState.queueBuilds([]*types.ActiveBuild{{
							ReplyId:           replyId,
							Idx:               i,
							FileDescription:   fileDescriptions[i],
							FileContent:       fileContents[i],
							FileContentTokens: fileContentTokens,
							Path:              file,
						}})
					}
					replyFiles = append(replyFiles, file)
					UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
						ap.Files = append(ap.Files, file)
					})
				}
			}
		}
	}
}

func (state *activeTellStreamState) storeAssistantReply() (*db.ConvoMessage, string, error) {
	currentOrgId := state.currentOrgId
	currentUserId := state.currentUserId
	planId := state.plan.Id
	branch := state.branch
	auth := state.auth
	replyNumTokens := state.replyNumTokens
	replyId := state.replyId
	convo := state.convo

	num := len(convo) + 1

	log.Printf("storing assistant reply | len(convo) %d | num %d\n", len(convo), num)

	activePlan := GetActivePlan(planId, branch)

	if activePlan == nil {
		return nil, "", fmt.Errorf("active plan not found")
	}

	assistantMsg := db.ConvoMessage{
		Id:      replyId,
		OrgId:   currentOrgId,
		PlanId:  planId,
		UserId:  currentUserId,
		Role:    openai.ChatMessageRoleAssistant,
		Tokens:  replyNumTokens,
		Num:     num,
		Message: activePlan.CurrentReplyContent,
	}

	commitMsg, err := db.StoreConvoMessage(&assistantMsg, auth.User.Id, branch, false)

	if err != nil {
		log.Printf("Error storing assistant message: %v\n", err)
		return nil, "", err
	}

	UpdateActivePlan(planId, branch, func(ap *types.ActivePlan) {
		ap.MessageNum = num
		ap.StoredReplyIds = append(ap.StoredReplyIds, replyId)
	})

	convo = append(convo, &assistantMsg)
	state.convo = convo

	return &assistantMsg, commitMsg, err
}

func (state *activeTellStreamState) onError(streamErr error, storeDesc bool, convoMessageId, commitMsg string) {
	log.Printf("\nStream error: %v\n", streamErr)

	planId := state.plan.Id
	branch := state.branch
	currentOrgId := state.currentOrgId
	summarizedToMessageId := state.summarizedToMessageId

	active := GetActivePlan(planId, branch)

	if active == nil {
		log.Printf("tellStream onError - Active plan not found for plan ID %s on branch %s\n", planId, branch)
		return
	}

	storeDescAndReply := func() error {
		ctx, cancelFn := context.WithCancel(context.Background())

		repoLockId, err := db.LockRepo(
			db.LockRepoParams{
				UserId:   state.currentUserId,
				OrgId:    state.currentOrgId,
				PlanId:   planId,
				Branch:   branch,
				Scope:    db.LockScopeWrite,
				Ctx:      ctx,
				CancelFn: cancelFn,
			},
		)

		if err != nil {
			log.Printf("Error locking repo for plan %s: %v\n", planId, err)
			return err
		} else {

			defer func() {
				err := db.DeleteRepoLock(repoLockId)
				if err != nil {
					log.Printf("Error unlocking repo for plan %s: %v\n", planId, err)
				}
			}()

			err := db.GitClearUncommittedChanges(state.currentOrgId, planId)
			if err != nil {
				log.Printf("Error clearing uncommitted changes for plan %s: %v\n", planId, err)
				return err
			}
		}

		storedMessage := false
		storedDesc := false

		if convoMessageId == "" {
			assistantMsg, msg, err := state.storeAssistantReply()
			if err == nil {
				convoMessageId = assistantMsg.Id
				commitMsg = msg
				storedMessage = true
			} else {
				log.Printf("Error storing assistant message after stream error: %v\n", err)
				return err
			}
		}

		if storeDesc && convoMessageId != "" {
			err := db.StoreDescription(&db.ConvoMessageDescription{
				OrgId:                 currentOrgId,
				PlanId:                planId,
				SummarizedToMessageId: summarizedToMessageId,
				MadePlan:              false,
				ConvoMessageId:        convoMessageId,
				BuildPathsInvalidated: map[string]bool{},
				Error:                 streamErr.Error(),
			})
			if err == nil {
				storedDesc = true
			} else {
				log.Printf("Error storing description after stream error: %v\n", err)
				return err
			}
		}

		if storedMessage || storedDesc {
			err := db.GitAddAndCommit(currentOrgId, planId, branch, commitMsg)
			if err != nil {
				log.Printf("Error committing after stream error: %v\n", err)
				return err
			}
		}

		return nil
	}

	storeDescAndReply()

	active.StreamDoneCh <- &shared.ApiError{
		Type:   shared.ApiErrorTypeOther,
		Status: http.StatusInternalServerError,
		Msg:    "Stream error: " + streamErr.Error(),
	}
}

func getCroppedChunk(uncropped, cropped, chunk string) string {
	uncroppedIdx := strings.Index(uncropped, chunk)
	if uncroppedIdx == -1 {
		return ""
	}
	croppedChunk := cropped[uncroppedIdx:]
	return croppedChunk
}
