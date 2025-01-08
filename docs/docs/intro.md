---
title: Intro
description: Gpt4cli is an open source, terminal-based AI coding engine that helps you work on complex, real-world development tasks with LLMs.
sidebar_position: 1
sidebar_label: Intro
---

# Introduction

Gpt4cli is an open source, terminal-based AI coding engine that helps you work on complex, real-world development tasks with LLMs.

It combines multiple agents to complete tasks that span many files and model responses. When you give Gpt4cli a task, it continues working automatically until the task is determined to be complete.   

## Use cases

- Build new projects from scratch.
- Add features to existing projects.
- Write tests.
- Write scripts.
- Fix bugs.
- Refactor.
- Work with unfamiliar technologies.
- Ask questions about code.
- Understand a codebase.
- Auto-debug a failing shell command (tests, type checks, scripts, etc.).

## What makes Gpt4cli different?

### Version Control

Gpt4cli gives the LLM its own version-controlled staging area/sandbox (separate from your project's git repo) where all of its proposed changes are accumulated. This allows you to:

- Iterate on your code and the LLM's plan at the same time without the changes becoming intertwined and difficult to disentangle.
- Review proposed changes across multiple files as a whole (rejecting any that aren't correct) to be sure that broken updates or hallucinations don't make it into your project files.
- Branch or rewind the LLM's plan in order to explore multiple paths or revert to the exact step where a task went off the rails.

### Context Management

Apart from version control, Gpt4cli also helps you manage what's in the LLM's context:

- Add files or directories to context from the terminal instead of copying and pasting or clicking around in a UI. 
- Files you add to context are kept up-to-date automatically so that the LLM is always using the latest version.
- Unlike IDE-based tools that automatically and opaquely load context in the background, Gpt4cli gives the developer precise control of what's in the LLM's context. You never have to wonder what's been loaded or whether it's up-to-date. This is crucial to getting good results and keeping a handle on costs, particularly when you want to go beyond auto-complete and work on larger tasks.

### Automated Debugging

Gpt4cli can repeatedly run a failing terminal command, continually making fixes based on the command's output and retrying until it succeeds.

## Models

By default, Gpt4cli uses a combination of Anthropic and OpenAI models; Anthropic models are served via OpenRouter.ai and OpenAI models are served directly from OpenAI.

That said, you can also use it with a wide range of other models from any OpenAI-compatible provider.

## Languages and Platforms

You can use Gpt4cli to work with any language or framework that the underlying LLM has been trained on. For the largest models, this includes just about any language or framework you can think of, though output quality will tend to be best for those that are more popular and therefore better represented in the training data.

Gpt4cli is cross-platform and easy to install. It supports Mac, Linux, FreeBSD, and Windows. It runs from a single binary with no dependencies.

## Hosting

Gpt4cli runs on a client-server architecture. The Gpt4cli server is open source and can be self-hosted. A cloud-hosted option is also offered for getting started as quickly as possible with minimal setup.

## Community

Join our growing community and help guide Gpt4cli's development.

- [GitHub](https://github.com/khulnasoft/gpt4cli) - post an [issue](https://github.com/khulnasoft/gpt4cli/issues), start a [discussion](https://github.com/khulnasoft/gpt4cli/discussions), or [fork and contribute.](https://github.com/khulnasoft/gpt4cli/fork)
- [Discord](https://discord.gg/khulnasoft) - ask for help, give feedback, share your use cases, or just hang out.
- [X](https://x.com/Gpt4cliAI) - follow for updates on new versions and other AI coding content.
- [YouTube](https://www.youtube.com/@gpt4cli-ny5ry) - subscribe to watch various tasks and projects get completed with Gpt4cli.
