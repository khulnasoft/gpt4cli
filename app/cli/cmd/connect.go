package cmd

import (
	"fmt"
	"gpt4cli/api"
	"gpt4cli/auth"
	"gpt4cli/lib"
	"gpt4cli/stream"
	streamtui "gpt4cli/stream_tui"
	"gpt4cli/term"
	"os"

	"github.com/spf13/cobra"
)

var connectCmd = &cobra.Command{
	Use:     "connect [stream-id-or-plan] [branch]",
	Aliases: []string{"conn"},
	Short:   "Connect to an active stream",
	// Long:  ``,
	Args: cobra.MaximumNArgs(2),
	Run:  connect,
}

func init() {
	RootCmd.AddCommand(connectCmd)

}

func connect(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputNoCurrentPlanErrorAndExit()
	}

	planId, branch, shouldContinue := lib.SelectActiveStream(args)

	if !shouldContinue {
		return
	}

	term.StartSpinner("")
	apiErr := api.Client.ConnectPlan(planId, branch, stream.OnStreamPlan)
	term.StopSpinner()

	if apiErr != nil {
		term.OutputErrorAndExit("Error connecting to stream: %v", apiErr)
	}

	go func() {
		err := streamtui.StartStreamUI("", false)

		if err != nil {
			term.OutputErrorAndExit("Error starting stream UI", err)
		}

		fmt.Println()
		term.PrintCmds("", "changes", "diff", "apply", "reject", "log")

		os.Exit(0)
	}()

	// Wait for the stream to finish
	select {}
}
