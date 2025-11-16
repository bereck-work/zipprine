package main

import (
	"fmt"
	"os"

	"zipprine/internal/ui"

	"github.com/charmbracelet/huh"
)

func main() {
	fmt.Println(ui.TitleStyle.Render("🗜️  Zipprine - Archive Like a Pro! 🚀"))
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
	case "exit":
		fmt.Println(ui.InfoStyle.Render("👋 Goodbye!"))
		return
	}

	fmt.Println()
	fmt.Println(ui.SuccessStyle.Render("✨ Operation completed successfully!"))
}