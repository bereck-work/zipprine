package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"zipprine/internal/models"
)

func TestCreateTar(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-tar-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content2"), 0644)

	// Create subdirectory
	subDir := filepath.Join(sourceDir, "subdir")
	os.Mkdir(subDir, 0755)
	os.WriteFile(filepath.Join(subDir, "file3.txt"), []byte("content3"), 0644)

	// Create TAR
	tarPath := filepath.Join(tmpDir, "test.tar")
	config := &models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       tarPath,
		ArchiveType:      models.TAR,
		CompressionLevel: 0,
	}

	err = createTar(config)
	if err != nil {
		t.Fatalf("createTar failed: %v", err)
	}

	// Verify TAR was created
	if _, err := os.Stat(tarPath); os.IsNotExist(err) {
		t.Error("TAR file was not created")
	}

	// Verify file size is reasonable
	info, _ := os.Stat(tarPath)
	if info.Size() < 100 {
		t.Error("TAR file seems too small")
	}
}

func TestCreateTarGz(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-targz-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("test content for compression"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("more test content"), 0644)

	// Create TAR.GZ with different compression levels
	testCases := []struct {
		name  string
		level int
	}{
		{"fast", 1},
		{"balanced", 5},
		{"best", 9},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			targzPath := filepath.Join(tmpDir, "test-"+tc.name+".tar.gz")
			config := &models.CompressConfig{
				SourcePath:       sourceDir,
				OutputPath:       targzPath,
				ArchiveType:      models.TARGZ,
				CompressionLevel: tc.level,
			}

			err = createTarGz(config)
			if err != nil {
				t.Fatalf("createTarGz failed: %v", err)
			}

			// Verify TAR.GZ was created
			if _, err := os.Stat(targzPath); os.IsNotExist(err) {
				t.Errorf("TAR.GZ file was not created")
			}
		})
	}
}

func TestCreateGzip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-gzip-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test file
	sourceFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("This is test content for GZIP compression")
	os.WriteFile(sourceFile, testContent, 0644)

	// Create GZIP
	gzipPath := filepath.Join(tmpDir, "test.txt.gz")
	config := &models.CompressConfig{
		SourcePath:       sourceFile,
		OutputPath:       gzipPath,
		ArchiveType:      models.GZIP,
		CompressionLevel: 5,
	}

	err = createGzip(config)
	if err != nil {
		t.Fatalf("createGzip failed: %v", err)
	}

	// Verify GZIP was created
	if _, err := os.Stat(gzipPath); os.IsNotExist(err) {
		t.Error("GZIP file was not created")
	}

	// Verify compressed file is smaller than original
	originalInfo, _ := os.Stat(sourceFile)
	compressedInfo, _ := os.Stat(gzipPath)
	if compressedInfo.Size() >= originalInfo.Size() {
		t.Log("Warning: Compressed file is not smaller (may be due to small test data)")
	}
}

func TestExtractTar(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-tar-extract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create and compress test files
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content2"), 0644)

	tarPath := filepath.Join(tmpDir, "test.tar")
	createTar(&models.CompressConfig{
		SourcePath:  sourceDir,
		OutputPath:  tarPath,
		ArchiveType: models.TAR,
	})

	// Extract
	destDir := filepath.Join(tmpDir, "dest")
	config := &models.ExtractConfig{
		ArchivePath:   tarPath,
		DestPath:      destDir,
		ArchiveType:   models.TAR,
		OverwriteAll:  true,
		PreservePerms: true,
	}

	err = extractTar(config)
	if err != nil {
		t.Fatalf("extractTar failed: %v", err)
	}

	// Verify files were extracted
	if _, err := os.Stat(filepath.Join(destDir, "file1.txt")); os.IsNotExist(err) {
		t.Error("file1.txt was not extracted")
	}
	if _, err := os.Stat(filepath.Join(destDir, "file2.txt")); os.IsNotExist(err) {
		t.Error("file2.txt was not extracted")
	}

	// Verify content
	content, _ := os.ReadFile(filepath.Join(destDir, "file1.txt"))
	if string(content) != "content1" {
		t.Errorf("Extracted content mismatch: got %q, want %q", string(content), "content1")
	}
}

func TestExtractTarGzDetailed(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-targz-extract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create and compress test files
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	testContent := "test content for tar.gz"
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte(testContent), 0644)

	targzPath := filepath.Join(tmpDir, "test.tar.gz")
	createTarGz(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       targzPath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	})

	// Extract
	destDir := filepath.Join(tmpDir, "dest")
	config := &models.ExtractConfig{
		ArchivePath:   targzPath,
		DestPath:      destDir,
		ArchiveType:   models.TARGZ,
		OverwriteAll:  true,
		PreservePerms: true,
	}

	err = extractTarGz(config)
	if err != nil {
		t.Fatalf("extractTarGz failed: %v", err)
	}

	// Verify file was extracted
	extractedFile := filepath.Join(destDir, "test.txt")
	if _, err := os.Stat(extractedFile); os.IsNotExist(err) {
		t.Error("test.txt was not extracted")
	}

	// Verify content
	content, _ := os.ReadFile(extractedFile)
	if string(content) != testContent {
		t.Errorf("Extracted content mismatch: got %q, want %q", string(content), testContent)
	}
}

func TestExtractGzip(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-gzip-extract-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create and compress test file
	sourceFile := filepath.Join(tmpDir, "test.txt")
	testContent := []byte("test content for gzip extraction")
	os.WriteFile(sourceFile, testContent, 0644)

	gzipPath := filepath.Join(tmpDir, "test.txt.gz")
	createGzip(&models.CompressConfig{
		SourcePath:       sourceFile,
		OutputPath:       gzipPath,
		ArchiveType:      models.GZIP,
		CompressionLevel: 5,
	})

	// Extract
	destDir := filepath.Join(tmpDir, "dest")
	os.Mkdir(destDir, 0755)
	config := &models.ExtractConfig{
		ArchivePath:   gzipPath,
		DestPath:      destDir,
		ArchiveType:   models.GZIP,
		OverwriteAll:  true,
		PreservePerms: true,
	}

	err = extractGzip(config)
	if err != nil {
		t.Fatalf("extractGzip failed: %v", err)
	}

	// Verify file was extracted (extractGzip removes .gz extension)
	extractedFile := filepath.Join(destDir, "test.txt")
	if _, err := os.Stat(extractedFile); os.IsNotExist(err) {
		t.Error("Extracted file was not created")
	}

	// Verify content
	content, _ := os.ReadFile(extractedFile)
	if string(content) != string(testContent) {
		t.Errorf("Extracted content mismatch: got %q, want %q", string(content), string(testContent))
	}
}

func TestAnalyzeTar(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-tar-analyze-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "file1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "file2.txt"), []byte("content2"), 0644)

	tarPath := filepath.Join(tmpDir, "test.tar")
	createTar(&models.CompressConfig{
		SourcePath:  sourceDir,
		OutputPath:  tarPath,
		ArchiveType: models.TAR,
	})

	// Analyze
	info, err := analyzeTar(tarPath, false)
	if err != nil {
		t.Fatalf("analyzeTar failed: %v", err)
	}

	// Verify results
	if info.Type != models.TAR {
		t.Errorf("Expected type TAR, got: %s", info.Type)
	}

	if info.FileCount < 2 {
		t.Errorf("Expected at least 2 files, got: %d", info.FileCount)
	}

	if len(info.Files) < 2 {
		t.Errorf("Expected at least 2 file entries, got: %d", len(info.Files))
	}

	if info.Checksum == "" {
		t.Error("Checksum should not be empty")
	}
}

func TestAnalyzeTarGz(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-targz-analyze-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archive
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test content"), 0644)

	targzPath := filepath.Join(tmpDir, "test.tar.gz")
	createTarGz(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       targzPath,
		ArchiveType:      models.TARGZ,
		CompressionLevel: 5,
	})

	// Analyze
	info, err := analyzeTar(targzPath, true)
	if err != nil {
		t.Fatalf("analyzeTar failed: %v", err)
	}

	// Verify results
	if info.Type != models.TARGZ {
		t.Errorf("Expected type TARGZ, got: %s", info.Type)
	}

	if info.FileCount < 1 {
		t.Errorf("Expected at least 1 file, got: %d", info.FileCount)
	}

	if info.CompressedSize == 0 {
		t.Error("CompressedSize should not be zero")
	}
}

func TestTarWithExcludePatterns(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-tar-exclude-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "include.txt"), []byte("include"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "exclude.log"), []byte("exclude"), 0644)
	os.WriteFile(filepath.Join(sourceDir, "also-include.txt"), []byte("include"), 0644)

	// Create TAR with exclude pattern
	tarPath := filepath.Join(tmpDir, "test.tar")
	config := &models.CompressConfig{
		SourcePath:   sourceDir,
		OutputPath:   tarPath,
		ArchiveType:  models.TAR,
		ExcludePaths: []string{"*.log"},
	}

	err = createTar(config)
	if err != nil {
		t.Fatalf("createTar failed: %v", err)
	}

	// Analyze to verify excluded files
	info, err := analyzeTar(tarPath, false)
	if err != nil {
		t.Fatalf("analyzeTar failed: %v", err)
	}

	// Should only have .txt files (may include directory entry)
	if info.FileCount < 2 {
		t.Errorf("Expected at least 2 files (excluded .log), got: %d", info.FileCount)
	}

	// Verify .log file is not in archive
	for _, file := range info.Files {
		if filepath.Ext(file.Name) == ".log" {
			t.Error("Excluded .log file found in archive")
		}
	}
}

func BenchmarkCreateTar(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "zipprine-bench-*")
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(sourceDir, "file"+string(rune('0'+i))+".txt"), []byte("content"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		tarPath := filepath.Join(tmpDir, "bench.tar")
		createTar(&models.CompressConfig{
			SourcePath:  sourceDir,
			OutputPath:  tarPath,
			ArchiveType: models.TAR,
		})
		os.Remove(tarPath)
	}
}

func BenchmarkCreateTarGz(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "zipprine-bench-*")
	defer os.RemoveAll(tmpDir)

	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	for i := 0; i < 10; i++ {
		os.WriteFile(filepath.Join(sourceDir, "file"+string(rune('0'+i))+".txt"), []byte("content"), 0644)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targzPath := filepath.Join(tmpDir, "bench.tar.gz")
		createTarGz(&models.CompressConfig{
			SourcePath:       sourceDir,
			OutputPath:       targzPath,
			ArchiveType:      models.TARGZ,
			CompressionLevel: 5,
		})
		os.Remove(targzPath)
	}
}
