package ui

import (
	"fmt"

	"zipprine/internal/fetcher"

	"github.com/charmbracelet/huh"
)

func RunRemoteFetchFlow() error {
	var url, destPath string
	var overwrite, preservePerms bool

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("🌐 Remote Archive URL").
				Description("HTTP/HTTPS URL to download archive from").
				Placeholder("https://example.com/archive.zip").
				Value(&url).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("URL cannot be empty")
					}
					if !fetcher.IsValidArchiveURL(s) {
						return fmt.Errorf("URL does not appear to point to a supported archive format")
					}
					return nil
				}),

			huh.NewInput().
				Title("📂 Destination Path").
				Description("Where to extract the downloaded archive - Tab for completions").
				Placeholder("/path/to/destination").
				Value(&destPath).
				Validate(func(s string) error {
					if s == "" {
						return fmt.Errorf("destination path cannot be empty")
					}
					return nil
				}).
				Suggestions(getDirCompletions("")),
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

	fmt.Println()
	fmt.Println(InfoStyle.Render("🌐 Fetching remote archive..."))

	if err := fetcher.FetchAndExtract(url, destPath, overwrite, preservePerms); err != nil {
		return fmt.Errorf("failed to fetch and extract: %w", err)
	}

	return nil
}
