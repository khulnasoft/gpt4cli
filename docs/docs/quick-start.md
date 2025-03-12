---
sidebar_position: 2
sidebar_label: Quickstart
---

# Quickstart

## Install Gpt4cli

```bash
curl -sL https://v2.khulnasoft.com/install.sh | bash
```

[Click here for more installation options.](./install.md)

Note that Windows is supported via [WSL](https://learn.microsoft.com/en-us/windows/wsl/about). Gpt4cli only works correctly on Windows in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell.

## Hosting Options

| Option                                | Description                                                                                                                                                                                                                                                 |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Gpt4cli Cloud (Integrated Models)** | • No separate accounts or API keys.<br/>• Easy multi-device usage.<br/>• Centralized billing and budgeting.<br/>• Quickest way to [get started.](https://app.khulnasoft.com/start?modelsMode=integrated)                                                        |
| **Gpt4cli Cloud (BYO API Key)**       | • Use Gpt4cli Cloud with your own [OpenRouter.ai](https://openrouter.ai) and [OpenAI](https://platform.openai.com) keys.<br/>                                                                                                                               |
| **Self-hosted/Local Mode**            | • Run Gpt4cli locally with Docker or host on your own server.<br/>• Use your own [OpenRouter.ai](https://openrouter.ai) and [OpenAI](https://platform.openai.com) keys.<br/>• Follow the [local-mode quickstart](./hosting/self-hosting/local-mode-quickstart.md) to get started. |

If you're going with a 'BYO API Key' option above (whether cloud or self-hosted), you'll need to set the `OPENROUTER_API_KEY` and `OPENAI_API_KEY` environment variables before continuing:

```bash
export OPENROUTER_API_KEY=...
export OPENAI_API_KEY=...
```

## Get Started

If you're starting on a new project, make a directory first:

```bash
mkdir your-project-dir
```

Now `cd` into your **project's directory.**

```bash
cd your-project-dir
```

Then just give a quick the REPL help text a quick read, and you're ready go. The REPL starts in _chat mode_ by default, which is good for fleshing out ideas before moving to implementation. Once the task is clear, Gpt4cli will prompt you to switch to _tell mode_ to make a detailed plan and start writing code.

```bash
gpt4cli
```

or for short:

```bash
g4c
```

☁️ _If you're using Gpt4cli Cloud, you'll be prompted at this point to start a trial._

Then just give the REPL help text a quick read, and you're ready go.
