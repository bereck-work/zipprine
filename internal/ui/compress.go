package ui

import (
	"fmt"
	"os"
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

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📁 Source Path").
				Description("Enter the path to compress (file or directory)").
				Placeholder("/path/to/source").
				Value(&sourcePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("source path cannot be empty")
					}
					if _, err := os.Stat(s); os.IsNotExist(err) {
						return fmt.Errorf("path does not exist")
					}
					return nil
				}),

			huh.NewInput().
				Title("💾 Output Path").
				Description("Where to save the archive").
				Placeholder("/path/to/output.zip").
				Value(&outputPath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("output path cannot be empty")
					}
					return nil
				}).Suggestions([]string{".zip", ".tar.gz", ".tar", ".gz"}),
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

	if err := archiver.Compress(config); err != nil {
		return err
	}

	fmt.Println(SuccessStyle.Render("✅ Archive created successfully!"))

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