---
sidebar_position: 5
sidebar_label: Version Control
---

# Version Control

Just about every aspect of a Gpt4cli plan is version-controlled, and anything that can happen during a plan creates a new version in the plan's history. This includes:

- Adding, removing, or updating context.
- When you send a prompt.
- When Gpt4cli responds.
- When Gpt4cli builds the plan's proposed updates to a file into a pending change.
- When pending changes are rejected.
- When pending changes are applied to your project.
- When models or model settings are updated.

## Viewing History

To see the history of your plan, use the `gpt4cli log` command:

```bash
gpt4cli log
```

## Rewinding

To rewind the plan to an earlier state, use the `gpt4cli rewind` command:

```bash
gpt4cli rewind # Rewind 1 step
gpt4cli rewind 3  # Rewind 3 steps
gpt4cli rewind a7c8d66  # Rewind to a specific step
```

## Preventing History Loss With Branches

Note that currently, there's no way to undo a `rewind` and recover any history that may have been cleared as a result. That said, you can use `rewind` without losing any history with [branches](./branches.md). Use `gpt4cli checkout` to a create a new branch before executing `rewind`, and the original branch will still include the history from before the `rewind`.

```bash
gpt4cli checkout undo-changes # create a new branch called 'undo-changes'
gpt4cli rewind ef883a # history is rewound in 'undo-changes' branch
gpt4cli checkout main # main branch still retains original history 
```

## Viewing Conversation

While the Gpt4cli history includes an entry for each message in the conversation, message content isn't included. To see the full conversation history, use the `gpt4cli convo` command:

```bash
gpt4cli convo
```

## Rewinding After `gpt4cli apply`

Like any other action that modifies a plan, running `gpt4cli apply` to apply pending changes to your project file creates a new version in the plan's history.

The `gpt4cli apply` action can be undone with `gpt4cli rewind`, but it's important to note that this will only make the changes pending again in the Gpt4cli sandbox. It **will not** undo the changes to your project files. You'll have to do that separately if desired. 
