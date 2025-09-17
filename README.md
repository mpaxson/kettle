# Kettle

A command-line tool for installing development tools and managing terminal configurations.

## Quick Start

```bash
wget -O ~/.local/bin/kettle https://github.com/mpaxson/kettle/releases/latest/download/kettle
chmod +x ~/.local/bin/kettle
kettle install
```

Ensure `~/.local/bin` is in your `PATH`.

now you will have tab completion and help docs for all your commands

### update

run `kettle update to download the latest version from github`

## What it does

Setting up a new development machine usually means installing the same tools over and over again - Go, various linters, terminal emulators, shell utilities, and getting all their configurations just right. Kettle automates this process by downloading the right binaries for your platform and setting up everything with sensible defaults.

## Why?

Managing updates to terminals as tools change isntallations processes differ, etc.. can be a pain. This allows a seperate .kettle.<zsh|fish|bash> profile to be updated and edited with per installation setup

I origininally was using a series of bash scripts in my dot files which became a mess, and I didn't have easy ways to use github releases api in a clean format for new tools i come across

Automatically adds tool completion to bash, zsh, fish when installing

- automatic updating of kettle, user prompting

### Kettle

Example commands include `kettle tools terminal ghostty bind-f1` this will create a ghostty tty script, run gsettings commands, and then bind a full screen view of ghostty to your terminal

## Features

### Tool Installation

- Downloads tools directly from GitHub releases
- Picks the right binary for your OS and architecture automatically
- Extracts archives (tar.gz, zip, .deb) and puts binaries where they belong
- Planned handling of .deb packages when available

### Language Support

- **Go**: Installs Go toolchain and sets up workspace
- **golangci-lint**: Go linting with shell completions
- More languages planned (Python, Rust, Node.js)

### Terminal Setup

- **Ghostty & Kitty**: Modern terminal emulator installation
- **zoxide**: Smart directory navigation
- **autoenv**: Automatic environment loading
- Manages shell profiles in one place (`~/.config/kettle/kettle.<bashrc|zshrc|fishrc`)

### Updates

- Checks versions and prompts before updating
- Replaces binaries safely
- Uses Charmbracelet libraries for nice terminal UI

## Why use it?

**Individual developers**: Stop spending hours setting up each new machine. Get a consistent environment quickly.

**Teams**: Everyone gets the same tool versions. New team members run one command instead of following a long setup document.

**Cross-platform**: Works the same on Linux, macOS, and Windows.

## How it works

Kettle looks at GitHub releases, ranks the available files by how well they match your system, downloads the best one, and extracts it if needed. It's built with Go using Cobra for CLI and handles errors properly.

The ranking system prefers:

1. Standalone binaries
2. Archive files (.tar.gz, .zip)
3. Package files (.deb)

It ignores source code archives completely.
