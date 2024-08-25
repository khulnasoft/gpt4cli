package types

import (
	"gpt4cli-server/db"

	"github.com/gpt4cli/gpt4cli/shared"
)

func HasPendingBuilds(planDescs []*db.ConvoMessageDescription) bool {
	apiDescs := make([]*shared.ConvoMessageDescription, len(planDescs))
	for i, desc := range planDescs {
		apiDescs[i] = desc.ToApi()
	}

	return shared.HasPendingBuilds(apiDescs)
}
