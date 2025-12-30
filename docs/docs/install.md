---
sidebar_position: 1
sidebar_label: Install
---

# Install Gpt4cli

## Quick Install

```bash
curl -sL https://gpt4cli.khulnasoft.com/install.sh | bash
```

## Manual install

Grab the appropriate binary for your platform from the latest [release](https://github.com/khulnasoft-lab/gpt4cli/releases) and put it somewhere in your `PATH`.

## Build from source

```bash
git clone https://github.com/khulnasoft-lab/gpt4cli.git
cd gpt4cli/app/cli
go build -ldflags "-X gpt4cli/version.Version=$(cat version.txt)"
mv gpt4cli /usr/local/bin # adapt as needed for your system
```

## Windows

Windows is supported via [WSL](https://learn.microsoft.com/en-us/windows/wsl/about).

Gpt4cli only works correctly in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell.

## Upgrading from v1 to v2

When you install the Gpt4cli v2 CLI with the quick install script, it will rename your existing `gpt4cli` command to `gpt4cli1` (and the `g4c` alias to `g4c1`). Gpt4cli v2 is designed to run *separately* from v1 rather than upgrading in place. [More details here.](./upgrading-v1-to-v2.md)