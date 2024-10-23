package lib

import "github.com/khulnasoft/gpt4cli/shared"

var buildPlanInlineFn func(maybeContexts []*shared.Context) (bool, error)

func SetBuildPlanInlineFn(fn func(maybeContexts []*shared.Context) (bool, error)) {
	buildPlanInlineFn = fn
}
