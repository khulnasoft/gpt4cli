---
sidebar_position: 1
sidebar_label: Install
---

# Install Gpt4cli

## Quick Install

```bash
curl -sL https://v2.khulnasoft.com/install.sh | bash
```

## Manual install

Grab the appropriate binary for your platform from the latest [release](https://github.com/khulnasoft/gpt4cli/releases) and put it somewhere in your `PATH`.

## Build from source

```bash
git clone https://github.com/khulnasoft/gpt4cli.git
git checkout v2
cd gpt4cli/app/cli
go build -ldflags "-X gpt4cli/version.Version=$(cat version.txt)"
mv gpt4cli /usr/local/bin # adapt as needed for your system
```

## Windows

Windows is supported via [WSL](https://learn.microsoft.com/en-us/windows/wsl/about).

Gpt4cli only works correctly in the WSL shell. It doesn't work in the Windows CMD prompt or PowerShell.
