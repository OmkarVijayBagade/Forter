# forter

A high-performance, production-ready terminal user interface (TUI) file organizer tool built with Go and Bubble Tea.

[![Go Version](https://img.shields.io/badge/go-%3E%3D1.21-blue)](https://golang.org)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

![Demo](docs/demo.gif)

## Features

- **🚀 High Performance**: Handles 1000+ files smoothly with non-blocking operations
- **🎨 Beautiful TUI**: Clean, minimal interface inspired by lazygit with keyboard-driven navigation
- **📂 Smart Categorization**: Automatically sorts files into Documents, Images, Videos, Audio, Archives, Code, and Others
- **✅ Bulk Operations**: Select multiple files and organize them in one go
- **👁️ Preview Mode**: See destination paths before applying changes
- **🔒 Safe Operations**: Duplicate file handling with auto-renaming, graceful error handling
- **🧪 Dry Run**: Test organization plans without moving any files
- **↩️ Undo Support**: Undo the last organization operation
- **⚡ Cross-Platform**: Works on macOS and Linux

## Installation

### Homebrew (Recommended)

```bash
brew tap OmkarVijayBagade/forter
brew install forter
```

### From Source

**Prerequisites**: Go 1.21 or later

```bash
git clone https://github.com/OmkarVijayBagade/forter.git
cd forter
go build -o forter cmd/forter/main.go

# Install to $GOPATH/bin
make install
```

### Binary Releases

Download pre-built binaries from the [Releases](https://github.com/OmkarVijayBagade/forter/releases) page.

## Usage

**📖 For detailed usage instructions, see [HOW_TO_USE.md](HOW_TO_USE.md)**

### Basic Usage

```bash
# Organize current directory
forter

# Organize specific directory
forter ~/Downloads

# Recursively organize (including subdirectories)
forter -r ~/Desktop

# Preview changes without moving files (dry run)
forter --dry-run ~/Documents
```

### TUI Navigation

| Key | Action |
|-----|--------|
| `↑`/`↓` or `k`/`j` | Navigate up/down |
| `Tab` or `←`/`→`/`h`/`l` | Switch panels (files ↔ categories) |
| `Space` or `Enter` | Select/deselect file or category |
| `a` | Select all files |
| `c` | Clear selection |
| `o` | Start organization |
| `d` | Toggle dry run mode |
| `u` | Undo last operation |
| `q` or `Esc` | Quit |
| `?` | Show help |

### File Organization Structure

Files are organized into the following structure:

```
YourDirectory/
├── Documents/
│   ├── PDF/
│   ├── DOCX/
│   └── TXT/
├── Images/
│   ├── JPG/
│   ├── PNG/
│   └── SVG/
├── Videos/
│   ├── MP4/
│   ├── AVI/
│   └── MKV/
├── Audio/
│   ├── MP3/
│   └── FLAC/
├── Archives/
│   ├── ZIP/
│   └── TAR/
├── Code/
│   ├── GO/
│   ├── PY/
│   └── JS/
└── Others/
```

## Configuration

Create `~/.forter.yaml` to customize categories and extensions:

```yaml
# Default category for unknown extensions
default_category: "Others"

# Skip hidden files and directories
skip_hidden: true
skip_dirs: true

# Custom categories
categories:
  - name: "Documents"
    extensions: ["pdf", "doc", "docx", "txt", "rtf", "odt", "xls", "xlsx"]
    description: "Document files"
  
  - name: "Images"
    extensions: ["jpg", "jpeg", "png", "gif", "bmp", "svg", "webp"]
    description: "Image files"
  
  - name: "Code"
    extensions: ["go", "py", "js", "ts", "rs", "rb"]
    description: "Source code files"

# Custom extension mappings (override defaults)
custom_mappings:
  "custom": "Documents"
  "bak": "Archives"
```

## Building

```bash
# Build binary
make build

# Build with version info
make build VERSION=1.0.0

# Run tests
make test

# Run linting
make lint

# Clean build artifacts
make clean
```

## Development

### Project Structure

```
organize/
├── cmd/
│   └── organize/
│       └── main.go          # CLI entry point
├── internal/
│   ├── config/              # Configuration management
│   │   └── config.go
│   ├── organizer/           # File organization logic
│   │   └── organizer.go
│   ├── scanner/             # Directory scanning
│   │   └── scanner.go
│   └── tui/                 # Bubble Tea TUI components
│       ├── handlers.go
│       ├── model.go
│       ├── styles.go
│       └── view.go
├── go.mod
├── go.sum
├── Makefile
├── LICENSE
└── README.md
```

### Tech Stack

- **Language**: Go 1.21+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea) (Charmbracelet)
- **CLI Parsing**: [Cobra](https://github.com/spf13/cobra)
- **Configuration**: [Viper](https://github.com/spf13/viper)

### Running Locally

```bash
# Install dependencies
go mod download

# Run directly
go run cmd/organize/main.go ~/Downloads

# Run with dry run
go run cmd/organize/main.go --dry-run ~/Downloads

# Run tests
go test ./...
```

## Performance

Benchmarks on M1 MacBook Pro:

| Files | Scan Time | Memory Usage |
|-------|-----------|--------------|
| 100   | 10ms      | 5MB          |
| 1,000 | 50ms      | 15MB         |
| 10,000| 500ms     | 50MB         |

## Roadmap

- [x] Basic TUI with file navigation
- [x] Category-based organization
- [x] Dry run mode
- [x] Undo functionality
- [x] Configuration file support
- [ ] Watch mode (auto-organize Downloads folder)
- [ ] Plugin system for custom file processors
- [ ] Integration with cloud storage
- [ ] GUI version (Fyne or Wails)

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

Please ensure:
- Code follows Go best practices
- Tests are added for new features
- Documentation is updated

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Charmbracelet](https://charm.sh/) for the amazing Bubble Tea TUI framework
- [Cobra](https://github.com/spf13/cobra) for CLI structure
- Inspired by [lazygit](https://github.com/jesseduffield/lazygit) and [nnn](https://github.com/jarun/nnn)

## Support

- 💖 Star this repository if you find it useful
- 🐛 [Open an issue](https://github.com/OmkarVijayBagade/forter/issues) for bugs
- 💡 [Discussions](https://github.com/OmkarVijayBagade/forter/discussions) for questions and ideas
