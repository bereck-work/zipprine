package ui

import (
	"os"
	"path/filepath"
	"strings"
)

// getPathCompletions returns file/directory path completions for autocomplete
func getPathCompletions(input string) []string {
	if input == "" {
		input = "."
	}

	if strings.HasPrefix(input, "~") {
		home, err := os.UserHomeDir()
		if err == nil {
			input = filepath.Join(home, input[1:])
		}
	}

	dir := filepath.Dir(input)
	pattern := filepath.Base(input)

	if strings.HasSuffix(input, string(filepath.Separator)) {
		dir = input
		pattern = ""
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
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

	if len(completions) > 15 {
		completions = completions[:15]
	}

	return completions
}

// getArchiveCompletions returns only archive file completions
func getArchiveCompletions(input string) []string {
	archiveExts := map[string]bool{
		".zip":    true,
		".tar":    true,
		".gz":     true,
		".tar.gz": true,
		".tgz":    true,
	}

	allCompletions := getPathCompletions(input)
	archiveCompletions := []string{}

	for _, path := range allCompletions {
		if strings.HasSuffix(path, string(filepath.Separator)) {
			archiveCompletions = append(archiveCompletions, path)
			continue
		}

		ext := filepath.Ext(path)
		if archiveExts[ext] {
			archiveCompletions = append(archiveCompletions, path)
			continue
		}

		// Check for .tar.gz
		if strings.HasSuffix(path, ".tar.gz") {
			archiveCompletions = append(archiveCompletions, path)
		}
	}

	return archiveCompletions
}

// getDirCompletions returns only directory completions
func getDirCompletions(input string) []string {
	allCompletions := getPathCompletions(input)
	dirCompletions := []string{}

	for _, path := range allCompletions {
		if strings.HasSuffix(path, string(filepath.Separator)) {
			dirCompletions = append(dirCompletions, path)
		}
	}

	return dirCompletions
}
