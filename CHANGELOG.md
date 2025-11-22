# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.3] - 2025-11-22

### Added

- **RAR Support**: Added extraction support for RAR archives (v4 and v5)
  - Magic byte detection for RAR files
  - Full extraction with permission preservation
  - Analysis capabilities for RAR archives
  - Note: RAR compression not supported due to proprietary format
- **Semantic Versioning**: Implemented proper semantic versioning system
  - Version module with Major.Minor.Patch format
  - `--version` flag to display version information
  - Version displayed in interactive TUI mode
- **Command-Line Interface (CLI)**: Added non-interactive CLI mode for automation
  - `--compress` flag for compression operations
  - `--extract` flag for extraction operations
  - `--analyze` flag for archive analysis
  - `--output` flag for specifying output paths
  - `--type` flag for archive type selection
  - `--level` flag for compression level control
  - `--overwrite` flag for overwrite control
  - `--preserve-perms` flag for permission preservation
  - `--exclude` and `--include` flags for filtering
  - `--verify` flag for integrity verification
  - `--help` flag for usage information
- **Remote URL Fetching**: Added ability to download and extract archives from URLs
  - `--url` flag for remote archive fetching
  - HTTP/HTTPS support
  - Progress tracking during download
  - Automatic format detection and extraction
  - Support for all archive formats via URL

### Changed

- Updated README.md with comprehensive documentation for new features
- Enhanced main.go to support both CLI and interactive modes
- Improved archive type detection to include RAR format

### Dependencies

- Added `github.com/nwaples/rardecode` v1.1.3 for RAR extraction support

## [0.x.x] - Previous Versions

Previous versions included:

- ZIP, TAR, TAR.GZ, and GZIP support
- Interactive TUI mode
- Batch operations
- Archive comparison
- Format conversion
- Compression levels
- Include/exclude patterns
- Integrity verification
