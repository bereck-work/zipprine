# 🗜️ Zipprine - Advanced Archive Manager

A beautiful, feature-rich TUI application for managing archives with style!

## ✨ Features

### 📦 Compression

- **Multiple formats**: ZIP, TAR, TAR.GZ, GZIP
- **Compression levels**: Fast, Balanced, Best
- **Smart filtering**: Include/exclude patterns with wildcards
- **Integrity verification**: SHA256 checksums and validation

### 📂 Extraction

- **Auto-detection**: Automatically detects archive type by magic bytes
- **Safe extraction**: Optional overwrite protection
- **Permission preservation**: Keep original file permissions
- **Progress tracking**: Real-time extraction feedback

### 🔍 Analysis

- **Detailed statistics**: File count, sizes, compression ratios
- **File listing**: View contents without extraction
- **Checksum verification**: SHA256 integrity checks
- **Format detection**: Magic byte analysis

## 🚀 Installation

```bash
# Clone the repository
git clone https://github.com/bereck-work/ziprine.git
cd ziprine

# Install dependencies
go mod download

# Build
go build -o ziprine ./cmd/ziprine

# Run
./ziprine
```

## 📖 Usage

Simply run `ziprine` and follow the interactive prompts!

### Compress Files

```bash
./ziprine
# Select: Compress files/folders
# Enter source path: /path/to/folder
# Choose format: ZIP, TAR.GZ, TAR, or GZIP
# Set compression level and filters
```

### Extract Archives

```bash
./ziprine
# Select: Extract archive
# Archive type is auto-detected!
# Choose destination and options
```

### Analyze Archives

```bash
./ziprine
# Select: Analyze archive
# View detailed statistics and file listing
```

## 🎨 Pattern Examples

**Exclude patterns**:

- `*.log` - Exclude all log files
- `node_modules` - Exclude node_modules directory
- `temp/*` - Exclude everything in temp folder
- `.git,__pycache__,*.tmp` - Multiple patterns

**Include patterns**:

- `*.go` - Only Go files
- `src/*,docs/*` - Only src and docs folders
- `*.md,*.txt` - Only markdown and text files

## 🏗️ Project Structure

```
ziprine/
├── cmd/ziprine/          # Main application entry
├── internal/
│   ├── archiver/         # Archive operations
│   ├── ui/              # TUI components
│   └── models/          # Data structures
└── pkg/fileutil/        # Utility functions
```

## 🛠️ Technologies

- **[Charm Bracelet Huh](https://github.com/charmbracelet/huh)** - Beautiful TUI forms
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling and colors
- **Go standard library** - Archive formats

## 📝 License

MIT License - Feel free to use and modify!

## 🤝 Contributing

Contributions are welcome! Feel free to open issues or submit PRs.
