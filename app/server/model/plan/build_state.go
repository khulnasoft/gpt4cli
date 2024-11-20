package plan

import (
	"gpt4cli-server/db"
	"gpt4cli-server/types"

	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/sashabaranov/go-openai"
)

const MaxBuildStreamErrorRetries = 3 // uses semi-exponential backoff so be careful with this

const FixSyntaxRetries = 2
const FixSyntaxEpochs = 2

type activeBuildStreamState struct {
	clients       map[string]*openai.Client
	auth          *types.ServerAuth
	currentOrgId  string
	currentUserId string
	plan          *db.Plan
	branch        string
	settings      *shared.PlanSettings
	modelContext  []*db.Context
	convo         []*db.ConvoMessage
}

type activeBuildStreamFileState struct {
	*activeBuildStreamState
	filePath           string
	convoMessageId     string
	build              *db.PlanBuild
	currentPlanState   *shared.CurrentPlanState
	activeBuild        *types.ActiveBuild
	preBuildState      string
	lineNumsNumRetry   int
	verifyFileNumRetry int
	fixFileNumRetry    int

	syntaxNumRetry int
	syntaxNumEpoch int

	isFixingSyntax bool
	isFixingOther  bool

	streamedChangesWithLineNums []*shared.StreamedChangeWithLineNums
	updated                     string

	verificationErrors string
	syntaxErrors       []string

	isNewFile bool
}
