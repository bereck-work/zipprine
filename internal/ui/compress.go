package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"zipprine/internal/archiver"
	"zipprine/internal/models"

	"github.com/charmbracelet/huh"
)

func RunCompressFlow() error {
	config := &models.CompressConfig{}

	var sourcePath, outputPath string
	var archiveTypeStr string
	var excludeInput, includeInput string
	var verify bool
	var compressionLevel string

	// Get current working directory
	cwd, _ := os.Getwd()

	// Function to get directory completions
	getDirCompletions := func(input string) []string {
		if input == "" {
			input = "."
		}

		// Expand home directory
		if strings.HasPrefix(input, "~") {
			home, err := os.UserHomeDir()
			if err == nil {
				input = filepath.Join(home, input[1:])
			}
		}

		// Get the directory and file pattern
		dir := filepath.Dir(input)
		pattern := filepath.Base(input)

		// If input ends with /, we want to list that directory
		if strings.HasSuffix(input, string(filepath.Separator)) {
			dir = input
			pattern = ""
		}

		// Read directory
		entries, err := os.ReadDir(dir)
		if err != nil {
			// If can't read, try parent directory
			entries, err = os.ReadDir(".")
			if err != nil {
				return []string{}
			}
			dir = "."
		}

		completions := []string{}
		for _, entry := range entries {
			name := entry.Name()
			
			// Skip hidden files unless explicitly requested
			if strings.HasPrefix(name, ".") && !strings.HasPrefix(pattern, ".") {
				continue
			}

			// Filter by pattern
			if pattern != "" && !strings.HasPrefix(strings.ToLower(name), strings.ToLower(pattern)) {
				continue
			}

			fullPath := filepath.Join(dir, name)
			if entry.IsDir() {
				fullPath += string(filepath.Separator)
			}

			completions = append(completions, fullPath)
		}

		// Limit to 10 suggestions
		if len(completions) > 10 {
			completions = completions[:10]
		}

		return completions
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📁 Source Path").
				Description("Enter the path to compress (file or directory) - Tab for completions").
				Placeholder(cwd).
				Value(&sourcePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("source path cannot be empty")
					}
					// Expand home directory
					if strings.HasPrefix(s, "~") {
						home, err := os.UserHomeDir()
						if err == nil {
							s = filepath.Join(home, s[1:])
						}
					}
					if _, err := os.Stat(s); os.IsNotExist(err) {
						return fmt.Errorf("path does not exist")
					}
					return nil
				}).
				Suggestions(getDirCompletions("")),

			huh.NewInput().
				Title("💾 Output Path").
				Description("Where to save (leave empty for auto-naming in current directory)").
				Placeholder("Auto: <source-name>.<type> in current directory").
				Value(&outputPath).
				Suggestions(getDirCompletions("")),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🎨 Archive Type").
				Description("Choose your compression format").
				Options(
					huh.NewOption("ZIP - Universal & Compatible 📦", "ZIP"),
					huh.NewOption("TAR.GZ - Linux Classic (Best Compression) 🐧", "TARGZ"),
					huh.NewOption("TAR - No Compression 📄", "TAR"),
					huh.NewOption("GZIP - Single File Compression 🔧", "GZIP"),
				).
				Value(&archiveTypeStr),

			huh.NewSelect[string]().
				Title("⚡ Compression Level").
				Description("Higher = smaller but slower").
				Options(
					huh.NewOption("Fast (Level 1)", "1"),
					huh.NewOption("Balanced (Level 5)", "5"),
					huh.NewOption("Best (Level 9)", "9"),
				).
				Value(&compressionLevel),
		),

		huh.NewGroup(
			huh.NewText().
				Title("🚫 Exclude Patterns").
				Description("Comma-separated patterns to exclude (e.g., *.log,node_modules,*.tmp)").
				Placeholder("*.log,temp/*,.git,__pycache__").
				Value(&excludeInput),

			huh.NewText().
				Title("✅ Include Patterns").
				Description("Comma-separated patterns to include (leave empty for all)").
				Placeholder("*.go,*.md,src/*").
				Value(&includeInput),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("🔐 Verify Archive Integrity").
				Description("Check the archive after creation?").
				Value(&verify).
				Affirmative("Yes please!").
				Negative("Skip it"),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	// Expand home directory in source path
	if strings.HasPrefix(sourcePath, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			sourcePath = filepath.Join(home, sourcePath[1:])
		}
	}

	// Make source path absolute
	if !filepath.IsAbs(sourcePath) {
		absPath, err := filepath.Abs(sourcePath)
		if err == nil {
			sourcePath = absPath
		}
	}

	// Auto-generate output path if not provided
	if outputPath == "" {
		sourceName := filepath.Base(sourcePath)
		
		// Remove trailing slashes
		sourceName = strings.TrimSuffix(sourceName, string(filepath.Separator))
		
		// Determine file extension based on archive type
		var extension string
		switch models.ArchiveType(archiveTypeStr) {
		case models.ZIP:
			extension = ".zip"
		case models.TARGZ:
			extension = ".tar.gz"
		case models.TAR:
			extension = ".tar"
		case models.GZIP:
			extension = ".gz"
		default:
			extension = ".zip"
		}

		// Create output path in current working directory
		outputPath = filepath.Join(cwd, sourceName+extension)
		
		fmt.Println(InfoStyle.Render(fmt.Sprintf("📝 Auto-generated output: %s", outputPath)))
	} else {
		// Expand home directory in output path
		if strings.HasPrefix(outputPath, "~") {
			home, err := os.UserHomeDir()
			if err == nil {
				outputPath = filepath.Join(home, outputPath[1:])
			}
		}

		// Make output path absolute if relative
		if !filepath.IsAbs(outputPath) {
			absPath, err := filepath.Abs(outputPath)
			if err == nil {
				outputPath = absPath
			}
		}
	}

	config.SourcePath = sourcePath
	config.OutputPath = outputPath
	config.ArchiveType = models.ArchiveType(archiveTypeStr)
	config.VerifyIntegrity = verify
	fmt.Sscanf(compressionLevel, "%d", &config.CompressionLevel)

	if excludeInput != "" {
		config.ExcludePaths = strings.Split(excludeInput, ",")
		for i := range config.ExcludePaths {
			config.ExcludePaths[i] = strings.TrimSpace(config.ExcludePaths[i])
		}
	}

	if includeInput != "" {
		config.IncludePaths = strings.Split(includeInput, ",")
		for i := range config.IncludePaths {
			config.IncludePaths[i] = strings.TrimSpace(config.IncludePaths[i])
		}
	}

	fmt.Println()
	fmt.Println(InfoStyle.Render("🎯 Starting compression..."))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("   Source: %s", config.SourcePath)))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("   Output: %s", config.OutputPath)))

	if err := archiver.Compress(config); err != nil {
		return err
	}

	fmt.Println(SuccessStyle.Render("✅ Archive created successfully!"))
	
	// Show file info
	fileInfo, err := os.Stat(config.OutputPath)
	if err == nil {
		sizeKB := float64(fileInfo.Size()) / 1024
		sizeMB := sizeKB / 1024
		if sizeMB >= 1 {
			fmt.Println(InfoStyle.Render(fmt.Sprintf("📦 Size: %.2f MB", sizeMB)))
		} else {
			fmt.Println(InfoStyle.Render(fmt.Sprintf("📦 Size: %.2f KB", sizeKB)))
		}
	}

	if config.VerifyIntegrity {
		fmt.Println(InfoStyle.Render("🔍 Verifying archive integrity..."))
		info, err := archiver.Analyze(config.OutputPath)
		if err != nil {
			return err
		}
		displayArchiveInfo(info)
	}

	return nil
}