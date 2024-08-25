---
sidebar_position: 3
sidebar_label: Quickstart
---

# Quickstart

## Install Gpt4cli

```bash
curl -sL https://gpt4cli.khulnasoft.com/install.sh | bash
```

[Click here for more installation options.](./install.md)

Note that Windows is supported via [WSL](https://learn.microsoft.com/en-us/windows/wsl/about). Gpt4cli only works correctly on Windows in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell.


## Set `OPENAI_API_KEY`

Gpt4cli uses OpenAI by default. If you don't have an OpenAI account, first [sign up here.](https://platform.openai.com/signup)

Then [generate an API key here](https://platform.openai.com/account/api-keys) and `export` it.

```bash
export OPENAI_API_KEY=...
```

## Create A Plan

If you're starting on a new project, make a directory first:

```bash
mkdir your-project-dir
```

Now `cd` into your **project's directory.** 

```bash
cd your-project-dir
```

For a new project, you might also want to initialize a git repo. Gpt4cli doesn't require that your project is in a git repo, but it does integrate well with git if you use it.

```bash
git init
```

Now **create your first plan** with `gpt4cli new`.

```bash
gpt4cli new
```

## Load In Context

Load any relevant files, directories, directory layouts, urls, or images **into the LLM's context** with `gpt4cli load`. You can also pipe in the results of a command.

```bash
gpt4cli load some-file.ts another-file.ts
gpt4cli load src/components -r # load a whole directory
gpt4cli load src --tree # load a directory layout (file names only)
gpt4cli load src/**/*.ts # load files matching a glob pattern
gpt4cli load https://raw.githubusercontent.com/khulnasoft/gpt4cli/main/README.md # load the text content of a url
gpt4cli load images/mockup.png # load an image
npm test | gpt4cli load # pipe in the output of a command
```

## Send A Prompt

Now **send your prompt.** You can pass it in as a file:

```bash
gpt4cli tell -f prompt.txt
```

Write it in vim:

```bash
gpt4cli tell # tell with no arguments opens vim so you can write your prompt there
```

Or pass it inline (use enter for line breaks):

```bash
gpt4cli tell "add a new line chart showing the number of foobars over time to components/charts.tsx"
```

Gpt4cli will make a plan for your task and then implement that plan in code. **The changes won't yet be applied to your project files.** Instead, they'll accumulate in Gpt4cli's sandbox. 

## Review The Changes

When Gpt4cli has finished with your task, you can review the proposed changes with the `gpt4cli diff` command, which shows them in `git diff` format:

```bash
gpt4cli diff
```

Or you can view them in Gpt4cli's changes TUI:

```bash
gpt4cli changes
```

## Iterate If Needed

If the proposed changes have issues or need more work, you have a few options:

### 1. Continue prompting.

You can send another prompt to continue updating or refining the plan.

```bash
gpt4cli tell "the line chart should be centered and have a width and height of 80% of the screen"
```

### 2. Rewind the plan.

You can use `gpt4cli rewind` to revert to an earlier step in the plan, load in new context or update the prompt as needed, then proceed from there with another `gpt4cli tell` or a `gpt4cli continue` (which continues from where the plan left off).

Use `gpt4cli log` for a list of all changes in a plan. You can rewind one step by running `gpt4cli rewind` with no arguments, go back a specific number of steps (`gpt4cli rewind 3`), or rewind to a specific change with a hash `gpt4cli rewind e7e06e0`.

Seeing the conversation history can also be helpful when rewinding, since `gpt4cli log` doesn't include conversation messages in its output. You can do that with `gpt4cli convo`.

### 3. Reject incorrect files.

While we're working hard to make file updates as reliable as possible, bad updates can still happen. If the plan's changes were applied incorrectly to a file, you can either [apply the changes](#apply-the-changes) and then fix the problems manually, *or* you can reject the updates to that file and then make the proposed changes yourself manually. 

To reject changes to a file (or multiple files), you can run `gpt4cli reject` with the file path(s):

```bash
gpt4cli reject components/charts.tsx
```

You can also reject changes using the `r` hotkey in the `gpt4cli changes` TUI.

Once the bad update is rejected, copy the changes from the plan's output or run `gpt4cli convo` to output the full conversation and copy them from there. Then apply the updates to that file yourself.

## Apply The Changes

Once you're happy (enough) with the plan's changes, you can apply them to your project files with `gpt4cli apply`:

```bash
gpt4cli apply
```

If you're in a git repository, Gpt4cli will give you the option of grouping the changes into a git commit with an automatically generated commit message.

You've now experienced the core workflow of Gpt4cli! While there are more commands and options available, those described above are what you'll be using most often. 

## CLI Help

After any gpt4cli command is run, commands that could make sense to run next will be suggested. You can learn to use Gpt4cli quickly by jumping in and following these suggestions.

You can get help on the CLI with `gpt4cli help` and a list of all commands with `gpt4cli help --all`. Get help on a specific command and its options with `gpt4cli [command] --help`.

## Aliases

You can use the `pdx` alias instead of `gpt4cli` to type a bit less, and most common commands have their own aliases as well.

Here are the same commands we went through above using aliases to minimize typing:

```bash
pdx new
pdx l some-file.ts another-file.ts # load
pdx t -f prompt.txt # tell
pdx diff
pdx ch # changes
pdx log
pdx rw e7e06e0 # rewind
pdx c # continue
pdx rj components/charts.tsx # reject
pdx ap # apply
```