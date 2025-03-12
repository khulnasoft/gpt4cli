package cmd

import (
	"fmt"
	"os"
	"gpt4cli-cli/term"
	"strings"

	"github.com/fatih/color"
	"github.com/spf13/cobra"
)

var helpShowAll bool

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use: `gpt4cli [command] [flags]`,
	// Short: "Gpt4cli: iterative development with AI",
	SilenceErrors: true,
	SilenceUsage:  true,
	Run: func(cmd *cobra.Command, args []string) {
		run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	// if no arguments were passed, start the repl
	if len(os.Args) == 1 ||
		(len(os.Args) == 2 && strings.HasPrefix(os.Args[1], "--") && os.Args[1] != "--help") ||
		(len(os.Args) == 3 && strings.HasPrefix(os.Args[1], "--") && os.Args[1] != "--help" && strings.HasPrefix(os.Args[2], "--") && os.Args[2] != "--help") {

		// Instead of directly calling replCmd.Run, parse the flags first
		replCmd.ParseFlags(os.Args[1:])
		replCmd.Run(replCmd, []string{})
		return
	}

	if err := RootCmd.Execute(); err != nil {
		// term.OutputErrorAndExit("Error executing root command: %v", err)
		// log.Fatalf("Error executing root command: %v", err)

		// output the error message to stderr
		term.OutputSimpleError("Error: %v", err)

		fmt.Println()

		color.New(color.Bold, color.BgGreen, color.FgHiWhite).Println(" Usage ")
		color.New(color.Bold).Println("  gpt4cli [command] [flags]")
		color.New(color.Bold).Println("  g4c [command] [flags]")
		fmt.Println()

		color.New(color.Bold, color.BgGreen, color.FgHiWhite).Println(" Help ")
		color.New(color.Bold).Println("  gpt4cli help # show basic usage")
		color.New(color.Bold).Println("  gpt4cli help --all # show all commands")
		color.New(color.Bold).Println("  gpt4cli [command] --help")
		fmt.Println()

		color.New(color.Bold, color.BgGreen, color.FgHiWhite).Println(" Common Commands ")
		color.New(color.Bold).Println("  gpt4cli new # create a new plan")
		color.New(color.Bold).Println("  gpt4cli tell # tell the plan what to do")
		color.New(color.Bold).Println("  gpt4cli continue # continue the current plan")
		color.New(color.Bold).Println("  gpt4cli settings # show plan settings")
		color.New(color.Bold).Println("  gpt4cli set # update plan settings")
		fmt.Println()

		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) {

}

func init() {
	var helpCmd = &cobra.Command{
		Use:     "help",
		Aliases: []string{"h"},
		Short:   "Display help for Gpt4cli",
		Long:    `Display help for Gpt4cli.`,
		Run: func(cmd *cobra.Command, args []string) {
			term.PrintCustomHelp(helpShowAll)
		},
	}

	RootCmd.AddCommand(helpCmd)

	// add an --all/-a flag
	helpCmd.Flags().BoolVarP(&helpShowAll, "all", "a", false, "Show all commands")
}
