package ui

import (
	"fmt"
	"os"

	"zipprine/internal/archiver"
	"zipprine/internal/models"

	"github.com/charmbracelet/huh"
)

func RunCompareFlow() error {
	var archive1Path, archive2Path string
	var showDetails bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📦 First Archive").
				Description("Path to the first archive").
				Value(&archive1Path).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path cannot be empty")
					}
					if _, err := os.Stat(s); os.IsNotExist(err) {
						return fmt.Errorf("archive does not exist")
					}
					return nil
				}).
				Suggestions(getPathCompletions("")),

			huh.NewInput().
				Title("📦 Second Archive").
				Description("Path to the second archive").
				Value(&archive2Path).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path cannot be empty")
					}
					if _, err := os.Stat(s); os.IsNotExist(err) {
						return fmt.Errorf("archive does not exist")
					}
					return nil
				}).
				Suggestions(getPathCompletions("")),

			huh.NewConfirm().
				Title("📋 Show Detailed Differences").
				Description("Display detailed file-by-file comparison").
				Value(&showDetails),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	fmt.Println(InfoStyle.Render("🔍 Analyzing archives..."))
	fmt.Println()

	// Detect archive types
	type1, err := archiver.DetectArchiveType(archive1Path)
	if err != nil {
		return fmt.Errorf("failed to detect first archive type: %w", err)
	}

	type2, err := archiver.DetectArchiveType(archive2Path)
	if err != nil {
		return fmt.Errorf("failed to detect second archive type: %w", err)
	}

	// Compare archives
	result, err := archiver.CompareArchives(archive1Path, archive2Path, type1, type2)
	if err != nil {
		return fmt.Errorf("failed to compare archives: %w", err)
	}

	// Display results
	fmt.Println(TitleStyle.Render("📊 Comparison Results"))
	fmt.Println()
	fmt.Println(result.Summary)
	fmt.Println()

	if showDetails {
		if len(result.OnlyInFirst) > 0 {
			fmt.Println(InfoStyle.Render("📁 Files only in first archive:"))
			for _, f := range result.OnlyInFirst {
				fmt.Printf("  • %s\n", f)
			}
			fmt.Println()
		}

		if len(result.OnlyInSecond) > 0 {
			fmt.Println(InfoStyle.Render("📁 Files only in second archive:"))
			for _, f := range result.OnlyInSecond {
				fmt.Printf("  • %s\n", f)
			}
			fmt.Println()
		}

		if len(result.Different) > 0 {
			fmt.Println(WarningStyle.Render("⚠️  Files that differ:"))
			for _, f := range result.Different {
				fmt.Printf("  • %s\n", f.Name)
				fmt.Printf("    Size: %d bytes → %d bytes\n", f.Size1, f.Size2)
				fmt.Printf("    ModTime: %s → %s\n", f.ModTime1, f.ModTime2)
			}
			fmt.Println()
		}

		if len(result.InBoth) > 0 && len(result.Different) == 0 {
			fmt.Println(SuccessStyle.Render("✅ All common files are identical!"))
		}
	}

	return nil
}

func RunConvertFlow() error {
	var sourcePath, destPath string
	var destTypeStr string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📦 Source Archive").
				Description("Path to the archive to convert").
				Value(&sourcePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path cannot be empty")
					}
					if _, err := os.Stat(s); os.IsNotExist(err) {
						return fmt.Errorf("archive does not exist")
					}
					return nil
				}).
				Suggestions(getPathCompletions("")),

			huh.NewInput().
				Title("💾 Destination Path").
				Description("Path for the converted archive").
				Value(&destPath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("path cannot be empty")
					}
					return nil
				}).
				Suggestions(getPathCompletions("")),

			huh.NewSelect[string]().
				Title("🎨 Destination Format").
				Options(
					huh.NewOption("ZIP", "ZIP"),
					huh.NewOption("TAR.GZ", "TARGZ"),
					huh.NewOption("TAR", "TAR"),
				).
				Value(&destTypeStr),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	fmt.Println(InfoStyle.Render("🔄 Converting archive..."))
	fmt.Println()

	// Detect source archive type
	sourceType, err := archiver.DetectArchiveType(sourcePath)
	if err != nil {
		return fmt.Errorf("failed to detect source archive type: %w", err)
	}

	destType := models.ArchiveType(destTypeStr)

	// Convert archive
	if err := archiver.ConvertArchive(sourcePath, destPath, sourceType, destType); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	fmt.Println(SuccessStyle.Render(fmt.Sprintf("✅ Successfully converted %s to %s", sourceType, destType)))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("📁 Output: %s", destPath)))

	return nil
}
