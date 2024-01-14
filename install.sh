#!/bin/sh

UPSTREAM_REPO="https://github.com/gekkowrld/memo"

isDirThere() {
	test -d "$1"
}

isFileThere() {
	test -f "$1"
}

isZeroString() {
	test -z "$1"
}

# Determine the user shell
# Determine where I am before running
# First check the expected files if they exist
if [[ $(isFileThere "$PWD/main.go") && $(isFileThere "$PWD/cmd/root.go") ]]; then
	# This means I'm in a directory, so I should continue running
	echo "Anything may happen, please check the script and make sure everything is OK"

	exit 0
else
	# Clone  a repository and then give the install.sh of the cloned
	# repo the work and then exit after being done.
	#
	# If it enters in an infinite loop, then I have no way of knowing
	echo "I'll have to get the files from $UPSTREAM_REPO"
	# Check if the required commands are available
	ANY_MISSING=false
	for r_command in "curl" "grep" "sort" "head" "tar"; do
		if ! command -v "$r_command" 2>&1 >/dev/null; then
			echo "$r_command not found, please install it"
			ANY_MISSING=true
		fi
	done

	# If any was found, then exit the program early on
	if $ANY_MISSING; then
		exit 1
	fi

	# Code gotten from:
	# https://stackoverflow.com/a/54608917
	MEMO_LATEST_VERSION=$(curl -s "$UPSTREAM_REPO/tags" | grep -Eo "$Version v[0-9]{1,2}.[0-9]{1,2}.[0-9]{1,2}" | sort -r | head -n1 | tr -d ' ')

	if [ $(isZeroString "$MEMO_LATEST_VERSION") ]; then
		echo "Error: Couldn't determine the latest version of Memo."
		exit 1
	fi

	# Get the tar.gz file and untar it
	TAR_TARGET="$UPSTREAM_REPO/archive/refs/tags/$MEMO_LATEST_VERSION.tar.gz"
	echo "curl -sL $TAR_TARGET"
	curl -sL $TAR_TARGET | tar xz
	SAVED_AS=$(echo "memo-$MEMO_LATEST_VERSION" | tr -d 'v')

	SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
	# Run the script in the file
	if [ -e "$SCRIPT_DIR/$SAVED_AS/install.sh" ]; then
		sh "$SCRIPT_DIR/$SAVED_AS/install.sh"
	else
		echo "Error: Couldn't find the install script ($SCRIPT_DIR/$SAVED_AS/install.sh)"
		exit 1
	fi

	exit 0
fi

USERSHELL=$(basename "$SHELL")

# Identify the users (bash, zsh, fish), this are the options
# offered by cobra-cli

BASH_USER=false
FISH_USER=false
ZSH_USER=false
SH_USER=false # Incase all three fails, use sh as the default.
# Should basically do some things?

case $USERSHELL in
"bash") BASH_USER=true ;;
"fish") FISH_USER=true ;;
"zsh") ZSH_USER=true ;;
*) SH_USER=true ;;
esac

# Check if the user set a diffrent config location
USER_CONFIG="$HOME/.config/memo"
ENV_USER_CONFIG=$(echo $GMEMOCONFLC)

if [ $ENV_USER_CONFIG ]; then
	USER_CONFIG=$ENV_USER_CONFIG
fi

# Check if the config directory even exists
# Create it if it doesn't

if [ ! $(isDirThere "$USER_CONFIG") ]; then
	mkdir -p "$USER_CONFIG"
fi

# Now check if the shell completion is set in the ~/.bashrc file
if $BASH_USER; then
	USER_CONFIG=$USER_CONFIG"/memo"
	# Assuming that the file structure is as is in the upstream

	if [ $(ls -R $PWD | grep "memo.bash") ]; then
		# Update the file if it exists
		# I don't expect it to be user maintained, so overriding it is ok
		cp -uv ./completion/memo.bash $USER_CONFIG
	fi

	if ! grep -q "source \"$USER_CONFIG\"" "$HOME/.bashrc"; then
		echo "# Source completion file for memo binary" >>"$HOME/.bashrc"
		echo "source \"$USER_CONFIG\"" >>"$HOME/.bashrc"
	fi
fi
