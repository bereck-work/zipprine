package fileutil

import (
	"fmt"
	"path/filepath"
	"strings"
)

// ShouldInclude determines if a file should be included based on exclude/include patterns
func ShouldInclude(path string, excludePaths, includePaths []string) bool {
	// Check exclude patterns first
	for _, pattern := range excludePaths {
		if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
			return false
		}
		// Also check if pattern is in the path
		if strings.Contains(path, pattern) {
			return false
		}
		// Handle directory patterns
		if strings.HasSuffix(pattern, "/*") {
			dirPattern := strings.TrimSuffix(pattern, "/*")
			if strings.Contains(path, dirPattern) {
				return false
			}
		}
	}

	// If include patterns are specified, check them
	if len(includePaths) > 0 {
		for _, pattern := range includePaths {
			if matched, _ := filepath.Match(pattern, filepath.Base(path)); matched {
				return true
			}
			if strings.Contains(path, pattern) {
				return true
			}
			// Handle directory patterns
			if strings.HasSuffix(pattern, "/*") {
				dirPattern := strings.TrimSuffix(pattern, "/*")
				if strings.Contains(path, dirPattern) {
					return true
				}
			}
		}
		return false
	}

	return true
}

// FormatBytes converts bytes to human-readable format
func FormatBytes(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}