package plan_exec

import "github.com/gpt4cli/gpt4cli/shared"

type ExecParams struct {
	CurrentPlanId        string
	CurrentBranch        string
	ApiKeys              map[string]string
	CheckOutdatedContext func(maybeContexts []*shared.Context) (bool, bool)
}
