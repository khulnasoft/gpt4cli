---
sidebar_position: 3
sidebar_label: Prompts
---

# Prompts

## Sending Prompts

To send a prompt to the model, use the `gpt4cli tell` command.

You can pass it in as a file with the `--file/-f` flag:

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

## Plan Stream TUI

After you send a prompt with `gpt4cli tell`, you'll see the plan stream TUI. The model's responses are streamed here. You'll see several hotkeys listed along the bottom row that allow you to stop the plan (s), send the plan to the background (b), scroll/page the streamed text, or jump to the beginning or end of the stream. If you're a vim user, you'll notice Gpt4cli's scrolling hotkeys are the same as vim's.

Note that scrolling the terminal window itself won't work while you're in the stream TUI. Use the scroll hotkeys instead. 

## Task Prompts

The most common prompt is a **task**. When you give Gpt4cli a task, it will first break down the task into steps, then it will proceed to implement each step in code. Gpt4cli will automatically continue sending model requests until the task is determined to be complete.

## Conversational Prompts

You can also ask Gpt4cli questions or chat with it about whatever you want:

```bash
gpt4cli tell "explain every function in lib/math.ts"
```

For conversational prompts, Gpt4cli will reply with just a single response and won't automatically continue.

Note that Gpt4cli can sometimes be a bit over-eager to interpret conversational prompts as tasks. Its built-in prompts generally give it a bias toward action. For example, responding to a prompt like this:

```bash
gpt4cli tell "did you forget to add the submit button?"
```

Gpt4cli is likely to respond by implementing the missing button rather than simply answering the question. Giving an additional nudge to clarify that you're only chatting and aren't implying that you want a task to be done can help.

```bash
gpt4cli tell "did you forget to add the submit button? don't add it--just answer yes or no."
```

## Stopping and Continuing

You can prevent Gpt4cli from automatically continuing for multiple responses by passing the `--stop/-s` flag:

```bash
gpt4cli tell -s "write tests for the charting helpers in lib/chart-helpers.ts"
```

Gpt4cli will then reply with just a single response. From there, you can continue if desired with the `continue` command. Like `tell`, `continue` can also accept a `--stop/-s` flag. Without the `--stop/-s` flag, `continue` will also cause Gpt4cli to continue automatically until the task is done. If you pass the `--stop/-s` flag, it will continue for just one more response.

```bash
gpt4cli continue -s
```

Apart from `--stop/-s` Gpt4cli's plan stream TUI also has an `s` hotkey that allows you to immediately stop a plan.

## Background Tasks

By default, `gpt4cli tell` opens the plan stream TUI and streams Gpt4cli's response(s) there, but you can also pass the `--bg` flag to run a task in the background instead.

You can learn more about using and interacting with background tasks [here](./background-tasks.md).

## Building Files

As Gpt4cli implements your task, files it creates or updates will appear in the `Building Plan` section of the plan stream TUI. Gpt4cli will **build** all changes proposed by the plan into a set of pending changesets for each affected file.

Keep in mind that initially, these changes **will not** be directly applied to your project files. Instead, they will be **pending** in Gpt4cli's version-controlled sandbox. This allows you to review the proposed changes or continue iterating and accumulating more changes. You can view the pending changes with `gpt4cli diff` (for git diff format) or `gpt4cli changes` (to view them in a TUI). Once you're happy with the changes, you can apply them to your project files with `gpt4cli apply`.

- [Learn more about reviewing changes.](./reviewing-changes.md)
- [Learn more about version control.](./version-control.md)

### Skipping builds / `gpt4cli build`

You can skip building files when you send a prompt by passing the `--no-build` flag to `gpt4cli tell` or `gpt4cli continue`. This can be useful if you want to ensure that a plan is on the right track before comitting the time and tokens to build files.

```bash
gpt4cli tell "implement sign up and sign in forms in src/components" --no-build
```

You can later build any changes that were implemented in the plan with the `gpt4cli build` command:

```bash
gpt4cli build
```

This will show a smaller version of the plan stream TUI that only includes the `Building Plan` section.

Like full plan streams, build streams can stopped with the `s` hotkey or sent to the background with the `b` hotkey. They can also be run fully in the background with the `--bg` flag:

```bash
gpt4cli build --bg
```

There's one more thing to keep in mind about builds. If you send a prompt with the `--no-build` flag:

```bash
gpt4cli tell "implement a forgot password email in src/emails" --no-build
```

Then you later send *another* prompt with `gpt4cli tell` or continue the plan with `gpt4cli continue` and you *don't* include the `--no-build` flag, any changes that were implemented previously but weren't built will immediately start building when the plan stream begins.

```bash
gpt4cli tell "now implement the UI portion of the forgot password flow" 
# the above will start building the changes proposed in the earlier prompt that was passed --no-build
```

## Iterating on a Plan

If you send a prompt:

```bash
gpt4cli tell "implement a fully working and production-ready tic tac toe game, including a computer-controlled AI, in html, css, and javascript"
```

And then you want to iterate on it, whether that's to add more functionality or correct something that went off track, you have a couple options.

### Continue the convo

The most straightforward way to continue iterating is to simply send another `gpt4cli tell` command:

```bash
gpt4cli tell "I plan to seek VC funding for this game, so please implement a dark mode toggle and give all buttons subtle gradient fills"
```

This is generally a good approach when you're happy with the current plan and want to extend it to add more functionality.

Note, you can view the full conversation history with the `gpt4cli convo` command:

```bash
gpt4cli convo
```

### Rewind and iterate

Another option is to use Gpt4cli's [version control](./version-control.md) features to rewind to the point just before your prompt was sent and then update it before sending the prompt again. 

You can use `gpt4cli log` to see the plan's history and determine which step to rewind to, then `gpt4cli rewind` with the appropriate hash to rewind to that step:

```bash
gpt4cli log # see the history
gpt4cli rewind accfe9 # rewind to right before your prompt
```

This approach works well in conjunction with **prompt files**. You write your prompts in files somewhere in your codebase, then pass those to `gpt4cli tell` using the `--file/-f` flag:

```bash
gpt4cli tell -f prompts/tic-tac-toe.txt
```

This makes it easy to continuously iterate on your prompt using `gpt4cli rewind` and `gpt4cli tell` until you get a result that you're happy with.

### Which is better?

There's not necessarily one right answer on whether to use an ongoing conversation or the `rewind` approach with prompt files for iteration. Here are a few things to consider when making the choice:

- Bad results tend to beget more bad results. Rewinding and iterating on the prompt is often more effective for correcting a wayward task than continuing to send more `tell` commands. Even if you are specifically prompting the model to *correct* a problem, having the wrong approach in its context will tend to bias it toward additional errors. Using `rewind` to the give the model a clean slate can work better in these scenarios.

- Iterating on a prompt file with the `rewind` approach until you find your way to an effective prompt has another benefit: you can keep the final version of the prompt that produced a given set of changes right alongside the changes themselves in your codebase. This can be helpful for other developers (or your future self) if you want to revisit a task later.

- A downside of the `rewind` approach is that it can involve re-running early steps of a plan over and over, which can be **a lot** more expensive than iterating with additional `tell` commands.