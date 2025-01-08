package lib

import "github.com/khulnasoft/gpt4cli/shared"

var buildPlanInlineFn func(autoConfirm bool, maybeContexts []*shared.Context) (bool, error)

func SetBuildPlanInlineFn(fn func(autoConfirm bool, maybeContexts []*shared.Context) (bool, error)) {
	buildPlanInlineFn = fn
}
