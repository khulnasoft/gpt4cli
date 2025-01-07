package cmd

import (
	"fmt"
	"gpt4cli/auth"
	"gpt4cli/lib"
	"gpt4cli/plan_exec"

	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/spf13/cobra"
)

var chatCmd = &cobra.Command{
	Use:     "chat [prompt]",
	Aliases: []string{"ch"},
	Short:   "Chat without making changes",
	// Long:  ``,
	Args: cobra.RangeArgs(0, 1),
	Run:  doChat,
}

func init() {
	RootCmd.AddCommand(chatCmd)

	initExecFlags(chatCmd, initExecFlagsParams{
		omitNoBuild: true,
		omitStop:    true,
		omitBg:      true,
		omitApply:   true,
		omitExec:    true,
	})

}

func doChat(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()
	mustSetPlanExecFlags(cmd)

	var apiKeys map[string]string
	if !auth.Current.IntegratedModelsMode {
		apiKeys = lib.MustVerifyApiKeys()
	}

	prompt := getTellPrompt(args)

	if prompt == "" {
		fmt.Println("🤷‍♂️ No prompt to send")
		return
	}

	plan_exec.TellPlan(plan_exec.ExecParams{
		CurrentPlanId: lib.CurrentPlanId,
		CurrentBranch: lib.CurrentBranch,
		ApiKeys:       apiKeys,
		CheckOutdatedContext: func(maybeContexts []*shared.Context) (bool, bool, error) {
			return lib.CheckOutdatedContextWithOutput(false, tellAutoContext, maybeContexts)
		},
	}, prompt, plan_exec.TellFlags{
		IsChatOnly:  true,
		AutoContext: tellAutoContext,
	})
}
