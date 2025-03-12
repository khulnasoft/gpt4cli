---
sidebar_position: 8
sidebar_label: Branches
---

# Branches

Branches in Gpt4cli allow you to easily try out multiple approaches to a task and see which gives you the best results. They work in conjunction with [version control](./version-control.md). Use cases include:

- Comparing different prompting strategies.
- Comparing results with different files in context.
- Comparing results with different models or model-settings.
- Using `gpt4cli rewind` without losing history (first check out a new branch, then rewind).

## Creating a Branch

To create a new branch, use the `gpt4cli checkout` command:

```bash
gpt4cli checkout new-branch
g4cd new-branch # alias
```

## Switching Branches

To switch to a different branch, also use the `gpt4cli checkout` command:

```bash
gpt4cli checkout existing-branch
```

## Listing Branches

To list all branches, use the `gpt4cli branches` command:

```bash
gpt4cli branches
```

## Deleting a Branch

To delete a branch, use the `gpt4cli delete-branch` command:

```bash
gpt4cli delete-branch branch-name
```
