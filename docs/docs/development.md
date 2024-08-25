---
sidebar_position: 10
sidebar_label: Development
---

# Development

To set up a development environment, first install dependencies:

- Go 1.21.3 - [install here](https://go.dev/doc/install)
- [reflex](https://github.com/cespare/reflex) 0.3.1 - for watching files and rebuilding in development. Install with `go install github.com/cespare/reflex@v0.3.1`
- PostgreSQL 14 - https://www.postgresql.org/download/

Make sure `$GOPATH` is in your $PATH

```bash
# print $GOPATH
echo $GOPATH

# if it's empty
export GOPATH=<path-to-go-folder>
```

Make sure the PostgreSQL server is running and create a database called `gpt4cli`.

Then make sure the following environment variables are set:

```bash
export DATABASE_URL=postgres://user:password@host:5432/gpt4cli?sslmode=disable # replace with your own database URL
export GOENV=development
```

Now from the root directory, run the script in `app/scripts/dev.sh`.

On Linux, you'll want to run this as `sudo` for copying the CLI to `/usr/local/bin` after it builds:

```bash
sudo ./app/scripts/dev.sh
```

You might also need sudo on MacOS if you don't have write permissions to `/usr/local/bin`, but this shouldn't be the case for most users. Assuming you have those write permissions, you can run the script without `sudo`:

```bash
./app/scripts/dev.sh
```

This creates watchers with `reflex` to rebuild both the server and the CLI when relevant files change.

The server runs on port 8080 by default.

After each build, the CLI is copied to `/usr/local/bin/gpt4cli-dev`so you can use it with just `gpt4cli-dev` in any directory. A `g4cd` alias is also created. Note the difference from the `gpt4cli` binary and `g4c` aliases which are installed for production usage—aliases are used for development to avoid overwriting the production install.

The output directory can be changed with the `GPT4CLI_DEV_CLI_OUT_DIR` environment variable. The binary name can be changed with `GPT4CLI_DEV_CLI_NAME` and the alias can be changed with `GPT4CLI_DEV_CLI_ALIAS`.

When running the Gpt4cli CLI, set `export GPT4CLI_ENV=development` to run in development mode, which connects to the development server by default.
