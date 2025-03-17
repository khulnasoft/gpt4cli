---
sidebar_position: 11
sidebar_label: Background Tasks
---

# Background Tasks

Gpt4cli allows you to run tasks in the background, helping you work on multiple tasks in parallel.

## Running a Task in the Background

To run a task in the background, use the `--bg` flag with `gpt4cli tell` or `gpt4cli continue`.

```bash
gpt4cli tell --bg "Add an update credit card form to 'src/components/billing'"
gpt4cli continue --bg
```

The gpt4cli stream TUI also has a `b` hotkey that allows you to send a streaming plan to the background.

## Listing Background Tasks

To list active and recently finished background tasks, use the `gpt4cli ps` command:

```bash
gpt4cli ps
```

## Connecting to a Background Task

To connect to a running background task and view its stream in the plan stream TUI, use the `gpt4cli connect` command:

```bash
gpt4cli connect
```

## Stopping a Background Task

To stop a running background task, use the `gpt4cli stop` command:

```bash
gpt4cli stop
```
