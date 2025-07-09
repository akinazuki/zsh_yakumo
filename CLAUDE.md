# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview

This is a zsh plugin called `zsh_yakumo` that provides AI-powered command completion for zsh shells. It integrates with OpenAI's API to suggest command completions based on the current buffer context. This is a Go rewrite of [zsh_codex](https://github.com/tom-doerr/zsh_codex) for easier installation and deployment.

## Architecture

The codebase follows a clean modular structure:

- **cmd/zsh_yakumo.go**: Main Go application that handles OpenAI API communication, request processing, and response parsing
- **zsh_yakumo.plugin.zsh**: Zsh plugin script that integrates with the shell, captures buffer context, and calls the Go binary
- **internal/defs/completion_request.go**: Go structs for OpenAI API request/response handling with JSON schema support
- **internal/logger/logger.go**: Structured logging with configurable levels and file output to `~/.config/zsh_yakumo.log`
- **dist/zsh_yakumo**: Compiled Go binary (built from cmd/zsh_yakumo.go)
- **zsh_yakumo.env**: Configuration template for OpenAI credentials

## Configuration

The plugin requires configuration in `~/.config/zsh_yakumo.env`:
- `OPENAI_TOKEN`: Required OpenAI API key
- `OPENAI_MODEL`: Optional model selection (defaults to "gpt-4o-mini")
- `LOG_LEVEL`: Optional logging level (DEBUG, INFO, WARN, ERROR - defaults to WARN)

## Build Commands

```bash
# Build the Go binary (primary build command)
go build -o dist/zsh_yakumo ./cmd/zsh_yakumo.go

# Install/update dependencies
go mod tidy

# Test the plugin (after building)
source zsh_yakumo.plugin.zsh
```

## Development Workflow

1. Modify Go source code in cmd/zsh_yakumo.go or internal/ packages
2. Rebuild the binary: `go build -o dist/zsh_yakumo ./cmd/zsh_yakumo.go`
3. Test the zsh plugin functionality by sourcing the plugin file or restarting zsh

## Plugin Installation & Usage

The plugin integrates with oh-my-zsh and requires:
- Adding `zsh_yakumo` to plugins array in `.zshrc`
- Key binding: `bindkey '^X' create_completion` (Ctrl+X triggers completions)
- Usage: Type comments like `# update brew packages` and press Ctrl+X for AI completions

## Key Components

- **Command Processing**: The Go binary reads stdin buffer and cursor position, sends context to OpenAI with structured JSON schema
- **Response Parsing**: Uses fastjson for efficient JSON parsing of OpenAI responses, extracting command from structured output
- **Buffer Management**: Handles zsh buffer manipulation including cursor positioning, text insertion, and prefix/suffix handling
- **Error Handling**: Comprehensive logging and graceful handling of API failures, malformed responses, and configuration issues
- **API Integration**: Uses OpenAI's completion API with structured output (JSON schema) to ensure consistent command completion format

## System Prompt

The system prompt emphasizes zsh expertise and preference for single quotes over double quotes to avoid shell expansion issues: "You are a zsh shell expert, please help me complete the following command, you should only output the completed command, no need to include any other explanation. Do not put completed command in a code block. When using quotes in commands, prefer single quotes over double quotes to avoid shell expansion issues."

## Dependencies

- `github.com/joho/godotenv`: Environment variable loading from config files
- `github.com/valyala/fastjson`: Fast JSON parsing for API responses
- Standard Go libraries for HTTP, JSON, and file operations