<h1 align="center">
 <a href="https://khulnasoft.com">
  <picture>
    <source media="(prefers-color-scheme: dark)" srcset="images/gpt4cli-logo-dark-v2.png"/>
    <source media="(prefers-color-scheme: light)" srcset="images/gpt4cli-logo-light-v2.png"/>
    <img width="400" src="images/gpt4cli-logo-dark-bg-v2.png"/>
 </a>
 <br />
</h1>
<br />

<div align="center">

<p align="center">
  <!-- Call to Action Links -->
  <a href="#install">
    <b>30-Second Install</b>
  </a>
   · 
  <a href="https://khulnasoft.com">
    <b>Website</b>
  </a>
   · 
  <a href="https://docs-v2.khulnasoft.com/">
    <b>Docs</b>
  </a>
   · 
  <!-- <a href="#more-examples-">
    <b>Examples</b>
  </a>
   ·  -->
  <a href="https://docs-v2.khulnasoft.com/hosting/self-hosting/local-mode-quickstart">
    <b>Local Self-Hosted Mode</b>
  </a>
</p>

<br>

[![Discord](https://img.shields.io/discord/1214825831973785600.svg?style=flat&logo=discord&label=Discord&refresh=1)](https://discord.gg/khulnasoft)
[![GitHub Repo stars](https://img.shields.io/github/stars/khulnasoft/gpt4cli?style=social)](https://github.com/khulnasoft/gpt4cli)
[![Twitter Follow](https://img.shields.io/twitter/follow/KhulnaSoft?style=social)](https://twitter.com/KhulnaSoft)

</div>

<p align="center">
  <!-- Badges -->
<a href="https://github.com/khulnasoft/gpt4cli/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome" /></a> <a href="https://github.com/khulnasoft/gpt4cli/releases?q=cli"><img src="https://img.shields.io/github/v/release/khulnasoft/gpt4cli?filter=cli*" alt="Release" /></a>
<a href="https://github.com/khulnasoft/gpt4cli/releases?q=server"><img src="https://img.shields.io/github/v/release/khulnasoft/gpt4cli?filter=server*" alt="Release" /></a>

  <!-- <a href="https://github.com/your_username/your_project/issues">
    <img src="https://img.shields.io/github/issues-closed/your_username/your_project.svg" alt="Issues Closed" />
  </a> -->

</p>

<br />

<div align="center">
<a href="https://trendshift.io/repositories/8994" target="_blank"><img src="https://trendshift.io/api/badge/repositories/8994" alt="khulnasoft%2Fgpt4cli | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
</div>

<br>

<h1 align="center" >
  An AI coding agent for large tasks and real world projects.<br/><br/>
</h1>

<!-- <h2 align="center">
  Designed for large tasks and real world projects.<br/><br/>
  </h2> -->
  <br/>

<div align="center">
  <a href="https://www.youtube.com/watch?v=SFSu2vNmlLk">
    <img src="images/gpt4cli-v2-yt.png" alt="Gpt4cli v2 Demo Video" width="800">
  </a>
</div>

<br/>

💻  Gpt4cli is a terminal-based AI development tool that can **plan and execute** large coding tasks that span many steps and touch dozens of files. It can handle up to 2M tokens of context directly (~100k per file), and can index directories with 20M tokens or more using tree-sitter project maps.

🔬  **A cumulative diff review sandbox** keeps AI-generated changes separate from your project files until they are ready to go. Command execution is controlled so you can easily roll back and debug. Gpt4cli helps you get the most out of AI without leaving behind a mess in your project.

🧠  **Combine the best models** from Anthropic, OpenAI, Google, and open source providers to build entire features and apps with a robust terminal-based workflow.

🚀  Gpt4cli is capable of <strong>full autonomy</strong>—it can load relevant files, plan and implement changes, execute commands, and automatically debug—but it's also highly flexible and configurable, giving developers fine-grained control and a step-by-step review process when needed.

💪  Gpt4cli is designed to be resilient to <strong>large projects and files</strong>. If you've found that others tools struggle once your project gets past a certain size or the changes are too complex, give Gpt4cli a shot.

## Smart context management that works in big projects

- 🐘 **2M token effective context window** with default model pack. Gpt4cli loads only what's needed for each step.
- 🗄️ **Reliable in large projects and files.** Easily generate, review, revise, and apply changes spanning dozens of files.
- 🗺️ **Fast project map generation** and syntax validation with tree-sitter. Supports 30+ languages.
- 💰 **Context caching** is used across the board for OpenAI and Anthropic models, reducing costs and latency.

## Tight control or full autonomy—it's up to you

- 🚦 **Configurable autonomy:** go from full auto mode to fine-grained control depending on the task.
- 🐞 **Automated debugging** of terminal commands (like builds, linters, tests, deployments, and scripts).

## Tools that help you get production-ready results

- 💬 **A project-aware chat mode** that helps you flesh out ideas before moving to implementation. Also great for asking questions and learning about a codebase.
- 🧠 **Easily try + combine models** from multiple providers. Curated model packs offer different tradeoffs of capability, cost, and speed, as well as open source and provider-specific packs.
- 🛡️ **Reliable file edits** that prioritize correctness. While most edits are quick and cheap, Gpt4cli validates both syntax and logic as needed, with multiple fallback layers when there are problems.
- 🔀 **Full-fledged version control** for every update to the plan, including branches for exploring multiple paths or comparing different models.
- 📂 **Git integration** with commit message generation and optional automatic commits.

## Dev-friendly, easy to install

- 🧑‍💻 **REPL mode** with fuzzy auto-complete for commands and file loading. Just run `gpt4cli` in any project to get started.
- 🛠️ **CLI interface** for scripting or piping data into context.
- 📦 **One-line, zero dependency CLI install**. Dockerized local mode for easily self-hosting the server. Cloud-hosting options for extra reliability and convenience.

<!-- <br/>
<br/> -->

<!-- Vimeo link is nicer on mobile than embedded video... downside is it navigates to vimeo in same tab (no way to add target=_blank) -->
<!-- https://github.com/khulnasoft/gpt4cli/assets/545350/c2ee3bcd-1512-493f-bdd5-e3a4ca534a36 -->

<!-- <a href="https://player.vimeo.com/video/926634577">
  <img src="images/gpt4cli-intro-vimeo.png" alt="Gpt4cli intro video" width="100%"/>
</a> -->

<!-- <br/>
<br/> -->

<!-- ## More examples  🎥

<h4>👉  <a href="https://www.youtube.com/watch?v=0ULjQx25S_Y">Building Pong in C/OpenGL with GPT-4o and Gpt4cli</a></h4>

<h4>👉  <a href="https://www.youtube.com/watch?v=rnlepfh7TN4">Fixing a tricky real-world bug in 5 minutes with Claude Opus 3 and Gpt4cli</a></h4>

<br/> -->

<!-- ## Learn more  🧐

<!-- - [Overview](#overview-) -->
<!-- - [Install](#install)
- [Hosting options](#hosting-options-)
- [Get started](#get-started-)
- [Docs](https://docs.khulnasoft.com/)
- [Build complex software](#build-complex-software-with-llms-)
- [Why Gpt4cli?](#why-gpt4cli-)
- [Discussion and discord](#discussion-and-discord-)
- [Contributors](#contributors-) -->
<br/>

<!-- ## Overview  📚

<p>Churn through your backlog, work with unfamiliar technologies, get unstuck, and <strong>spend less time on the boring stuff.</strong></p>

<p>Gpt4cli is a <strong>reliable and developer-friendly</strong> AI coding agent in your terminal. It can plan out and complete <strong>large tasks</strong> that span many files and steps.</p>

<p>Designed for <strong>real-world use-cases</strong>, Gpt4cli can help you build a new app quickly, add new features to an existing codebase, write tests and scripts, understand code, fix bugs, and automatically debug failing commands (like tests, typecheckers, linters, etc.). </p>

<br/> -->

## Workflow  🔄

<img src="images/gpt4cli-workflow.png" alt="Gpt4cli workflow" width="100%"/>

## Install  📥

```bash
curl -sL https://khulnasoft.com/install.sh | bash
```

**Note:** Windows is supported via [WSL](https://learn.microsoft.com/en-us/windows/wsl/install). Gpt4cli only works correctly on Windows in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell.

[More installation options.](https://docs-v2.khulnasoft.com/install)

## Hosting  ⚖️

| Option                                | Description                                                                                                                                                                                                                                                 |
| ------------------------------------- | ----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| **Gpt4cli Cloud (Integrated Models)** | • No separate accounts or API keys.<br/>• Easy multi-device usage.<br/>• Centralized billing and budgeting.<br/>• Quickest way to [get started.](https://app.khulnasoft.com/start?modelsMode=integrated)                                                        |
| **Gpt4cli Cloud (BYO API Key)**       | • Use Gpt4cli Cloud with your own [OpenRouter.ai](https://openrouter.ai) and [OpenAI](https://platform.openai.com) keys.<br/>• [Get started](https://app.khulnasoft.com/start?modelsMode=byo)                                                                   |
| **Self-hosted/Local Mode**            | • Run Gpt4cli locally with Docker or host on your own server.<br/>• Use your own [OpenRouter.ai](https://openrouter.ai) and [OpenAI](https://platform.openai.com) keys.<br/>• Follow the [local-mode quickstart](./hosting/self-hosting.md) to get started. |

## Provider keys  🔑

If you're going with a 'BYO API Key' option above (whether cloud or self-hosted), you'll need to set the `OPENROUTER_API_KEY` and `OPENAI_API_KEY` environment variables before continuing:

```bash
export OPENROUTER_API_KEY=...
export OPENAI_API_KEY=...
```

<br/>

## Get started  🚀

First, `cd` into a **project directory** where you want to get something done or chat about the project. Make a new directory first with `mkdir your-project-dir` if you're starting on a new project.

```bash
cd your-project-dir
```

For a new project, you might also want to initialize a git repo. Gpt4cli doesn't require that your project is in a git repo, but it does integrate well with git if you use it.

```bash
git init
```

Now start the Gpt4cli REPL in your project:

```bash
gpt4cli
```

or for short:

```bash
g4c
```

☁️ _If you're using Gpt4cli Cloud, you'll be prompted at this point to start a trial._

Then just give the REPL help text a quick read, and you're ready go. The REPL starts in _chat mode_ by default, which is good for fleshing out ideas before moving to implementation. Once the task is clear, Gpt4cli will prompt you to switch to _tell mode_ to make a detailed plan and start writing code.

<br/>

## Docs  🛠️

### [👉  Full documentation.](https://docs-v2.khulnasoft.com/)

<br/>

## Discussion and discord  💬

Please feel free to give your feedback, ask questions, report a bug, or just hang out:

- [Discord](https://discord.gg/khulnasoft)
- [Discussions](https://github.com/khulnasoft/gpt4cli/discussions)
- [Issues](https://github.com/khulnasoft/gpt4cli/issues)

## Follow and subscribe

- [Follow @KhulnaSoft](https://x.com/KhulnaSoft)
- [Follow @Danenania](https://x.com/Danenania) (Gpt4cli's creator)
- [Subscribe on YouTube](https://x.com/KhulnaSoft)

<br/>

## Contributors  👥

⭐️  Please star, fork, explore, and contribute to Gpt4cli. There's a lot of work to do and so much that can be improved.

[Here's an overview on setting up a development environment.](https://docs-v2.khulnasoft.com/development)
