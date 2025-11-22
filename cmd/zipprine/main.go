package main

import (
	"fmt"
	"os"

	"zipprine/internal/cli"
	"zipprine/internal/ui"
	"zipprine/internal/version"

	"github.com/charmbracelet/huh"
)

func main() {
	// Try CLI mode first
	if cli.Run() {
		return
	}

	fmt.Println(ui.TitleStyle.Render("Zipprine - TUI Archiver"))
	fmt.Println(ui.InfoStyle.Render("Version: " + version.Version()))
	fmt.Println()

	var operation string

	mainMenu := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("🎯 What would you like to do?").
				Options(
					huh.NewOption("📦 Compress files/folders", "compress"),
					huh.NewOption("📂 Extract archive", "extract"),
					huh.NewOption("🔍 Analyze archive", "analyze"),
					huh.NewOption("🌐 Fetch from URL", "remote-fetch"),
					huh.NewOption("📚 Batch compress", "batch-compress"),
					huh.NewOption("📂 Batch extract", "batch-extract"),
					huh.NewOption("🔄 Convert archive format", "convert"),
					huh.NewOption("⚖️  Compare archives", "compare"),
					huh.NewOption("🚪 Exit", "exit"),
				).
				Value(&operation),
		),
	).WithTheme(huh.ThemeCatppuccin())

	if err := mainMenu.Run(); err != nil {
		fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
		os.Exit(1)
	}

	switch operation {
	case "compress":
		if err := ui.RunCompressFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "extract":
		if err := ui.RunExtractFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "analyze":
		if err := ui.RunAnalyzeFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "remote-fetch":
		if err := ui.RunRemoteFetchFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "batch-compress":
		if err := ui.RunBatchCompressFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "batch-extract":
		if err := ui.RunBatchExtractFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "convert":
		if err := ui.RunConvertFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "compare":
		if err := ui.RunCompareFlow(); err != nil {
			fmt.Println(ui.ErrorStyle.Render("❌ Error: " + err.Error()))
			os.Exit(1)
		}
	case "exit":
		fmt.Println(ui.InfoStyle.Render("👋 Goodbye!"))
		return
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("✨ Operation completed successfully!"))
}