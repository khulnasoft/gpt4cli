package cmd

import (
	"gpt4cli/auth"
	"gpt4cli/term"

	"github.com/spf13/cobra"
)

var signInCmd = &cobra.Command{
	Use:   "sign-in",
	Short: "Sign in to a Gpt4cli account",
	Args:  cobra.NoArgs,
	Run:   signIn,
}

func init() {
	RootCmd.AddCommand(signInCmd)

	signInCmd.Flags().String("code", "", "Sign in code from the Gpt4cli web UI")
}

func signIn(cmd *cobra.Command, args []string) {
	code, err := cmd.Flags().GetString("code")
	if err != nil {
		term.OutputErrorAndExit("Error getting code: %v", err)
	}

	if code != "" {
		err = auth.SignInWithCode(code, "")

		if err != nil {
			term.OutputErrorAndExit("Error signing in: %v", err)
		}

		return
	}

	err = auth.SelectOrSignInOrCreate()

	if err != nil {
		term.OutputErrorAndExit("Error signing in: %v", err)
	}
}
