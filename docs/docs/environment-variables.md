---
sidebar_position: 11
sidebar_label: Environment Variables
---

# Environment Variables

This is an overview of all the environment variables that can be used with Gpt4cli.

## CLI

### LLM Providers

```bash
OPENAI_API_BASE= # Your OpenAI server, such as http://localhost:1234/v1 Defaults to empty.
OPENAI_API_KEY= # Your OpenAI key.

# optional - set API keys for any other providers you're using
export OPENROUTER_API_KEY= # Your OpenRouter.ai API key.
export TOGETHER_API_KEY = # Your Together.ai API key.
# etc.
```

### Upgrades

```bash
GPT4CLI_SKIP_UPGRADE= # Set this to '1' to skip the auto-upgrade check when running the CLI.
```

### Development

Check out the [Development Guide](./development.md) for more details.

```bash
GPT4CLI_ENV=development # Set this to 'development' to default to the local development server instead of Gpt4cli Cloud when working on Gpt4cli itself.
GPT4CLI_API_HOST= # Defaults to http://localhost:8080 if GPT4CLI_ENV is development, otherwise it's https://api.gpt4cli.khulnasoft.com—override this to use a different host.
GPT4CLI_OUT_DIR=/usr/local/bin # Where the development binary should be output when using dev.sh
GPT4CLI_DEV_CLI_OUT_DIR=/usr/local/bin # Where the development binary should be output when using dev.sh
GPT4CLI_DEV_CLI_NAME=gpt4cli-dev # The name of the development binary when using dev.sh
GPT4CLI_DEV_CLI_ALIAS=g4cd # The alias for the development binary when using dev.sh
GOPATH= # This should be already set to your Go folder if you've installed Golang.
```

## Server

Check out the [Self-Hosting Guide](./hosting/self-hosting.md) for more details.

### General

```bash
GOENV=development # Whether to run in development or production mode. Must be 'development' or 'production'
GPT4CLI_BASE_DIR= # The base directory to read and write files. Defaults to '$HOME/gpt4cli-server' in development mode, '/gpt4cli-server' in production.
PORT=8080 # The port the server listens on. Defaults to 8080.
```

### docker-compose

For self-hosting with docker-compose, default environment variables are set in `app/_env`. This file should be copied to `app/.env` before running the server. You can override any of these defaults in `.env`. 

```bash
GPT4CLI_DATA_DIR=/var/lib/gpt4cli/data # When using docker-compose, this is the directory *on your machine* that the Gpt4cli server will use to store data—it will be mounted to the Docker container as a volume.
PGDATA_DIR=/var/lib/postgresql/data # Where PostgreSQL should store its data.

# Database Credentials
POSTGRES_DATABASE=gpt4cli # Your postgres database.
POSTGRES_USER=gpt4cli # Your postgres user.
POSTGRES_PASSWORD=gpt4cli # Your postgres password.
```

### Other methods

If you're *not* using docker-compose, you'll need a `DATABASE_URL` environment variable that points to a PostgreSQL database. For example, if you're running PostgreSQL locally, you might set it to something like this:

```bash
DATABASE_URL=postgres://gpt4cli:<password>@<host>:<port>/gpt4cli?sslmode=disable
```

### SMTP

If you're running in production mode (with `GOENV=production`, typically on a remote server), you'll need SMTP credentials:

```bash
SMTP_HOST= # Your SMTP host.
SMTP_PORT= # Set this to 1025 e.g. if you are using mailhog.
SMTP_USER= # SMTP username.
SMTP_PASSWORD= # SMTP password.
```
