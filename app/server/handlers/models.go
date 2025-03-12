package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"gpt4cli-server/db"

	shared "gpt4cli-shared"

	"github.com/gorilla/mux"
)

func CreateCustomModelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for CreateCustomModelHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	var model shared.AvailableModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		log.Printf("Error decoding request body: %v\n", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if os.Getenv("IS_CLOUD") != "" && model.Provider == shared.ModelProviderCustom {
		http.Error(w, "Custom model providers are not supported on Gpt4cli Cloud", http.StatusBadRequest)
		return
	}

	dbModel := &db.AvailableModel{
		Id:                    model.Id,
		OrgId:                 auth.OrgId,
		Provider:              model.Provider,
		CustomProvider:        model.CustomProvider,
		BaseUrl:               model.BaseUrl,
		ModelName:             model.ModelName,
		Description:           model.Description,
		MaxTokens:             model.MaxTokens,
		ApiKeyEnvVar:          model.ApiKeyEnvVar,
		HasImageSupport:       model.HasImageSupport,
		DefaultMaxConvoTokens: model.DefaultMaxConvoTokens,
		MaxOutputTokens:       model.MaxOutputTokens,
		ReservedOutputTokens:  model.ReservedOutputTokens,
	}

	if err := db.CreateCustomModel(dbModel); err != nil {
		log.Printf("Error creating custom model: %v\n", err)
		http.Error(w, "Failed to create custom model: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Println("Successfully created custom model")
}

func ListCustomModelsHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for ListCustomModelsHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	models, err := db.ListCustomModels(auth.OrgId)
	if err != nil {
		log.Printf("Error fetching custom models: %v\n", err)
		http.Error(w, "Failed to fetch custom models: "+err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(models)

	log.Println("Successfully fetched custom models")
}

func DeleteAvailableModelHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for DeleteAvailableModelHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	modelId := mux.Vars(r)["modelId"]
	if err := db.DeleteAvailableModel(modelId); err != nil {
		log.Printf("Error deleting custom model: %v\n", err)
		http.Error(w, "Failed to delete custom model: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Println("Successfully deleted custom model")
}

func CreateModelPackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for CreateModelPackHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	var ms shared.ModelPack
	if err := json.NewDecoder(r.Body).Decode(&ms); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	dbMs := &db.ModelPack{
		OrgId:       auth.OrgId,
		Name:        ms.Name,
		Description: ms.Description,
		PlanSummary: ms.PlanSummary,
		Builder:     ms.Builder,
		Namer:       ms.Namer,
		CommitMsg:   ms.CommitMsg,
		ExecStatus:  ms.ExecStatus,
		Architect:   ms.Architect,
		Coder:       ms.Coder,
	}

	if err := db.CreateModelPack(dbMs); err != nil {
		log.Printf("Error creating model pack: %v\n", err)
		http.Error(w, "Failed to create model pack: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)

	log.Println("Successfully created model pack")
}

func ListModelPacksHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for ListModelPacksHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	sets, err := db.ListModelPacks(auth.OrgId)
	if err != nil {
		log.Printf("Error fetching model packs: %v\n", err)
		http.Error(w, "Failed to fetch model packs: "+err.Error(), http.StatusInternalServerError)
		return
	}

	var apiPacks []*shared.ModelPack

	for _, mp := range sets {
		apiPacks = append(apiPacks, mp.ToApi())
	}

	json.NewEncoder(w).Encode(apiPacks)

	log.Println("Successfully fetched model packs")
}

func DeleteModelPackHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Received request for DeleteModelPackHandler")

	auth := Authenticate(w, r, true)
	if auth == nil {
		return
	}

	setId := mux.Vars(r)["setId"]

	log.Printf("Deleting model pack with id: %s\n", setId)

	if err := db.DeleteModelPack(setId); err != nil {
		log.Printf("Error deleting model pack: %v\n", err)
		http.Error(w, "Failed to delete model pack: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	log.Println("Successfully deleted model pack")
}
