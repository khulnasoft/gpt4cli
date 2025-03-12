package lib

import (
	"bufio"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"gpt4cli-cli/api"
	"gpt4cli-cli/auth"
	"gpt4cli-cli/fs"
	"gpt4cli-cli/term"
	"gpt4cli-cli/types"
	"gpt4cli-cli/url"
	"strings"
	"sync"

	shared "gpt4cli-shared"

	"github.com/fatih/color"
)

func MustLoadContext(resources []string, params *types.LoadContextParams) {
	// startTime := time.Now()
	showElapsed := func(msg string) {
		// elapsed := time.Since(startTime)
		// log.Println(msg, "elapsed: %s\n", elapsed)
	}

	if params.DefsOnly {
		// while caching is set up to work with multiple map paths, it can end up in a partially loaded state if token limits are exceeded, so better to just load one at a time
		if len(resources) > 1 {
			term.OutputErrorAndExit("Please load a single map directory at a time")
		}

		term.LongSpinnerWithWarning("🗺️  Building project map...", "🗺️  This can take a while in larger projects...")
	} else if params.NamesOnly {
		term.LongSpinnerWithWarning("🌳 Loading directory tree...", "🌳 This can take a while in larger projects...")
	} else {
		term.StartSpinner("📥 Loading context...")
	}

	onErr := func(err error) {
		term.StopSpinner()
		term.OutputErrorAndExit("Failed to load context: %v", err)
	}

	var loadContextReq shared.LoadContextRequest

	fileInfo, err := os.Stdin.Stat()
	if err != nil {
		onErr(fmt.Errorf("failed to stat stdin: %v", err))
	}

	var apiKeys map[string]string
	var openAIBase string

	if !auth.Current.IntegratedModelsMode {
		if params.Note != "" || fileInfo.Mode()&os.ModeNamedPipe != 0 {
			apiKeys = MustVerifyApiKeysSilent()
			openAIBase = os.Getenv("OPENAI_API_BASE")
			if openAIBase == "" {
				openAIBase = os.Getenv("OPENAI_ENDPOINT")
			}
		}
	}

	if params.Note != "" {
		loadContextReq = append(loadContextReq, &shared.LoadContextParams{
			ContextType: shared.ContextNoteType,
			Body:        params.Note,
			ApiKeys:     apiKeys,
			OpenAIBase:  openAIBase,
			OpenAIOrgId: os.Getenv("OPENAI_ORG_ID"),
			AutoLoaded:  params.AutoLoaded,
		})
	}

	if fileInfo.Mode()&os.ModeNamedPipe != 0 {
		reader := bufio.NewReader(os.Stdin)
		pipedData, err := io.ReadAll(reader)
		if err != nil {
			onErr(fmt.Errorf("failed to read piped data: %v", err))
		}

		if len(pipedData) > 0 {
			loadContextReq = append(loadContextReq, &shared.LoadContextParams{
				ContextType: shared.ContextPipedDataType,
				Body:        string(pipedData),
				ApiKeys:     apiKeys,
				OpenAIBase:  openAIBase,
				OpenAIOrgId: os.Getenv("OPENAI_ORG_ID"),
				AutoLoaded:  params.AutoLoaded,
			})
		}
	}

	var inputUrls []string
	var inputFilePaths []string

	if len(resources) > 0 {
		for _, resource := range resources {
			// so far resources are either files or urls
			if url.IsValidURL(resource) {
				inputUrls = append(inputUrls, resource)
			} else {
				if strings.HasPrefix(resource, "."+string(os.PathSeparator)) {
					resource = resource[2:]
				}

				inputFilePaths = append(inputFilePaths, resource)
			}
		}
	}

	var contextMu sync.Mutex

	errCh := make(chan error)
	ignoredPaths := make(map[string]string)
	mapFilesSkippedTooLarge := []struct {
		Path string
		Size int64
	}{}
	mapFilesSkippedAfterSizeLimit := []string{}

	numRoutines := 0

	// filter out already loaded contexts
	alreadyLoadedByComposite := make(map[string]*shared.Context)
	existingContexts, apiErr := api.Client.ListContext(CurrentPlanId, CurrentBranch)
	if apiErr != nil {
		onErr(fmt.Errorf("failed to list contexts: %v", apiErr.Msg))
	}

	existsByComposite := make(map[string]*shared.Context)
	for _, context := range existingContexts {
		switch context.ContextType {
		case shared.ContextFileType, shared.ContextDirectoryTreeType, shared.ContextMapType, shared.ContextImageType:
			existsByComposite[strings.Join([]string{string(context.ContextType), context.FilePath}, "|")] = context
		case shared.ContextURLType:
			existsByComposite[strings.Join([]string{string(context.ContextType), context.Url}, "|")] = context
		}
	}

	var cachedMapPaths map[string]bool
	var cachedMapLoadRes *shared.LoadContextResponse

	if len(inputFilePaths) > 0 {
		mapInputsByPath := map[string]shared.FileMapInputs{}
		toLoadMapPaths := []string{}

		var mapSize int64

		if params.DefsOnly {
			for _, inputFilePath := range inputFilePaths {
				composite := strings.Join([]string{string(shared.ContextMapType), inputFilePath}, "|")
				if existsByComposite[composite] != nil {
					alreadyLoadedByComposite[composite] = existsByComposite[composite]
					continue
				}

				mapInputsByPath[inputFilePath] = shared.FileMapInputs{}
				toLoadMapPaths = append(toLoadMapPaths, inputFilePath)
			}

			var uncachedMapPaths []string

			res, err := api.Client.LoadCachedFileMap(CurrentPlanId, CurrentBranch, shared.LoadCachedFileMapRequest{
				FilePaths: toLoadMapPaths,
			})

			if err != nil {
				onErr(fmt.Errorf("error checking cached file map: %v", err))
			}

			if res.LoadRes != nil {
				if res.LoadRes.MaxTokensExceeded {
					term.StopSpinner()
					overage := res.LoadRes.TotalTokens - res.LoadRes.MaxTokens

					term.OutputErrorAndExit("Update would add %d 🪙 and exceed token limit (%d) by %d 🪙\n", res.LoadRes.TokensAdded, res.LoadRes.MaxTokens, overage)
				}

				cachedMapLoadRes = res.LoadRes
				cachedMapPaths = res.CachedByPath

				for _, path := range toLoadMapPaths {
					if !cachedMapPaths[path] {
						uncachedMapPaths = append(uncachedMapPaths, path)
					}
				}
			} else {
				uncachedMapPaths = toLoadMapPaths
			}

			toLoadMapPaths = uncachedMapPaths
			inputFilePaths = toLoadMapPaths

			showElapsed("Checked cached maps")
		}

		if len(inputFilePaths) > 0 {
			baseDir := fs.GetBaseDirForFilePaths(inputFilePaths)

			showElapsed("Got base dir")

			paths, err := fs.GetProjectPaths(baseDir)
			if err != nil {
				onErr(fmt.Errorf("failed to get project paths: %v", err))
			}

			showElapsed("Got project paths")

			// log.Println(spew.Sdump(paths))

			// fmt.Println("active paths", len(paths.ActivePaths))
			// fmt.Println("all paths", len(paths.AllPaths))
			// fmt.Println("ignored paths", len(paths.IgnoredPaths))

			// spew.Dump(paths.IgnoredPaths)
			// spew.Dump(paths.ActivePaths)

			if !params.ForceSkipIgnore {
				var filteredPaths []string
				for _, inputFilePath := range inputFilePaths {
					if _, ok := paths.ActivePaths[inputFilePath]; !ok {
						ignored, reason, err := fs.IsIgnored(paths, inputFilePath, baseDir)
						if err != nil {
							onErr(fmt.Errorf("failed to check if %s is ignored: %v", inputFilePath, err))
						}
						if ignored {
							ignoredPaths[inputFilePath] = reason
						}
					} else {
						filteredPaths = append(filteredPaths, inputFilePath)
					}
				}
				inputFilePaths = filteredPaths

				showElapsed("Filtered paths")
			}

			if params.NamesOnly {
				for _, inputFilePath := range inputFilePaths {
					composite := strings.Join([]string{string(shared.ContextDirectoryTreeType), inputFilePath}, "|")
					if existsByComposite[composite] != nil {
						alreadyLoadedByComposite[composite] = existsByComposite[composite]
						continue
					}

					numRoutines++
					go func(inputFilePath string) {
						flattenedPaths, err := ParseInputPaths(ParseInputPathsParams{
							FileOrDirPaths: []string{inputFilePath},
							BaseDir:        baseDir,
							ProjectPaths:   paths,
							LoadParams:     params,
						})
						if err != nil {
							errCh <- fmt.Errorf("failed to parse input paths: %v", err)
							return
						}

						if !params.ForceSkipIgnore {
							var filteredPaths []string
							for _, path := range flattenedPaths {
								if _, ok := paths.ActivePaths[path]; ok {
									filteredPaths = append(filteredPaths, path)
								} else {
									ignored, reason, err := fs.IsIgnored(paths, path, baseDir)
									if err != nil {
										errCh <- fmt.Errorf("failed to check if %s is ignored: %v", path, err)
										return
									}
									if ignored {
										ignoredPaths[path] = reason
										ignoredPaths[path] = paths.IgnoredPaths[path]
									}
								}
							}
							flattenedPaths = filteredPaths
						}

						body := strings.Join(flattenedPaths, "\n")

						name := inputFilePath
						if name == "." {
							name = "cwd"
						}
						if name == ".." {
							name = "parent"
						}

						contextMu.Lock()
						defer contextMu.Unlock()
						loadContextReq = append(loadContextReq, &shared.LoadContextParams{
							ContextType:     shared.ContextDirectoryTreeType,
							Name:            name,
							Body:            body,
							FilePath:        inputFilePath,
							ForceSkipIgnore: params.ForceSkipIgnore,
							AutoLoaded:      params.AutoLoaded,
						})

						errCh <- nil
					}(inputFilePath)
				}

			} else {
				flattenedPaths, err := ParseInputPaths(ParseInputPathsParams{
					FileOrDirPaths: inputFilePaths,
					BaseDir:        baseDir,
					ProjectPaths:   paths,
					LoadParams:     params,
				})
				if err != nil {
					onErr(fmt.Errorf("failed to parse input paths: %v", err))
				}

				showElapsed("Parsed input paths")

				// // Dump flattenedPaths to JSON file for debugging
				// debugData, err := json.MarshalIndent(flattenedPaths, "", "  ")
				// if err != nil {
				// 	onErr(fmt.Errorf("failed to marshal flattened paths: %v", err))
				// 	return
				// }

				// if err := os.WriteFile("flattened_paths_debug.json", debugData, 0644); err != nil {
				// 	onErr(fmt.Errorf("failed to write debug file: %v", err))
				// 	return
				// }

				if !params.ForceSkipIgnore {
					var filteredPaths []string
					for _, path := range flattenedPaths {
						if _, ok := paths.ActivePaths[path]; ok {
							filteredPaths = append(filteredPaths, path)
						} else {
							ignored, reason, err := fs.IsIgnored(paths, path, baseDir)
							if err != nil {
								onErr(fmt.Errorf("failed to check if %s is ignored: %v", path, err))
							}
							if ignored {
								ignoredPaths[path] = reason
							}
						}
					}
					flattenedPaths = filteredPaths

					showElapsed("Filtered paths")
				}

				// Add this check for the number of files (after filtering out ignored/irrelevant paths)
				var numPaths int
				if params.DefsOnly {
					for _, path := range flattenedPaths {
						if shared.HasFileMapSupport(path) {
							numPaths++
						}
					}
					showElapsed("Counted map paths")
				} else {
					numPaths = len(flattenedPaths)
				}

				if (params.DefsOnly || params.NamesOnly) && numPaths > shared.MaxContextMapPaths {
					onErr(fmt.Errorf("too many files to load (found %d, limit is %d)", numPaths, shared.MaxContextMapPaths))
				} else if !params.DefsOnly && !params.NamesOnly && numPaths > shared.MaxContextCount {
					onErr(fmt.Errorf("too many files to load (found %d, limit is %d)", numPaths, shared.MaxContextCount))
				}

				inputFilePaths = flattenedPaths

				for _, path := range flattenedPaths {
					var mapInputPath string
					if params.DefsOnly {
						for _, inputPath := range toLoadMapPaths {
							// Clean and make absolute paths for comparison
							absPath, err := filepath.Abs(path)
							if err != nil {
								continue
							}
							absInputPath, err := filepath.Abs(inputPath)
							if err != nil {
								continue
							}

							// Check if paths are equal or if path is under inputPath
							if absPath == absInputPath || strings.HasPrefix(absPath+string(filepath.Separator), absInputPath+string(filepath.Separator)) {
								mapInputPath = inputPath
								break
							}
						}

						if mapInputPath == "" {
							continue // not a child of any input path
						}

						// add empty entry for non-supported file types to show file tree
						// if !shared.HasFileMapSupport(path) {
						// 	// not a tree-sitter supported file type
						// 	continue
						// }

						if _, ok := mapInputsByPath[mapInputPath]; !ok {
							mapInputsByPath[mapInputPath] = shared.FileMapInputs{}
						}
					}

					var contextType shared.ContextType
					isImage := shared.IsImageFile(path)
					if isImage {
						contextType = shared.ContextImageType
					} else if params.DefsOnly {
						contextType = shared.ContextMapType
					} else {
						contextType = shared.ContextFileType
					}

					if !params.DefsOnly {
						composite := strings.Join([]string{string(contextType), path}, "|")

						if existsByComposite[composite] != nil {
							alreadyLoadedByComposite[composite] = existsByComposite[composite]
							continue
						}
					}

					numRoutines++
					go func(path string) {
						var fileContent []byte
						var size int64
						// File size check
						fileInfo, err := os.Stat(path)
						if err != nil {
							errCh <- fmt.Errorf("failed to get file info for %s: %v", path, err)
							return
						}

						size = fileInfo.Size()

						if !params.DefsOnly && size > shared.MaxContextBodySize {
							errCh <- fmt.Errorf("file %s exceeds size limit (size %.2f MB, limit %d MB)", path, float64(fileInfo.Size())/1024/1024, int(shared.MaxContextBodySize)/1024/1024)
							return
						}

						fileContent, err = os.ReadFile(path)
						if err != nil {
							errCh <- fmt.Errorf("failed to read the file %s: %v", path, err)
							return
						}

						contextMu.Lock()
						defer contextMu.Unlock()

						if params.DefsOnly {
							if size > shared.MaxContextBodySize {
								mapFilesSkippedTooLarge = append(mapFilesSkippedTooLarge, struct {
									Path string
									Size int64
								}{Path: path, Size: size})
								errCh <- nil
								return
							}

							var mapSizeExceeded bool
							if mapSize+size > shared.MaxContextMapInputSize {
								mapSizeExceeded = true
								mapFilesSkippedAfterSizeLimit = append(mapFilesSkippedAfterSizeLimit, path)
							}

							content := string(fileContent)
							if mapSizeExceeded {
								content = ""
							}

							mapInputsByPath[mapInputPath][path] = content
							mapSize += size
						} else if isImage {
							loadContextReq = append(loadContextReq, &shared.LoadContextParams{
								ContextType: shared.ContextImageType,
								Name:        path,
								Body:        base64.StdEncoding.EncodeToString(fileContent),
								FilePath:    path,
								ImageDetail: params.ImageDetail,
								AutoLoaded:  params.AutoLoaded,
							})
						} else {
							loadContextReq = append(loadContextReq, &shared.LoadContextParams{
								ContextType: shared.ContextFileType,
								Name:        path,
								Body:        string(fileContent),
								FilePath:    path,
								AutoLoaded:  params.AutoLoaded,
							})
						}

						errCh <- nil
					}(path)
				}

				showElapsed("Got map input paths")

				if params.DefsOnly {
					for _, inputPath := range toLoadMapPaths {
						var name string
						if inputPath == "." {
							name = "cwd"
						} else if inputPath == ".." {
							name = "parent"
						} else {
							name = inputPath
						}

						loadContextReq = append(loadContextReq, &shared.LoadContextParams{
							ContextType: shared.ContextMapType,
							Name:        name,
							MapInputs:   mapInputsByPath[inputPath],
							FilePath:    inputPath,
							AutoLoaded:  params.AutoLoaded,
						})
					}
				}
			}
		}
	}

	if len(inputUrls) > 0 {
		for _, u := range inputUrls {
			composite := strings.Join([]string{string(shared.ContextURLType), u}, "|")
			if existsByComposite[composite] != nil {
				alreadyLoadedByComposite[composite] = existsByComposite[composite]
				continue
			}

			numRoutines++
			go func(u string) {
				body, err := url.FetchURLContent(u)
				if err != nil {
					errCh <- fmt.Errorf("failed to fetch content from URL %s: %v", u, err)
					return
				}

				name := url.SanitizeURL(u)
				// show the first 20 characters, then ellipsis then the last 20 characters of 'name'
				if len(name) > 40 {
					name = name[:20] + "⋯" + name[len(name)-20:]
				}

				contextMu.Lock()
				defer contextMu.Unlock()

				loadContextReq = append(loadContextReq, &shared.LoadContextParams{
					ContextType: shared.ContextURLType,
					Name:        name,
					Body:        body,
					Url:         u,
					AutoLoaded:  params.AutoLoaded,
				})

				errCh <- nil
			}(u)
		}
	}

	for i := 0; i < numRoutines; i++ {
		err := <-errCh
		if err != nil {
			onErr(err)
		}
	}

	showElapsed("Loaded reqs")

	filesToLoad := map[string]string{}
	for _, context := range loadContextReq {
		if context.ContextType == shared.ContextFileType {
			filesToLoad[context.FilePath] = context.Body
		}
	}

	hasConflicts, err := checkContextConflicts(filesToLoad)

	showElapsed("Checked conflicts")

	if err != nil {
		onErr(fmt.Errorf("failed to check context conflicts: %v", err))
	}

	if len(loadContextReq)+len(cachedMapPaths) == 0 {
		term.StopSpinner()
		fmt.Println("🤷‍♂️ No context loaded")

		didOutputReason := false
		if len(alreadyLoadedByComposite) > 0 {
			printAlreadyLoadedMsg(alreadyLoadedByComposite)
			didOutputReason = true
		}
		if len(ignoredPaths) > 0 && !params.SkipIgnoreWarning {
			printIgnoredMsg()
			didOutputReason = true
		}

		if !didOutputReason {
			fmt.Println()
			fmt.Printf("Use %s to load a file or URL:", color.New(color.BgCyan, color.FgHiWhite).Sprint(" gpt4cli load [file-path|url] "))
			fmt.Println()
			fmt.Println("gpt4cli load file.c file.h")
			fmt.Println("gpt4cli load https://github.com/some-org/some-repo/README.md")

			fmt.Println()
			fmt.Printf("%s with the --recursive/-r flag:\n", color.New(color.Bold, term.ColorHiCyan).Sprint("Load a whole directory"))
			fmt.Println("gpt4cli load app/src -r")

			fmt.Println()
			fmt.Printf("%s with the --tree flag:\n", color.New(color.Bold, term.ColorHiCyan).Sprint("Load a directory layout (file names only)"))

			fmt.Println()
			fmt.Printf("%s file paths are relative to the current directory\n", color.New(color.Bold, term.ColorHiYellow).Sprint("Note:"))

			fmt.Println()
			fmt.Printf("%s with the -n flag:\n", color.New(color.Bold, term.ColorHiCyan).Sprint("Load a note"))
			fmt.Println("gpt4cli load -n 'Some note here'")

			fmt.Println()
			fmt.Printf("%s from any command:\n", color.New(color.Bold, term.ColorHiCyan).Sprint("Pipe data in"))
			fmt.Println("npm test | gpt4cli load")
		}

		os.Exit(0)
	}

	var res *shared.LoadContextResponse
	if cachedMapLoadRes != nil {
		res = cachedMapLoadRes
	} else {
		res, apiErr = api.Client.LoadContext(CurrentPlanId, CurrentBranch, loadContextReq)

		if apiErr != nil {
			onErr(fmt.Errorf("failed to load context: %v", apiErr.Msg))
		}
	}

	showElapsed("Made reqs")

	term.StopSpinner()

	if res.MaxTokensExceeded {
		overage := res.TotalTokens - res.MaxTokens
		term.OutputErrorAndExit("Update would add %d 🪙 and exceed token limit (%d) by %d 🪙\n", res.TokensAdded, res.MaxTokens, overage)
	}

	if hasConflicts {
		term.StartSpinner("🏗️  Starting build...")
		_, err := buildPlanInlineFn(false, nil)

		if err != nil {
			onErr(fmt.Errorf("failed to build plan: %v", err))
		}

		fmt.Println()
	}

	fmt.Println("✅ " + res.Msg)

	if len(alreadyLoadedByComposite) > 0 {
		printAlreadyLoadedMsg(alreadyLoadedByComposite)
	}

	if len(ignoredPaths) > 0 && !params.SkipIgnoreWarning {
		printIgnoredMsg()
	}

	if len(mapFilesSkippedTooLarge) > 0 || len(mapFilesSkippedAfterSizeLimit) > 0 {
		printSkippedMapFilesMsg(mapFilesSkippedTooLarge, mapFilesSkippedAfterSizeLimit)
	}
}

func AutoLoadContextFiles(ctx context.Context, files []string) (string, error) {
	loadContextReqs := shared.LoadContextRequest{}
	errCh := make(chan error, len(files))
	var mu sync.Mutex

	for _, path := range files {
		go func(path string) {
			body, err := os.ReadFile(path)
			if err != nil {
				errCh <- fmt.Errorf("failed to read file %s: %v", path, err)
				return
			}

			mu.Lock()
			defer mu.Unlock()

			loadContextReqs = append(loadContextReqs, &shared.LoadContextParams{
				ContextType: shared.ContextFileType,
				FilePath:    path,
				Name:        path,
				Body:        string(body),
				AutoLoaded:  true,
			})

			errCh <- nil
		}(path)
	}

	for i := 0; i < len(files); i++ {
		err := <-errCh
		if err != nil {
			return "", fmt.Errorf("failed to load context: %v", err)
		}
	}

	res, apiErr := api.Client.AutoLoadContext(ctx, CurrentPlanId, CurrentBranch, loadContextReqs)

	if apiErr != nil {
		return "", fmt.Errorf("failed to load context: %v", apiErr.Msg)
	}

	if res.MaxTokensExceeded {
		overage := res.TotalTokens - res.MaxTokens
		return "", fmt.Errorf("update would add %d 🪙 and exceed token limit (%d) by %d 🪙", res.TokensAdded, res.MaxTokens, overage)
	}

	return res.Msg, nil
}

func MustLoadAutoContextMap() {
	// fmt.Println("Select a base directory to load context from. Press enter to use current directory (.), otherwise use a relative path like 'src' or 'lib'.")
	// fmt.Println()

	// baseDir, err := term.GetUserStringInputWithDefault("Base directory for context:", ".")

	// if err != nil {
	// 	term.OutputErrorAndExit("Error: %v", err)
	// }

	MustLoadContext([]string{"."}, &types.LoadContextParams{
		DefsOnly:          true,
		SkipIgnoreWarning: true,
		AutoLoaded:        true,
	})
}

func printAlreadyLoadedMsg(alreadyLoadedByComposite map[string]*shared.Context) {
	fmt.Println()
	pronoun := "they're"
	if len(alreadyLoadedByComposite) == 1 {
		pronoun = "it's"
	}
	fmt.Printf("🙅‍♂️ Skipped because %s already in context:\n", pronoun)
	for _, context := range alreadyLoadedByComposite {
		_, icon := context.TypeAndIcon()

		fmt.Printf("  • %s %s\n", icon, context.Name)
	}
}

func printIgnoredMsg() {
	fmt.Println()
	fmt.Println("ℹ️  " + color.New(color.FgWhite).Sprint("Due to .gitignore or .gpt4cliignore, some paths weren't loaded.\nUse --force / -f to load ignored paths."))
}

func printSkippedMapFilesMsg(mapFilesSkippedTooLarge []struct {
	Path string
	Size int64
}, mapFilesSkippedAfterSizeLimit []string) {
	fmt.Println()
	if len(mapFilesSkippedTooLarge) > 0 {
		fmt.Println("ℹ️  These files were skipped because they're too large to map:")
		for _, file := range mapFilesSkippedTooLarge {
			fmt.Printf("  • %s - %d MB\n", file.Path, file.Size/1024/1024)
		}
	}
	if len(mapFilesSkippedAfterSizeLimit) > 0 {
		fmt.Println("ℹ️  These files were skipped because the total map size limit was exceeded:")
		for _, file := range mapFilesSkippedAfterSizeLimit {
			fmt.Printf("  • %s\n", file)
		}
		fmt.Println()
		fmt.Println("They will still be included in the map as paths in the project, but no maps will be generated for them.")
	}
}
