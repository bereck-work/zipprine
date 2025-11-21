# 🗜️ Zipprine - TUI zipping tool

Zipprine is a modern TUI application for managing archives with support for multiple formats including ZIP, TAR, TAR.GZ, and GZIP.

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

### 📚 Batch Operations

- **Batch compression**: Compress multiple files/folders at once
- **Batch extraction**: Extract multiple archives simultaneously
- **Parallel processing**: Speed up operations with concurrent workers
- **Progress tracking**: Real-time feedback for each operation

### ⚖️ Archive Comparison

- **Compare two archives**: Find differences between archives
- **Detailed reports**: See files unique to each archive
- **File differences**: Identify files that differ in size or modification time
- **Cross-format support**: Compare different archive formats

### 🔄 Archive Conversion

- **Format conversion**: Convert between ZIP, TAR, TAR.GZ formats
- **Preserve contents**: Maintains file structure and permissions
- **Automatic extraction**: Seamless conversion process

## 🚀 Installation

```bash
# Clone the repository
git clone https://gitlab.com/bereckobrian/zipprine.git
cd zipprine

# Build for your platform
make build

# Run it
./build/zipprine
```

## 📖 Usage

Just run `zipprine` and follow the interactive menu:

**Compress** - Choose files/folders, pick a format (ZIP, TAR, TAR.GZ, GZIP), and set your preferences

**Extract** - Point to an archive and choose where to extract (format is auto-detected)

**Analyze** - View detailed stats about any archive without extracting it

**Batch Operations** - Compress or extract multiple files at once with optional parallel processing

**Compare** - Find differences between two archives

**Convert** - Change archive formats while preserving structure

## 🔨 Building

```bash
# Build for your platform
make build

# Build for all platforms (Linux, macOS, Windows)
make build-all

# Install to $GOPATH/bin
make install

# Or use Docker
docker build -t zipprine .
```

## 🧪 Testing

We have comprehensive test coverage:

```bash
# Run tests
make test

# With coverage report
make test-coverage

# Run benchmarks
make bench
```

**Coverage:** 77.9% for archiver, 96.6% for utilities

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

## 🛠️ Technologies

- **[Charm Bracelet Huh](https://github.com/charmbracelet/huh)** - Beautiful TUI forms
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling and colors
- **Go standard library** - Archive formats

## 📝 License

MIT License
