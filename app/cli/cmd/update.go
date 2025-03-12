package cmd

import (
	"fmt"
	"gpt4cli-cli/api"
	"gpt4cli-cli/auth"
	"gpt4cli-cli/fs"
	"gpt4cli-cli/lib"
	"gpt4cli-cli/term"

	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:     "update ",
	Aliases: []string{"u"},
	Short:   "Update outdated context",
	Args:    cobra.MaximumNArgs(1),
	Run:     update,
}

func init() {
	RootCmd.AddCommand(updateCmd)

}

func update(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	term.StartSpinner("")

	contexts, apiErr := api.Client.ListContext(lib.CurrentPlanId, lib.CurrentBranch)

	if apiErr != nil {
		term.StopSpinner()
		term.OutputErrorAndExit("failed to list context: %s", apiErr)
	}

	paths, err := fs.GetProjectPaths(fs.ProjectRoot)

	if err != nil {
		term.OutputErrorAndExit("error getting project paths: %v", err)
	}

	outdated, err := lib.CheckOutdatedContext(contexts, paths)

	if err != nil {
		term.StopSpinner()
		term.OutputErrorAndExit("failed to check outdated context: %s", err)
	}

	if len(outdated.UpdatedContexts) == 0 {
		term.StopSpinner()
		fmt.Println("✅ Context is up to date")
		return
	}

	lib.UpdateContextWithOutput(lib.UpdateContextParams{
		Contexts:    contexts,
		OutdatedRes: *outdated,
		Req:         outdated.Req,
	})
}
