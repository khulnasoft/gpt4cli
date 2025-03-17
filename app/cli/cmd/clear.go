package cmd

import (
	"fmt"
	"gpt4cli-cli/api"
	"gpt4cli-cli/auth"
	"gpt4cli-cli/lib"
	"gpt4cli-cli/term"

	shared "gpt4cli-shared"

	"github.com/spf13/cobra"
)

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all context",
	Long:  `Clear all context.`,
	Run:   clearAllContext,
}

func clearAllContext(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputNoCurrentPlanErrorAndExit()
	}

	term.StartSpinner("")
	contexts, err := api.Client.ListContext(lib.CurrentPlanId, lib.CurrentBranch)
	term.StopSpinner()

	if err != nil {
		term.OutputErrorAndExit("Error retrieving context: %v", err)
	}

	deleteIds := map[string]bool{}

	for _, context := range contexts {
		deleteIds[context.Id] = true
	}

	if len(deleteIds) > 0 {
		res, err := api.Client.DeleteContext(lib.CurrentPlanId, lib.CurrentBranch, shared.DeleteContextRequest{
			Ids: deleteIds,
		})

		if err != nil {
			term.OutputErrorAndExit("Error deleting context: %v", err)
		}

		fmt.Println("✅ " + res.Msg)
	} else {
		fmt.Println("🤷‍♂️ No context removed")
	}

}

func init() {
	RootCmd.AddCommand(clearCmd)
}
