package handlers

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"gpt4cli-server/db"
	"gpt4cli-server/hooks"
	"gpt4cli-server/host"
	modelPlan "gpt4cli-server/model/plan"
	"gpt4cli-server/types"
	"time"

	"github.com/gorilla/mux"
	"github.com/khulnasoft/gpt4cli/shared"
)

const TrialMaxReplies = 10

func TellPlanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for TellPlanHandler", "ip:", host.Ip)

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	vars := mux.Vars(r)
	planId := vars["planId"]
	branch := vars["branch"]

	log.Println("planId: ", planId)

	plan := authorizePlanExecUpdate(w, planId, auth)
	if plan == nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer func() {
		log.Println("Closing request body")
		r.Body.Close()
	}()

	var requestBody shared.TellPlanRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	_, apiErr := hooks.ExecHook(hooks.WillTellPlan, hooks.HookParams{
		Auth: auth,
		Plan: plan,
	})
	if apiErr != nil {
		writeApiError(w, *apiErr)
		return
	}

	clients := initClients(
		initClientsParams{
			w:           w,
			auth:        auth,
			apiKey:      requestBody.ApiKey,
			apiKeys:     requestBody.ApiKeys,
			endpoint:    requestBody.Endpoint,
			openAIBase:  requestBody.OpenAIBase,
			openAIOrgId: requestBody.OpenAIOrgId,
			plan:        plan,
		},
	)
	err = modelPlan.Tell(clients, plan, branch, auth, &requestBody)

	if err != nil {
		log.Printf("Error telling plan: %v\n", err)
		http.Error(w, "Error telling plan", http.StatusInternalServerError)
		return
	}

	if requestBody.ConnectStream {
		startResponseStream(w, auth, planId, branch, false)
	}

	log.Println("Successfully processed request for TellPlanHandler")
}

func BuildPlanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for BuildPlanHandler", "ip:", host.Ip)
	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	vars := mux.Vars(r)
	planId := vars["planId"]
	branch := vars["branch"]

	log.Println("planId: ", planId)
	plan := authorizePlanExecUpdate(w, planId, auth)
	if plan == nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer func() {
		log.Println("Closing request body")
		r.Body.Close()
	}()

	var requestBody shared.BuildPlanRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	clients := initClients(
		initClientsParams{
			w:           w,
			auth:        auth,
			apiKey:      requestBody.ApiKey,
			apiKeys:     requestBody.ApiKeys,
			endpoint:    requestBody.Endpoint,
			openAIBase:  requestBody.OpenAIBase,
			openAIOrgId: requestBody.OpenAIOrgId,
			plan:        plan,
		},
	)
	numBuilds, err := modelPlan.Build(clients, plan, branch, auth)

	if err != nil {
		log.Printf("Error building plan: %v\n", err)
		http.Error(w, "Error building plan", http.StatusInternalServerError)
		return
	}

	if numBuilds == 0 {
		log.Println("No builds were executed")
		http.Error(w, shared.NoBuildsErr, http.StatusNotFound)
		return
	}

	if requestBody.ConnectStream {
		startResponseStream(w, auth, planId, branch, false)
	}

	log.Println("Successfully processed request for BuildPlanHandler")
}

func ConnectPlanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for ConnectPlanHandler", "ip:", host.Ip)

	vars := mux.Vars(r)
	planId := vars["planId"]
	branch := vars["branch"]
	log.Println("planId: ", planId)
	log.Println("branch: ", branch)
	active := modelPlan.GetActivePlan(planId, branch)
	isProxy := r.URL.Query().Get("proxy") == "true"

	if active == nil {
		if isProxy {
			log.Println("No active plan on proxied request")
			http.Error(w, "No active plan", http.StatusNotFound)
			return
		}

		log.Println("No active plan -- proxying request")

		proxyActivePlanMethod(w, r, planId, branch, "connect")
		return
	}

	auth := Authenticate(w, r, true)
	if auth == nil {
		log.Println("No auth")
		return
	}

	plan := authorizePlan(w, planId, auth)
	if plan == nil {
		log.Println("No plan")
		return
	}

	startResponseStream(w, auth, planId, branch, true)

	log.Println("Successfully processed request for ConnectPlanHandler")
}

func StopPlanHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for StopPlanHandler", "ip:", host.Ip)

	vars := mux.Vars(r)
	planId := vars["planId"]
	branch := vars["branch"]
	log.Println("planId: ", planId)
	log.Println("branch: ", branch)
	active := modelPlan.GetActivePlan(planId, branch)
	isProxy := r.URL.Query().Get("proxy") == "true"

	if active == nil {
		if isProxy {
			log.Println("No active plan on proxied request")
			http.Error(w, "No active plan", http.StatusNotFound)
			return
		}
		proxyActivePlanMethod(w, r, planId, branch, "stop")
		return
	}

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	if authorizePlan(w, planId, auth) == nil {
		return
	}

	log.Println("Sending stream aborted message to client")

	active.Stream(shared.StreamMessage{
		Type: shared.StreamMessageAborted,
	})

	// give some time for stream message to be processed before canceling
	log.Println("Sleeping for 100ms before canceling")
	time.Sleep(100 * time.Millisecond)

	var err error
	ctx, cancel := context.WithCancel(context.Background())
	unlockFn := LockRepo(w, r, auth, db.LockScopeWrite, ctx, cancel, true)
	if unlockFn == nil {
		return
	} else {
		defer func() {
			(*unlockFn)(err)

			if err == nil {
				err = modelPlan.Stop(planId, branch, auth.User.Id, auth.OrgId)

				if err != nil {
					log.Printf("Error stopping plan: %v\n", err)
					http.Error(w, "Error stopping plan", http.StatusInternalServerError)
					return
				}

				log.Println("Successfully processed request for StopPlanHandler")
			}
		}()
	}

	log.Println("Stopping plan")
	err = modelPlan.StorePartialReply(planId, branch, auth.User.Id, auth.OrgId)

	if err != nil {
		log.Printf("Error storing partial reply: %v\n", err)
		http.Error(w, "Error storing partial reply", http.StatusInternalServerError)
		return
	}
}

func RespondMissingFileHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for RespondMissingFileHandler", "ip:", host.Ip)

	vars := mux.Vars(r)
	planId := vars["planId"]
	branch := vars["branch"]
	log.Println("planId: ", planId)
	log.Println("branch: ", branch)
	isProxy := r.URL.Query().Get("proxy") == "true"

	active := modelPlan.GetActivePlan(planId, branch)
	if active == nil {
		if isProxy {
			log.Println("No active plan on proxied request")
			http.Error(w, "No active plan", http.StatusNotFound)
			return
		}

		proxyActivePlanMethod(w, r, planId, branch, "respond_missing_file")
		return
	}

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	plan := authorizePlan(w, planId, auth)
	if plan == nil {
		return
	}

	// read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestBody shared.RespondMissingFileRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	log.Println("missing file choice:", requestBody.Choice)

	if requestBody.Choice == shared.RespondMissingFileChoiceLoad {
		log.Println("loading missing file")
		res, dbContexts := loadContexts(w, r, auth, &shared.LoadContextRequest{
			&shared.LoadContextParams{
				ContextType: shared.ContextFileType,
				Name:        requestBody.FilePath,
				FilePath:    requestBody.FilePath,
				Body:        requestBody.Body,
			},
		}, plan, branch, nil)
		if res == nil {
			return
		}

		dbContext := dbContexts[0]

		log.Println("loaded missing file:", dbContext.FilePath)

		modelPlan.UpdateActivePlan(planId, branch, func(activePlan *types.ActivePlan) {
			activePlan.Contexts = append(activePlan.Contexts, dbContext)
			activePlan.ContextsByPath[dbContext.FilePath] = dbContext
		})
	}

	// This will resume model stream
	log.Println("Resuming model stream")
	active.MissingFileResponseCh <- requestBody.Choice

	log.Println("Successfully processed request for RespondMissingFileHandler")
}

func AutoLoadContextHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for AutoLoadContextHandler", "ip:", host.Ip)

	vars := mux.Vars(r)
	planId := vars["planId"]
	branch := vars["branch"]
	log.Println("planId: ", planId)
	log.Println("branch: ", branch)

	isProxy := r.URL.Query().Get("proxy") == "true"

	active := modelPlan.GetActivePlan(planId, branch)
	if active == nil {
		if isProxy {
			log.Println("No active plan on proxied request")
			http.Error(w, "No active plan", http.StatusNotFound)
			return
		}

		proxyActivePlanMethod(w, r, planId, branch, "auto_load_context")
		return
	}

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	plan := authorizePlan(w, planId, auth)
	if plan == nil {
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("Error reading request body: %v\n", err)
		http.Error(w, "Error reading request body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var requestBody shared.LoadContextRequest
	if err := json.Unmarshal(body, &requestBody); err != nil {
		log.Printf("Error parsing request body: %v\n", err)
		http.Error(w, "Error parsing request body", http.StatusBadRequest)
		return
	}

	log.Println("AutoLoadContextHandler - loading contexts")
	res, dbContexts := loadContexts(w, r, auth, &requestBody, plan, branch, nil)

	if res == nil {
		return
	}

	log.Println("AutoLoadContextHandler - updating active plan")

	modelPlan.UpdateActivePlan(planId, branch, func(activePlan *types.ActivePlan) {
		activePlan.Contexts = append(activePlan.Contexts, dbContexts...)
		for _, dbContext := range dbContexts {
			activePlan.ContextsByPath[dbContext.FilePath] = dbContext
		}
	})

	var apiContexts []*shared.Context
	for _, dbContext := range dbContexts {
		apiContexts = append(apiContexts, dbContext.ToApi())
	}

	msg := shared.SummaryForLoadContext(apiContexts, res.TokensAdded, res.TotalTokens)
	markdownRes := shared.LoadContextResponse{
		TokensAdded:       res.TokensAdded,
		TotalTokens:       res.TotalTokens,
		MaxTokensExceeded: res.MaxTokensExceeded,
		MaxTokens:         res.MaxTokens,
		Msg:               msg,
	}

	bytes, err := json.Marshal(markdownRes)
	if err != nil {
		log.Printf("Error marshalling response: %v\n", err)
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}

	w.Write(bytes)

	active.AutoLoadContextCh <- struct{}{}

	log.Println("Successfully processed request for AutoLoadContextHandler")
}

func authorizePlanExecUpdate(w http.ResponseWriter, planId string, auth *types.ServerAuth) *db.Plan {
	plan := authorizePlan(w, planId, auth)
	if plan == nil {
		return nil
	}

	if plan.OwnerId != auth.User.Id && !auth.HasPermission(shared.PermissionUpdateAnyPlan) {
		log.Println("User does not have permission to update plan")
		http.Error(w, "User does not have permission to update plan", http.StatusForbidden)
		return nil
	}

	return plan
}
