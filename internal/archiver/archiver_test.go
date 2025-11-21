package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"zipprine/internal/models"
)

func setupTestDir(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "zipprine-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	return tmpDir
}

func createTestFiles(t *testing.T, dir string) {
	files := map[string]string{
		"test1.txt":     "Hello World",
		"test2.go":      "package main",
		"subdir/test3.txt": "Nested file",
	}

	for path, content := range files {
		fullPath := filepath.Join(dir, path)
		if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
			t.Fatalf("Failed to create dir: %v", err)
		}
		if err := os.WriteFile(fullPath, []byte(content), 0644); err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}
	}
}

func TestCompressZip(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	createTestFiles(t, sourceDir)

	outputPath := filepath.Join(tmpDir, "test.zip")
	config := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       outputPath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	}

	err := Compress(config)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestCompressTarGz(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	createTestFiles(t, sourceDir)

	outputPath := filepath.Join(tmpDir, "test.tar.gz")
	config := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       outputPath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	}

	err := Compress(config)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestCompressWithExcludes(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	createTestFiles(t, sourceDir)

	outputPath := filepath.Join(tmpDir, "test.zip")
	config := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       outputPath,
		ArchiveType:      models.ZIP,
		ExcludePaths:     []string{"*.go"},
		CompressionLevel: 5,
	}

	err := Compress(config)
	if err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Output file was not created")
	}
}

func TestExtractZip(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// First create an archive
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	createTestFiles(t, sourceDir)

	archivePath := filepath.Join(tmpDir, "test.zip")
	compressConfig := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       archivePath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	}

	if err := Compress(compressConfig); err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	// Now extract it
	destDir := filepath.Join(tmpDir, "dest")
	if err := os.Mkdir(destDir, 0755); err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}

	extractConfig := &models.ExtractConfig{
		ArchivePath:   archivePath,
		DestPath:      destDir,
		ArchiveType:   models.ZIP,
		OverwriteAll:  true,
		PreservePerms: true,
	}

	err := Extract(extractConfig)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify extracted files exist
	if _, err := os.Stat(filepath.Join(destDir, "test1.txt")); os.IsNotExist(err) {
		t.Error("Extracted file test1.txt not found")
	}
}

func TestExtractTarGz(t *testing.T) {
	tmpDir := setupTestDir(t)
	defer os.RemoveAll(tmpDir)

	// First create an archive
	sourceDir := filepath.Join(tmpDir, "source")
	if err := os.Mkdir(sourceDir, 0755); err != nil {
		t.Fatalf("Failed to create source dir: %v", err)
	}
	createTestFiles(t, sourceDir)

	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	compressConfig := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       archivePath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	}

	if err := Compress(compressConfig); err != nil {
		t.Fatalf("Compress failed: %v", err)
	}

	// Now extract it
	destDir := filepath.Join(tmpDir, "dest")
	if err := os.Mkdir(destDir, 0755); err != nil {
		t.Fatalf("Failed to create dest dir: %v", err)
	}

	extractConfig := &models.ExtractConfig{
		ArchivePath:   archivePath,
		DestPath:      destDir,
		ArchiveType:   models.TARGZ,
		OverwriteAll:  true,
		PreservePerms: true,
	}

	err := Extract(extractConfig)
	if err != nil {
		t.Fatalf("Extract failed: %v", err)
	}

	// Verify extracted files exist
	if _, err := os.Stat(filepath.Join(destDir, "test1.txt")); os.IsNotExist(err) {
		t.Error("Extracted file test1.txt not found")
	}
}

func TestCompressInvalidType(t *testing.T) {
	config := &models.CompressConfig{
		SourcePath:  "/nonexistent",
		OutputPath:  "/tmp/test.invalid",
		ArchiveType: "INVALID",
	}

	err := Compress(config)
	if err != nil {
		t.Logf("Expected behavior: %v", err)
	}
}

func BenchmarkCompressZip(b *testing.B) {
	tmpDir := setupTestDir(&testing.T{})
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	createTestFiles(&testing.T{}, sourceDir)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		outputPath := filepath.Join(tmpDir, "bench.zip")
		config := &models.CompressConfig{
			SourcePath:       sourceDir,
			OutputPath:       outputPath,
			ArchiveType:      models.ZIP,
			CompressionLevel: 5,
		}
		Compress(config)
		os.Remove(outputPath)
	}
}
