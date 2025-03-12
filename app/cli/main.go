package main

import (
	"log"
	"os"
	"path/filepath"
	"gpt4cli-cli/api"
	"gpt4cli-cli/auth"
	"gpt4cli-cli/cmd"
	"gpt4cli-cli/fs"
	"gpt4cli-cli/lib"
	"gpt4cli-cli/plan_exec"
	"gpt4cli-cli/term"
	"gpt4cli-cli/types"
	"gpt4cli-cli/ui"

	shared "gpt4cli-shared"

	"gopkg.in/natefinch/lumberjack.v2"
)

func init() {
	// inter-package dependency injections to avoid circular imports
	auth.SetApiClient(api.Client)

	auth.SetOpenUnauthenticatedCloudURLFn(ui.OpenUnauthenticatedCloudURL)
	auth.SetOpenAuthenticatedURLFn(ui.OpenAuthenticatedURL)

	term.SetOpenAuthenticatedURLFn(ui.OpenAuthenticatedURL)
	term.SetOpenUnauthenticatedCloudURLFn(ui.OpenUnauthenticatedCloudURL)
	term.SetConvertTrialFn(auth.ConvertTrial)

	lib.SetBuildPlanInlineFn(func(autoConfirm bool, maybeContexts []*shared.Context) (bool, error) {
		var apiKeys map[string]string
		if !auth.Current.IntegratedModelsMode {
			apiKeys = lib.MustVerifyApiKeys()
		}
		return plan_exec.Build(plan_exec.ExecParams{
			CurrentPlanId: lib.CurrentPlanId,
			CurrentBranch: lib.CurrentBranch,
			ApiKeys:       apiKeys,
			CheckOutdatedContext: func(maybeContexts []*shared.Context, projectPaths *types.ProjectPaths) (bool, bool, error) {
				return lib.CheckOutdatedContextWithOutput(true, autoConfirm, maybeContexts, projectPaths)
			},
		}, types.BuildFlags{})
	})

	// set up a rotating file logger
	logger := &lumberjack.Logger{
		Filename:   filepath.Join(fs.HomeGpt4cliDir, "gpt4cli.log"),
		MaxSize:    10,   // megabytes before rotation
		MaxBackups: 3,    // number of backups to keep
		MaxAge:     28,   // days to keep old logs
		Compress:   true, // compress rotated files
	}

	// Set the output of the logger
	log.SetOutput(logger)
	log.SetFlags(log.LstdFlags | log.Lmicroseconds | log.Lshortfile)

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
