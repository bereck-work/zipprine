package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"zipprine/internal/models"
)

func TestCreateRarNotSupported(t *testing.T) {
	config := &models.CompressConfig{
		SourcePath:  "/tmp/test",
		OutputPath:  "/tmp/test.rar",
		ArchiveType: models.RAR,
	}

	err := createRar(config)
	if err == nil {
		t.Error("Expected error for RAR compression, got nil")
	}

	expectedMsg := "RAR compression is not supported"
	if err != nil && err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("Expected error message to start with %q, got %q", expectedMsg, err.Error())
	}
}

func TestRARDetectionByExtension(t *testing.T) {
	// Create a temporary file with .rar extension
	tempDir := t.TempDir()
	rarFile := filepath.Join(tempDir, "test.rar")

	// Create a file with RAR magic bytes
	file, err := os.Create(rarFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write RAR magic bytes: "Rar!" (0x52 0x61 0x72 0x21)
	rarMagic := []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}
	if _, err := file.Write(rarMagic); err != nil {
		t.Fatalf("Failed to write magic bytes: %v", err)
	}
	file.Close()

	// Test detection
	archiveType, err := DetectArchiveType(rarFile)
	if err != nil {
		t.Fatalf("Failed to detect archive type: %v", err)
	}

	if archiveType != models.RAR {
		t.Errorf("Expected RAR archive type, got %s", archiveType)
	}
}

func TestRARDetectionByMagicBytes(t *testing.T) {
	// Create a temporary file without .rar extension
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test.bin")

	file, err := os.Create(testFile)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	// Write RAR magic bytes
	rarMagic := []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00}
	if _, err := file.Write(rarMagic); err != nil {
		t.Fatalf("Failed to write magic bytes: %v", err)
	}
	file.Close()

	// Test detection by magic bytes
	archiveType, err := DetectArchiveType(testFile)
	if err != nil {
		t.Fatalf("Failed to detect archive type: %v", err)
	}

	if archiveType != models.RAR {
		t.Errorf("Expected RAR archive type by magic bytes, got %s", archiveType)
	}
}

func TestExtractRarInvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	invalidRar := filepath.Join(tempDir, "invalid.rar")

	// Create an invalid RAR file
	if err := os.WriteFile(invalidRar, []byte("not a rar file"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	config := &models.ExtractConfig{
		ArchivePath: invalidRar,
		DestPath:    filepath.Join(tempDir, "output"),
		ArchiveType: models.RAR,
	}

	err := extractRar(config)
	if err == nil {
		t.Error("Expected error for invalid RAR file, got nil")
	}
}

func TestExtractRarNonExistentFile(t *testing.T) {
	config := &models.ExtractConfig{
		ArchivePath: "/nonexistent/file.rar",
		DestPath:    "/tmp/output",
		ArchiveType: models.RAR,
	}

	err := extractRar(config)
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}

func TestAnalyzeRarInvalidFile(t *testing.T) {
	tempDir := t.TempDir()
	invalidRar := filepath.Join(tempDir, "invalid.rar")

	// Create an invalid RAR file
	if err := os.WriteFile(invalidRar, []byte("not a rar file"), 0644); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := analyzeRar(invalidRar)
	if err == nil {
		t.Error("Expected error for invalid RAR file, got nil")
	}
}

func TestAnalyzeRarNonExistentFile(t *testing.T) {
	_, err := analyzeRar("/nonexistent/file.rar")
	if err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}
}
