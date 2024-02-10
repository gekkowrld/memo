#!/bin/python3
# -*- coding: utf-8 -*-

UPSTREAM_REPO = "https://github.com/gekkowrld/memo"
PRISM_CSS = "https://gist.github.com/gekkowrld/11622b8dd1adcb783d6fa23ed1883338/raw/7ebdac23526b21e4e30486357de9aed2aa84b815/prism.css"
PRISM_JS = "https://gist.github.com/gekkowrld/11622b8dd1adcb783d6fa23ed1883338/raw/7ebdac23526b21e4e30486357de9aed2aa84b815/prism.js"
PRISM_GIST = "https://gist.github.com/gekkowrld/11622b8dd1adcb783d6fa23ed1883338"


def check_if_in_repo():
    """Check if in repo"""
    import os

    if not os.path.exists(os.path.join(os.getcwd(), "main.go")):
        print("You are not in a git repository. I'll have to clone it from upstream")
        clone_from_upstream()
        # Move into the repo
        os.chdir("memo")


def clone_from_upstream():
    """Clone from upstream"""
    # Check if gh (for network sake and slow internet) is installed and use it
    import subprocess

    clone_repo = False
    try:
        subprocess.run(["gh", "repo", "clone", UPSTREAM_REPO, "memo"])
        clone_repo = True
    except FileNotFoundError:
        pass

    if not clone_repo:
        try:
            subprocess.run(["git", "clone", UPSTREAM_REPO, "memo"])
            clone_repo = True
        except FileNotFoundError:
            pass

    if not clone_repo:
        print("You should install either git or gh to clone the repo")
        sys.exit(1)


def install():
    """Compile and install

    Compile the golang code, installing the dependancy and then install the
        binary to the $GOPATH/bin directory.
    Then copy the assets to the $XDG_DATA_HOME directory.
    Write the initial config to the $XDG_CONFIG_HOME directory.
    """
    import os
    import shutil
    import subprocess
    import sys

    # Check if in repo
    check_if_in_repo()

    # Check if go is installed
    try:
        subprocess.run(["go", "version"])
    except FileNotFoundError:
        print("Go is not installed. Please install it and try again.")
        sys.exit(1)

    # Download the prism files (for language highlighting)
    download_prism = False
    try:
        subprocess.run(["wget", "-c", PRISM_JS, "-O", "assets/prism.js"])
        subprocess.run(["wget", "-c", PRISM_CSS, "-O", "assets/prism.css"])
        download_prism = True
    except FileNotFoundError:
        pass

    if not download_prism:
        try:
            # Use curl to download the files
            subprocess.run(["curl", "-o", "assets/prism.js", PRISM_JS])
            subprocess.run(["curl", "-o", "assets/prism.css", PRISM_CSS])
            download_prism = True
        except FileNotFoundError:
            pass

    if not download_prism:
        print(
            f"Couldn't download the files for language highlighting, find them here: {PRISM_GIST}\n"
        )

    # Compile the code with no debug info and optimizations
    print("Building the go app with verbose and optimizations ON\n")
    subprocess.run(
        ["go", "build", "-x", "-gcflags", "-l=4 -m", "-ldflags", "-s -w", "-o", "memo"]
    )

    # Install the binary
    # Check if GOPATH is set
    if "GOPATH" not in os.environ:
        print(
            f"GOPATH is not set. I'll install the binary in {os.path.expanduser('~/.local/bin')}"
        )
        # Check if $HOME/.local/bin exists
        if not os.path.exists(os.path.expanduser("~/.local/bin")):
            # Create it
            os.makedirs(os.path.expanduser("~/.local/bin"))
        # Install
        shutil.move("memo", os.path.expanduser("~/.local/bin"))
    else:
        # Install
        shutil.move("memo", os.environ["GOPATH"])

    # Copy the assets
    # Check if XDG_DATA_HOME is set
    if "XDG_DATA_HOME" not in os.environ:
        print(
            f"XDG_DATA_HOME is not set. I'll put assets in {os.path.expanduser('~/.local/share')}"
        )
        # Check if $HOME/.local/share exists
        if not os.path.exists(os.path.expanduser("~/.local/share")):
            # Create it
            os.makedirs(os.path.expanduser("~/.local/share"))
        # Copy
        shutil.copytree("assets", os.path.expanduser("~/.local/share/memo"))
    else:
        # Copy
        shutil.copytree("assets", os.environ["XDG_DATA_HOME"] + "/memo")

    # Check if there is a config, create it if not write the default values
    """Default config values (~/.config/memo/config.toml)
    memodir = "~/.memo"
    editor = "vim"
    listfgcolour = ""
    listbgcolour = ""
    displaywidth = 0
    editconfig = false
    git = false
    """

    config = {
        "memodir": "~/.memo",
        "editor": "vim",
        "listfgcolour": "",
        "listbgcolour": "",
        "displaywidth": 0,
        "editconfig": False,
        "git": False,
        "staticfiles": "~/.local/share/memo/",
    }

    # Check if XDG_CONFIG_HOME is set
    if "XDG_CONFIG_HOME" not in os.environ:
        print(
            f"XDG_CONFIG_HOME is not set. I'll install in {os.path.expanduser('~/.config')}"
        )
        # Check if $HOME/.config exists
        if not os.path.exists(os.path.expanduser("~/.config")):
            # Create it
            os.makedirs(os.path.expanduser("~/.config"))
        # Check if $HOME/.config/memo exists
        if not os.path.exists(os.path.expanduser("~/.config/memo")):
            # Create it
            os.makedirs(os.path.expanduser("~/.config/memo"))
        # Check if $HOME/.config/memo/config.toml exists
        if not os.path.exists(os.path.expanduser("~/.config/memo/config.toml")):
            # Create it
            with open(os.path.expanduser("~/.config/memo/config.toml"), "w") as f:
                for key, value in config.items():
                    f.write(f"{key} = {value}\n")
    else:
        # Check if $XDG_CONFIG_HOME/memo exists
        if not os.path.exists(os.environ["XDG_CONFIG_HOME"] + "/memo"):
            # Create it
            os.makedirs(os.environ["XDG_CONFIG_HOME"] + "/memo")
        # Check if $XDG_CONFIG_HOME/memo/config.toml exists
        if not os.path.exists(os.environ["XDG_CONFIG_HOME"] + "/memo/config.toml"):
            # Create it
            with open(os.environ["XDG_CONFIG_HOME"] + "/memo/config.toml", "w") as f:
                for key, value in config.items():
                    f.write(f"{key} = {value}\n")

    # Check if ~./bashrc exists and check if the file contains:
    # source ~/.local/share/memo/memo.bash
    # or the xdg data directory.
    # If not, then add the line to it, and update the file
    # location.
    # If the file exists, update it and tell the user of the update.

    # Define the path to .bashrc and the source command
    bashrc_path = os.path.expanduser("~/.bashrc")
    memo_p = "/memo/memo.bash"
    source_command = (
        f"source {os.environ.get('XDG_DATA_HOME', '~/.local/share')}{memo_p}\n"
    )

    # Check if .bashrc exists
    if os.path.exists(bashrc_path):
        # Read the contents of .bashrc
        with open(bashrc_path, "r") as bashrc_file:
            bashrc_contents = bashrc_file.readlines()

        # Check if the source command is already in .bashrc
        if any(source_command in line for line in bashrc_contents):
            print(".bashrc is already configured with memo.")
        else:
            # Append the source command to .bashrc
            with open(bashrc_path, "a") as bashrc_file:
                bashrc_file.write(source_command)
            print("Updated .bashrc to configure memo.")

        # Update the memo.bash file.
        shutil.copy(
            "completion/memo.bash",
            os.environ.get(f"XDG_DATA_HOME{memo_p}", "~/.local/share{memo_p}"),
        )


if __name__ == "__main__":
    # Check if in repo (this will run install() if in repo)
    install()
