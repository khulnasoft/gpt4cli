g4c-1: package cmd
g4c-2: 
g4c-3: import (
g4c-4: 	"fmt"
g4c-5: 	"path/filepath"
g4c-6: 	"gpt4cli/api"
g4c-7: 	"gpt4cli/auth"
g4c-8: 	"gpt4cli/lib"
g4c-9: 	"gpt4cli/term"
g4c-10: 
g4c-11: 	"github.com/gpt4cli/gpt4cli/shared"
g4c-12: 	"github.com/spf13/cobra"
g4c-13: )
g4c-14: 
g4c-15: var contextRmCmd = &cobra.Command{
g4c-16: 	Use:     "rm",
g4c-17: 	Aliases: []string{"remove", "unload"},
g4c-18: 	Short:   "Remove context",
g4c-19: 	Long:    `Remove context by index, name, or glob.`,
g4c-20: 	Args:    cobra.MinimumNArgs(1),
g4c-21: 	Run:     contextRm,
g4c-22: }
g4c-23: 
g4c-24: func contextRm(cmd *cobra.Command, args []string) {
g4c-25: 	auth.MustResolveAuthWithOrg()
g4c-26: 	lib.MustResolveProject()
g4c-27: 
g4c-28: 	if lib.CurrentPlanId == "" {
g4c-29: 		fmt.Println("🤷‍♂️ No current plan")
g4c-30: 		return
g4c-31: 	}
g4c-32: 
g4c-33: 	term.StartSpinner("")
g4c-34: 	contexts, err := api.Client.ListContext(lib.CurrentPlanId, lib.CurrentBranch)
g4c-35: 
g4c-36: 	if err != nil {
g4c-37: 		term.OutputErrorAndExit("Error retrieving context: %v", err)
g4c-38: 	}
g4c-39: 
g4c-40: 	deleteIds := map[string]bool{}
g4c-41: 
g4c-42: 	for i, context := range contexts {
g4c-43: 		for _, id := range args {
g4c-44: 			if fmt.Sprintf("%d", i+1) == id || context.Name == id || context.FilePath == id || context.Url == id {
g4c-45: 				deleteIds[context.Id] = true
g4c-46: 				break
g4c-47: 			} else if context.FilePath != "" {
g4c-48: 				// Check if id is a glob pattern
g4c-49: 				matched, err := filepath.Match(id, context.FilePath)
g4c-50: 				if err != nil {
g4c-51: 					term.OutputErrorAndExit("Error matching glob pattern: %v", err)
g4c-52: 				}
g4c-53: 				if matched {
g4c-54: 					deleteIds[context.Id] = true
g4c-55: 					break
g4c-56: 				}
g4c-57: 
g4c-58: 				// Check if id is a parent directory
g4c-59: 				parentDir := context.FilePath
g4c-60: 				for parentDir != "." && parentDir != "/" && parentDir != "" {
g4c-61: 					if parentDir == id {
g4c-62: 						deleteIds[context.Id] = true
g4c-63: 						break
g4c-64: 					}
g4c-65: 					parentDir = filepath.Dir(parentDir) // Move up one directory
g4c-66: 				}
g4c-67: 
g4c-68: 			}
g4c-69: 		}
g4c-70: 	}
g4c-71: 
g4c-72: 	if len(deleteIds) > 0 {
g4c-73: 		res, err := api.Client.DeleteContext(lib.CurrentPlanId, lib.CurrentBranch, shared.DeleteContextRequest{
g4c-74: 			Ids: deleteIds,
g4c-75: 		})
g4c-76: 		term.StopSpinner()
g4c-77: 
g4c-78: 		if err != nil {
g4c-79: 			term.OutputErrorAndExit("Error deleting context: %v", err)
g4c-80: 		}
g4c-81: 
g4c-82: 		fmt.Println("✅ " + res.Msg)
g4c-83: 	} else {
g4c-84: 		term.StopSpinner()
g4c-85: 		fmt.Println("🤷‍♂️ No context removed")
g4c-86: 	}
g4c-87: }
g4c-88: 
g4c-89: func init() {
g4c-90: 	RootCmd.AddCommand(contextRmCmd)
g4c-91: }
g4c-92: 