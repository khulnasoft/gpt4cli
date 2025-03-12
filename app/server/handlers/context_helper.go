package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"gpt4cli-server/db"
	"gpt4cli-server/model"
	"gpt4cli-server/types"

	shared "gpt4cli-shared"
)

type loadContextsParams struct {
	w                http.ResponseWriter
	r                *http.Request
	auth             *types.ServerAuth
	loadReq          *shared.LoadContextRequest
	plan             *db.Plan
	branchName       string
	cachedMapsByPath map[string]*db.CachedMap
	autoLoaded       bool
}

func loadContexts(
	params loadContextsParams,
) (*shared.LoadContextResponse, []*db.Context) {
	w := params.w
	r := params.r
	auth := params.auth
	loadReq := params.loadReq
	plan := params.plan
	branchName := params.branchName
	cachedMapsByPath := params.cachedMapsByPath
	autoLoaded := params.autoLoaded

	log.Printf("Starting loadContexts with %d contexts, cachedMapsByPath: %v, autoLoaded: %v", len(*loadReq), cachedMapsByPath != nil, autoLoaded)

	// check file count and size limits
	totalFiles := 0
	mapFilesCount := 0
	for _, context := range *loadReq {
		totalFiles++
		if context.ContextType == shared.ContextMapType {
			mapFilesCount++
			log.Printf("Found map file: %s with %d map inputs", context.FilePath, len(context.MapInputs))
		}

		if totalFiles > shared.MaxContextCount {
			log.Printf("Error: Too many contexts to load (found %d, limit is %d)\n", totalFiles, shared.MaxContextCount)
			http.Error(w, fmt.Sprintf("Too many contexts to load (found %d, limit is %d)", totalFiles, shared.MaxContextCount), http.StatusBadRequest)
			return nil, nil
		}

		fileSize := int64(len(context.Body))
		if fileSize > shared.MaxContextBodySize {
			log.Printf("Error: Context %s exceeds size limit (size %.2f MB, limit %d MB)\n", context.Name, float64(fileSize)/1024/1024, int(shared.MaxContextBodySize)/1024/1024)
			http.Error(w, fmt.Sprintf("Context %s exceeds size limit (size %.2f MB, limit %d MB)", context.Name, float64(fileSize)/1024/1024, int(shared.MaxContextBodySize)/1024/1024), http.StatusBadRequest)
			return nil, nil
		}
	}

	if mapFilesCount > 0 {
		log.Printf("Processing %d map files out of %d total contexts", mapFilesCount, totalFiles)
	}

	var err error

	var settings *shared.PlanSettings
	var clients map[string]model.ClientInfo

	for _, context := range *loadReq {
		if context.ContextType == shared.ContextPipedDataType || context.ContextType == shared.ContextNoteType || context.ContextType == shared.ContextImageType {

			settings, err = db.GetPlanSettings(plan, true)

			if err != nil {
				log.Printf("Error getting plan settings: %v\n", err)
				http.Error(w, "Error getting plan settings: "+err.Error(), http.StatusInternalServerError)
				return nil, nil
			}

			clients = initClients(
				initClientsParams{
					w:           w,
					auth:        auth,
					apiKeys:     context.ApiKeys,
					openAIBase:  context.OpenAIBase,
					openAIOrgId: context.OpenAIOrgId,
					plan:        plan,
				},
			)

			break
		}
	}

	// ensure image compatibility if we're loading an image
	for _, context := range *loadReq {
		if context.ContextType == shared.ContextImageType {
			if !settings.ModelPack.Planner.BaseModelConfig.HasImageSupport {
				log.Printf("Error loading context: %s does not support images in context\n", settings.ModelPack.Planner.BaseModelConfig.ModelName)
				http.Error(w, fmt.Sprintf("Error loading context: %s does not support images in context", settings.ModelPack.Planner.BaseModelConfig.ModelName), http.StatusBadRequest)
				return nil, nil
			}
		}
	}

	// get name for piped data or notes if present
	num := 0
	errCh := make(chan error, len(*loadReq))
	for _, context := range *loadReq {
		if context.ContextType == shared.ContextPipedDataType {
			num++

			go func(context *shared.LoadContextParams) {
				name, err := model.GenPipedDataName(r.Context(), auth, plan, settings, clients, context.Body)

				if err != nil {
					errCh <- fmt.Errorf("error generating name for piped data: %v", err)
					return
				}

				context.Name = name
				errCh <- nil
			}(context)
		} else if context.ContextType == shared.ContextNoteType {
			num++

			go func(context *shared.LoadContextParams) {
				name, err := model.GenNoteName(r.Context(), auth, plan, settings, clients, context.Body)

				if err != nil {
					errCh <- fmt.Errorf("error generating name for note: %v", err)
					return
				}

				context.Name = name
				errCh <- nil
			}(context)
		}
	}
	if num > 0 {
		for i := 0; i < num; i++ {
			err := <-errCh
			if err != nil {
				log.Printf("Error: %v\n", err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return nil, nil
			}
		}
	}

	ctx, cancel := context.WithCancel(r.Context())

	var loadRes *shared.LoadContextResponse
	var dbContexts []*db.Context

	err = db.ExecRepoOperation(db.ExecRepoOperationParams{
		OrgId:          auth.OrgId,
		UserId:         auth.User.Id,
		PlanId:         plan.Id,
		Branch:         branchName,
		Reason:         "load contexts",
		Scope:          db.LockScopeWrite,
		Ctx:            ctx,
		CancelFn:       cancel,
		ClearRepoOnErr: true,
	}, func(repo *db.GitRepo) error {
		log.Printf("Calling db.LoadContexts with %d contexts, %d cached maps", len(*loadReq), len(cachedMapsByPath))
		for path := range cachedMapsByPath {
			log.Printf("Using cached map for path: %s", path)
		}

		res, dbContextsRes, err := db.LoadContexts(ctx, db.LoadContextsParams{
			OrgId:            auth.OrgId,
			Plan:             plan,
			BranchName:       branchName,
			Req:              loadReq,
			UserId:           auth.User.Id,
			CachedMapsByPath: cachedMapsByPath,
			AutoLoaded:       autoLoaded,
		})

		if err != nil {
			return err
		}

		loadRes = res
		dbContexts = dbContextsRes

		log.Printf("db.LoadContexts completed successfully, loaded %d contexts", len(dbContexts))

		// Log information about loaded map contexts
		mapContextsCount := 0
		for _, context := range dbContexts {
			if context.ContextType == shared.ContextMapType {
				mapContextsCount++
				log.Printf("Loaded map context: %s, path: %s, tokens: %d", context.Name, context.FilePath, context.NumTokens)
			}
		}
		if mapContextsCount > 0 {
			log.Printf("Successfully loaded %d map contexts out of %d total contexts", mapContextsCount, len(dbContexts))
		}

		if loadRes.MaxTokensExceeded {
			return nil
		}

		err = repo.GitAddAndCommit(branchName, res.Msg)

		if err != nil {
			return fmt.Errorf("error committing changes: %v", err)
		}

		return nil
	})

	if err != nil {
		log.Printf("Error loading contexts: %v\n", err)
		http.Error(w, "Error loading contexts: "+err.Error(), http.StatusInternalServerError)
		return nil, nil
	}

	if loadRes.MaxTokensExceeded {
		log.Printf("The total number of tokens (%d) exceeds the maximum allowed (%d)", loadRes.TotalTokens, loadRes.MaxTokens)
		bytes, err := json.Marshal(loadRes)

		if err != nil {
			log.Printf("Error marshalling response: %v\n", err)
			http.Error(w, "Error marshalling response: "+err.Error(), http.StatusInternalServerError)
			return nil, nil
		}

		w.Write(bytes)
		return nil, nil
	}

	return loadRes, dbContexts
}
