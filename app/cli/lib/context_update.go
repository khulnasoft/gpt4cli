package lib

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"gpt4cli-cli/api"
	"gpt4cli-cli/fs"
	"gpt4cli-cli/term"
	"gpt4cli-cli/types"
	"gpt4cli-cli/url"
	"strconv"
	"strings"
	"sync"

	shared "gpt4cli-shared"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"
)

func CheckOutdatedContextWithOutput(quiet, autoConfirm bool, maybeContexts []*shared.Context, projectPaths *types.ProjectPaths) (contextOutdated, updated bool, err error) {
	if !quiet {
		term.StartSpinner("🔬 Checking context...")
	}

	var contexts []*shared.Context

	if maybeContexts != nil {
		contexts = maybeContexts
	} else {
		res, err := api.Client.ListContext(CurrentPlanId, CurrentBranch)
		if err != nil {
			term.StopSpinner()
			return false, false, fmt.Errorf("failed to list context: %s", err)
		}
		contexts = res
	}

	outdatedRes, err := CheckOutdatedContext(contexts, projectPaths)
	if err != nil {
		term.StopSpinner()
		return false, false, fmt.Errorf("failed to check outdated context: %s", err)
	}

	if !quiet {
		term.StopSpinner()
	}

	if len(outdatedRes.UpdatedContexts) == 0 && len(outdatedRes.RemovedContexts) == 0 {
		if !quiet {
			fmt.Println("✅ Context is up to date")
		}
		return false, false, nil
	}
	if len(outdatedRes.UpdatedContexts) > 0 {
		types := []string{}
		if outdatedRes.NumFiles > 0 {
			lbl := "file"
			if outdatedRes.NumFiles > 1 {
				lbl = "files"
			}
			lbl = strconv.Itoa(outdatedRes.NumFiles) + " " + lbl
			types = append(types, lbl)
		}
		if outdatedRes.NumUrls > 0 {
			lbl := "url"
			if outdatedRes.NumUrls > 1 {
				lbl = "urls"
			}
			lbl = strconv.Itoa(outdatedRes.NumUrls) + " " + lbl
			types = append(types, lbl)
		}
		if outdatedRes.NumTrees > 0 {
			lbl := "directory tree"
			if outdatedRes.NumTrees > 1 {
				lbl = "directory trees"
			}
			lbl = strconv.Itoa(outdatedRes.NumTrees) + " " + lbl
			types = append(types, lbl)
		}
		if outdatedRes.NumMaps > 0 {
			lbl := "map"
			if outdatedRes.NumMaps > 1 {
				lbl = "maps"
			}
			lbl = strconv.Itoa(outdatedRes.NumMaps) + " " + lbl
			types = append(types, lbl)
		}

		var msg string
		if len(types) <= 2 {
			msg += strings.Join(types, " and ")
		} else {
			for i, add := range types {
				if i == len(types)-1 {
					msg += ", and " + add
				} else {
					msg += ", " + add
				}
			}
		}

		phrase := "have been"
		if len(outdatedRes.UpdatedContexts) == 1 {
			phrase = "has been"
		}

		if !quiet {
			term.StopSpinner()

			color.New(term.ColorHiCyan, color.Bold).Printf("%s in context %s modified 👇\n\n", msg, phrase)

			tableString := tableForContextOutdated(outdatedRes.UpdatedContexts, outdatedRes.TokenDiffsById)
			fmt.Println(tableString)
		}
	}

	if len(outdatedRes.RemovedContexts) > 0 {
		types := []string{}
		if outdatedRes.NumFilesRemoved > 0 {
			lbl := "file"
			if outdatedRes.NumFilesRemoved > 1 {
				lbl = "files"
			}
			lbl = strconv.Itoa(outdatedRes.NumFilesRemoved) + " " + lbl
			types = append(types, lbl)
		}
		if outdatedRes.NumTreesRemoved > 0 {
			lbl := "directory tree"
			if outdatedRes.NumTreesRemoved > 1 {
				lbl = "directory trees"
			}
			lbl = strconv.Itoa(outdatedRes.NumTreesRemoved) + " " + lbl
			types = append(types, lbl)
		}

		var msg string
		if len(types) <= 2 {
			msg += strings.Join(types, " and ")
		} else {
			for i, add := range types {
				if i == len(types)-1 {
					msg += ", and " + add
				} else {
					msg += ", " + add
				}
			}
		}

		phrase := "have been"
		if len(outdatedRes.RemovedContexts) == 1 {
			phrase = "has been"
		}

		if !quiet {
			term.StopSpinner()

			color.New(term.ColorHiCyan, color.Bold).Printf("%s in context %s removed 👇\n\n", msg, phrase)

			tableString := tableForContextOutdated(outdatedRes.RemovedContexts, outdatedRes.TokenDiffsById)
			fmt.Println(tableString)
		}
	}

	confirmed := autoConfirm

	if !autoConfirm {
		confirmed, err = term.ConfirmYesNo("Update context now?")

		if err != nil {
			term.OutputErrorAndExit("failed to get user input: %s", err)
		}
	}

	if confirmed {
		_, err := UpdateContextWithOutput(UpdateContextParams{
			Contexts:    contexts,
			OutdatedRes: *outdatedRes,
			Req:         outdatedRes.Req,
		})
		if err != nil {
			return false, false, fmt.Errorf("error updating context: %v", err)
		}
		return true, true, nil
	} else {
		return true, false, nil
	}

}

type UpdateContextParams struct {
	Contexts    []*shared.Context
	OutdatedRes types.ContextOutdatedResult
	Req         map[string]*shared.UpdateContextParams
}

type UpdateContextResult struct {
	HasConflicts bool
	Msg          string
}

func UpdateContextWithOutput(params UpdateContextParams) (UpdateContextResult, error) {
	term.StartSpinner("🔄 Updating context...")

	updateRes, err := UpdateContext(params)

	if err != nil {
		return UpdateContextResult{}, err
	}

	term.StopSpinner()

	fmt.Println("✅ " + updateRes.Msg)

	return updateRes, nil
}

func UpdateContext(params UpdateContextParams) (UpdateContextResult, error) {
	req := params.Req

	var hasConflicts bool
	var msg string

	contextsById := map[string]*shared.Context{}
	for _, context := range params.Contexts {
		contextsById[context.Id] = context
	}
	deleteIds := map[string]bool{}
	for _, context := range params.OutdatedRes.RemovedContexts {
		deleteIds[context.Id] = true
	}

	filesToLoad := map[string]string{}
	for id := range req {
		context := contextsById[id]
		if context.ContextType == shared.ContextFileType {
			filesToLoad[context.FilePath] = context.Body
		}
	}
	for id := range deleteIds {
		context := contextsById[id]
		if context.ContextType == shared.ContextFileType {
			filesToLoad[context.FilePath] = ""
		}
	}

	var err error
	hasConflicts, err = checkContextConflicts(filesToLoad)
	if err != nil {
		return UpdateContextResult{}, fmt.Errorf("failed to check context conflicts: %v", err)
	}

	if len(req) > 0 {
		// log.Println("updating context")
		res, apiErr := api.Client.UpdateContext(CurrentPlanId, CurrentBranch, req)
		if apiErr != nil {
			return UpdateContextResult{}, fmt.Errorf("failed to update context: %v", apiErr)
		}
		// log.Println("updated context")
		// log.Println("res.Msg", res.Msg)
		msg = res.Msg
	} else {
	}

	if len(deleteIds) > 0 {
		res, apiErr := api.Client.DeleteContext(CurrentPlanId, CurrentBranch, shared.DeleteContextRequest{
			Ids: deleteIds,
		})
		if apiErr != nil {
			return UpdateContextResult{}, fmt.Errorf("failed to delete contexts: %v", apiErr)
		}
		msg += " " + res.Msg
	}

	return UpdateContextResult{
		HasConflicts: hasConflicts,
		Msg:          msg,
	}, nil
}

func CheckOutdatedContext(maybeContexts []*shared.Context, projectPaths *types.ProjectPaths) (*types.ContextOutdatedResult, error) {
	return checkOutdatedAndMaybeUpdateContext(false, maybeContexts, projectPaths)
}

func checkOutdatedAndMaybeUpdateContext(doUpdate bool, maybeContexts []*shared.Context, projectPaths *types.ProjectPaths) (*types.ContextOutdatedResult, error) {
	var contexts []*shared.Context

	if maybeContexts == nil {
		var apiErr *shared.ApiError
		contexts, apiErr = api.Client.ListContext(CurrentPlanId, CurrentBranch)
		if apiErr != nil {
			return nil, fmt.Errorf("error retrieving context: %v", apiErr)
		}
	} else {
		contexts = maybeContexts
	}

	totalTokens := 0
	for _, context := range contexts {
		totalTokens += context.NumTokens
	}

	var errs []error

	req := shared.UpdateContextRequest{}
	var updatedContexts []*shared.Context
	var tokenDiffsById = map[string]int{}
	var numFiles int
	var numUrls int
	var numTrees int
	var numMaps int
	var numFilesRemoved int
	var numTreesRemoved int
	var mu sync.Mutex
	var wg sync.WaitGroup
	contextsById := map[string]*shared.Context{}
	deleteIds := map[string]bool{}

	paths := projectPaths

	for _, context := range contexts {
		contextsById[context.Id] = context

		if context.ContextType == shared.ContextFileType {
			wg.Add(1)
			go func(context *shared.Context) {
				defer wg.Done()

				mu.Lock()
				defer mu.Unlock()

				if _, err := os.Stat(context.FilePath); os.IsNotExist(err) {
					deleteIds[context.Id] = true
					numFilesRemoved++
					tokenDiffsById[context.Id] = -context.NumTokens
					return
				}

				fileContent, err := os.ReadFile(context.FilePath)

				if err != nil {
					errs = append(errs, fmt.Errorf("failed to read the file %s: %v", context.FilePath, err))
					return
				}

				hash := sha256.Sum256(fileContent)
				sha := hex.EncodeToString(hash[:])

				if sha != context.Sha {
					// log.Println()
					// log.Println("context.FilePath", context.FilePath)
					// log.Println("context.Sha", context.Sha, "sha", sha)
					// log.Println("fileContent", string(fileContent))
					// log.Println()

					body := string(fileContent)

					numTokens := shared.GetNumTokensEstimate(body)
					tokenDiffsById[context.Id] = numTokens - context.NumTokens

					numFiles++
					updatedContexts = append(updatedContexts, context)

					req[context.Id] = &shared.UpdateContextParams{
						Body: body,
					}
				}
			}(context)

		} else if context.ContextType == shared.ContextDirectoryTreeType {
			wg.Add(1)
			go func(context *shared.Context) {
				defer wg.Done()

				// check if the directory tree exists
				if _, err := os.Stat(context.FilePath); os.IsNotExist(err) {
					mu.Lock()
					defer mu.Unlock()
					deleteIds[context.Id] = true
					numTreesRemoved++
					tokenDiffsById[context.Id] = -context.NumTokens
					return
				}

				baseDir := fs.GetBaseDirForFilePaths([]string{context.FilePath})

				flattenedPaths, err := ParseInputPaths(ParseInputPathsParams{
					FileOrDirPaths: []string{context.FilePath},
					BaseDir:        baseDir,
					ProjectPaths:   paths,
					LoadParams: &types.LoadContextParams{
						NamesOnly:       true,
						ForceSkipIgnore: context.ForceSkipIgnore,
					},
				})

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					errs = append(errs, fmt.Errorf("failed to get the directory tree %s: %v", context.FilePath, err))
					return
				}

				if !context.ForceSkipIgnore {
					if paths == nil {
						errs = append(errs, fmt.Errorf("project paths are nil"))
						return
					}

					var filteredPaths []string
					for _, path := range flattenedPaths {
						if _, ok := paths.ActivePaths[path]; ok {
							filteredPaths = append(filteredPaths, path)
						}
					}
					flattenedPaths = filteredPaths
				}

				body := strings.Join(flattenedPaths, "\n")
				bytes := []byte(body)

				hash := sha256.Sum256(bytes)
				sha := hex.EncodeToString(hash[:])

				if sha != context.Sha {
					numTokens := shared.GetNumTokensEstimate(body)
					tokenDiffsById[context.Id] = numTokens - context.NumTokens

					numTrees++
					updatedContexts = append(updatedContexts, context)
					req[context.Id] = &shared.UpdateContextParams{
						Body: body,
					}
				}
			}(context)

		} else if context.ContextType == shared.ContextMapType {
			wg.Add(1)
			go func(context *shared.Context) {
				defer wg.Done()

				mu.Lock()
				defer mu.Unlock()

				var removedMapPaths []string

				// Check if any input files have changed
				var updatedInputs = make(shared.FileMapInputs)
				var updatedInputShas = map[string]string{}

				for path, currentSha := range context.MapShas {
					bytes, err := os.ReadFile(path)
					if err != nil {
						if os.IsNotExist(err) {
							removedMapPaths = append(removedMapPaths, path)
							continue
						}

						errs = append(errs, fmt.Errorf("failed to read map file %s: %v", path, err))
						return
					}

					hash := sha256.Sum256(bytes)
					newSha := hex.EncodeToString(hash[:])

					if newSha != currentSha {
						// log.Println("path", path, "newSha", newSha, "currentSha", currentSha)
						// fmt.Println("path", path, "newSha", newSha, "currentSha", currentSha)
						content := string(bytes)
						updatedInputs[path] = content
						updatedInputShas[path] = newSha
					}
				}

				// Check if new files were added
				baseDir := fs.GetBaseDirForFilePaths([]string{context.FilePath})

				flattenedPaths, err := ParseInputPaths(ParseInputPathsParams{
					FileOrDirPaths: []string{context.FilePath},
					BaseDir:        baseDir,
					ProjectPaths:   paths,
					LoadParams:     &types.LoadContextParams{Recursive: true},
				})
				if err != nil {
					errs = append(errs, fmt.Errorf("failed to get the directory tree %s: %v", context.FilePath, err))
					return
				}

				var filteredPaths []string
				for _, inputFilePath := range flattenedPaths {
					if _, ok := paths.ActivePaths[inputFilePath]; ok {
						filteredPaths = append(filteredPaths, inputFilePath)
					}
				}
				flattenedPaths = filteredPaths

				for _, path := range flattenedPaths {
					if !shared.HasFileMapSupport(path) {
						continue
					}

					if _, ok := context.MapShas[path]; !ok {
						// log.Println("path", path, "not in context.MapShas")

						bytes, err := os.ReadFile(path)
						if err != nil {
							errs = append(errs, fmt.Errorf("failed to read map file %s: %v", path, err))
							return
						}

						content := string(bytes)
						updatedInputs[path] = content
						hash := sha256.Sum256(bytes)
						sha := hex.EncodeToString(hash[:])
						updatedInputShas[path] = sha
					} else {
					}
				}
				// If any files changed, get new map

				if len(updatedInputs) > 0 || len(removedMapPaths) > 0 {
					// Check total map input size and paths before making API call
					if len(updatedInputs) > shared.MaxContextMapPaths {
						errs = append(errs, fmt.Errorf("total map paths limit exceeded (found %d, limit %d)", len(updatedInputs), shared.MaxContextMapPaths))
						return
					}

					var totalMapSize int64
					for _, input := range updatedInputs {
						totalMapSize += int64(len(input))
					}

					if totalMapSize > shared.MaxContextMapInputSize {
						errs = append(errs, fmt.Errorf("total map size limit exceeded (size %.2f MB, limit %d MB)", float64(totalMapSize)/1024/1024, int(shared.MaxContextMapInputSize)/1024/1024))
						return
					}

					updatedParts := make(shared.FileMapBodies)
					for k, v := range context.MapParts {
						updatedParts[k] = v
					}
					var updatedMapBodies shared.FileMapBodies
					if len(updatedInputs) > 0 {
						mapRes, apiErr := api.Client.GetFileMap(shared.GetFileMapRequest{
							MapInputs: updatedInputs,
						})
						if apiErr != nil {
							errs = append(errs, fmt.Errorf("failed to get file map: %v", apiErr))
							return
						}
						updatedMapBodies = mapRes.MapBodies

						// Update map parts with new content
						for path, body := range mapRes.MapBodies {
							updatedParts[path] = body

							prevTokens := context.MapTokens[path]
							numTokens := mapRes.MapBodies.TokenEstimateForPath(path)

							if numTokens != prevTokens {
								tokenDiffsById[context.Id] += numTokens - prevTokens
							}
						}
					}

					if len(removedMapPaths) > 0 {
						for _, path := range removedMapPaths {
							delete(updatedParts, path)
							tokenDiffsById[context.Id] -= context.MapTokens[path]
						}
					}

					numMaps++
					updatedContexts = append(updatedContexts, context)
					req[context.Id] = &shared.UpdateContextParams{
						MapBodies:       updatedMapBodies,
						InputShas:       updatedInputShas,
						RemovedMapPaths: removedMapPaths,
					}

				} else {
				}
			}(context)
		} else if context.ContextType == shared.ContextURLType {
			wg.Add(1)
			go func(context *shared.Context) {
				defer wg.Done()
				body, err := url.FetchURLContent(context.Url)

				mu.Lock()
				defer mu.Unlock()

				if err != nil {
					errs = append(errs, fmt.Errorf("failed to fetch the URL %s: %v", context.Url, err))
					return
				}

				hash := sha256.Sum256([]byte(body))
				sha := hex.EncodeToString(hash[:])

				if sha != context.Sha {
					numTokens := shared.GetNumTokensEstimate(body)
					tokenDiffsById[context.Id] = numTokens - context.NumTokens

					numUrls++
					updatedContexts = append(updatedContexts, context)
					req[context.Id] = &shared.UpdateContextParams{
						Body: body,
					}
				}

			}(context)
		}
	}

	wg.Wait()

	if len(errs) > 0 {
		return nil, fmt.Errorf("failed to check context outdated: %v", errs)
	}

	var totalContextCount int
	var totalBodySize int64

	for _, context := range contexts {
		totalContextCount++
		totalBodySize += int64(len(context.Body))
	}

	for _, context := range updatedContexts {
		if req[context.Id] != nil {
			totalBodySize += int64(len(req[context.Id].Body)) - int64(len(context.Body))
		}
	}

	if totalContextCount > shared.MaxContextCount {
		return nil, fmt.Errorf("too many contexts to update (found %d, limit is %d)", totalContextCount, shared.MaxContextCount)
	}

	if totalBodySize > shared.MaxContextBodySize {
		return nil, fmt.Errorf("total context body size exceeds limit (size %.2f MB, limit %d MB)", float64(totalBodySize)/1024/1024, int(shared.MaxContextBodySize)/1024/1024)
	}

	var removedContexts []*shared.Context
	for id := range deleteIds {
		removedContexts = append(removedContexts, contextsById[id])
	}

	var outdatedRes types.ContextOutdatedResult
	var msg string
	var hasConflicts bool

	if len(req) == 0 && len(deleteIds) == 0 {
		// log.Println("return context is up to date res")
		return &types.ContextOutdatedResult{
			Msg: "Context is up to date",
		}, nil
	} else {

		outdatedRes = types.ContextOutdatedResult{
			UpdatedContexts: updatedContexts,
			RemovedContexts: removedContexts,
			TokenDiffsById:  tokenDiffsById,
			NumFiles:        numFiles,
			NumUrls:         numUrls,
			NumTrees:        numTrees,
			NumMaps:         numMaps,
			NumFilesRemoved: numFilesRemoved,
			NumTreesRemoved: numTreesRemoved,
			Req:             req,
		}

		if doUpdate {

			res, err := UpdateContext(UpdateContextParams{
				Contexts:    contexts,
				OutdatedRes: outdatedRes,
				Req:         req,
			})
			if err != nil {
				return nil, fmt.Errorf("failed to update context: %v", err)
			}
			hasConflicts = res.HasConflicts
			msg = res.Msg
			outdatedRes.Msg = msg
		} else {
			tokensDiff := 0
			for _, diff := range tokenDiffsById {
				tokensDiff += diff
			}
			total := totalTokens + tokensDiff
			outdatedRes.Msg = shared.SummaryForUpdateContext(shared.SummaryForUpdateContextParams{
				NumFiles:    numFiles,
				NumTrees:    numTrees,
				NumUrls:     numUrls,
				NumMaps:     numMaps,
				TokensDiff:  tokensDiff,
				TotalTokens: total,
			})
		}
	}

	if hasConflicts {
		term.StartSpinner("🏗️  Starting build...")
		_, err := buildPlanInlineFn(false, nil) // don't pass in outdated contexts -- nil value causes them to be refetched, which is what we want since they were just updated

		if err != nil {
			return nil, fmt.Errorf("failed to build plan: %v", err)
		}

		fmt.Println()
	}

	return &outdatedRes, nil
}

func tableForContextOutdated(updatedContexts []*shared.Context, tokenDiffsById map[string]int) string {
	if len(updatedContexts) == 0 {
		return ""
	}

	tableString := &strings.Builder{}
	table := tablewriter.NewWriter(tableString)
	table.SetHeader([]string{"Name", "Type", "🪙"})
	table.SetAutoWrapText(false)

	for _, context := range updatedContexts {
		t, icon := context.TypeAndIcon()
		diff := tokenDiffsById[context.Id]
		diffStr := "+" + strconv.Itoa(diff)
		tableColor := tablewriter.FgHiGreenColor

		if diff < 0 {
			diffStr = strconv.Itoa(diff)
			tableColor = tablewriter.FgHiRedColor
		}

		row := []string{
			" " + icon + " " + context.Name,
			t,
			diffStr,
		}

		table.Rich(row, []tablewriter.Colors{
			{tableColor, tablewriter.Bold},
			{tableColor},
			{tableColor},
		})
	}

	table.Render()

	return tableString.String()
}
