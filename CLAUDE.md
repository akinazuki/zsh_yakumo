# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a zsh plugin called `zsh_yakumo` that provides AI-powered command completion for zsh shells. It integrates with OpenAI's API to suggest command completions based on the current buffer context. This is a Go rewrite of [zsh_codex](https://github.com/tom-doerr/zsh_codex) for easier installation and deployment.

## Architecture

- **cmd/zsh_yakumo.go**: Core Go application that handles OpenAI API communication and command processing
- **zsh_yakumo.plugin.zsh**: Zsh plugin script that integrates with the shell and calls the Go binary
- **internal/defs/completion_request.go**: Go structs for OpenAI API request/response handling
- **internal/logger/logger.go**: Logging functionality with environment-based configuration
- **dist/zsh_yakumo**: Compiled Go binary (built from cmd/zsh_yakumo.go)
- **zsh_yakumo.env**: Configuration template for OpenAI credentials

## Configuration

The plugin requires configuration in `~/.config/zsh_yakumo.env`:
- `OPENAI_TOKEN`: Required OpenAI API key
- `OPENAI_MODEL`: Optional model selection (defaults to "gpt-4o-mini")

## Build Commands

```bash
# Build the Go binary (correct command for new structure)
go build -o dist/zsh_yakumo ./cmd/zsh_yakumo.go

# Install dependencies
go mod tidy
```

## Development Workflow

1. Modify Go source code in cmd/zsh_yakumo.go or internal/ packages
2. Rebuild the binary: `go build -o dist/zsh_yakumo ./cmd/zsh_yakumo.go`
3. Test the zsh plugin functionality by sourcing the plugin file

## Plugin Installation

The plugin integrates with oh-my-zsh and requires:
- Adding `zsh_yakumo` to plugins array in `.zshrc`
- Key binding: `bindkey '^X' create_completion` (Ctrl+X triggers completions)

## Key Components

- **Command Processing**: The Go binary reads stdin buffer and cursor position, sends context to OpenAI, and returns completion suggestions
- **Response Parsing**: Uses fastjson for efficient JSON parsing of OpenAI responses
- **Buffer Management**: Handles zsh buffer manipulation including cursor positioning and text insertion
- **Error Handling**: Graceful handling of API failures and malformed responses

## API Integration

The plugin uses OpenAI's completion API with structured output (JSON schema) to ensure consistent command completion format. The system prompt specifically instructs the AI to act as a zsh shell expert.