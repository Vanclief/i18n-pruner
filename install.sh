#!/usr/bin/env bash
set -e

echo Installing i18n-pruner...

TARGET_DIR="$HOME/.i18n-pruner/bin"
SERVER="https://github.com/Vanclief/i18n-pruner/raw/master/bin"

# Detect the platform (similar to $OSTYPE)
OS="$(uname)"
if [[ "$OS" == "Linux" ]]; then
	# Linux
	FILENAME="i18n-pruner-linux"
elif [[ "$OS" == "Darwin" ]]; then
	# MacOS, should validate if Intel or ARM
	UNAMEM="$(uname -m)"
	if [[ "$UNAMEM" == "x86_64" ]]; then
		FILENAME="i18n-pruner-mac-intel"
	else
		FILENAME="i18n-pruner-mac-arm64"
	fi
else
	echo "unrecognized OS: $OS"
	echo "Exiting..."
	exit 1
fi

# Check if ~/.i18n-pruner/bin exists, if not create it
if [[ ! -e "${TARGET_DIR}" ]]; then
	mkdir -p "${TARGET_DIR}"
fi

# Download the appropriate binary
echo "Downloading $SERVER/$FILENAME..."
curl -# -L "${SERVER}/${FILENAME}" -o "${TARGET_DIR}/i18n-pruner"
chmod +x "${TARGET_DIR}/i18n-pruner"
echo "Installed under ${TARGET_DIR}/i18n-pruner"

# Update PATH depending on the shell
if [[ "${SHELL}" == *"zsh"* ]]; then
	SHELL_PROFILE="$HOME/.zshrc"
elif [[ "${SHELL}" == *"bash"* ]]; then
	SHELL_PROFILE="$HOME/.bashrc"
else
	echo "Shell not supported. Please add the line 'export PATH=\"\$HOME/.i18n-pruner/bin:\$PATH\"' to your shell profile manually."
	exit 1
fi

echo 'export PATH="$HOME/.i18n-pruner/bin:$PATH"' >>${SHELL_PROFILE}
source ${SHELL_PROFILE}

# Confirmation
echo "i18n-pruner successfully installed!"
