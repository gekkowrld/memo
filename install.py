#!/bin/python3
# -*- coding: utf-8 -*-

UPSTREAM_REPO = "https://github.com/gekkowrld/memo"


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
    # Check if gh is installed and use it
    import subprocess

    try:
        subprocess.run(["gh", "repo", "clone", UPSTREAM_REPO, "memo"])
    except FileNotFoundError:
        # Use git
        subprocess.run(["git", "clone", UPSTREAM_REPO])


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

    # Compile the code with no debug info and optimizations
    subprocess.run(["go", "build", "-ldflags", "-s -w", "-o", "memo"])

    # Install the binary
    # Check if GOPATH is set
    if "GOPATH" not in os.environ:
        print(
            f"GOPATH is not set. I'll install in {os.path.expanduser('~/.local/bin')}"
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
            f"XDG_DATA_HOME is not set. I'll install in {os.path.expanduser('~/.local/share')}"
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


if __name__ == "__main__":
    # Check if in repo (this will run install() if in repo)
    install()
