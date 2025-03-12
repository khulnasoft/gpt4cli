package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"gpt4cli-server/db"
	"gpt4cli-server/syntax/file_map"
	"runtime"
	"sync"

	shared "gpt4cli-shared"

	"github.com/gorilla/mux"
)

func GetFileMapHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for GetFileMapHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		log.Println("GetFileMapHandler: auth failed")
		return
	}

	var req shared.GetFileMapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
		return
	}

	log.Println("GetFileMapHandler: checking limits")

	if len(req.MapInputs) > shared.MaxContextMapPaths {
		http.Error(w, fmt.Sprintf("Too many files to map: %d (max %d)", len(req.MapInputs), shared.MaxContextMapPaths), http.StatusBadRequest)
		return
	}

	totalSize := 0
	for _, input := range req.MapInputs {
		totalSize += len(input)
	}

	// Allow a little extra space for empty file maps after the total size limit is exceeded
	// On the client, once the total size limit is exceeded, we send empty file maps for remaining files
	if totalSize > shared.MaxContextMapInputSize+10000 {
		http.Error(w, fmt.Sprintf("Max map size exceeded: %d (max %d)", totalSize, shared.MaxContextMapInputSize), http.StatusBadRequest)
		return
	}

	maps := make(shared.FileMapBodies)

	// Use half of available CPUs
	cpus := runtime.NumCPU()
	log.Printf("GetFileMapHandler: Available CPUs: %d", cpus)
	maxWorkers := cpus / 2
	if maxWorkers < 1 {
		maxWorkers = 1 // Ensure at least one worker
	}
	log.Printf("GetFileMapHandler: Max workers: %d", maxWorkers)
	sem := make(chan struct{}, maxWorkers)
	wg := sync.WaitGroup{}
	var mu sync.Mutex

	log.Printf("GetFileMapHandler: len(req.MapInputs): %d", len(req.MapInputs))

	for path, input := range req.MapInputs {
		if !shared.HasFileMapSupport(path) {
			mu.Lock()
			maps[path] = "[NO MAP]"
			mu.Unlock()
			continue
		}

		wg.Add(1)
		sem <- struct{}{}
		go func(path string, input string) {
			defer wg.Done()
			defer func() { <-sem }()
			fileMap, err := file_map.MapFile(r.Context(), path, []byte(input))
			if err != nil {
				// Skip files that can't be parsed, just log the error
				log.Printf("Error mapping file %s: %v", path, err)
				return
			}
			mu.Lock()
			maps[path] = fileMap.String()
			mu.Unlock()
		}(path, input)
	}
	wg.Wait()

	log.Printf("GetFileMapHandler: len(maps): %d", len(maps))

	resp := shared.GetFileMapResponse{
		MapBodies: maps,
	}

	respBytes, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("GetFileMapHandler success - writing response bytes: %d", len(respBytes))

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBytes)
}

func LoadCachedFileMapHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for LoadCachedFileMapHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	vars := mux.Vars(r)
	planId := vars["planId"]
	branchName := vars["branch"]
	log.Println("planId: ", planId, "branchName: ", branchName)

	plan := authorizePlan(w, planId, auth)

	if plan == nil {
		return
	}

	var req shared.LoadCachedFileMapRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, fmt.Sprintf("Error decoding request: %v", err), http.StatusBadRequest)
		return
	}

	cachedMetaByPath := map[string]*shared.Context{}
	cachedMapsByPath := map[string]*db.CachedMap{}
	var mu sync.Mutex
	errCh := make(chan error, len(req.FilePaths))

	for _, path := range req.FilePaths {
		go func(path string) {
			cachedContext, err := db.GetCachedMap(plan.OrgId, plan.ProjectId, path)
			if err != nil {
				errCh <- fmt.Errorf("error getting cached map: %v", err)
				return
			}
			if cachedContext != nil {
				mu.Lock()
				cachedMetaByPath[path] = cachedContext.ToMeta().ToApi()
				cachedMapsByPath[path] = &db.CachedMap{
					MapParts:  cachedContext.MapParts,
					MapShas:   cachedContext.MapShas,
					MapTokens: cachedContext.MapTokens,
				}
				mu.Unlock()
			}
			errCh <- nil
		}(path)
	}

	for range req.FilePaths {
		err := <-errCh
		if err != nil {
			http.Error(w, fmt.Sprintf("Error getting cached map: %v", err), http.StatusInternalServerError)
			return
		}
	}

	resp := shared.LoadCachedFileMapResponse{}

	var loadRes *shared.LoadContextResponse
	if len(cachedMetaByPath) == 0 {
		log.Println("no cached maps found")
	} else {
		log.Println("cached map found")

		cachedByPath := map[string]bool{}
		for _, cachedContext := range cachedMetaByPath {
			cachedByPath[cachedContext.FilePath] = true
		}
		resp.CachedByPath = cachedByPath

		var loadReq shared.LoadContextRequest
		for _, cachedContext := range cachedMetaByPath {
			loadReq = append(loadReq, &shared.LoadContextParams{
				ContextType: shared.ContextMapType,
				Name:        cachedContext.Name,
				FilePath:    cachedContext.FilePath,
				Body:        cachedContext.Body,
			})
		}

		loadRes, _ = loadContexts(loadContextsParams{
			w:                w,
			r:                r,
			auth:             auth,
			loadReq:          &loadReq,
			plan:             plan,
			branchName:       branchName,
			cachedMapsByPath: cachedMapsByPath,
		})

		if loadRes == nil {
			log.Println("LoadCachedFileMapHandler - loadRes is nil")
			return
		}

		resp.LoadRes = loadRes
	}

	bytes, err := json.Marshal(resp)
	if err != nil {
		log.Printf("Error marshalling response: %v", err)
		http.Error(w, fmt.Sprintf("Error marshalling response: %v", err), http.StatusInternalServerError)
		return
	}

	w.Write(bytes)
}
