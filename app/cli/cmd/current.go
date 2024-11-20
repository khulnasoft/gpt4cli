package cmd

import (
	"fmt"
	"gpt4cli/api"
	"gpt4cli/auth"
	"gpt4cli/format"
	"gpt4cli/lib"
	"gpt4cli/term"
	"os"
	"strconv"

	"github.com/fatih/color"
	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:     "current",
	Aliases: []string{"cu"},
	Short:   "Get the current plan",
	Run:     current,
}

func init() {
	RootCmd.AddCommand(currentCmd)
}

func current(cmd *cobra.Command, args []string) {
	auth.MustResolveAuthWithOrg()
	lib.MaybeResolveProject()

	if lib.CurrentPlanId == "" {
		term.OutputNoCurrentPlanErrorAndExit()
	}

	term.StartSpinner("")
	plan, err := api.Client.GetPlan(lib.CurrentPlanId)
	term.StopSpinner()

	if err != nil {
		term.OutputErrorAndExit("Error getting plan: %v", err)
		return
	}

	currentBranchesByPlanId, err := api.Client.GetCurrentBranchByPlanId(lib.CurrentProjectId, shared.GetCurrentBranchByPlanIdRequest{
		CurrentBranchByPlanId: map[string]string{
			lib.CurrentPlanId: lib.CurrentBranch,
		},
	})

	if err != nil {
		term.OutputErrorAndExit("Error getting current branches: %v", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetAutoWrapText(false)
	table.SetHeader([]string{"Current Plan", "Updated", "Created" /*"Branches",*/, "Branch", "Context", "Convo"})

	name := color.New(color.Bold, term.ColorHiGreen).Sprint(plan.Name)
	branch := currentBranchesByPlanId[lib.CurrentPlanId]

	row := []string{
		name,
		format.Time(plan.UpdatedAt),
		format.Time(plan.CreatedAt),
		// strconv.Itoa(plan.ActiveBranches),
		lib.CurrentBranch,
		strconv.Itoa(branch.ContextTokens) + " 🪙",
		strconv.Itoa(branch.ConvoTokens) + " 🪙",
	}

	style := []tablewriter.Colors{
		{tablewriter.FgGreenColor, tablewriter.Bold},
	}

	table.Rich(row, style)

	table.Render()
	fmt.Println()
	term.PrintCmds("", "tell", "ls", "plans")

}
