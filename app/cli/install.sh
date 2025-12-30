#!/usr/bin/env bash
#
# Gpt4cli Quick Installer
# This script installs the latest (or requested) version of Gpt4cli for supported platforms.
# See: https://docs.gpt4cli.khulnasoft.com

set -euo pipefail  # Exit on error, treat unset vars as errors, fail on pipeline errors


# Print usage/help
if [[ "${1:-}" =~ ^(-h|--help)$ ]]; then
  echo "Usage: bash install.sh"
  echo "Installs the latest (or requested) version of Gpt4cli for supported platforms."
  echo "Set GPT4CLI_VERSION to install a specific version."
  exit 0
fi

# Check for required commands
for cmd in curl tar tput uname grep mktemp; do
  if ! command -v "$cmd" >/dev/null 2>&1; then
    echo "ERROR: Required command '$cmd' not found. Please install it and retry." >&2
    exit 1
  fi
done

PLATFORM=""
ARCH=""
VERSION=""
RELEASES_URL="https://github.com/khulnasoft-lab/gpt4cli/releases/download"

# Ensure cleanup happens on exit and on specific signals
trap cleanup EXIT
trap cleanup INT TERM

# Use a secure temp directory
INSTALL_TMPDIR="$(mktemp -d -t gpt4cli_install_XXXXXX)"

cleanup () {
  # Only remove if it exists
  [ -d "$INSTALL_TMPDIR" ] && rm -rf "$INSTALL_TMPDIR"
}

# Set platform
case "$(uname -s)" in
 Darwin)
   PLATFORM='darwin'
   ;;

 Linux)
   PLATFORM='linux'
   ;;

 FreeBSD)
   PLATFORM='freebsd'
   ;;

 CYGWIN*|MINGW*|MSYS*)
   PLATFORM='windows'
   ;;

 *)
   echo "Platform may or may not be supported. Will attempt to install."
   PLATFORM='linux'
   ;;
esac

if [[ "$PLATFORM" == "windows" ]]; then
  echo "🚨 Windows is only supported via WSL. It doesn't work in the Windows CMD prompt or PowerShell."
  echo "How to install WSL 👉 https://learn.microsoft.com/en-us/windows/wsl/about"
  exit 1
fi

# Set arch (support more variants)
case "$(uname -m)" in
  x86_64|amd64)
    ARCH="amd64"
    ;;
  arm64|aarch64)
    ARCH="arm64"
    ;;
  armv7l)
    ARCH="armv7"
    ;;
  *)
    echo "Unknown architecture: $(uname -m)" >&2
    exit 1
    ;;
esac

# Detect Docker (portable)
if grep -qE '/docker/|/lxc/' /proc/1/cgroup 2>/dev/null || [ -f /.dockerenv ]; then
  IS_DOCKER=true
else
  IS_DOCKER=false
fi

# Set Version
if [[ -z "${GPT4CLI_VERSION}" ]]; then
  VERSION=$(curl -sL https://gpt4cli.khulnasoft.com/v2/cli-version.txt)
else
  VERSION=$GPT4CLI_VERSION
  echo "Using custom version $VERSION"
fi


welcome_gpt4cli () {
  echo ""
  printf '%*s' "$(tput cols)" '' | tr ' ' -
  echo ""
  echo "🚀 Gpt4cli v$VERSION • Quick Install"
  echo ""
  printf '%*s' "$(tput cols)" '' | tr ' ' -
  echo ""
}

download_gpt4cli () {
  ENCODED_TAG="cli%2Fv${VERSION}"
  url="${RELEASES_URL}/${ENCODED_TAG}/gpt4cli_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"

  cd "$INSTALL_TMPDIR"

  echo "📥 Downloading Gpt4cli tarball"
  echo ""
  echo "👉 $url"
  echo ""
  curl -s -L -o gpt4cli.tar.gz "${url}"

  tar zxf gpt4cli.tar.gz 1> /dev/null

  should_sudo=false

  if [ "$PLATFORM" = "darwin" ] || $IS_DOCKER ; then
    if [[ -d /usr/local/bin ]]; then
      if ! mv gpt4cli /usr/local/bin/ 2>/dev/null; then
        echo "Permission denied when attempting to move Gpt4cli to /usr/local/bin."
        if command -v sudo 2>/dev/null; then
          should_sudo=true
          echo "Attempting to use sudo to complete installation."
          sudo mv gpt4cli /usr/local/bin/
          if sudo mv gpt4cli /usr/local/bin/; then
            echo "✅ Gpt4cli is installed in /usr/local/bin"
            echo ""
          else
            echo "Failed to install Gpt4cli using sudo. Please manually move Gpt4cli to a directory in your PATH."
            exit 1
          fi
        else
          echo "sudo not found. Please manually move Gpt4cli to a directory in your PATH."
          exit 1
        fi
      else
        echo "✅ Gpt4cli is installed in /usr/local/bin"
      fi
    else
      echo >&2 'Error: /usr/local/bin does not exist. Create this directory with appropriate permissions, then re-install.'
      exit 1
    fi
  else
    if [ "$(id -u)" -eq 0 ]; then
      mv gpt4cli /usr/local/bin/
    elif command -v sudo 1>/dev/null 2>&1; then
      sudo mv gpt4cli /usr/local/bin/
      should_sudo=true
    else
      echo "ERROR: This script must be run as root or be able to sudo to complete the installation."
      exit 1
    fi
    echo "✅ Gpt4cli is installed in /usr/local/bin"
  fi

  # create 'g4c' alias, but don't overwrite existing g4c command
  if ! command -v g4c >/dev/null 2>&1; then
    echo "🎭 Creating g4c alias..."
    LOC=$(command -v gpt4cli)
    BIN_DIR="$(dirname "$LOC")"

    if [ "$should_sudo" = true ]; then
      if sudo ln -s "$LOC" "$BIN_DIR/g4c"; then
        echo "✅ Successfully created 'g4c' alias with sudo."
      else
        echo "⚠️ Failed to create 'g4c' alias even with sudo. Please create it manually."
      fi
    else
      if ln -s "$LOC" "$BIN_DIR/g4c"; then
        echo "✅ Successfully created 'g4c' alias."
      else
        echo "⚠️ Failed to create 'g4c' alias. Please create it manually."
      fi
    fi
  fi
}

check_existing_installation () {
  if command -v gpt4cli >/dev/null 2>&1; then
    existing_version=$(gpt4cli version 2>/dev/null || echo "unknown")
    # Check if version starts with 1.x.x
    if [[ "$existing_version" =~ ^1\. ]]; then
      echo "Found existing Gpt4cli v1.x installation ($existing_version). Renaming to 'gpt4cli1' before installing v2..."
      
      # Get the location of existing binary
      existing_binary=$(command -v gpt4cli)
      binary_dir="$(dirname "$existing_binary")"
      
      # Rename gpt4cli to gpt4cli1
      if ! mv "$existing_binary" "${binary_dir}/gpt4cli1" 2>/dev/null; then
        if command -v sudo 1>/dev/null 2>&1; then
          sudo mv "$existing_binary" "${binary_dir}/gpt4cli1"
        else
          echo "Cannot rename old gpt4cli without sudo. Please rename manually."
        fi
      fi
      
      # Rename g4c to g4c1 if it exists
      if [ -L "${binary_dir}/g4c" ]; then
        if ! mv "${binary_dir}/g4c" "${binary_dir}/g4c1" 2>/dev/null; then
          if command -v sudo 1>/dev/null 2>&1; then
            sudo mv "${binary_dir}/g4c" "${binary_dir}/g4c1"
          else
            echo "Cannot rename old g4c without sudo. Please rename manually."
          fi
        fi
        echo "Renamed 'g4c' alias to 'g4c1'"
      fi
      
      echo "Your v1.x installation is now accessible as 'gpt4cli1' and 'g4c1'"
    fi
  fi
}

welcome_gpt4cli
check_existing_installation
download_gpt4cli

echo ""
echo "🎉 Installation complete"
echo ""
printf '%*s' "$(tput cols)" '' | tr ' ' -
echo ""
echo "⚡️ Run 'gpt4cli' or 'g4c' in any project directory and start building!"
echo ""
printf '%*s' "$(tput cols)" '' | tr ' ' -
echo ""
echo "📚 Need help? 👉 https://docs.gpt4cli.khulnasoft.com"
echo ""
echo "👋 Join a community of AI builders 👉 https://discord.gg/khulnasoft"
echo ""
printf '%*s' "$(tput cols)" '' | tr ' ' -
echo ""

