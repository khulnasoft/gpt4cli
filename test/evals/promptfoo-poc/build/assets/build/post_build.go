g4c-1: package cmd
g4c-2: 
g4c-3: import (
g4c-4: 	"fmt"
g4c-5: 	"path/filepath"
g4c-6: 	"gpt4cli/api"
g4c-7: 	"gpt4cli/auth"
g4c-8: 	"gpt4cli/lib"
g4c-9: 	"gpt4cli/term"
g4c-10: 	"strconv"
g4c-11: 	"strings"
g4c-12: 
g4c-13: 	"github.com/khulnasoft/gpt4cli/shared"
g4c-14: 	"github.com/spf13/cobra"
g4c-15: )
g4c-16: 
g4c-17: func parseRange(arg string) ([]int, error) {
g4c-18: 	var indices []int 
g4c-19: 	parts := strings.Split(arg, "-")
g4c-20: 	if len(parts) == 2 {
g4c-21: 		start, err := strconv.Atoi(parts[0])
g4c-22: 		if err != nil {
g4c-23: 			return nil, err
g4c-24: 		}
g4c-25: 		end, err := strconv.Atoi(parts[1])
g4c-26: 		if err != nil {
g4c-27: 			return nil, err
g4c-28: 		}
g4c-29: 		for i := start; i <= end; i++ {
g4c-30: 			indices = append(indices, i)
g4c-31: 		}
g4c-32: 	} else {
g4c-33: 		index, err := strconv.Atoi(arg)
g4c-34: 		if err != nil {
g4c-35: 			return nil, err
g4c-36: 		}
g4c-37: 		indices = append(indices, index)
g4c-38: 	}
g4c-39: 	return indices, nil
g4c-40: }
g4c-41: 
g4c-42: func contextRm(cmd *cobra.Command, args []string) {
g4c-43: 	auth.MustResolveAuthWithOrg()
g4c-44: 	lib.MustResolveProject()
g4c-45: 
g4c-46: 	if lib.CurrentPlanId == "" {
g4c-47: 		fmt.Println("🤷‍♂️ No current plan")
g4c-48: 		return
g4c-49: 	}
g4c-50: 
g4c-51: 	term.StartSpinner("")
g4c-52: 	contexts, err := api.Client.ListContext(lib.CurrentPlanId, lib.CurrentBranch)
g4c-53: 
g4c-54: 	if err != nil {
g4c-55: 		term.OutputErrorAndExit("Error retrieving context: %v", err)
g4c-56: 	}
g4c-57: 
g4c-58: 	deleteIds := map[string]bool{}
g4c-59: 
g4c-60: 	for _, arg := range args {
g4c-61: 		indices, err := parseRange(arg)
g4c-62: 		if err != nil {
g4c-63: 			term.OutputErrorAndExit("Error parsing range: %v", err)
g4c-64: 		}
g4c-65: 
g4c-66: 		for _, index := range indices {
g4c-67: 			if index > 0 && index <= len(contexts) {
g4c-68: 				context := contexts[index-1]
g4c-69: 				deleteIds[context.Id] = true
g4c-70: 			}
g4c-71: 		}
g4c-72: 	}
g4c-73: 
g4c-74: 	for i, context := range contexts {
g4c-75: 		for _, id := range args {
g4c-76: 			if fmt.Sprintf("%d", i+1) == id || context.Name == id || context.FilePath == id || context.Url == id {
g4c-77: 				deleteIds[context.Id] = true
g4c-78: 				break
g4c-79: 			} else if context.FilePath != "" {
g4c-80: 				// Check if id is a glob pattern
g4c-81: 				matched, err := filepath.Match(id, context.FilePath)
g4c-82: 				if err != nil {
g4c-83: 					term.OutputErrorAndExit("Error matching glob pattern: %v", err)
g4c-84: 				}
g4c-85: 				if matched {
g4c-86: 					deleteIds[context.Id] = true
g4c-87: 					break
g4c-88: 				}
g4c-89: 
g4c-90: 				// Check if id is a parent directory
g4c-91: 				parentDir := context.FilePath
g4c-92: 				for parentDir != "." && parentDir != "/" && parentDir != "" {
g4c-93: 					if parentDir == id {
g4c-94: 						deleteIds[context.Id] = true
g4c-95: 						break
g4c-96: 					}
g4c-97: 					parentDir = filepath.Dir(parentDir) // Move up one directory
g4c-98: 				}
g4c-99: 			}
g4c-100: 		}
g4c-101: 	}
g4c-102: 
g4c-103: 	if len(deleteIds) > 0 {
g4c-104: 		res, err := api.Client.DeleteContext(lib.CurrentPlanId, lib.CurrentBranch, shared.DeleteContextRequest{
g4c-105: 			Ids: deleteIds,
g4c-106: 		})
g4c-107: 		term.StopSpinner()
g4c-108: 
g4c-109: 		if err != nil {
g4c-110: 			term.OutputErrorAndExit("Error deleting context: %v", err)
g4c-111: 		}
g4c-112: 
g4c-113: 		fmt.Println("✅ " + res.Msg)
g4c-114: 	} else {
g4c-115: 		term.StopSpinner()
g4c-116: 		fmt.Println("🤷‍♂️ No context removed")
g4c-117: 	}
g4c-118: }
g4c-119: 
g4c-120: func init() {
g4c-121: 	RootCmd.AddCommand(contextRmCmd)
g4c-122: }
g4c-123: 