package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"zipprine/internal/models"
)

func TestCompareArchives(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-compare-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create first archive
	dir1 := filepath.Join(tmpDir, "source1")
	os.Mkdir(dir1, 0755)
	os.WriteFile(filepath.Join(dir1, "common.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(dir1, "only1.txt"), []byte("only in first"), 0644)

	archive1 := filepath.Join(tmpDir, "archive1.zip")
	err = Compress(&models.CompressConfig{
		SourcePath:       dir1,
		OutputPath:       archive1,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})
	if err != nil {
		t.Fatalf("Failed to create first archive: %v", err)
	}

	// Create second archive
	dir2 := filepath.Join(tmpDir, "source2")
	os.Mkdir(dir2, 0755)
	os.WriteFile(filepath.Join(dir2, "common.txt"), []byte("content2"), 0644)
	os.WriteFile(filepath.Join(dir2, "only2.txt"), []byte("only in second"), 0644)

	archive2 := filepath.Join(tmpDir, "archive2.zip")
	err = Compress(&models.CompressConfig{
		SourcePath:       dir2,
		OutputPath:       archive2,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})
	if err != nil {
		t.Fatalf("Failed to create second archive: %v", err)
	}

	// Compare archives
	result, err := CompareArchives(archive1, archive2, models.ZIP, models.ZIP)
	if err != nil {
		t.Fatalf("Failed to compare archives: %v", err)
	}

	// Verify results
	if len(result.OnlyInFirst) != 1 || result.OnlyInFirst[0] != "only1.txt" {
		t.Errorf("Expected only1.txt in OnlyInFirst, got: %v", result.OnlyInFirst)
	}

	if len(result.OnlyInSecond) != 1 || result.OnlyInSecond[0] != "only2.txt" {
		t.Errorf("Expected only2.txt in OnlyInSecond, got: %v", result.OnlyInSecond)
	}

	if len(result.InBoth) != 1 || result.InBoth[0] != "common.txt" {
		t.Errorf("Expected common.txt in InBoth, got: %v", result.InBoth)
	}

	if result.Summary == "" {
		t.Error("Summary should not be empty")
	}
}

func TestCompareIdenticalArchives(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-compare-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content"), 0644)

	// Create two identical archives
	archive1 := filepath.Join(tmpDir, "archive1.zip")
	archive2 := filepath.Join(tmpDir, "archive2.zip")

	config := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       archive1,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	}
	Compress(config)

	config.OutputPath = archive2
	Compress(config)

	// Compare archives
	result, err := CompareArchives(archive1, archive2, models.ZIP, models.ZIP)
	if err != nil {
		t.Fatalf("Failed to compare archives: %v", err)
	}

	// Verify results
	if len(result.OnlyInFirst) != 0 {
		t.Errorf("Expected no files only in first, got: %v", result.OnlyInFirst)
	}

	if len(result.OnlyInSecond) != 0 {
		t.Errorf("Expected no files only in second, got: %v", result.OnlyInSecond)
	}

	if len(result.InBoth) != 2 {
		t.Errorf("Expected 2 files in both, got: %d", len(result.InBoth))
	}
}

func TestAnalyzeArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-analyze-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content2"), 0644)

	archivePath := filepath.Join(tmpDir, "test.zip")
	err = Compress(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       archivePath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})
	if err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	// Analyze archive
	info, err := AnalyzeArchive(archivePath, models.ZIP)
	if err != nil {
		t.Fatalf("Failed to analyze archive: %v", err)
	}

	// Verify results
	if info.Type != models.ZIP {
		t.Errorf("Expected type ZIP, got: %s", info.Type)
	}

	if info.FileCount != 2 {
		t.Errorf("Expected 2 files, got: %d", info.FileCount)
	}

	if len(info.Files) != 2 {
		t.Errorf("Expected 2 file entries, got: %d", len(info.Files))
	}
}

func TestGetArchiveStats(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-stats-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test content"), 0644)

	archivePath := filepath.Join(tmpDir, "test.tar.gz")
	err = Compress(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       archivePath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	})
	if err != nil {
		t.Fatalf("Failed to create archive: %v", err)
	}

	// Get stats
	stats, err := GetArchiveStats(archivePath, models.TARGZ)
	if err != nil {
		t.Fatalf("Failed to get archive stats: %v", err)
	}

	if stats == "" {
		t.Error("Stats should not be empty")
	}

	// Stats should contain key information
	if !contains(stats, "Type") || !contains(stats, "Files") || !contains(stats, "Checksum") {
		t.Errorf("Stats missing key information: %s", stats)
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
