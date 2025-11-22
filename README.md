# 🗜️ Zipprine - TUI/CLI Archiving Tool

Zipprine is a modern TUI/CLI application for managing archives with support for multiple formats including ZIP, TAR, TAR.GZ, GZIP, and RAR (extraction only).

**Version:** 1.0.3

## ✨ Features

### 📦 Compression

- **Multiple formats**: ZIP, TAR, TAR.GZ, GZIP
- **Compression levels**: Fast, Balanced, Best
- **Smart filtering**: Include/exclude patterns with wildcards
- **Integrity verification**: SHA256 checksums and validation
- **CLI mode**: Non-interactive command-line interface for automation

### 📂 Extraction

- **Auto-detection**: Automatically detects archive type by magic bytes
- **RAR support**: Extract RAR archives (v4 and v5)
- **Remote fetching**: Download and extract archives from URLs
- **Safe extraction**: Optional overwrite protection
- **Permission preservation**: Keep original file permissions
- **Progress tracking**: Real-time extraction feedback

### 🔍 Analysis

- **Detailed statistics**: File count, sizes, compression ratios
- **File listing**: View contents without extraction
- **Checksum verification**: SHA256 integrity checks
- **Format detection**: Magic byte analysis (including RAR)

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

### 🌐 Remote Archive Fetching

- **URL download**: Fetch archives from HTTP/HTTPS URLs
- **Auto-extract**: Automatically extract downloaded archives
- **Progress tracking**: Real-time download progress
- **Format detection**: Supports all archive formats via URL

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

### Interactive Mode (TUI)

Just run `zipprine` without arguments to launch the interactive menu:

```bash
zipprine
```

**Compress** - Choose files/folders, pick a format (ZIP, TAR, TAR.GZ, GZIP), and set your preferences

**Extract** - Point to an archive and choose where to extract (format is auto-detected, supports RAR)

**Analyze** - View detailed stats about any archive without extracting it

**Batch Operations** - Compress or extract multiple files at once with optional parallel processing

**Compare** - Find differences between two archives

**Convert** - Change archive formats while preserving structure

### Command-Line Mode (CLI)

For automation and scripting, use CLI flags:

```bash
# Compress a directory
zipprine --compress /path/to/source --output archive.zip --type zip

# Extract an archive (auto-detects format)
zipprine --extract archive.tar.gz --output /path/to/dest

# Extract a RAR archive
zipprine --extract archive.rar --output /path/to/dest

# Analyze an archive
zipprine --analyze archive.zip

# Download and extract from URL
zipprine --url https://example.com/archive.zip --output /path/to/dest

# Compress with exclusions
zipprine --compress /project --output project.tar.gz --type tar.gz --exclude '*.log,*.tmp'

# Show version
zipprine --version

# Show help
zipprine --help
```

#### CLI Options

- `--compress <path>` - Compress files/folders at the specified path
- `--extract <path>` - Extract archive at the specified path
- `--analyze <path>` - Analyze archive at the specified path
- `--output <path>` - Output path for compression or extraction
- `--type <type>` - Archive type: zip, tar, tar.gz, gzip, rar (default: zip)
- `--level <1-9>` - Compression level: 1=fast, 6=balanced, 9=best (default: 6)
- `--overwrite` - Overwrite existing files during extraction
- `--preserve-perms` - Preserve file permissions (default: true)
- `--exclude <patterns>` - Comma-separated patterns to exclude
- `--include <patterns>` - Comma-separated patterns to include
- `--verify` - Verify archive integrity after compression
- `--url <url>` - Download and extract archive from remote URL
- `--version` - Show version information
- `--help` - Show help message

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

## 📚 Supported Formats

### Compression (Create Archives)

- **ZIP** - Universal format, works everywhere
- **TAR** - Unix standard, no compression
- **TAR.GZ** - Compressed TAR, best for Linux
- **GZIP** - Single file compression

### Extraction (Read Archives)

- **ZIP** - Full support
- **TAR** - Full support
- **TAR.GZ** - Full support
- **GZIP** - Full support
- **RAR** - Extraction only (RAR v4 and v5)

**Note:** RAR compression is not supported due to proprietary format restrictions. Use ZIP or TAR.GZ for creating archives.

## 🛠️ Technologies

- **[Charm Bracelet Huh](https://github.com/charmbracelet/huh)** - Beautiful TUI forms
- **[Lipgloss](https://github.com/charmbracelet/lipgloss)** - Styling and colors
- **[rardecode](https://github.com/nwaples/rardecode)** - RAR extraction support
- **Go standard library** - Archive formats and HTTP client

## 📝 License

MIT License
