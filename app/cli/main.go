package main

import (
	"gpt4cli/api"
	"gpt4cli/auth"
	"gpt4cli/cmd"
	"gpt4cli/fs"
	"gpt4cli/lib"
	"gpt4cli/plan_exec"
	"gpt4cli/term"
	"log"
	"os"
	"path/filepath"

	"github.com/khulnasoft/gpt4cli/shared"
)

func init() {
	// inter-package dependency injections to avoid circular imports
	auth.SetApiClient(api.Client)
	lib.SetBuildPlanInlineFn(func(maybeContexts []*shared.Context) (bool, error) {
		apiKeys := lib.MustVerifyApiKeys()
		return plan_exec.Build(plan_exec.ExecParams{
			CurrentPlanId: lib.CurrentPlanId,
			CurrentBranch: lib.CurrentBranch,
			ApiKeys:       apiKeys,
			CheckOutdatedContext: func(maybeContexts []*shared.Context) (bool, bool) {
				return lib.MustCheckOutdatedContext(true, maybeContexts)
			},
		}, false)
	})

	// set up a file logger
	// TODO: log rotation

	file, err := os.OpenFile(filepath.Join(fs.HomeGpt4cliDir, "gpt4cli.log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		term.OutputErrorAndExit("Error opening log file: %v", err)
	}

	// Set the output of the logger to the file
	log.SetOutput(file)

	// log.Println("Starting Gpt4cli - logging initialized")
}

func main() {
	checkForUpgrade()

	// Manually check for help flags at the root level
	if len(os.Args) == 2 && (os.Args[1] == "-h" || os.Args[1] == "--help") {
		// Display your custom help here
		term.PrintCustomHelp(true)
		os.Exit(0)
	}

	cmd.Execute()
}
