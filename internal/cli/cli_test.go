package cli

import (
	"testing"

	"zipprine/internal/models"
)

func TestParseArchiveType(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected models.ArchiveType
	}{
		{"ZIP lowercase", "zip", models.ZIP},
		{"ZIP uppercase", "ZIP", models.ZIP},
		{"TAR", "tar", models.TAR},
		{"TAR.GZ", "tar.gz", models.TARGZ},
		{"TARGZ", "targz", models.TARGZ},
		{"TGZ", "tgz", models.TARGZ},
		{"GZIP", "gzip", models.GZIP},
		{"GZ", "gz", models.GZIP},
		{"RAR", "rar", models.RAR},
		{"RAR uppercase", "RAR", models.RAR},
		{"AUTO", "auto", models.AUTO},
		{"Unknown defaults to ZIP", "unknown", models.ZIP},
		{"Empty defaults to ZIP", "", models.ZIP},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseArchiveType(tt.input)
			if result != tt.expected {
				t.Errorf("parseArchiveType(%q) = %v; want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestParseArchiveTypeCaseInsensitive(t *testing.T) {
	// Test that parsing is case-insensitive
	inputs := []string{"zip", "ZIP", "Zip", "ZiP"}
	for _, input := range inputs {
		result := parseArchiveType(input)
		if result != models.ZIP {
			t.Errorf("parseArchiveType(%q) = %v; want %v", input, result, models.ZIP)
		}
	}
}

func TestParseArchiveTypeAllFormats(t *testing.T) {
	// Ensure all supported formats are handled
	formats := map[string]models.ArchiveType{
		"zip":    models.ZIP,
		"tar":    models.TAR,
		"tar.gz": models.TARGZ,
		"gzip":   models.GZIP,
		"rar":    models.RAR,
		"auto":   models.AUTO,
	}

	for input, expected := range formats {
		result := parseArchiveType(input)
		if result != expected {
			t.Errorf("parseArchiveType(%q) = %v; want %v", input, result, expected)
		}
	}
}

func TestParseArchiveTypeTarGzVariants(t *testing.T) {
	// Test all variants of TAR.GZ
	variants := []string{"tar.gz", "targz", "tgz"}
	for _, variant := range variants {
		result := parseArchiveType(variant)
		if result != models.TARGZ {
			t.Errorf("parseArchiveType(%q) = %v; want %v", variant, result, models.TARGZ)
		}
	}
}

func TestParseArchiveTypeGzipVariants(t *testing.T) {
	// Test all variants of GZIP
	variants := []string{"gzip", "gz"}
	for _, variant := range variants {
		result := parseArchiveType(variant)
		if result != models.GZIP {
			t.Errorf("parseArchiveType(%q) = %v; want %v", variant, result, models.GZIP)
		}
	}
}
