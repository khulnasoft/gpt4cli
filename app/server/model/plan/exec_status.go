package plan

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"gpt4cli-server/model"
	"gpt4cli-server/model/prompts"
	"gpt4cli-server/types"
	"gpt4cli-server/utils"
	"strings"

	shared "gpt4cli-shared"

	"github.com/sashabaranov/go-openai"
)

// controls the number steps to spent trying to finish a single subtask
// if a subtask is not finished in this number of steps, we'll give up and mark it done
// necessary to prevent infinite loops
const MaxPreviousMessages = 4

type execStatusShouldContinueResult struct {
	subtaskFinished bool
}

func (state *activeTellStreamState) execStatusShouldContinue(currentMessage string, ctx context.Context) (execStatusShouldContinueResult, *shared.ApiError) {
	auth := state.auth
	plan := state.plan
	settings := state.settings
	clients := state.clients
	config := settings.ModelPack.ExecStatus

	// Check subtask completion
	if state.currentSubtask != nil {
		completionMarker := fmt.Sprintf("**%s** has been completed", state.currentSubtask.Title)
		log.Printf("[ExecStatus] Checking for subtask completion marker: %q", completionMarker)
		log.Printf("[ExecStatus] Current subtask: %q", state.currentSubtask.Title)

		if strings.Contains(currentMessage, completionMarker) {
			log.Printf("[ExecStatus] ✓ Subtask completion marker found")
			return execStatusShouldContinueResult{
				subtaskFinished: true,
			}, nil

			// NOTE: tried using an LLM to verify "suspicious" subtask completions, but in practice led to too many extra LLM calls and disagreement cycles between agent roles (it's finished. no it's note! etc.)
			// now just going back to trusting the completion marker... basically it's better to err on the side of marking tasks done.

			// var potentialProblem bool

			// if len(state.chunkProcessor.replyOperations) == 0 {
			// 	log.Printf("[ExecStatus] ✗ Subtask completion marker found, but there are no operations to execute")
			// 	potentialProblem = true
			// } else {
			// wroteToPaths := map[string]bool{}
			// for _, op := range state.chunkProcessor.replyOperations {
			// 	if op.Type == shared.OperationTypeFile {
			// 		wroteToPaths[op.Path] = true
			// 	}
			// }

			// for _, path := range state.currentSubtask.UsesFiles {
			// 	if !wroteToPaths[path] {
			// 		log.Printf("[ExecStatus] ✗ Subtask completion marker found, but the operations did not write to the file %q from the 'Uses' list", path)
			// 		potentialProblem = true
			// 		break
			// 	}
			// }
			// }

			// if !potentialProblem {
			// 	log.Printf("[ExecStatus] ✓ Subtask completion marker found and no potential problem - will mark as completed")

			// 	return execStatusShouldContinueResult{
			// 		subtaskFinished: true,
			// 	}, nil
			// } else if state.currentSubtask.NumTries >= 1 {
			// 	log.Printf("[ExecStatus] ✓ Subtask completion marker found, but the operations are questionable -- marking it done anyway since it's the second try and we can't risk an infinite loop")

			// 	return execStatusShouldContinueResult{
			// 		subtaskFinished: true,
			// 	}, nil
			// } else {
			// 	log.Printf("[ExecStatus] ✗ Subtask completion marker found, but the operations are questionable -- will verify with LLM call")
			// }
		} else {
			log.Printf("[ExecStatus] ✗ No subtask completion marker found in message")
		}

		log.Println("[ExecStatus] Current subtasks state:")
		for i, task := range state.subtasks {
			log.Printf("[ExecStatus] Task %d: %q (finished=%v)", i+1, task.Title, task.IsFinished)
		}
	}

	log.Println("Checking if plan should continue based on exec status")

	fullSubtask := state.currentSubtask.Title
	fullSubtask += "\n\n" + state.currentSubtask.Description

	previousMessages := []string{}
	for _, msg := range state.convo {
		if msg.Subtask != nil && msg.Subtask.Title == state.currentSubtask.Title {
			previousMessages = append(previousMessages, msg.Message)
		}
	}

	if len(previousMessages) >= MaxPreviousMessages {
		log.Printf("[ExecStatus] ✗ Max previous messages reached - will mark as completed and move on to next subtask")
		return execStatusShouldContinueResult{
			subtaskFinished: true,
		}, nil
	}

	prompt := prompts.GetExecStatusFinishedSubtask(prompts.GetExecStatusFinishedSubtaskParams{
		UserPrompt:                 state.userPrompt,
		CurrentSubtask:             fullSubtask,
		CurrentMessage:             currentMessage,
		PreviousMessages:           previousMessages,
		PreferredModelOutputFormat: config.BaseModelConfig.PreferredModelOutputFormat,
	})

	messages := []types.ExtendedChatMessage{
		{
			Role: openai.ChatMessageRoleSystem,
			Content: []types.ExtendedChatMessagePart{
				{
					Type: openai.ChatMessagePartTypeText,
					Text: prompt,
				},
			},
		},
	}

	modelRes, err := model.ModelRequest(ctx, model.ModelRequestParams{
		Clients:        clients,
		Auth:           auth,
		Plan:           plan,
		ModelConfig:    &config,
		Purpose:        "Check if task finished",
		Messages:       messages,
		ModelStreamId:  state.modelStreamId,
		ConvoMessageId: state.replyId,
	})

	if err != nil {
		log.Printf("[ExecStatus] Error in model call: %v", err)
		return execStatusShouldContinueResult{}, nil
	}

	content := modelRes.Content

	var reasoning string
	var subtaskFinished bool

	if config.BaseModelConfig.PreferredModelOutputFormat == shared.ModelOutputFormatXml {
		reasoning = utils.GetXMLContent(content, "reasoning")
		subtaskFinishedStr := utils.GetXMLContent(content, "subtaskFinished")
		subtaskFinished = subtaskFinishedStr == "true"

		if reasoning == "" || subtaskFinishedStr == "" {
			return execStatusShouldContinueResult{}, &shared.ApiError{
				Type:   shared.ApiErrorTypeOther,
				Status: http.StatusInternalServerError,
				Msg:    "Missing required XML tags in response",
			}
		}
	} else {

		if content == "" {
			log.Printf("[ExecStatus] No function response found in model output")
			return execStatusShouldContinueResult{}, nil
		}

		var res types.ExecStatusResponse
		if err := json.Unmarshal([]byte(content), &res); err != nil {
			log.Printf("[ExecStatus] Failed to parse response: %v", err)
			return execStatusShouldContinueResult{}, nil
		}

		reasoning = res.Reasoning
		subtaskFinished = res.SubtaskFinished
	}

	log.Printf("[ExecStatus] Decision: subtaskFinished=%v, reasoning=%v",
		subtaskFinished, reasoning)

	return execStatusShouldContinueResult{
		subtaskFinished: subtaskFinished,
	}, nil
}
