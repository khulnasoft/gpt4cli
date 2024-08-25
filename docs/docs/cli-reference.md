---
sidebar_position: 8
sidebar_label: CLI Reference
---

# CLI Reference

All Gpt4cli CLI commands and their options.

## Usage

```bash
gpt4cli [command] [flags]
pdx [command] [flags] # 'pdx' is an alias for 'gpt4cli'
```

## Help

Built-in help.

```bash
gpt4cli help
pdx h # alias
```

`--all/-a`: List all commands.

For help on a specific command, use:

```bash
gpt4cli [command] --help
```

## Plans

### new

Start a new plan.

```bash
gpt4cli new
gpt4cli new -n new-plan # with name
```

`--name/-n`: Name of the new plan. The name is generated automatically after first prompt if no name is specified on creation.

### plans

List plans. Output includes index, when each plan was last updated, the current branch of each plan, the number of tokens in context, and the number of tokens in the conversation (prior to summarization).

Includes full details on plans in current directory. Also includes names of plans in parent directories and child directories.

```bash
gpt4cli plans
gpt4cli plans --archived # list archived plans only

pdx pl # alias
```

`--archived/-a`: List archived plans only.

### current

Show current plan. Output includes when the plan was last updated and created, the current branch, the number of tokens in context, and the number of tokens in the conversation (prior to summarization).

```bash
gpt4cli current
pdx cu # alias
```

### cd

Set current plan by name or index.

```bash
gpt4cli cd # select from a list of plans
gpt4cli cd some-plan # by name
gpt4cli cd 4 # by index in `gpt4cli plans`
```

With no arguments, Gpt4cli prompts you with a list of plans to select from.

With one argument, Gpt4cli selects a plan by name or by index in the `gpt4cli plans` list.

### delete-plan

Delete a plan by name or index.

```bash
gpt4cli delete-plan # select from a list of plans
gpt4cli delete-plan some-plan # by name
gpt4cli delete-plan 4 # by index in `gpt4cli plans`

pdx dp # alias
```

With no arguments, Gpt4cli prompts you with a list of plans to select from.

With one argument, Gpt4cli deletes a plan by name or by index in the `gpt4cli plans` list.

### rename

Rename the current plan.

```bash
gpt4cli rename # prompt for new name
gpt4cli rename new-name # set new name
```

With no arguments, Gpt4cli prompts you for a new name.

With one argument, Gpt4cli sets the new name.

### archive

Archive a plan.

```bash
gpt4cli archive # select from a list of plans
gpt4cli archive some-plan # by name
gpt4cli archive 4 # by index in `gpt4cli plans`

pdx arc # alias
```

With no arguments, Gpt4cli prompts you with a list of plans to select from.

With one argument, Gpt4cli archives a plan by name or by index in the `gpt4cli plans` list.

### unarchive

Unarchive a plan.

```bash
gpt4cli unarchive # select from a list of archived plans
gpt4cli unarchive some-plan # by name
gpt4cli unarchive 4 # by index in `gpt4cli plans --archived`
pdx unarc # alias
```

## Context

### load

Load files, directories, directory layouts, URLs, notes, images, or piped data into context.

```bash
gpt4cli load component.ts # single file
gpt4cli load component.ts action.ts reducer.ts # multiple files
gpt4cli load lib -r # loads lib and all its subdirectories
gpt4cli load tests/**/*.ts # loads all .ts files in tests and its subdirectories
gpt4cli load . --tree # loads the layout of the current directory and its subdirectories (file names only)
gpt4cli load https://redux.js.org/usage/writing-tests # loads the text-only content of the url
npm test | gpt4cli load # loads the output of `npm test`
gpt4cli load -n 'add logging statements to all the code you generate.' # load a note into context
gpt4cli load ui-mockup.png # load an image into context

pdx l component.ts # alias
```

`--recursive/-r`: Load an entire directory and all its subdirectories.

`--tree`: Load directory tree layout with file names only.

`--note/-n`: Load a note into context.

`--force/-f`: Load files even when ignored by .gitignore or .gpt4cliignore.

`--detail/-d`: Image detail level when loading an image (high or low)—default is high. See https://platform.openai.com/docs/guides/vision/low-or-high-fidelity-image-understanding for more info.

### ls

List everything in the current plan's context. Output includes index, name, type, token size, when the context added, and when the context was last updated.

```bash
gpt4cli ls

gpt4cli list-context # longer alias
```

### rm

Remove context by index, range, name, or glob.

```bash
gpt4cli rm some-file.ts # by name
gpt4cli rm app/**/*.ts # by glob pattern
gpt4cli rm 4 # by index in `gpt4cli ls`
plandx rm 2-4 # by range of indices

gpt4cli remove # longer alias
gpt4cli unload # longer alias
```

### update

Update any outdated context.

```bash
gpt4cli update
pdx u # alias
```

### clear

Remove all context.

```bash
gpt4cli clear
```

## Control

### tell

Describe a task, ask a question, or chat.

```bash
gpt4cli tell -f prompt.txt # from file
gpt4cli tell # open vim to write prompt
gpt4cli tell "add a cancel button to the left of the submit button" # inline

pdx t # alias
```

`--file/-f`: File path containing prompt.

`--stop/-s`: Stop after a single model response (don't auto-continue).

`--no-build/-n`: Don't build proposed changes into pending file updates.

`--bg`: Run task in the background.

### continue

Continue the plan.

```bash
gpt4cli continue

pdx c # alias
```

`--stop/-s`: Stop after a single model response (don't auto-continue).

`--no-build/-n`: Don't build proposed changes into pending file updates.

`--bg`: Run task in the background.

### build

Build any unbuilt pending changes from the plan conversation.

```bash
gpt4cli build
pdx b # alias
```

`--bg`: Build in the background.

## Changes

### diff

Review pending changes in 'git diff' format.

```bash
gpt4cli diff
```

### changes

Review pending changes in a TUI.

```bash
gpt4cli changes
```

### apply

Apply pending changes to project files.

```bash
gpt4cli apply
pdx ap # alias
```

`--yes/-y`: Skip confirmation.

### reject

Reject pending changes to one or more project files.

```bash
gpt4cli reject file.ts # one file
gpt4cli reject file.ts another-file.ts # multiple files
gpt4cli reject --all # all pending files

pdx rj file.ts # alias
```

`--all/-a`: Reject all pending files.

## History

### log

Show plan history.

```bash
gpt4cli log

gpt4cli history # alias
gpt4cli logs # alias
```

### rewind

Rewind to a previous state.

```bash
gpt4cli rewind # rewind 1 step
gpt4cli rewind 3 # rewind 3 steps
gpt4cli rewind a7c8d66 # rewind to a specific step from `gpt4cli log`
```

With no arguments, Gpt4cli rewinds one step.

With one argument, Gpt4cli rewinds the specified number of steps (if an integer is passed) or rewinds to the specified step (if a hash from `gpt4cli log` is passed).

### convo

Show the current plan's conversation.

```bash
gpt4cli convo
gpt4cli convo 1 # show a specific message
gpt4cli convo 1-5 # show a range of messages
gpt4cli convo 3- # show all messages from 3 to the end
```

`--plain/-p`: Output conversation in plain text with no ANSI codes.

### summary

Show the latest summary of the current plan.

```bash
gpt4cli summary
```

`--plain/-p`: Output summary in plain text with no ANSI codes.

## Branches

### branches

List plan branches. Output includes index, name, when the branch was last updated, the number of tokens in context, and the number of tokens in the conversation (prior to summarization).

```bash
gpt4cli branches
pdx br # alias
```

### checkout

Checkout or create a branch.

```bash
gpt4cli checkout # select from a list of branches or prompt to create a new branch
gpt4cli checkout some-branch # checkout by name or create a new branch with that name

pdx co # alias
```

### delete-branch

Delete a branch by name or index.

```bash
gpt4cli delete-branch # select from a list of branches
gpt4cli delete-branch some-branch # by name
gpt4cli delete-branch 4 # by index in `gpt4cli branches`

pdx db # alias
```

With no arguments, Gpt4cli prompts you with a list of branches to select from.

With one argument, Gpt4cli deletes a branch by name or by index in the `gpt4cli branches` list.

## Background Tasks / Streams

### ps

List active and recently finished plan streams. Output includes stream ID, plan name, branch name, when the stream was started, and the stream's status (active, finished, stopped, errored, or waiting for a missing file to be selected).

```bash
gpt4cli ps
```

### connect

Connect to an active plan stream.

```bash
gpt4cli connect # select from a list of active streams
gpt4cli connect a4de # by stream ID in `gpt4cli ps`
gpt4cli connect some-plan main # by plan name and branch name
```

With no arguments, Gpt4cli prompts you with a list of active streams to select from.

With one argument, Gpt4cli connects to a stream by stream ID in the `gpt4cli ps` list.

With two arguments, Gpt4cli connects to a stream by plan name and branch name.

### stop

Stop an active plan stream.

```bash
gpt4cli stop # select from a list of active streams
gpt4cli stop a4de # by stream ID in `gpt4cli ps`
gpt4cli stop some-plan main # by plan name and branch name
```

With no arguments, Gpt4cli prompts you with a list of active streams to select from.

With one argument, Gpt4cli connects to a stream by stream ID in the `gpt4cli ps` list.

With two arguments, Gpt4cli connects to a stream by plan name and branch name.

## Models

### models

Show current plan models and model settings.

```bash
gpt4cli models
```

### models default

Show org-wide default models and model settings for new plans.

```bash
gpt4cli models default
```

### models available

Show available models.

```bash
gpt4cli models available # show all available models
gpt4cli models available --custom # show available custom models only
```

`--custom`: Show available custom models only.

### set-model

Update current plan models or model settings.

```bash
gpt4cli set-model # select from a list of models and settings
gpt4cli set-model planner openai/gpt-4 # set the model for a role
gpt4cli set-model gpt-4-turbo-latest # set the current plan's model pack by name (sets all model roles at once—see `model-packs` below)
gpt4cli set-model builder temperature 0.1 # set a model setting for a role
gpt4cli set-model max-tokens 4000 # set the planner model overall token limit to 4000
gpt4cli set-model max-convo-tokens 20000  # set how large the conversation can grow before Gpt4cli starts using summaries
```

With no arguments, Gpt4cli prompts you to select from a list of models and settings.

With arguments, can take one of the following forms:

- `gpt4cli set-model [role] [model]`: Set the model for a role.
- `gpt4cli set-model [model-pack]`: Set the current plan's model pack by name.
- `gpt4cli set-model [role] [setting] [value]`: Set a model setting for a role.
- `gpt4cli set-model [setting] [value]`: Set a model setting for the current plan.

Models are specified as `provider/model-name`, e.g. `openai/gpt-4`, `openrouter/anthropic/claude-opus-3`, `together/mistralai/Mixtral-8x22B-Instruct-v0.1`, etc.

See all the model roles [here](./models/roles.md).

Model role settings:

- `temperature`: Higher temperature means more randomness, which can produce more creativity but also more errors.
- `top-p`: Top-p sampling is a way to prevent the model from generating improbable text by only considering the most likely tokens.

Plan settings:

- `max-tokens`: The overall token limit for the planner model.
- `max-convo-tokens`: How large the conversation can grow before Gpt4cli starts using summaries.
- `reserved-output-tokens`: The number of tokens reserved for output from the model.

### set-model default

Update org-wide default model settings for new plans.

```bash
gpt4cli set-model default # select from a list of models and settings
gpt4cli set-model default planner openai/gpt-4 # set the model for a role
# etc. — same options as `set-model` above
```

Works exactly the same as `set-model` above, but sets the default model settings for all new plans instead of only the current plan.

### models add

Add a custom model.

```bash
gpt4cli models add
```

Gpt4cli will prompt you for all required information to add a custom model.

### models delete

Delete a custom model.

```bash
gpt4cli models delete # select from a list of custom models
gpt4cli models delete some-model # by name
gpt4cli models delete 4 # by index in `gpt4cli models available --custom`
```

### model-packs

Show all available model packs.

```bash
gpt4cli model-packs
```

### model-packs create

Create a new custom model pack.

```bash
gpt4cli model-packs create
```

Gpt4cli will prompt you for all required information to create a custom model pack.

### model-packs delete

Delete a custom model pack.

```bash
gpt4cli model-packs delete
gpt4cli model-packs delete some-model-pack # by name
gpt4cli model-packs delete 4 # by index in `gpt4cli model-packs --custom`
```

## Account Management

### sign-in

Sign in, accept an invite, or create an account.

```bash
gpt4cli sign-in
```

Gpt4cli will prompt you for all required information to sign in, accept an invite, or create an account.

### invite

Invite a user to join your org.

```bash
gpt4cli invite # prompt for email, name, and role
gpt4cli invite name@domain.com 'Full Name' member # invite with email, name, and role 
```

Users can be invited as `member`, `admin`, or `owner`.

### revoke

Revoke an invite or remove a user from your org.

```bash
gpt4cli revoke # select from a list of users and invites
gpt4cli revoke name@domain.com # by email
```

### users

List users and pending invites in your org.

```bash
gpt4cli users
```

