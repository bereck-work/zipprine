package fileutil

import (
	"testing"
)

func TestShouldInclude(t *testing.T) {
	tests := []struct {
		name         string
		path         string
		excludePaths []string
		includePaths []string
		expected     bool
	}{
		{
			name:         "no filters - should include",
			path:         "/path/to/file.txt",
			excludePaths: []string{},
			includePaths: []string{},
			expected:     true,
		},
		{
			name:         "exclude pattern match",
			path:         "/path/to/file.log",
			excludePaths: []string{"*.log"},
			includePaths: []string{},
			expected:     false,
		},
		{
			name:         "exclude directory pattern",
			path:         "/path/node_modules/file.js",
			excludePaths: []string{"node_modules"},
			includePaths: []string{},
			expected:     false,
		},
		{
			name:         "exclude with wildcard directory",
			path:         "/path/temp/file.txt",
			excludePaths: []string{"temp/*"},
			includePaths: []string{},
			expected:     false,
		},
		{
			name:         "include pattern match",
			path:         "/path/to/file.go",
			excludePaths: []string{},
			includePaths: []string{"*.go"},
			expected:     true,
		},
		{
			name:         "include pattern no match",
			path:         "/path/to/file.txt",
			excludePaths: []string{},
			includePaths: []string{"*.go"},
			expected:     false,
		},
		{
			name:         "exclude takes precedence",
			path:         "/path/to/file.log",
			excludePaths: []string{"*.log"},
			includePaths: []string{"*.log"},
			expected:     false,
		},
		{
			name:         "multiple exclude patterns",
			path:         "/path/to/file.tmp",
			excludePaths: []string{"*.log", "*.tmp", "*.cache"},
			includePaths: []string{},
			expected:     false,
		},
		{
			name:         "include directory pattern",
			path:         "/path/src/main.go",
			excludePaths: []string{},
			includePaths: []string{"src/*"},
			expected:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ShouldInclude(tt.path, tt.excludePaths, tt.includePaths)
			if result != tt.expected {
				t.Errorf("ShouldInclude(%q, %v, %v) = %v; want %v",
					tt.path, tt.excludePaths, tt.includePaths, result, tt.expected)
			}
		})
	}
}

func TestFormatBytes(t *testing.T) {
	tests := []struct {
		name     string
		bytes    int64
		expected string
	}{
		{"zero bytes", 0, "0 B"},
		{"bytes", 512, "512 B"},
		{"kilobytes", 1024, "1.0 KB"},
		{"megabytes", 1048576, "1.0 MB"},
		{"gigabytes", 1073741824, "1.0 GB"},
		{"terabytes", 1099511627776, "1.0 TB"},
		{"mixed KB", 1536, "1.5 KB"},
		{"mixed MB", 2621440, "2.5 MB"},
		{"large value", 5368709120, "5.0 GB"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatBytes(tt.bytes)
			if result != tt.expected {
				t.Errorf("FormatBytes(%d) = %q; want %q", tt.bytes, result, tt.expected)
			}
		})
	}
}

func BenchmarkShouldInclude(b *testing.B) {
	excludePaths := []string{"*.log", "*.tmp", "node_modules", ".git"}
	includePaths := []string{"*.go", "*.md"}
	path := "/path/to/some/file.go"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ShouldInclude(path, excludePaths, includePaths)
	}
}

func BenchmarkFormatBytes(b *testing.B) {
	bytes := int64(1073741824) // 1 GB

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		FormatBytes(bytes)
	}
}
