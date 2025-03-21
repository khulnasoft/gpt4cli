package model

import (
	"context"
	"fmt"
	"gpt4cli-server/db"
	"gpt4cli-server/model/prompts"
	"time"

	"github.com/khulnasoft/gpt4cli/shared"
	"github.com/sashabaranov/go-openai"
)

type PlanSummaryParams struct {
	Conversation                []*openai.ChatCompletionMessage
	LatestConvoMessageId        string
	LatestConvoMessageCreatedAt time.Time
	NumMessages                 int
	OrgId                       string
	PlanId                      string
}

func PlanSummary(client *openai.Client, config shared.ModelRoleConfig, params PlanSummaryParams, ctx context.Context) (*db.ConvoSummary, error) {
	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleSystem,
			Content: prompts.Identity,
		},
	}

	for _, message := range params.Conversation {
		messages = append(messages, *message)
	}

	messages = append(messages, openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleUser,
		Content: prompts.PlanSummary,
	})

	fmt.Println("summarizing messages:")
	// spew.Dump(messages)

	resp, err := CreateChatCompletionWithRetries(
		client,
		ctx,
		openai.ChatCompletionRequest{
			Model:       config.BaseModelConfig.ModelName,
			Messages:    messages,
			Temperature: config.Temperature,
			TopP:        config.TopP,
		},
	)

	if err != nil {
		fmt.Println("PlanSummary err:", err)

		return nil, err
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no response from GPT")
	}

	content := resp.Choices[0].Message.Content

	// log.Println("Plan summary content:")
	// log.Println(content)

	return &db.ConvoSummary{
		OrgId:                       params.OrgId,
		PlanId:                      params.PlanId,
		Summary:                     "## Summary of the plan so far:\n\n" + content,
		Tokens:                      resp.Usage.CompletionTokens,
		LatestConvoMessageId:        params.LatestConvoMessageId,
		LatestConvoMessageCreatedAt: params.LatestConvoMessageCreatedAt,
		NumMessages:                 params.NumMessages,
	}, nil

}
