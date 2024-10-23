package cmd

import (
	"fmt"
	"gpt4cli/auth"
	"gpt4cli/lib"
	"gpt4cli/plan_exec"
	"gpt4cli/term"

	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/spf13/cobra"
)

var buildBg bool

var buildCmd = &cobra.Command{
	Use:     "build",
	Aliases: []string{"b"},
	Short:   "Build pending changes",
	// Long:  ``,
	Args: cobra.NoArgs,
	Run:  build,
}

func init() {
	RootCmd.AddCommand(buildCmd)
	buildCmd.Flags().BoolVar(&buildBg, "bg", false, "Execute autonomously in the background")
}

func build(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputNoCurrentPlanErrorAndExit()
	}

	apiKeys := lib.MustVerifyApiKeys()

	didBuild, err := plan_exec.Build(plan_exec.ExecParams{
		CurrentPlanId: lib.CurrentPlanId,
		CurrentBranch: lib.CurrentBranch,
		ApiKeys:       apiKeys,
		CheckOutdatedContext: func(maybeContexts []*shared.Context) (bool, bool) {
			return lib.MustCheckOutdatedContext(false, maybeContexts)
		},
	}, buildBg)

	if err != nil {
		term.OutputErrorAndExit("Error building plan: %v", err)
	}

	if !didBuild {
		fmt.Println()
		term.PrintCmds("", "log", "tell", "continue")
		return
	}

	if buildBg {
		fmt.Println("🏗️ Building plan in the background")
		fmt.Println()
		term.PrintCmds("", "ps", "connect", "stop")
	} else {
		fmt.Println()
		term.PrintCmds("", "changes", "apply", "reject", "log")
	}
}
