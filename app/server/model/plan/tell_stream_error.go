package plan

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"gpt4cli-server/db"
	"strconv"
	"time"

	shared "gpt4cli-shared"
)

type onErrorParams struct {
	streamErr      error
	streamApiErr   *shared.ApiError
	storeDesc      bool
	convoMessageId string
	commitMsg      string
	canRetry       bool
}

type onErrorResult struct {
	shouldContinueMainLoop bool
	shouldReturn           bool
}

func (state *activeTellStreamState) onError(params onErrorParams) onErrorResult {
	log.Printf("\nStream error: %v\n", params.streamErr)
	streamErr := params.streamErr
	storeDesc := params.storeDesc
	convoMessageId := params.convoMessageId
	commitMsg := params.commitMsg

	planId := state.plan.Id
	branch := state.branch
	currentOrgId := state.currentOrgId
	summarizedToMessageId := state.summarizedToMessageId

	active := GetActivePlan(planId, branch)
	numRetries := state.execTellPlanParams.numErrorRetry

	if active == nil {
		log.Printf("tellStream onError - Active plan not found for plan ID %s on branch %s\n", planId, branch)
		return onErrorResult{
			shouldReturn: true,
		}
	}

	canRetry := params.canRetry

	if canRetry {
		if numRetries >= NumTellStreamRetries {
			log.Printf("tellStream onError - Max retries reached for plan ID %s on branch %s\n", planId, branch)

			canRetry = false
		}
	}

	if canRetry {
		// stop stream via context (ensures we stop child streams too)
		active.CancelModelStreamFn()

		active.ResetModelCtx()

		retryDelaySeconds := 1 * numRetries * (numRetries / 2)

		log.Printf("tellStream onError - Retry %d/%d - Retrying stream in %d seconds", numRetries+1, NumTellStreamRetries, retryDelaySeconds)
		time.Sleep(time.Duration(retryDelaySeconds) * time.Second)

		params := state.execTellPlanParams
		params.numErrorRetry = numRetries + 1

		execTellPlan(params)
		return onErrorResult{
			shouldReturn: true,
		}
	}

	storeDescAndReply := func() error {
		ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)

		err := db.ExecRepoOperation(db.ExecRepoOperationParams{
			OrgId:    currentOrgId,
			UserId:   state.currentUserId,
			PlanId:   planId,
			Branch:   branch,
			Scope:    db.LockScopeWrite,
			Ctx:      ctx,
			CancelFn: cancelFn,
			Reason:   "store desc and reply",
		}, func(repo *db.GitRepo) error {
			storedMessage := false
			storedDesc := false

			if convoMessageId == "" {
				hasUnfinishedSubtasks := false
				for _, subtask := range state.subtasks {
					if !subtask.IsFinished {
						hasUnfinishedSubtasks = true
						break
					}
				}

				assistantMsg, msg, err := state.storeAssistantReply(repo, storeAssistantReplyParams{
					flags: shared.ConvoMessageFlags{
						CurrentStage:          state.currentStage,
						HasUnfinishedSubtasks: hasUnfinishedSubtasks,
					},
					subtask:       nil,
					addedSubtasks: nil,
				})
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
					WroteFiles:            false,
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
				err := repo.GitAddAndCommit(branch, commitMsg)
				if err != nil {
					log.Printf("Error committing after stream error: %v\n", err)
					return err
				}
			}

			return nil
		})

		if err != nil {
			log.Printf("Error storing description and reply after stream error: %v\n", err)
			return err
		}

		return nil
	}

	storeDescAndReply() // best effort to store description and reply, ignore errors

	if params.streamApiErr != nil {
		active.StreamDoneCh <- params.streamApiErr
	} else {
		msg := "Stream error: " + streamErr.Error()
		if params.canRetry && numRetries >= NumTellStreamRetries {
			msg += " | Failed after " + strconv.Itoa(numRetries) + " retries."
		}

		active.StreamDoneCh <- &shared.ApiError{
			Type:   shared.ApiErrorTypeOther,
			Status: http.StatusInternalServerError,
			Msg:    msg,
		}
	}

	return onErrorResult{
		shouldContinueMainLoop: true,
	}
}

func (state *activeTellStreamState) onActivePlanMissingError() {
	planId := state.plan.Id
	branch := state.branch
	log.Printf("Active plan not found for plan ID %s on branch %s\n", planId, branch)
	state.onError(onErrorParams{
		streamErr: fmt.Errorf("active plan not found for plan ID %s on branch %s", planId, branch),
		storeDesc: true,
	})
}
