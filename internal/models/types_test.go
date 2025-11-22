package models

import (
	"testing"
)

func TestArchiveTypeConstants(t *testing.T) {
	tests := []struct {
		name     string
		archType ArchiveType
		expected string
	}{
		{"ZIP type", ZIP, "ZIP"},
		{"TARGZ type", TARGZ, "TAR.GZ"},
		{"TAR type", TAR, "TAR"},
		{"GZIP type", GZIP, "GZIP"},
		{"RAR type", RAR, "RAR"},
		{"AUTO type", AUTO, "AUTO"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.archType) != tt.expected {
				t.Errorf("ArchiveType = %q; want %q", tt.archType, tt.expected)
			}
		})
	}
}

func TestCompressConfigValidation(t *testing.T) {
	config := &CompressConfig{
		SourcePath:       "/test/path",
		OutputPath:       "/test/output.zip",
		ArchiveType:      ZIP,
		ExcludePaths:     []string{"*.log"},
		IncludePaths:     []string{"*.go"},
		VerifyIntegrity:  true,
		CompressionLevel: 5,
	}

	if config.SourcePath != "/test/path" {
		t.Errorf("SourcePath = %q; want %q", config.SourcePath, "/test/path")
	}
	if config.ArchiveType != ZIP {
		t.Errorf("ArchiveType = %q; want %q", config.ArchiveType, ZIP)
	}
	if config.CompressionLevel != 5 {
		t.Errorf("CompressionLevel = %d; want %d", config.CompressionLevel, 5)
	}
}

func TestExtractConfigValidation(t *testing.T) {
	config := &ExtractConfig{
		ArchivePath:   "/test/archive.zip",
		DestPath:      "/test/dest",
		ArchiveType:   ZIP,
		OverwriteAll:  true,
		PreservePerms: true,
	}

	if config.ArchivePath != "/test/archive.zip" {
		t.Errorf("ArchivePath = %q; want %q", config.ArchivePath, "/test/archive.zip")
	}
	if !config.OverwriteAll {
		t.Error("OverwriteAll should be true")
	}
	if !config.PreservePerms {
		t.Error("PreservePerms should be true")
	}
}

func TestArchiveInfo(t *testing.T) {
	info := &ArchiveInfo{
		Type:             ZIP,
		FileCount:        10,
		TotalSize:        1024000,
		CompressedSize:   512000,
		CompressionRatio: 50.0,
		Files: []FileInfo{
			{Name: "test.txt", Size: 1024, IsDir: false, ModTime: "2024-01-01"},
		},
		Checksum: "abc123",
	}

	if info.FileCount != 10 {
		t.Errorf("FileCount = %d; want %d", info.FileCount, 10)
	}
	if info.CompressionRatio != 50.0 {
		t.Errorf("CompressionRatio = %f; want %f", info.CompressionRatio, 50.0)
	}
	if len(info.Files) != 1 {
		t.Errorf("Files length = %d; want %d", len(info.Files), 1)
	}
}