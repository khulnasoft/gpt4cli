package types

import (
	shared "gpt4cli-shared"
)

type ChangesWithLineNums struct {
	Comments []struct {
		Txt       string `json:"txt"`
		Reference bool   `json:"reference"`
	}
	Changes []*shared.StreamedChangeWithLineNums `json:"changes"`
}
