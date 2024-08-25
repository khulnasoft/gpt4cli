# gpt4cli

<div align="">

<p align="">
  <!-- Call to Action Links -->
  <a href="#install">
    <b>30-Second Install</b>
  </a>
   · 
  <a href="https://gpt4cli.khulnasoft.com">
    <b>Website</b>
  </a>
   · 
  <a href="https://gpt4cli.khulnasoft.com/">
    <b>Docs</b>
  </a>
   · 
  <a href="#more-examples-">
    <b>Examples</b>
  </a>
   · 
  <a href="https://gpt4cli.khulnasoft.com/hosting/self-hosting">
    <b>Self-Hosting</b>
  </a>
   <!-- · 
  <a href="./guides/DEVELOPMENT.md">
    <b>Development</b>
  </a> -->
  <!--  · 
  <a href="https://discord.gg/khulnasoft">
    <b>Discord</b>
  </a>  
   · 
  <a href="#weekly-office-hours-">
    <b>Office Hours</b>
  </a>  
  -->
</p>

<br>

[![Discord](https://img.shields.io/discord/1144758149232988282.svg?style=flat&logo=discord&label=Discord&refresh=1)](https://discord.gg/khulnasoft)
[![GitHub Repo stars](https://img.shields.io/github/stars/khulnasoft/gpt4cli?style=social)](https://github.com/khulnasoft/gpt4cli)
[![Twitter Follow](https://img.shields.io/twitter/follow/khulnasoft?style=social)](https://twitter.com/khulnasoft)

</div>

<p align="">
  <!-- Badges -->
<a href="https://github.com/khulnasoft/gpt4cli/pulls"><img src="https://img.shields.io/badge/PRs-welcome-brightgreen.svg" alt="PRs Welcome" /></a> <a href="https://github.com/khulnasoft/gpt4cli/releases?q=cli"><img src="https://img.shields.io/github/v/release/khulnasoft/gpt4cli?filter=cli*" alt="Release" /></a>
<a href="https://github.com/khulnasoft/gpt4cli/releases?q=server"><img src="https://img.shields.io/github/v/release/khulnasoft/gpt4cli?filter=server*" alt="Release" /></a>

  <!-- <a href="https://github.com/your_username/your_project/issues">
    <img src="https://img.shields.io/github/issues-closed/your_username/your_project.svg" alt="Issues Closed" />
  </a> -->

</p>

<br />

<div align="">
<a href="https://trendshift.io/repositories/8994" target="_blank"><img src="https://trendshift.io/api/badge/repositories/8994" alt="khulnasoft%2Fgpt4cli | Trendshift" style="width: 250px; height: 55px;" width="250" height="55"/></a>
</div>

<br>

<h3 align="">AI driven development in your terminal.<br/>Build entire features and apps with a robust workflow.</h3>

<br/>
<br/>

<!-- Vimeo link is nicer on mobile than embedded video... downside is it navigates to vimeo in same tab (no way to add target=_blank) -->
<!-- https://github.com/khulnasoft/gpt4cli/assets/545350/c2ee3bcd-1512-493f-bdd5-e3a4ca534a36 -->

<a href="https://player.vimeo.com/video/926634577">
  <img src="images/gpt4cli-intro-vimeo.png" alt="Gpt4cli intro video" width="100%"/>
</a>

<br/>
<br/>


## Learn more about Gpt4cli  🧐

- [Overview](#overview-)
- [Install](#install)
- [Get started](#get-started-)
- [Docs](https://gpt4cli.khulnasoft.com/)
- [Build complex software](#build-complex-software-with-llms-)
- [Why Gpt4cli?](#why-gpt4cli-)
- [Roadmap](#roadmap-%EF%B8%8F)
- [Discussion and discord](#discussion-and-discord-)
- [Contributors](#contributors-)
<br/>

## Overview  📚

<p>Churn through your backlog, work with unfamiliar technologies, get unstuck, and <strong>spend less time on the boring stuff.</strong></p>

<p>Gpt4cli is a <strong>reliable and developer-friendly</strong> AI coding agent in your terminal. It can plan out and complete <strong>large tasks</strong> that span many files and steps.</p>
 
<p>Designed for <strong>real-world use-cases</strong>, Gpt4cli can help you build a new app quickly, add new features to an existing codebase, write tests and scripts, understand code, and fix bugs. </p>

<br/>

## Install  📥

```bash
curl -sL https://raw.githubusercontent.com/khulnasoft/gpt4cli/main/app/cli/install.sh | bash
```

**Note:** Windows is supported via [WSL](https://learn.microsoft.com/en-us/windows/wsl/install). Gpt4cli only works correctly on Windows in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell.

[More installation options.](https://gpt4cli.khulnasoft.com/install)

<br/>

## Get started  🚀

Gpt4cli uses OpenAI by default. If you don't have an OpenAI account, first [sign up here.](https://platform.openai.com/signup)

Then [generate an API key here](https://platform.openai.com/account/api-keys) and `export` it.

```bash
export OPENAI_API_KEY=...
```


Now `cd` into your **project's directory.** Make a new directory first with `mkdir your-project-dir` if you're starting on a new project.

```bash
cd your-project-dir
```


Then **start your first plan** with `gpt4cli new`.

```bash
gpt4cli new
```


Load any relevant files, directories, directory layouts, urls, or images **into the LLM's context** with `gpt4cli load`.

```bash
gpt4cli load some-file.ts another-file.ts
gpt4cli load src/components -r # load a whole directory
gpt4cli load src --tree # load a directory layout (file names only)
gpt4cli load src/**/*.ts # load files matching a glob pattern
gpt4cli load https://raw.githubusercontent.com/khulnasoft/gpt4cli/main/README.md # load the text content of a url
gpt4cli load images/mockup.png # load an image
```


Now **send your prompt.** You can pass it in as a file:

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

Gpt4cli will make a plan for your task and then implement that plan in code. **The changes won't yet be applied to your project files.** Instead, they'll accumulate in Gpt4cli's sandbox.

To learn about reviewing changes, iterating on the plan, and applying changes to your project, **[continue with the full quickstart.](https://gpt4cli.khulnasoft.com/quick-start#review-the-changes)**

<br/>

## Docs  🛠️

### [👉  Full documentation.](https://gpt4cli.khulnasoft.com/)


<br/>

## Build complex software with LLMs  🌟

⚡️  Changes are accumulated in a protected sandbox so that you can review them before automatically applying them to your project files. Built-in version control allows you to easily go backwards and try a different approach. Branches allow you to try multiple approaches and compare the results.

📑  Manage context efficiently in the terminal. Easily add files or entire directories to context, and keep them updated automatically as you work so that models always have the latest state of your project.

🧠  By default, Gpt4cli relies on the OpenAI API and requires an `OPENAI_API_KEY` environment variable. You can also use it with a wide range of other models, including Anthropic Claude, Google Gemini, Mixtral, Llama and many more via OpenRouter.ai, Together.ai, or any other OpenAI-compatible provider.

✅  Gpt4cli supports Mac, Linux, FreeBSD, and Windows. It runs from a single binary with no dependencies.

<br/>

## Why Gpt4cli?  🤔

🏗️  Go beyond autocomplete to build complex functionality with AI.<br>
🚫  Stop the mouse-centered, copy-pasting madness of coding with ChatGPT.<br>
⚡️  Ensure the model always has the latest versions of files in context.<br>
🪙  Retain granular control over what's in the model's context and how many tokens you're using.<br>
⏪  Rewind, iterate, and retry as needed until you get your prompt just right.<br>
🌱  Explore multiple approaches with branches.<br>
🔀  Run tasks in the background or work on multiple tasks in parallel.<br>
🎛️  Try different models and temperatures, then compare results.<br>

<br/>

## Roadmap  🗺️

🧠  Support for open source models, Google Gemini, and Anthropic Claude in addition to OpenAI  ✅ released<br>
🖼️  Support for multi-modal models—add images and screenshots to context ✅ released<br>
🤝  Plan sharing and team collaboration<br>
🖥️  VSCode and JetBrains extensions<br>
📦  Community plugins and modules<br>
🔌  Github integration<br>
🌐  Web dashboard and GUI<br>
🔐  SOC2 compliance<br>
🛩️  Fine-tuned models<br>

This list will grow and be prioritized based on your feedback.

<br/>

## Discussion and discord  💬

Speaking of feedback, feel free to give yours, ask questions, report a bug, or just hang out:

- [Discord](https://discord.gg/khulnasoft)
- [Discussions](https://github.com/khulnasoft/gpt4cli/discussions)
- [Issues](https://github.com/khulnasoft/gpt4cli/issues)

<br/>

## Contributors  👥

⭐️  Please star, fork, explore, and contribute to Gpt4cli. There's a lot of work to do and so much that can be improved.

Work on tests, evals, prompts, and bug fixes is especially appreciated.

[Here's an overview on setting up a development environment.](https://gpt4cli.khulnasoft.com/development)


