package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"zipprine/internal/models"
)

func TestDetectArchiveTypeByExtension(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-detect-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	testCases := []struct {
		name         string
		filename     string
		expectedType models.ArchiveType
	}{
		{"zip_extension", "test.zip", models.ZIP},
		{"tar_extension", "test.tar", models.TAR},
		{"targz_extension", "test.tar.gz", models.TARGZ},
		{"tgz_extension", "test.tgz", models.TARGZ},
		{"gz_extension", "test.gz", models.GZIP},
		{"uppercase_zip", "test.ZIP", models.ZIP},
		{"uppercase_tar", "test.TAR", models.TAR},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create empty file with the extension
			testFile := filepath.Join(tmpDir, tc.filename)
			os.WriteFile(testFile, []byte("dummy content"), 0644)

			detectedType, err := DetectArchiveType(testFile)
			if err != nil {
				t.Fatalf("DetectArchiveType failed: %v", err)
			}

			if detectedType != tc.expectedType {
				t.Errorf("Expected %s, got %s", tc.expectedType, detectedType)
			}
		})
	}
}

func TestDetectArchiveTypeByMagicBytes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-magic-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create actual archives to test magic byte detection
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test content"), 0644)

	testCases := []struct {
		name         string
		archiveType  models.ArchiveType
		extension    string
		expectedType models.ArchiveType
	}{
		{"zip_magic", models.ZIP, ".unknown", models.ZIP},
		{"tar_magic", models.TAR, ".unknown", models.TAR},
		{"targz_magic", models.TARGZ, ".unknown", models.TARGZ},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create archive with wrong extension to force magic byte detection
			archivePath := filepath.Join(tmpDir, "test"+tc.extension)
			
			config := &models.CompressConfig{
				SourcePath:       sourceDir,
				OutputPath:       archivePath,
				ArchiveType:      tc.archiveType,
				CompressionLevel: 5,
			}

			err := Compress(config)
			if err != nil {
				t.Fatalf("Failed to create archive: %v", err)
			}

			// Detect type by magic bytes
			detectedType, err := DetectArchiveType(archivePath)
			if err != nil {
				t.Fatalf("DetectArchiveType failed: %v", err)
			}

			if detectedType != tc.expectedType {
				t.Errorf("Expected %s, got %s", tc.expectedType, detectedType)
			}
		})
	}
}

func TestDetectZipMagicBytes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-zip-magic-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a real ZIP file
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("content"), 0644)

	zipPath := filepath.Join(tmpDir, "test.noext")
	createZip(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       zipPath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})

	// Detect without extension
	detectedType, err := DetectArchiveType(zipPath)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.ZIP {
		t.Errorf("Expected ZIP, got %s", detectedType)
	}
}

func TestDetectTarMagicBytes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-tar-magic-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a real TAR file
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("content"), 0644)

	tarPath := filepath.Join(tmpDir, "test.noext")
	createTar(&models.CompressConfig{
		SourcePath:  sourceDir,
		OutputPath:  tarPath,
		ArchiveType: models.TAR,
	})

	// Detect without extension
	detectedType, err := DetectArchiveType(tarPath)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.TAR {
		t.Errorf("Expected TAR, got %s", detectedType)
	}
}

func TestDetectTarGzMagicBytes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-targz-magic-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a real TAR.GZ file
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("content"), 0644)

	targzPath := filepath.Join(tmpDir, "test.noext")
	createTarGz(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       targzPath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	})

	// Detect without extension
	detectedType, err := DetectArchiveType(targzPath)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.TARGZ {
		t.Errorf("Expected TARGZ, got %s", detectedType)
	}
}

func TestDetectGzipMagicBytes(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-gzip-magic-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a real GZIP file (not tar.gz)
	sourceFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(sourceFile, []byte("test content"), 0644)

	gzipPath := filepath.Join(tmpDir, "test.noext")
	createGzip(&models.CompressConfig{
		SourcePath:       sourceFile,
		OutputPath:       gzipPath,
		ArchiveType:      models.GZIP,
		CompressionLevel: 5,
	})

	// Detect without extension
	detectedType, err := DetectArchiveType(gzipPath)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.GZIP {
		t.Errorf("Expected GZIP, got %s", detectedType)
	}
}

func TestDetectNonExistentFile(t *testing.T) {
	// DetectArchiveType returns type based on extension without checking file existence
	// This is by design - it detects type, not validates existence
	detectedType, err := DetectArchiveType("/nonexistent/file.zip")
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if detectedType != models.ZIP {
		t.Errorf("Expected ZIP type for .zip extension, got: %s", detectedType)
	}
	
	// Test with file that has no extension and doesn't exist
	_, err = DetectArchiveType("/nonexistent/file")
	if err == nil {
		t.Error("Expected error when trying to read magic bytes from non-existent file")
	}
}

func TestDetectUnknownFormat(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-unknown-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a file with unknown format
	unknownFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(unknownFile, []byte("just plain text"), 0644)

	detectedType, err := DetectArchiveType(unknownFile)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.AUTO {
		t.Errorf("Expected AUTO for unknown format, got %s", detectedType)
	}
}

func TestAnalyzeZipArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-analyze-zip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content2"), 0644)

	zipPath := filepath.Join(tmpDir, "test.zip")
	createZip(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       zipPath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})

	// Analyze
	info, err := Analyze(zipPath)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if info.Type != models.ZIP {
		t.Errorf("Expected type ZIP, got: %s", info.Type)
	}

	if info.FileCount != 2 {
		t.Errorf("Expected 2 files, got: %d", info.FileCount)
	}
}

func TestAnalyzeTarArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-analyze-tar-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("content"), 0644)

	tarPath := filepath.Join(tmpDir, "test.tar")
	createTar(&models.CompressConfig{
		SourcePath:  sourceDir,
		OutputPath:  tarPath,
		ArchiveType: models.TAR,
	})

	// Analyze
	info, err := Analyze(tarPath)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if info.Type != models.TAR {
		t.Errorf("Expected type TAR, got: %s", info.Type)
	}

	if info.FileCount < 1 {
		t.Errorf("Expected at least 1 file, got: %d", info.FileCount)
	}
}

func TestAnalyzeTarGzArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-analyze-targz-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("content"), 0644)

	targzPath := filepath.Join(tmpDir, "test.tar.gz")
	createTarGz(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       targzPath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	})

	// Analyze
	info, err := Analyze(targzPath)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if info.Type != models.TARGZ {
		t.Errorf("Expected type TARGZ, got: %s", info.Type)
	}

	if info.FileCount < 1 {
		t.Errorf("Expected at least 1 file, got: %d", info.FileCount)
	}
}

func TestAnalyzeGzipArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-analyze-gzip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test GZIP file
	sourceFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(sourceFile, []byte("test content"), 0644)

	gzipPath := filepath.Join(tmpDir, "test.txt.gz")
	createGzip(&models.CompressConfig{
		SourcePath:       sourceFile,
		OutputPath:       gzipPath,
		ArchiveType:      models.GZIP,
		CompressionLevel: 5,
	})

	// Analyze
	info, err := Analyze(gzipPath)
	if err != nil {
		t.Fatalf("Analyze failed: %v", err)
	}

	if info.Type != models.GZIP {
		t.Errorf("Expected type GZIP, got: %s", info.Type)
	}

	if info.FileCount != 1 {
		t.Errorf("Expected 1 file for GZIP, got: %d", info.FileCount)
	}

	if info.CompressedSize == 0 {
		t.Error("CompressedSize should not be zero")
	}
}

func TestAnalyzeNonExistentArchive(t *testing.T) {
	_, err := Analyze("/nonexistent/archive.zip")
	if err == nil {
		t.Error("Expected error for non-existent archive")
	}
}

func TestDetectWithSmallFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-small-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a very small file
	smallFile := filepath.Join(tmpDir, "small.dat")
	os.WriteFile(smallFile, []byte("x"), 0644)

	// Should not crash
	detectedType, err := DetectArchiveType(smallFile)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.AUTO {
		t.Logf("Detected type: %s", detectedType)
	}
}

func TestDetectEmptyFile(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-empty-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create an empty file
	emptyFile := filepath.Join(tmpDir, "empty.dat")
	os.WriteFile(emptyFile, []byte{}, 0644)

	// Should not crash
	detectedType, err := DetectArchiveType(emptyFile)
	if err != nil {
		t.Fatalf("DetectArchiveType failed: %v", err)
	}

	if detectedType != models.AUTO {
		t.Logf("Detected type: %s", detectedType)
	}
}

func BenchmarkDetectArchiveType(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "zipprine-bench-*")
	defer os.RemoveAll(tmpDir)

	// Create a test ZIP file
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("content"), 0644)

	zipPath := filepath.Join(tmpDir, "test.zip")
	createZip(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       zipPath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		DetectArchiveType(zipPath)
	}
}
