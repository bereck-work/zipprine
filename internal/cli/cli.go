package cli

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"zipprine/internal/archiver"
	"zipprine/internal/fetcher"
	"zipprine/internal/models"
	"zipprine/internal/version"
)

// Run executes the CLI mode
func Run() bool {
	// Define flags
	compress := flag.String("compress", "", "Compress files/folders (source path)")
	extract := flag.String("extract", "", "Extract archive (archive path)")
	analyze := flag.String("analyze", "", "Analyze archive (archive path)")
	output := flag.String("output", "", "Output path for compression or extraction")
	archiveType := flag.String("type", "zip", "Archive type (zip, tar, tar.gz, gzip, rar)")
	level := flag.Int("level", 6, "Compression level (1=fast, 6=balanced, 9=best)")
	overwrite := flag.Bool("overwrite", false, "Overwrite existing files during extraction")
	preservePerms := flag.Bool("preserve-perms", true, "Preserve file permissions during extraction")
	exclude := flag.String("exclude", "", "Comma-separated list of patterns to exclude")
	include := flag.String("include", "", "Comma-separated list of patterns to include")
	verify := flag.Bool("verify", false, "Verify archive integrity after compression")
	remoteURL := flag.String("url", "", "Remote URL to download and extract archive from")
	showVersion := flag.Bool("version", false, "Show version information")
	help := flag.Bool("help", false, "Show help information")

	flag.Parse()

	// Show version
	if *showVersion {
		fmt.Println(version.FullVersion())
		return true
	}

	// Show help
	if *help {
		printHelp()
		return true
	}

	// Check if any CLI flags were provided
	if flag.NFlag() == 0 {
		return false // No flags, use interactive mode
	}

	// Handle remote URL fetching
	if *remoteURL != "" {
		if *output == "" {
			fmt.Println("❌ Error: --output is required when using --url")
			os.Exit(1)
		}

		if !fetcher.IsValidArchiveURL(*remoteURL) {
			fmt.Println("⚠️  Warning: URL does not appear to point to a supported archive format")
		}

		if err := fetcher.FetchAndExtract(*remoteURL, *output, *overwrite, *preservePerms); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✨ Remote archive fetched and extracted successfully!")
		return true
	}

	// Handle compression
	if *compress != "" {
		if *output == "" {
			fmt.Println("❌ Error: --output is required for compression")
			os.Exit(1)
		}

		archType := parseArchiveType(*archiveType)
		if archType == models.RAR {
			fmt.Println("❌ Error: RAR compression is not supported (proprietary format)")
			os.Exit(1)
		}

		config := &models.CompressConfig{
			SourcePath:       *compress,
			OutputPath:       *output,
			ArchiveType:      archType,
			CompressionLevel: *level,
			VerifyIntegrity:  *verify,
		}

		if *exclude != "" {
			config.ExcludePaths = strings.Split(*exclude, ",")
		}
		if *include != "" {
			config.IncludePaths = strings.Split(*include, ",")
		}

		fmt.Printf("📦 Compressing %s to %s (%s)...\n", *compress, *output, archType)
		if err := archiver.Compress(config); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✨ Compression completed successfully!")
		return true
	}

	// Handle extraction
	if *extract != "" {
		if *output == "" {
			fmt.Println("❌ Error: --output is required for extraction")
			os.Exit(1)
		}

		// Detect archive type if not specified or set to auto
		archType := parseArchiveType(*archiveType)
		if archType == models.AUTO || *archiveType == "" {
			detectedType, err := archiver.DetectArchiveType(*extract)
			if err != nil {
				fmt.Printf("❌ Error detecting archive type: %v\n", err)
				os.Exit(1)
			}
			archType = detectedType
			fmt.Printf("🔍 Detected archive type: %s\n", archType)
		}

		config := &models.ExtractConfig{
			ArchivePath:   *extract,
			DestPath:      *output,
			ArchiveType:   archType,
			OverwriteAll:  *overwrite,
			PreservePerms: *preservePerms,
		}

		fmt.Printf("📂 Extracting %s to %s...\n", *extract, *output)
		if err := archiver.Extract(config); err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✨ Extraction completed successfully!")
		return true
	}

	// Handle analysis
	if *analyze != "" {
		info, err := archiver.Analyze(*analyze)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("\n📊 Archive Analysis")
		fmt.Println("==================")
		fmt.Printf("Type:              %s\n", info.Type)
		fmt.Printf("File Count:        %d\n", info.FileCount)
		fmt.Printf("Total Size:        %d bytes\n", info.TotalSize)
		fmt.Printf("Compressed Size:   %d bytes\n", info.CompressedSize)
		if info.CompressionRatio > 0 {
			fmt.Printf("Compression Ratio: %.2f%%\n", info.CompressionRatio*100)
		}
		if info.Checksum != "" {
			fmt.Printf("Checksum (SHA256): %s\n", info.Checksum)
		}
		fmt.Println("\n📁 Files:")
		for i, file := range info.Files {
			if i >= 20 {
				fmt.Printf("... and %d more files\n", len(info.Files)-20)
				break
			}
			fmt.Printf("  - %s (%d bytes)\n", file.Name, file.Size)
		}
		return true
	}

	// If we get here, no valid operation was specified
	fmt.Println("❌ Error: No valid operation specified. Use --help for usage information.")
	os.Exit(1)
	return true
}

func parseArchiveType(typeStr string) models.ArchiveType {
	switch strings.ToLower(typeStr) {
	case "zip":
		return models.ZIP
	case "tar":
		return models.TAR
	case "tar.gz", "targz", "tgz":
		return models.TARGZ
	case "gzip", "gz":
		return models.GZIP
	case "rar":
		return models.RAR
	case "auto":
		return models.AUTO
	default:
		return models.ZIP
	}
}

func printHelp() {
	fmt.Println(version.FullVersion())
	fmt.Println("\n🗜️  A modern TUI/CLI archiving tool with support for multiple formats")
	fmt.Println("\nUSAGE:")
	fmt.Println("  Interactive mode (default):")
	fmt.Println("    zipprine")
	fmt.Println("\n  Command-line mode:")
	fmt.Println("    zipprine [OPTIONS]")
	fmt.Println("\nOPTIONS:")
	fmt.Println("  --compress <path>       Compress files/folders at the specified path")
	fmt.Println("  --extract <path>        Extract archive at the specified path")
	fmt.Println("  --analyze <path>        Analyze archive at the specified path")
	fmt.Println("  --output <path>         Output path for compression or extraction")
	fmt.Println("  --type <type>           Archive type: zip, tar, tar.gz, gzip, rar (default: zip)")
	fmt.Println("  --level <1-9>           Compression level: 1=fast, 6=balanced, 9=best (default: 6)")
	fmt.Println("  --overwrite             Overwrite existing files during extraction")
	fmt.Println("  --preserve-perms        Preserve file permissions (default: true)")
	fmt.Println("  --exclude <patterns>    Comma-separated patterns to exclude")
	fmt.Println("  --include <patterns>    Comma-separated patterns to include")
	fmt.Println("  --verify                Verify archive integrity after compression")
	fmt.Println("  --url <url>             Download and extract archive from remote URL")
	fmt.Println("  --version               Show version information")
	fmt.Println("  --help                  Show this help message")
	fmt.Println("\nEXAMPLES:")
	fmt.Println("  # Compress a directory")
	fmt.Println("  zipprine --compress /path/to/source --output archive.zip --type zip")
	fmt.Println("\n  # Extract an archive")
	fmt.Println("  zipprine --extract archive.tar.gz --output /path/to/dest")
	fmt.Println("\n  # Analyze an archive")
	fmt.Println("  zipprine --analyze archive.zip")
	fmt.Println("\n  # Download and extract from URL")
	fmt.Println("  zipprine --url https://example.com/archive.zip --output /path/to/dest")
	fmt.Println("\n  # Compress with exclusions")
	fmt.Println("  zipprine --compress /project --output project.tar.gz --type tar.gz --exclude '*.log,*.tmp'")
	fmt.Println("\nSUPPORTED FORMATS:")
	fmt.Println("  Compression: ZIP, TAR, TAR.GZ, GZIP")
	fmt.Println("  Extraction:  ZIP, TAR, TAR.GZ, GZIP, RAR")
	fmt.Println("\nNOTE:")
	fmt.Println("  RAR compression is not supported due to proprietary format.")
	fmt.Println("  RAR extraction is supported for reading existing RAR archives.")
}
