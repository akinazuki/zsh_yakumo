# zsh_yakumo

A Zsh plugin that integrates with OpenAI's GPT models to provide intelligent command completion and suggestions
based on [zsh_codex](https://github.com/tom-doerr/zsh_codex), but rewritten in Go.  

## Why Go?
the [zsh_codex](https://github.com/tom-doerr/zsh_codex) is written in Python, it's very difficult to install and use on some systems.
So I decided to rewrite it in Go, which is a easy to install and use language.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/akinazuki/zsh_yakumo.git ~/.oh-my-zsh/custom/plugins/zsh_yakumo
   ```

2. Build the plugin:
   ```bash
   cd ~/.oh-my-zsh/custom/plugins/zsh_yakumo
   go build -o dist/ .
   ```

   Ensure you have Go installed. If not, you can install it at [here](https://go.dev/doc/install).

3. Add the plugin to your `.zshrc` file:
   ```bash
   plugins=(zsh_yakumo)
   bindkey '^X' create_completion
   ```

4. Create a file called `zsh_yakumo.env` in `~/.config`. Example
   ```ini
   OPENAI_MODEL="gpt-4o-mini"
   OPENAI_TOKEN="YOUR_OPENAI_API_KEY_HERE"
   ```
