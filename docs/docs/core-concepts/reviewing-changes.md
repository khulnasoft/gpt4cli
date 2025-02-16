---
sidebar_position: 4
sidebar_label: Pending Changes
---

# Pending Changes

When you give Gpt4cli a task, by default the changes aren't applied directly to your project files. Instead, they are accumulated in Gpt4cli's version-controlled **sandbox** so that you can review them first.

## `gpt4cli diffs` / `gpt4cli diffs --ui`

When Gpt4cli has finished with your task, you can review the proposed changes with the `gpt4cli diff` command, which shows them in `git diff` format:

```bash
gpt4cli diff
```

`--plain/-p`: Outputs the conversation in plain text with no ANSI codes.

You can also view them in a local browser UI with the `gpt4cli diffs --ui` command:

```bash
gpt4cli diffs --ui
```

If you pass the `--side-by-side/-s` flag alongside `--ui`, the diffs will be shown in a side-by-side view rather than line-by-line.

## Rejecting Files

While we're working hard to make file updates as reliable as possible, bad updates can still happen. If the plan's changes were applied incorrectly to a file, you can either [apply the changes](#apply-the-changes) and then fix the problems manually, *or* you can reject the updates to that file and then make the proposed changes yourself manually. 

To reject changes to a file (or multiple files), you can run `gpt4cli reject` with one or more file paths:

```bash
gpt4cli reject some-file.py
```

You can reject *all* currently pending files by passing no arguments to the reject command (you'll then be prompted to confirm the rejection):

```bash
gpt4cli reject
```

Once the bad update is rejected, copy the changes from the plan's output or run `gpt4cli convo` to output the full conversation and copy them from there. Then apply the updates to that file yourself.

## Apply The Changes

Once you're happy with the plan's changes, you can apply them to your project files with `gpt4cli apply`:

```bash
gpt4cli apply
```

If you're in a git repository, Gpt4cli will give you the option of grouping the changes into a git commit with an automatically generated commit message. Any uncommitted changes that were present in your working directory beforehand will be unaffected.

You can skip the `gpt4cli apply` confirmation with the `-y` flag.

You can pass the `--commit/-c` flag to commit the changes to git after applying them without confirmation.

You can pass the `--no-commit/-n` flag to prevent the changes from being committed to git after applying them without confirmation.

## Auto-Applying Changes

If you want to skip the review step and automatically apply the changes from a plan immediately after it's complete, you can pass the `--apply/-a` flag to `gpt4cli tell`, `gpt4cli continue`, or `gpt4cli build`.

If you do this, you can also pass the `--commit/-c` flag to commit the automatically applied changes to git.