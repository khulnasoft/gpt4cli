---
sidebar_position: 3
sidebar_label: Context
---

# Context Management

Context in Gpt4cli refers to files, directories, URLs, images, notes, or piped in data that the LLM uses to understand and work on your project. Context is always associated with a [plan](./plans.md).

Changes to context are [version controlled](./version-control.md) and can be [branched](./branches.md).

## Automatic vs. Manual

As of v2, Gpt4cli loads context automatically by default. When a new plan is created, a [project map](#loading-project-maps) is generated and loaded into context. The LLM then uses this map to select relevant context before planning a task or responding to a message.

### Tradeoffs

Automatic context loading makes Gpt4cli more powerful and easier to work with, but there are tradeoffs in terms of cost, focus, and output speed. If you're trying to minimize costs or you know for sure that only one or two files are relevant to a task, you might prefer to load context manually.

### Setting Manual Mode

You can use manual context loading by:

- Using `set-auto` to choose a lower [autonomy level](./autonomy.md) that has auto-load-context disabled (like `plus` or `basic`).

```bash
gpt4cli set-auto plus
gpt4cli set-auto basic
gpt4cli set-auto default plus # set the default value for all new plans
```

- Starting a new REPL or a new plan with the `--plus` or `--basic` flags, which will automatically set the config option to the chosen autonomy level.

```bash
gpt4cli --plus
gpt4cli new --basic
```

- Setting the `auto-load-context` [config option](./configuration.md) to `false`:

```bash
gpt4cli set-config auto-load-context false
gpt4cli set-config default auto-load-context false # set the default value for all new plans
```

### Smart Context Window Management

Another new context management feature in v2 is smart context window management. When making a plan with multiple steps, Gpt4cli will determine which files are relevant to each step. Only those files will be loaded into context during implementation.

When combined with automatic context loading, this effectively creates a sliding context window that grows and shrinks as needed throughout the plan.

Smart context can also be used when you're managing context manually. To give an example: say you've manually loaded a directory with 10 files in it, and you need to make some updates to each one of them. Without smart context, each step of the implementation will load all 10 files into context. But if you use smart context, only the one or two files that are edited in each step will be loaded.

Smart context is enabled in the `plus` autonomy level and above. You can also toggle it with `set-config`:

```bash
gpt4cli set-config smart-context true
gpt4cli set-config smart-context false
gpt4cli set-config default smart-context false # set the default value for all new plans
```

### Automatic Context Updates

When you make your own changes to files in context separately from Gpt4cli, those files need to be updated before the plan can continue. Previously, Gpt4cli would prompt you to update context every time a file was changed. This is now automatic by default.

Automatic updates are enabled in the `plus` autonomy level and above. You can also toggle them with `set-config`:

```bash
gpt4cli set-config auto-update-context true
gpt4cli set-config auto-update-context false
gpt4cli set-config default auto-update-context false # set the default value for all new plans
```

### Autonomy Matrix

Here are the different autonomy levels as they relate to context management config options:

|                       | `none` | `basic` | `plus` | `semi` | `full` |
| --------------------- | ------ | ------- | ------ | ------ | ------ |
| `auto-load-context`   | ❌     | ❌      | ❌     | ✅     | ✅     |
| `smart-context`       | ❌     | ❌      | ✅     | ✅     | ✅     |
| `auto-update-context` | ❌     | ❌      | ✅     | ✅     | ✅     |

### Mixing Automatic and Manual Context

You can manually load additional context even if automatic loading is enabled. The way this additional context is handled works somewhat differently.

First, consider how automatic context loading works across each stage of a plan:

#### Automatic context loading (no manual context added)

1. **Context loading:** Only the project map is initially loaded. The map, along with your prompt, is used to select relevant context.
2. **Planning:** Only context selected in step 1 is loaded.
3. **Implementation:** Smart context (if enabled) filters context again, loading only what's directly relevant to each step.

Here's how it changes when you load manual context on top:

#### Automatic loading + manual context

1. **Context loading:** Your manually loaded context is **always included** alongside the project map.
2. **Planning:** Manually loaded context is always loaded, whether or not it's selected by the map-based context selection step.
3. **Implementation:** Smart context (if enabled) filters all context again (both manual and automatic), loading only what's directly relevant to each implementation step.

Loading files manually when using automatic context loading can sometimes be useful when you **know** certain files are relevant and don't want to risk the LLM leaving them out, or when the LLM is struggling to select the right context. If there are files that can help the LLM select the right context, like READMEs or documentation that describes the structure of the project, those can also be good candidates for manual loading.

Another use for manual context loading is for context types that can't be loaded automatically, like images, URLs, notes, or piped data (for now Gpt4cli can only automatically load project files).

## Manually Loading Context

To load files, directories, directory layouts, urls, images, notes, or piped data into a plan's context, use the `gpt4cli load` command.

### Loading Files

You can pass `load` one or more file paths. File paths are relative to the current directory in your terminal.

```bash
gpt4cli load component.ts # single file
gpt4cli load component.ts action.ts reducer.ts # multiple files
g4c l component.ts # alias
```

You can also load multiple files using glob patterns:

```bash
gpt4cli load tests/**/*.ts # loads all .ts files in 'tests' and its subdirectories
gpt4cli load * # loads all files in the current directory
```

You can load context from parent or sibling directories if needed by using `..` in your load paths.

```bash
gpt4cli load ../file.go # loads file.go from parent directory
gpt4cli load ../sibling-dir/test.go # loads test.go from sibling directory
```

### Loading Directories

You can load an entire directory with the `--recursive/-r` flag:

```bash
gpt4cli load lib -r # loads lib, all its files and all its subdirectories
gpt4cli load * -r # loads all files in the current directory and all its subdirectories
```

### Loading Files and Directories in the REPL

In the [Gpt4cli REPL](../repl.md), you can use the shortcut `@` plus a relative file path to load a file or directory into context.

```bash
@component.ts # loads component.ts
@lib # loads lib directory, and all its files and subdirectories
```

### Loading Directory Layouts

There are tasks where it's helpful for the LLM to the know the structure of your project or sections of your project, but it doesn't necessarily need to the see the content of every file. In that case, you can pass in a directory with the `--tree` flag to load in the directory layout. It will include just the names of all included files and subdirectories (and each subdirectory's files and subdirectories, and so on).

```bash
gpt4cli load . --tree # loads the layout of the current directory and its subdirectories (file names only)
gpt4cli load src/components --tree # loads the layout of the src/components directory
```

### Loading Project Maps

Gpt4cli can create a **project map** for any directory using [tree-sitter](https://tree-sitter.github.io/tree-sitter). This shows all the top-level symbols, like variables, functions, classes, etc. in each file. 30+ languages are supported. For non-supported languages, files are still listed without symbols so that the model is aware of their existence.

Maps are mainly used for selecting context during automatic context loading, but can also be used with manual context management in order to improve output. Maps make it much more likely that an LLM will, for example, use an existing function in your project (and call it correctly) rather than generating a new one that does the same thing.

```bash
gpt4cli load . --map
```

### Loading URLs

Gpt4cli can load the text content of URLs, which can be useful for adding relevant documentation, blog posts, discussions, and the like.

```bash
gpt4cli load https://redux.js.org/usage/writing-tests # loads the text-only content of the url
```

### Loading Images

Gpt4cli can load images into context.

```bash
gpt4cli load ui-mockup.png
```

For the default GPT-4o model, png, jpeg, non-animated gif, and webp formats are supported. For other models, support for images in general, and particular formats specifically, will depend on the model.

### Loading Notes

You can add notes to context, which are just simple strings.

```bash
gpt4cli load -n 'add logging statements to all the code you generate.' # load a note into context
```

Notes can be useful as 'sticky' explanations or instructions that will tend to have more prominence throughout a long conversation than normal prompts. That's because long conversations are summarized to stay below a token limit, which can cause some details from your prompts to be dropped along the way. This doesn't happen if you use notes.

### Piping Into Context

You can pipe the results of other commands into context:

```bash
npm test | gpt4cli load # loads the output of `npm test`
```

### Ignoring files

If you're in a git repo, Gpt4cli respects `.gitignore` and won't load any files that you're ignoring. You can also add a `.gpt4cliignore` file with ignore patterns to any directory.

You can force Gpt4cli to load ignored files with the `--force/-f` flag:

```bash
gpt4cli load .env --force # loads the .env file even if it's in .gitignore or .gpt4cliignore
```

## Viewing Context

To list everything in context, use the `gpt4cli ls` command:

```bash
gpt4cli ls
```

You can also see the content of any context item with the `gpt4cli show` command:

```bash
gpt4cli show component.ts # show the content of component.ts
gpt4cli show 2 # show the content of the second item in the `gpt4cli ls` list
```

## Removing Context

To remove selectively remove context, use the `gpt4cli rm` command:

```bash
gpt4cli rm component.ts # remove by name
gpt4cli rm 2 # remove by number in the `gpt4cli ls` list
gpt4cli rm 2-5 # remove a range of indices
gpt4cli rm lib/**/*.js # remove by glob pattern
gpt4cli rm lib # remove whole directory
```

## Clearing Context

To clear all context, use the `gpt4cli clear` command:

```bash
gpt4cli clear
```

## Updating Context

If files, directory layouts, or URLs in context are modified outside of Gpt4cli, they will need to be updated next time you send a prompt.

Whether they'll be updated automatically or you'll be prompted to update them depends on the `auto-update-context` config option.

You can also update any outdated files with the `update` command.

```bash
gpt4cli update # update files in context
```
