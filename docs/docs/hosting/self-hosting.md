---
sidebar_position: 2
sidebar_label: Self-Hosting
---

# Self-Hosting

Gpt4cli is open source and uses a client-server architecture. The server can be self-hosted. You can run either run it locally or on a cloud server that you control. 

## Quickstart

The self-hosting quickstart requires git, docker, and docker-compose. It's designed for local use with a single user. If you instead want to run Gpt4cli on a remote server with multiple users or orgs, continue on to the [Requirements](#requirements) section below.

1. Run the server in development mode: 

```bash
git clone https://github.com/khulnasoft/gpt4cli.git
cd gpt4cli/app
./start_local.sh
```

2. Install the Gpt4cli CLI if you haven't already:

```bash
curl -sL https://gpt4cli.khulnasoft.com/install.sh | bash
```

3. Then run:

```bash
gpt4cli sign-in
```

4. Follow the prompts from there to create a new account on your self-hosted server. From there, check out the more general [CLI quickstart](../quick-start.md) to get fully up and running.

## Requirements

The Gpt4cli server requires a PostgreSQL database (ideally v14), a persistent file system, and git.

## Development vs. Production

The Gpt4cli server can be run in development or production mode. The main differences are how authentication pins and emails are handled, and the default path for the persistent file system.

Development mode is designed for local usage with a single user. Email isn't enabled. Authentication pins are copied to the clipboard instead of sent via email, and a system notification will pop up to let you know that the pin is ready to paste. The pin will also be output to the console. In development mode, the default base directory is `$HOME/gpt4cli-server`.

Production mode is designed for multiple users or organizations. Email is enabled and SMTP environment variables are required. Authentication pins are sent via email. In production mode, the default base directory is `/gpt4cli-server`.

Development or production mode is set with the `GOENV` environment variable. It should be set to either `development` or `production`.

In both development and production mode, the server runs on port 8080 by default. This can be changed with the `PORT` environment variable.

## docker-compose

For local usage in development mode, you can skip setting up PostgreSQL if you use `docker-compose` with the included `docker-compose.yml` file:

```bash
git clone https://github.com/khulnasoft/gpt4cli.git
cd gpt4cli/app
cp _env .env
# edit .env to override any default environment variables
docker compose build
docker compose up
```

## Other Methods

### PostgreSQL Database

If you aren't using docker-compose, you'll need a PostgreSQL database. You can run the following SQL to create a user and database, replacing `user` and `password` with your own values:

```sql
CREATE USER 'user' WITH PASSWORD 'password';
CREATE DATABASE 'gpt4cli' OWNER 'user';
GRANT ALL PRIVILEGES ON DATABASE 'gpt4cli' TO 'user';
```

### Environment Variables

Set `GOENV` to either `development` or `production` as described above in the [Development vs. Production](#development-vs-production) section:

```bash
export GOENV=development
```

or
  
```bash
export GOENV=production
```

You'll also need a `DATABASE_URL`:

```bash
export DATABASE_URL=postgres://user:password@host:5432/gpt4cli # replace with your own database URL
```

If you're running in production mode, you'll need to connect to SMTP to send emails. Set the following environment variables:

```bash
export SMTP_HOST=smtp.example.com
export SMTP_PORT=587
export SMTP_USER=user
export SMTP_PASSWORD=password
export SMTP_FROM=user@example.com # optional, if not set then SMTP_USER is used
```

In either development or production mode, the base directory for the persistent file system can optionally be overridden with the `GPT4CLI_BASE_DIR` environment variable:

```bash
export GPT4CLI_BASE_DIR=~/some-dir/gpt4cli-server
```

### Using Docker Build

The server can be run from a Dockerfile at `app/Dockerfile.server`:

```bash
git clone https://github.com/khulnasoft/gpt4cli.git
VERSION=$(cat app/server/version.txt) # or use the version you want
git checkout server/v$VERSION
cd gpt4cli/app
mkdir ~/gpt4cli-server # or another directory where you want to store files
docker build -t gpt4cli-server -f Dockerfile.server .
docker run -p 8080:8080 \
  -v ~/gpt4cli-server:/gpt4cli-server \
  -e DATABASE_URL \
  -e GOENV \
  -e SMTP_HOST \ 
  -e SMTP_PORT \
  -e SMTP_USER \
  -e SMTP_PASSWORD \
  gpt4cli-server
```

The SMTP environment variables above are only required if you're running in [production mode](#development-vs-production).

### Run From Source

You can also run the server from source:

```bash
git clone https://github.com/khulnasoft/gpt4cli.git
cd gpt4cli/
VERSION=$(cat app/server/version.txt) # or use the version you want
git checkout server/v$VERSION
cd app/server
export GPT4CLI_BASE_DIR=~/gpt4cli-server # or another directory where you want to store files
go run main.go
```

## Health Check

You can check if the server is running by sending a GET request to `/health`. If all is well, it will return a 200 status code.

## Create a New Account

Once the server is running and you've [installed the Gpt4cli CLI](../install.md) on your local development machine, you can create a new account by running `gpt4cli sign-in`: 

```bash
gpt4cli sign-in # follow the prompts to create a new account on your self-hosted server
```

## Note On Local CLI Files

If you use the Gpt4cli CLI and then for some reason you reset the database or use a new one, you'll need to remove the local files that the CLI creates in directories where you used Gpt4cli in order to start fresh. Otherwise, the CLI will attempt to authenticate with an account that doesn't exist in the new database and you'll get errors. This could also happen if you use Gpt4cli Cloud and then switch to self-hosting.

To resolve this, remove the following in any directory you used the CLI in:

- `.gpt4cli-dev` directory if you ran the CLI with `GPT4CLI_ENV=development`
- `.gpt4cli` directory otherwise

Then run `gpt4cli sign-in` again to create a new account.

If you're still having trouble with accounts, you can also remove the following from your $HOME directory to fully reset them:

- `.gpt4cli-home-dev` directory if you ran the CLI with `GPT4CLI_ENV=development`
- `.gpt4cli-home` directory otherwise