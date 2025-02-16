package plan

import (
	"gpt4cli-server/db"
	"gpt4cli-server/types"

	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/sashabaranov/go-openai"
)

type activeTellStreamState struct {
	clients                map[string]*openai.Client
	req                    *shared.TellPlanRequest
	auth                   *types.ServerAuth
	currentOrgId           string
	currentUserId          string
	plan                   *db.Plan
	branch                 string
	iteration              int
	replyId                string
	modelContext           []*db.Context
	convo                  []*db.ConvoMessage
	currentPlanState       *shared.CurrentPlanState
	missingFileResponse    shared.RespondMissingFileChoice
	summaries              []*db.ConvoSummary
	summarizedToMessageId  string
	latestSummaryTokens    int
	userPrompt             string
	promptMessage          *openai.ChatCompletionMessage
	replyParser            *types.ReplyParser
	replyNumTokens         int
	messages               []openai.ChatCompletionMessage
	tokensBeforeConvo      int
	totalRequestTokens     int
	settings               *shared.PlanSettings
	currentReplyNumRetries int
	subtasks               []*db.Subtask
	currentSubtask         *db.Subtask
}
