"use strict";(self.webpackChunkdocs=self.webpackChunkdocs||[]).push([[567],{5226:e=>{e.exports=JSON.parse('{"version":{"pluginId":"default","version":"current","label":"Next","banner":null,"badge":false,"noIndex":false,"className":"docs-version-current","isLast":true,"docsSidebars":{"tutorialSidebar":[{"type":"link","label":"Intro","href":"/intro","docId":"intro","unlisted":false},{"type":"link","label":"Install","href":"/install","docId":"install","unlisted":false},{"type":"link","label":"Quickstart","href":"/quick-start","docId":"quick-start","unlisted":false},{"type":"category","label":"Core Concepts","collapsible":true,"collapsed":false,"items":[{"type":"link","label":"Plans","href":"/core-concepts/plans","docId":"core-concepts/plans","unlisted":false},{"type":"link","label":"Context","href":"/core-concepts/context-management","docId":"core-concepts/context-management","unlisted":false},{"type":"link","label":"Prompts","href":"/core-concepts/prompts","docId":"core-concepts/prompts","unlisted":false},{"type":"link","label":"Pending Changes","href":"/core-concepts/reviewing-changes","docId":"core-concepts/reviewing-changes","unlisted":false},{"type":"link","label":"Version Control","href":"/core-concepts/version-control","docId":"core-concepts/version-control","unlisted":false},{"type":"link","label":"Branches","href":"/core-concepts/branches","docId":"core-concepts/branches","unlisted":false},{"type":"link","label":"Conversations","href":"/core-concepts/conversations","docId":"core-concepts/conversations","unlisted":false},{"type":"link","label":"Background Tasks","href":"/core-concepts/background-tasks","docId":"core-concepts/background-tasks","unlisted":false},{"type":"link","label":"Collaboration / Orgs","href":"/core-concepts/orgs","docId":"core-concepts/orgs","unlisted":false}]},{"type":"category","label":"Models","collapsible":true,"collapsed":false,"items":[{"type":"link","label":"Providers","href":"/models/model-providers","docId":"models/model-providers","unlisted":false},{"type":"link","label":"Roles","href":"/models/roles","docId":"models/roles","unlisted":false},{"type":"link","label":"Settings","href":"/models/model-settings","docId":"models/model-settings","unlisted":false}]},{"type":"category","label":"Hosting","collapsible":true,"collapsed":false,"items":[{"type":"link","label":"Cloud","href":"/hosting/cloud","docId":"hosting/cloud","unlisted":false},{"type":"link","label":"Self-Hosting","href":"/hosting/self-hosting","docId":"hosting/self-hosting","unlisted":false}]},{"type":"link","label":"CLI Reference","href":"/cli-reference","docId":"cli-reference","unlisted":false},{"type":"link","label":"Security","href":"/security","docId":"security","unlisted":false},{"type":"link","label":"Development","href":"/development","docId":"development","unlisted":false},{"type":"link","label":"Environment Variables","href":"/environment-variables","docId":"environment-variables","unlisted":false}]},"docs":{"cli-reference":{"id":"cli-reference","title":"CLI Reference","description":"All Gpt4cli CLI commands and their options.","sidebar":"tutorialSidebar"},"core-concepts/background-tasks":{"id":"core-concepts/background-tasks","title":"Background Tasks","description":"Gpt4cli allows you to run tasks in the background, helping you work on multiple tasks in parallel.","sidebar":"tutorialSidebar"},"core-concepts/branches":{"id":"core-concepts/branches","title":"Branches","description":"Branches in Gpt4cli allow you to easily try out multiple approaches to a task and see which gives you the best results. They work in conjunction with version control. Use cases include:","sidebar":"tutorialSidebar"},"core-concepts/context-management":{"id":"core-concepts/context-management","title":"Context Management","description":"Context in Gpt4cli refers to files, directories, URLs, images, notes, or piped in data that the LLM uses to understand and work on your project. Context is always associated with a plan","sidebar":"tutorialSidebar"},"core-concepts/conversations":{"id":"core-concepts/conversations","title":"Conversations","description":"Each time you send a prompt to Gpt4cli or Gpt4cli responds, the plan\'s conversation is updated. Conversations are version controlled and can be branched.","sidebar":"tutorialSidebar"},"core-concepts/orgs":{"id":"core-concepts/orgs","title":"Collaboration and Orgs","description":"While so far Gpt4cli is mainly focused on a single-user experience, we plan to add features for sharing, collaboration, and team management in the future, and some groundwork has already been done. Orgs are the basis for collaboration in Gpt4cli.","sidebar":"tutorialSidebar"},"core-concepts/plans":{"id":"core-concepts/plans","title":"Plans","description":"A plan in Gpt4cli is similar to a conversation in ChatGPT. It might only include a single prompt and model response that executes one small task, or it could represent a long back and forth with the model that generates dozens of files and builds a whole feature or an entire app.","sidebar":"tutorialSidebar"},"core-concepts/prompts":{"id":"core-concepts/prompts","title":"Prompts","description":"Sending Prompts","sidebar":"tutorialSidebar"},"core-concepts/reviewing-changes":{"id":"core-concepts/reviewing-changes","title":"Pending Changes","description":"When you give Gpt4cli a task, the changes aren\'t applied directly to your project files. Instead, they are accumulated in Gpt4cli\'s version-controlled sandbox so that you can review them first.","sidebar":"tutorialSidebar"},"core-concepts/version-control":{"id":"core-concepts/version-control","title":"Version Control","description":"Just about every aspect of a Gpt4cli plan is version-controlled, and anything that can happen during a plan creates a new version in the plan\'s history. This includes:","sidebar":"tutorialSidebar"},"development":{"id":"development","title":"Development","description":"To set up a development environment, first install dependencies:","sidebar":"tutorialSidebar"},"environment-variables":{"id":"environment-variables","title":"Environment Variables","description":"This is an overview of all the environment variables that can be used with Gpt4cli.","sidebar":"tutorialSidebar"},"hosting/cloud":{"id":"hosting/cloud","title":"Gpt4cli Cloud","description":"Overview","sidebar":"tutorialSidebar"},"hosting/self-hosting":{"id":"hosting/self-hosting","title":"Self-Hosting","description":"Gpt4cli is open source and uses a client-server architecture. The server can be self-hosted. You can run either run it locally or on a cloud server that you control.","sidebar":"tutorialSidebar"},"install":{"id":"install","title":"Install Gpt4cli","description":"Quick Install","sidebar":"tutorialSidebar"},"intro":{"id":"intro","title":"Intro","description":"Gpt4cli is an open source, terminal-based AI coding engine that helps you work on complex, real-world development tasks with LLMs.","sidebar":"tutorialSidebar"},"models/model-providers":{"id":"models/model-providers","title":"Model Providers","description":"By default, Gpt4cli uses OpenAI models, but you can use models from any provider that provides an OpenAI-compatible API, like OpenRouter.ai (Anthropic, Gemini, and open source models), Together.ai (open source models), Replicate, Ollama, and more.","sidebar":"tutorialSidebar"},"models/model-settings":{"id":"models/model-settings","title":"Model Settings","description":"Gpt4cli gives you a number of ways to control the models and models settings used in your plans. Changes to models and model settings are version controlled and can be branched.","sidebar":"tutorialSidebar"},"models/roles":{"id":"models/roles","title":"Model Roles","description":"Gpt4cli has multiple roles that are used for different aspects of its functionality. Each role can have its model and settings changed independently. These are the roles:","sidebar":"tutorialSidebar"},"quick-start":{"id":"quick-start","title":"Quickstart","description":"Install Gpt4cli","sidebar":"tutorialSidebar"},"security":{"id":"security","title":"Security","description":"Ignoring Sensitive Files","sidebar":"tutorialSidebar"}}}}')}}]);