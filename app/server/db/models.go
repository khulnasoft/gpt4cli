package db

import (
	"fmt"
)

func CreateCustomModel(model *AvailableModel) error {
	query := `INSERT INTO custom_models (org_id, provider, custom_provider, base_url, model_name, description, max_tokens, api_key_env_var, default_max_convo_tokens, max_output_tokens, reserved_output_tokens, preferred_output_format, has_image_support) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	RETURNING id, created_at, updated_at`

	err := Conn.QueryRow(query, model.OrgId, model.Provider, model.CustomProvider, model.BaseUrl, model.ModelName, model.Description, model.MaxTokens, model.ApiKeyEnvVar, model.DefaultMaxConvoTokens, model.MaxOutputTokens, model.ReservedOutputTokens, model.PreferredOutputFormat, model.HasImageSupport).Scan(&model.Id, &model.CreatedAt, &model.UpdatedAt)
	if err != nil {
		return fmt.Errorf("error inserting new custom model: %v", err)
	}

	return nil
}

func ListCustomModels(orgId string) ([]*AvailableModel, error) {
	var models []*AvailableModel

	query := `SELECT * FROM custom_models WHERE org_id = $1`

	err := Conn.Select(&models, query, orgId)

	if err != nil {
		return nil, fmt.Errorf("error fetching custom models: %v", err)
	}

	return models, nil
}

func DeleteAvailableModel(modelId string) error {
	query := `DELETE FROM custom_models WHERE id = $1`
	_, err := Conn.Exec(query, modelId)

	if err != nil {
		return fmt.Errorf("error deleting custom model: %v", err)
	}

	return nil
}

func CreateModelPack(ms *ModelPack) error {
	query := `INSERT INTO model_sets (org_id, name, description, planner, plan_summary, builder, namer, commit_msg, exec_status, context_loader, coder) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	RETURNING id, created_at`

	err := Conn.QueryRow(query, ms.OrgId, ms.Name, ms.Description, ms.Planner, ms.PlanSummary, ms.Builder, ms.Namer, ms.CommitMsg, ms.ExecStatus, ms.Architect, ms.Coder).Scan(&ms.Id, &ms.CreatedAt)

	if err != nil {
		return fmt.Errorf("error inserting new model pack: %v", err)
	}

	return nil
}

func ListModelPacks(orgId string) ([]*ModelPack, error) {
	var modelPacks []*ModelPack

	query := `SELECT * FROM model_sets WHERE org_id = $1`
	err := Conn.Select(&modelPacks, query, orgId)

	if err != nil {
		return nil, fmt.Errorf("error fetching model packs: %v", err)
	}

	return modelPacks, nil
}

func DeleteModelPack(setId string) error {
	query := `DELETE FROM model_sets WHERE id = $1`
	_, err := Conn.Exec(query, setId)

	if err != nil {
		return fmt.Errorf("error deleting model pack: %v", err)
	}

	return nil
}
