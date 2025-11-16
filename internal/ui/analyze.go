package ui

import (
	"fmt"
	"os"

	"zipprine/internal/archiver"
	"zipprine/internal/models"

	"github.com/charmbracelet/huh"
)

func RunAnalyzeFlow() error {
	var archivePath string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📦 Archive Path").
				Description("Path to the archive to analyze").
				Placeholder("/path/to/archive.zip").
				Value(&archivePath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("archive path cannot be empty")
					}
					if _, err := os.Stat(s); os.IsNotExist(err) {
						return fmt.Errorf("archive does not exist")
					}
					return nil
				}),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	fmt.Println()
	fmt.Println(InfoStyle.Render("🔍 Analyzing archive..."))

	info, err := archiver.Analyze(archivePath)
	if err != nil {
		return err
	}

	displayArchiveInfo(info)
	return nil
}

func displayArchiveInfo(info *models.ArchiveInfo) {
	fmt.Println()
	fmt.Println(HeaderStyle.Render("📊 Archive Information"))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("  🎨 Type: %s", info.Type)))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("  📁 Files: %d", info.FileCount)))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("  💾 Uncompressed: %.2f MB", float64(info.TotalSize)/(1024*1024))))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("  📦 Compressed: %.2f MB", float64(info.CompressedSize)/(1024*1024))))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("  🎯 Ratio: %.1f%%", info.CompressionRatio)))
	fmt.Println(InfoStyle.Render(fmt.Sprintf("  🔒 SHA256: %s...", info.Checksum[:16])))

	if len(info.Files) > 0 && len(info.Files) <= 20 {
		fmt.Println()
		fmt.Println(HeaderStyle.Render("📝 File List"))
		for _, f := range info.Files {
			icon := "📄"
			if f.IsDir {
				icon = "📁"
			}
			fmt.Println(InfoStyle.Render(fmt.Sprintf("  %s %s (%.2f KB)", icon, f.Name, float64(f.Size)/1024)))
		}
	}
}