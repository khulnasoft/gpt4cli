package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"gpt4cli/auth"
	"gpt4cli/lib"
	"gpt4cli/plan_exec"
	"gpt4cli/term"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/spf13/cobra"
)

const DebugDefaultTries = 5

var debugCmd = &cobra.Command{
	Use:     "debug [tries] <cmd>",
	Aliases: []string{"db"},
	Short:   "Debug a failing command with Gpt4cli",
	Args:    cobra.MinimumNArgs(1),
	Run:     doDebug,
}

func init() {
	RootCmd.AddCommand(debugCmd)
	debugCmd.Flags().BoolVarP(&autoCommit, "commit", "c", false, "Commit changes to git on each try")
}

func doDebug(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MustResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputNoCurrentPlanErrorAndExit()
	}

	// Parse tries and command
	tries := DebugDefaultTries
	cmdArgs := args

	// Check if first arg is tries count
	if val, err := strconv.Atoi(args[0]); err == nil {
		if val <= 0 {
			term.OutputErrorAndExit("Tries must be greater than 0")
		}
		tries = val
		cmdArgs = args[1:]
		if len(cmdArgs) == 0 {
			term.OutputErrorAndExit("No command specified")
		}
	}

	var apiKeys map[string]string
	if !auth.Current.IntegratedModelsMode {
		apiKeys = lib.MustVerifyApiKeys()
	}

	// Get current working directory
	cwd, err := os.Getwd()
	if err != nil {
		term.OutputErrorAndExit("Failed to get working directory: %v", err)
	}

	cmdStr := strings.Join(cmdArgs, " ")

	// Execute command and handle retries
	for attempt := 0; attempt < tries; attempt++ {
		term.StartSpinner("")
		// Use shell to handle operators like && and |
		execCmd := exec.Command("sh", "-c", cmdStr)
		execCmd.Dir = cwd
		execCmd.Env = os.Environ()

		output, err := execCmd.CombinedOutput()

		term.StopSpinner()

		outputStr := string(output)
		if outputStr == "" && err != nil {
			// If no output but error occurred, include error in output
			outputStr = err.Error()
		}

		if outputStr != "" {
			fmt.Println(outputStr)
		}

		if err == nil {
			if attempt == 0 {
				fmt.Printf("✅ Command %s succeeded on first try\n", color.New(color.Bold, term.ColorHiCyan).Sprintf(cmdStr))
			} else {
				lbl := "attempts"
				if attempt == 1 {
					lbl = "attempt"
				}
				fmt.Printf("✅ Command %s succeeded after %d fix %s\n", color.New(color.Bold, term.ColorHiCyan).Sprintf(cmdStr), attempt, lbl)
			}
			return
		}

		if attempt == tries-1 {
			fmt.Printf("Command failed after %d tries\n", tries)
			os.Exit(1)
		}

		// Prepare prompt for TellPlan
		exitErr, ok := err.(*exec.ExitError)
		status := -1
		if ok {
			status = exitErr.ExitCode()
		}

		prompt := fmt.Sprintf("'%s' failed with exit status %d. Output:\n\n%s\n\n--\n\n",
			strings.Join(cmdArgs, " "), status, string(output))

		plan_exec.TellPlan(plan_exec.ExecParams{
			CurrentPlanId: lib.CurrentPlanId,
			CurrentBranch: lib.CurrentBranch,
			ApiKeys:       apiKeys,
			CheckOutdatedContext: func(maybeContexts []*shared.Context) (bool, bool, error) {
				return lib.CheckOutdatedContextWithOutput(false, true, maybeContexts)
			},
		}, prompt, plan_exec.TellFlags{IsUserDebug: true})

		flags := lib.ApplyFlags{
			AutoConfirm: true,
			AutoCommit:  autoCommit,
			NoCommit:    !autoCommit,
			NoExec:      true,
		}

		lib.MustApplyPlan(
			lib.CurrentPlanId,
			lib.CurrentBranch,
			flags,
			plan_exec.GetOnApplyExecFail(flags),
		)
	}
}
