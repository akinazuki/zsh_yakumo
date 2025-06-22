# 🚀 ZSH Yakumo

<div align="center">

[![Go Version](https://img.shields.io/badge/Go-1.23+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)](LICENSE)
[![Zsh](https://img.shields.io/badge/Shell-Zsh-89e051?style=for-the-badge&logo=gnu-bash)](https://www.zsh.org/)

*AI-powered command completion for your terminal* ✨

</div>

---

## 📖 Overview

**ZSH Yakumo** is an intelligent Zsh plugin that leverages OpenAI's GPT models to provide context-aware command completion and suggestions. Born from [zsh_codex](https://github.com/tom-doerr/zsh_codex) but completely rewritten in Go for better performance and easier deployment.

### ✨ Features

- 🤖 **AI-Powered Completions**: Get intelligent command suggestions based on context
- ⚡ **Fast & Lightweight**: Written in Go for optimal performance
- 🔧 **Easy Installation**: Single binary, no complex dependencies
- 📝 **Comment-Based Completion**: Type comments and get corresponding commands
- 🛡️ **Robust Error Handling**: Comprehensive logging and error management
- 🎛️ **Configurable Logging**: Environment-based log level control

### 🎯 Why Choose ZSH Yakumo?

| Feature | ZSH Yakumo (Go) | Original zsh_codex (Python) |
|---------|-----------------|------------------------------|
| **Installation** | ✅ Single binary | ❌ Complex Python dependencies |
| **Performance** | ✅ Fast startup | ❌ Slower Python execution |
| **Deployment** | ✅ Self-contained | ❌ Requires Python environment |
| **Maintenance** | ✅ Easy updates | ❌ Dependency conflicts |

---

## 🚀 Quick Start

### Prerequisites

- **Zsh shell** (oh-my-zsh recommended)
- **Go 1.23+** ([Install Go](https://go.dev/doc/install))
- **OpenAI API Key** ([Get one here](https://platform.openai.com/api-keys))

### 📦 Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/akinazuki/zsh_yakumo.git ~/.oh-my-zsh/custom/plugins/zsh_yakumo
   ```

2. **Build the plugin**
   ```bash
   cd ~/.oh-my-zsh/custom/plugins/zsh_yakumo
   go build -o dist/ .
   ```

3. **Configure your shell**
   
   Add to your `~/.zshrc`:
   ```bash
   plugins=(... zsh_yakumo)
   
   # Bind Ctrl+X for completion
   bindkey '^X' create_completion
   ```

4. **Set up configuration**
   
   Create `~/.config/zsh_yakumo.env`:
   ```bash
   # Required: Your OpenAI API key
   OPENAI_TOKEN="sk-your-openai-api-key-here"
   
   # Optional: Model selection (default: gpt-4o-mini)
   OPENAI_MODEL="gpt-4o-mini"
   
   # Optional: Log level (default: WARN)
   LOG_LEVEL="INFO"
   ```

5. **Reload your shell**
   ```bash
   source ~/.zshrc
   ```

---

## 🎮 Usage

### Basic Completion

Type a comment describing what you want to do, then press `Ctrl+X`:

```bash
# update brew packages
# → brew update && brew upgrade

# list all running processes
# → ps aux

# find large files in current directory
# → find . -type f -size +100M -exec ls -lh {} \;
```

### Advanced Examples

```bash
# create a git branch and switch to it
# → git checkout -b new-feature

# compress a directory to tar.gz
# → tar -czf archive.tar.gz directory/

# monitor system resources
# → top -o cpu
```

---

## ⚙️ Configuration

### Environment Variables

| Variable | Default | Description |
|----------|---------|-------------|
| `OPENAI_TOKEN` | *required* | Your OpenAI API key |
| `OPENAI_MODEL` | `gpt-4o-mini` | OpenAI model to use |
| `LOG_LEVEL` | `WARN` | Logging level: `DEBUG`, `INFO`, `WARN`, `ERROR` |

### Log Levels

- **DEBUG**: Detailed execution logs, API requests/responses
- **INFO**: General information about operations
- **WARN**: Warning messages and potential issues
- **ERROR**: Error messages only

### Custom Key Binding

Change the completion key binding in your `~/.zshrc`:

```bash
# Use Ctrl+Space instead of Ctrl+X
bindkey '^@' create_completion

# Use Ctrl+G
bindkey '^G' create_completion
```

---

## 🛠️ Development

### Building from Source

```bash
# Clone and build
git clone https://github.com/akinazuki/zsh_yakumo.git
cd zsh_yakumo
go mod tidy
go build -o dist/ .
```

### Project Structure

```
zsh_yakumo/
├── main.go                 # Main application logic
├── internal/
│   ├── logger/            # Logging functionality
│   └── defs/              # Type definitions
├── dist/                  # Built binaries
├── zsh_yakumo.plugin.zsh  # Zsh plugin script
└── zsh_yakumo.env         # Configuration template
```

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/amazing-feature`
3. Commit changes: `git commit -m 'Add amazing feature'`
4. Push to branch: `git push origin feature/amazing-feature`
5. Open a Pull Request

---

## 🐛 Troubleshooting

### Common Issues

**No completions generated**
- Check your OpenAI API key in `~/.config/zsh_yakumo.env`
- Verify your API key has sufficient credits
- Set `LOG_LEVEL=DEBUG` to see detailed logs

**Plugin not loading**
- Ensure the plugin is in your `plugins=()` array in `~/.zshrc`
- Verify the binary exists: `ls ~/.oh-my-zsh/custom/plugins/zsh_yakumo/dist/`
- Try reloading your shell: `source ~/.zshrc`

**Key binding not working**
- Check your key binding: `bindkey | grep create_completion`
- Try a different key combination
- Ensure no conflicts with existing bindings

### Debug Mode

Enable debug logging to troubleshoot issues:

```bash
LOG_LEVEL=DEBUG zsh  # Start a new shell with debug logging
```

Check logs at: `~/.config/zsh_yakumo.log`

---

## 📜 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 🙏 Acknowledgments

- [zsh_codex](https://github.com/tom-doerr/zsh_codex) - Original Python implementation
- [OpenAI](https://openai.com) - For providing the AI models
- [oh-my-zsh](https://ohmyz.sh/) - Zsh framework

---

<div align="center">

**Made with ❤️ and Go**

[Report Bug](https://github.com/akinazuki/zsh_yakumo/issues) · [Request Feature](https://github.com/akinazuki/zsh_yakumo/issues) · [Documentation](https://github.com/akinazuki/zsh_yakumo/wiki)

</div>
