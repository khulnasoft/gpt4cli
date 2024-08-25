---
sidebar_position: 9
sidebar_label: Security
---

# Security
## Ignoring Sensitive Files

Gpt4cli respects `.gitignore` and won't load any files that you're ignoring unless you use the `--force/-f` flag with `gpt4cli load`. You can also add a `.gpt4cliignore` file with ignore patterns to any directory.

## API Key Security

Gpt4cli is a bring-your-own-API-key tool. On the Gpt4cli server, whether that's Gpt4cli Cloud or a self-hosted server, API keys are only stored ephemerally in RAM while they are in active use. They are never written to disk, logged, or stored in a database. As soon as a plan stream ends, the API key is removed from memory and no longer exists anywhere on the Gpt4cli server.

It's also up to you to manage your API keys securely. Try to avoid storing them in multiple places, exposing them to third party services, or sending them around in plain text.