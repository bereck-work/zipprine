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

func RunBatchCompressFlow() error {
	var sourcePaths string
	var outputDir string
	var archiveTypeStr string
	var compressionLevel string
	var parallel bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("📁 Source Paths").
				Description("Enter paths separated by commas (e.g., /path1,/path2,/path3)").
				Value(&sourcePaths).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("source paths cannot be empty")
					}
					paths := strings.Split(s, ",")
					for _, p := range paths {
						p = strings.TrimSpace(p)
						if _, err := os.Stat(p); os.IsNotExist(err) {
							return fmt.Errorf("path does not exist: %s", p)
						}
					}
					return nil
				}),

			huh.NewInput().
				Title("💾 Output Directory").
				Description("Directory to save all archives").
				Value(&outputDir).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("output directory cannot be empty")
					}
					return nil
				}).
				Suggestions(getPathCompletions("")),
		),

		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🎨 Archive Type").
				Options(
					huh.NewOption("ZIP", "ZIP"),
					huh.NewOption("TAR.GZ", "TARGZ"),
					huh.NewOption("TAR", "TAR"),
				).
				Value(&archiveTypeStr),

			huh.NewSelect[string]().
				Title("⚡ Compression Level").
				Options(
					huh.NewOption("Fast (Level 1)", "1"),
					huh.NewOption("Balanced (Level 5)", "5"),
					huh.NewOption("Best (Level 9)", "9"),
				).
				Value(&compressionLevel),

			huh.NewConfirm().
				Title("🚀 Parallel Processing").
				Description("Process archives in parallel (faster for multiple files)").
				Value(&parallel),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	// Parse paths
	paths := strings.Split(sourcePaths, ",")
	configs := make([]*models.CompressConfig, 0, len(paths))

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create configs for each path
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		basename := filepath.Base(path)
		ext := getExtension(models.ArchiveType(archiveTypeStr))
		outputPath := filepath.Join(outputDir, basename+ext)

		level := 5
		fmt.Sscanf(compressionLevel, "%d", &level)

		configs = append(configs, &models.CompressConfig{
			SourcePath:       path,
			OutputPath:       outputPath,
			ArchiveType:      models.ArchiveType(archiveTypeStr),
			CompressionLevel: level,
		})
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("📦 Batch compressing %d items...", len(configs))))
	fmt.Println()

	// Create batch config
	batchConfig := &archiver.BatchCompressConfig{
		Configs:    configs,
		Parallel:   parallel,
		MaxWorkers: 4,
		OnProgress: func(index, total int, filename string) {
			fmt.Printf("  [%d/%d] Processing: %s\n", index, total, filepath.Base(filename))
		},
		OnError: func(index int, filename string, err error) {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("  ❌ Failed: %s - %v", filepath.Base(filename), err)))
		},
		OnComplete: func(index int, filename string) {
			fmt.Println(SuccessStyle.Render(fmt.Sprintf("  ✅ Completed: %s", filepath.Base(filename))))
		},
	}

	errors := archiver.BatchCompress(batchConfig)

	// Count successes
	successCount := 0
	for _, err := range errors {
		if err == nil {
			successCount++
		}
	}

	fmt.Println()
	fmt.Println(SuccessStyle.Render(fmt.Sprintf("✨ Batch complete: %d/%d successful", successCount, len(configs))))

	return nil
}

func RunBatchExtractFlow() error {
	var archivePaths string
	var outputDir string
	var parallel bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewText().
				Title("📦 Archive Paths").
				Description("Enter archive paths separated by commas").
				Value(&archivePaths).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("archive paths cannot be empty")
					}
					paths := strings.Split(s, ",")
					for _, p := range paths {
						p = strings.TrimSpace(p)
						if _, err := os.Stat(p); os.IsNotExist(err) {
							return fmt.Errorf("archive does not exist: %s", p)
						}
					}
					return nil
				}),

			huh.NewInput().
				Title("💾 Output Directory").
				Description("Directory to extract all archives").
				Value(&outputDir).
				Suggestions(getPathCompletions("")),

			huh.NewConfirm().
				Title("🚀 Parallel Processing").
				Description("Extract archives in parallel").
				Value(&parallel),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	// Parse paths
	paths := strings.Split(archivePaths, ",")
	configs := make([]*models.ExtractConfig, 0, len(paths))

	// Ensure output directory exists
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create configs for each archive
	for _, path := range paths {
		path = strings.TrimSpace(path)
		if path == "" {
			continue
		}

		// Detect archive type
		archiveType, err := archiver.DetectArchiveType(path)
		if err != nil {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("⚠️  Skipping %s: %v", path, err)))
			continue
		}

		basename := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
		destPath := filepath.Join(outputDir, basename)

		configs = append(configs, &models.ExtractConfig{
			ArchivePath:   path,
			DestPath:      destPath,
			ArchiveType:   archiveType,
			OverwriteAll:  true,
			PreservePerms: true,
		})
	}

	fmt.Println(InfoStyle.Render(fmt.Sprintf("📂 Batch extracting %d archives...", len(configs))))
	fmt.Println()

	// Create batch config
	batchConfig := &archiver.BatchExtractConfig{
		Configs:    configs,
		Parallel:   parallel,
		MaxWorkers: 4,
		OnProgress: func(index, total int, filename string) {
			fmt.Printf("  [%d/%d] Extracting: %s\n", index, total, filepath.Base(filename))
		},
		OnError: func(index int, filename string, err error) {
			fmt.Println(ErrorStyle.Render(fmt.Sprintf("  ❌ Failed: %s - %v", filepath.Base(filename), err)))
		},
		OnComplete: func(index int, filename string) {
			fmt.Println(SuccessStyle.Render(fmt.Sprintf("  ✅ Completed: %s", filepath.Base(filename))))
		},
	}

	errors := archiver.BatchExtract(batchConfig)

	// Count successes
	successCount := 0
	for _, err := range errors {
		if err == nil {
			successCount++
		}
	}

	fmt.Println()
	fmt.Println(SuccessStyle.Render(fmt.Sprintf("✨ Batch complete: %d/%d successful", successCount, len(configs))))

	return nil
}

func getExtension(archiveType models.ArchiveType) string {
	switch archiveType {
	case models.ZIP:
		return ".zip"
	case models.TARGZ:
		return ".tar.gz"
	case models.TAR:
		return ".tar"
	case models.GZIP:
		return ".gz"
	default:
		return ".archive"
	}
}
