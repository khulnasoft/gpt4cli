---
sidebar_position: 4
sidebar_label: Pending Changes
---

# Pending Changes

When you give Gpt4cli a task, the changes aren't applied directly to your project files. Instead, they are accumulated in Gpt4cli's version-controlled **sandbox** so that you can review them first.

## `gpt4cli diffs` / Changes TUI

When Gpt4cli has finished with your task, you can review the proposed changes with the `gpt4cli diff` command, which shows them in `git diff` format:

```bash
gpt4cli diff
```

`--plain/-p`: Outputs the conversation in plain text with no ANSI codes.

 You can also view them in Gpt4cli's changes TUI:

```bash
gpt4cli changes
```

## Rejecting Files

While we're working hard to make file updates as reliable as possible, bad updates can still happen. If the plan's changes were applied incorrectly to a file, you can either [apply the changes](#apply-the-changes) and then fix the problems manually, *or* you can reject the updates to that file and then make the proposed changes yourself manually. 

To reject changes to a file (or multiple files), you can run `gpt4cli reject` with one ore more file paths:

```bash
gpt4cli reject 
```

You can also reject changes using the `r` hotkey in the `gpt4cli changes` TUI.

Once the bad update is rejected, copy the changes from the plan's output or run `gpt4cli convo` to output the full conversation and copy them from there. Then apply the updates to that file yourself.

## Apply The Changes

Once you're happy with the plan's changes, you can apply them to your project files with `gpt4cli apply`:

```bash
gpt4cli apply
```

If you're in a git repository, Gpt4cli will give you the option of grouping the changes into a git commit with an automatically generated commit message. Any uncommitted changes that were present in your working directory beforehand will be unaffected.

You can skip the `gpt4cli apply` confirmation with the `-y` flag.
