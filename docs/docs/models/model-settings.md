---
sidebar_position: 3
sidebar_label: Settings
---

# Model Settings

Gpt4cli gives you a number of ways to control the models and models settings used in your plans. Changes to models and model settings are [version controlled](../core-concepts/version-control.md) and can be [branched](../core-concepts/branches.md).

## `models` and `set-model`

You can see the current plan's models and model settings with the `models` command and change them with the `set-model` command.

```bash
gpt4cli models # show the current AI models and model settings
gpt4cli models available # show all available models
gpt4cli set-model # select from a list of models, model packs, and settings
gpt4cli set-model planner openrouter/anthropic/claude-3.7-sonnet # set the main planner model to Claude Sonnet 3.7 from OpenRouter.ai
gpt4cli set-model builder temperature 0.1 # set the builder model's temperature to 0.1
gpt4cli set-model max-tokens 4000 # set the planner model overall token limit to 4000
gpt4cli set-model max-convo-tokens 20000  # set how large the conversation can grow before Gpt4cli starts using summaries
```

## Model Defaults 

`set-model` updates model settings for the current plan. If you want to change the default model settings for all new plans, use `set-model default`.

```bash
gpt4cli models default # show the default model settings
gpt4cli set-model default # select from a list of models and settings
gpt4cli set-model default planner openai/gpt-4o # set the default planner model to OpenAI gpt-4o
```

## Model Packs

Instead of changing models for each role one by one, model packs let you switch out all roles at once. It's the recommended way to manage models.

You can list available model packs with `model-packs`:

```bash
gpt4cli model-packs # list all available model packs
```

You can create your own model packs with `model-packs create`, list built-in and custom model packs with `model-packs`, and remove custom model packs with `model-packs delete`.

```bash
gpt4cli set-model # select from a list of model packs for the current plan
gpt4cli set-model default # select from a list of model packs to set as the default for all new plans
gpt4cli set-model anthropic-claude-3.5-sonnet-gpt-4o # set the current plan's model pack by name
gpt4cli set-model default Mixtral-8x22b/Mixtral-8x7b/gpt-4o # set the default model pack for all new plans

gpt4cli model-packs # list built-in and custom model packs
gpt4cli model-packs create # create a new custom model pack
gpt4cli model-packs --custom # list only custom model packs
```

## Custom Models

Use `models add` to add a custom model and use any provider that is compatible with OpenAI, including OpenRouter.ai, Together.ai, Ollama, Replicate, and more.

```bash
gpt4cli models add # add a custom model
gpt4cli models available --custom # show all available custom models
gpt4cli models delete # delete a custom model
```
