#!/usr/bin/env bash
set -e

echo Installing i18n-pruner...

TARGET_DIR="$HOME/.i18n-pruner/bin"
TARGET_DIR_BIN="$TARGET_DIR/i18n-pruner"

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
curl -# -L "${SERVER}/${FILENAME}" -o "${TARGET_DIR_BIN}"
chmod +x "${TARGET_DIR_BIN}"
echo "Installed under ${TARGET_DIR_BIN}"

# Store the correct profile file (i.e. .profile for bash or .zshenv for ZSH).
case $SHELL in
*/zsh)
	PROFILE=${ZDOTDIR-"$HOME"}/.zshenv
	PREF_SHELL=zsh
	;;
*/bash)
	PROFILE=$HOME/.bashrc
	PREF_SHELL=bash
	;;
*/fish)
	PROFILE=$HOME/.config/fish/config.fish
	PREF_SHELL=fish
	;;
*/ash)
	PROFILE=$HOME/.profile
	PREF_SHELL=ash
	;;
*)
	echo "could not detect shell, manually add ${TARGET_DIR_BIN} to your PATH."
	exit 1
	;;
esac

# Only add if it isn't already in PATH.
if [[ ":$PATH:" != *":${TARGET_DIR_BIN}:"* ]]; then
	# Add the directory to the path and ensure the old PATH variables remain.
	echo >>$PROFILE && echo "export PATH=\"\$PATH:$TARGET_DIR_BIN\"" >>$PROFILE
fi

echo && echo "Detected your preferred shell is ${PREF_SHELL} and added i18n-pruner to PATH. Run 'source ${PROFILE}' or start a new terminal session to use i18n-pruner."

# Confirmation
echo "i18n-pruner successfully installed!"
