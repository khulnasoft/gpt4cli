#!/usr/bin/env bash

set -e

SCRIPT_DIR=$( cd -- "$( dirname -- "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )

PLATFORM=
ARCH=
VERSION=
RELEASES_URL="https://github.com/khulnasoft/gpt4cli/releases/download"

 # Ensure cleanup happens on exit and on specific signals
trap cleanup EXIT
trap cleanup INT TERM

cleanup () {
  cd "${SCRIPT_DIR}"
  rm -rf gpt4cli_install_tmp
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

# Set arch
if [[ "$(uname -m)" == 'x86_64' ]]; then
  ARCH="amd64"
elif [[ "$(uname -m)" == 'arm64' || "$(uname -m)" == 'aarch64' ]]; then
  ARCH="arm64"
fi

if [[ "$(cat /proc/1/cgroup 2> /dev/null | grep docker | wc -l)" > 0 ]] || [ -f /.dockerenv ]; then
  IS_DOCKER=true
else
  IS_DOCKER=false
fi

# Set Version
if [[ -z "${GPT4CLI_VERSION}" ]]; then
  VERSION=$(curl -sL https://v2.khulnasoft.com/v2/cli-version.txt)
else
  VERSION=$GPT4CLI_VERSION
  echo "Using custom version $VERSION"
fi


welcome_gpt4cli () {
  echo ""
  echo "$(printf '%*s' "$(tput cols)" '' | tr ' ' -)"
  echo ""
  echo "🚀 Gpt4cli v$VERSION • Quick Install"
  echo ""
  echo "$(printf '%*s' "$(tput cols)" '' | tr ' ' -)"
  echo ""
}

download_gpt4cli () {
  ENCODED_TAG="cli%2Fv${VERSION}"

  url="${RELEASES_URL}/${ENCODED_TAG}/gpt4cli_${VERSION}_${PLATFORM}_${ARCH}.tar.gz"

  mkdir -p gpt4cli_install_tmp
  cd gpt4cli_install_tmp

  echo "📥 Downloading Gpt4cli tarball"
  echo ""
  echo "👉 $url"
  echo ""
  curl -s -L -o gpt4cli.tar.gz "${url}"

  tar zxf gpt4cli.tar.gz 1> /dev/null

  if [ "$PLATFORM" == "darwin" ] || $IS_DOCKER ; then
    if [[ -d /usr/local/bin ]]; then
      if ! mv gpt4cli /usr/local/bin/ 2>/dev/null; then
        echo "Permission denied when attempting to move Gpt4cli to /usr/local/bin."
        if hash sudo 2>/dev/null; then
          echo "Attempting to use sudo to complete installation."
          sudo mv gpt4cli /usr/local/bin/
          if [[ $? -eq 0 ]]; then
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
  elif [ "$PLATFORM" == "windows" ]; then
    # ensure $HOME/bin exists (it's in PATH but not present in default git-bash install)
    mkdir "$HOME/bin" 2> /dev/null
    mv gpt4cli.exe "$HOME/bin/"
    echo "✅ Gpt4cli is installed in '$HOME/bin'"
  else
    if [ $UID -eq 0 ]
    then
      # we are root
      mv gpt4cli /usr/local/bin/
    elif hash sudo 2>/dev/null;
    then
      # not root, but can sudo
      sudo mv gpt4cli /usr/local/bin/
    else
      echo "ERROR: This script must be run as root or be able to sudo to complete the installation."
      exit 1
    fi

    echo "✅ Gpt4cli is installed in /usr/local/bin"
  fi

  # create 'g4c' alias, but don't ovewrite existing g4c command
  if [ ! -x "$(command -v g4c)" ]; then
    echo "creating g4c alias"
    LOC=$(which gpt4cli)
    BIN_DIR=$(dirname $LOC)
    error_msg=$(ln -s "$LOC" "$BIN_DIR/g4c" 2>&1) || { echo "⚠️ Failed to create 'g4c' alias for Gpt4cli. Error: $error_msg. Please create it manually if needed."; }
  fi
}

check_existing_installation () {
  if command -v gpt4cli >/dev/null 2>&1; then
    existing_version=$(gpt4cli version 2>/dev/null || echo "unknown")
    # Check if version starts with 1.x.x
    if [[ "$existing_version" =~ ^1\. ]]; then
      echo "Found existing Gpt4cli v1.x installation ($existing_version). Renaming to 'gpt4cli1' before installing v2..."
      
      # Get the location of existing binary
      existing_binary=$(which gpt4cli)
      binary_dir=$(dirname "$existing_binary")
      
      # Rename gpt4cli to gpt4cli1
      if ! mv "$existing_binary" "${binary_dir}/gpt4cli1" 2>/dev/null; then
        sudo mv "$existing_binary" "${binary_dir}/gpt4cli1"
      fi
      
      # Rename g4c to g4c1 if it exists
      if [ -L "${binary_dir}/g4c" ]; then
        if ! mv "${binary_dir}/g4c" "${binary_dir}/g4c1" 2>/dev/null; then
          sudo mv "${binary_dir}/g4c" "${binary_dir}/g4c1"
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
echo "$(printf '%*s' "$(tput cols)" '' | tr ' ' -)"
echo ""
echo "⚡️ Run 'gpt4cli' or 'g4c' in any project directory and start building!"
echo ""
echo "$(printf '%*s' "$(tput cols)" '' | tr ' ' -)"
echo ""
echo "📚 Need help? 👉 https://docs.khulnasoft.com"
echo ""
echo "👋 Join a community of AI builders 👉 https://discord.gg/khulnasoft"
echo ""
echo "$(printf '%*s' "$(tput cols)" '' | tr ' ' -)"
echo ""

