package ui

import (
	"fmt"
	"os"

	"zipprine/internal/archiver"
	"zipprine/internal/models"

	"github.com/charmbracelet/huh"
)

func RunExtractFlow() error {
	config := &models.ExtractConfig{}

	var archivePath, destPath string
	var overwrite, preservePerms bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("📦 Archive Path").
				Description("Path to the archive file").
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

			huh.NewInput().
				Title("📂 Destination Path").
				Description("Where to extract files").
				Placeholder("/path/to/destination").
				Value(&destPath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("destination path cannot be empty")
					}
					return nil
				}),
		),

		huh.NewGroup(
			huh.NewConfirm().
				Title("⚠️  Overwrite Existing Files").
				Description("Replace files if they already exist?").
				Value(&overwrite).
				Affirmative("Yes, overwrite").
				Negative("No, skip"),

			huh.NewConfirm().
				Title("🔒 Preserve Permissions").
				Description("Keep original file permissions?").
				Value(&preservePerms).
				Affirmative("Yes").
				Negative("No"),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := form.Run(); err != nil {
		return err
	}

	config.ArchivePath = archivePath
	config.DestPath = destPath
	config.OverwriteAll = overwrite
	config.PreservePerms = preservePerms

	fmt.Println()
	fmt.Println(InfoStyle.Render("🔍 Detecting archive type..."))

	detectedType, err := archiver.DetectArchiveType(archivePath)
	if err != nil {
		return err
	}
	config.ArchiveType = detectedType

	fmt.Println(SuccessStyle.Render(fmt.Sprintf("✅ Detected: %s", detectedType)))
	fmt.Println(InfoStyle.Render("📂 Extracting files..."))

	if err := archiver.Extract(config); err != nil {
		return err
	}

	fmt.Println(SuccessStyle.Render("✅ Extraction completed!"))
	return nil
}