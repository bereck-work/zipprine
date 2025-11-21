package archiver

import (
	"os"
	"path/filepath"
	"testing"

	"zipprine/internal/models"
)

func TestBatchCompressSequential(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-batch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directories
	dir1 := filepath.Join(tmpDir, "source1")
	dir2 := filepath.Join(tmpDir, "source2")
	os.Mkdir(dir1, 0755)
	os.Mkdir(dir2, 0755)

	// Create test files
	os.WriteFile(filepath.Join(dir1, "test1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(dir2, "test2.txt"), []byte("content2"), 0644)

	// Create batch config
	configs := []*models.CompressConfig{
		{
			SourcePath:       dir1,
			OutputPath:       filepath.Join(tmpDir, "archive1.zip"),
			ArchiveType:      models.ZIP,
			CompressionLevel: 5,
		},
		{
			SourcePath:       dir2,
			OutputPath:       filepath.Join(tmpDir, "archive2.zip"),
			ArchiveType:      models.ZIP,
			CompressionLevel: 5,
		},
	}

	batchConfig := &BatchCompressConfig{
		Configs:    configs,
		Parallel:   false,
		MaxWorkers: 2,
	}

	errors := BatchCompress(batchConfig)

	// Check for errors
	for i, err := range errors {
		if err != nil {
			t.Errorf("Batch compress failed for config %d: %v", i, err)
		}
	}

	// Verify archives were created
	if _, err := os.Stat(filepath.Join(tmpDir, "archive1.zip")); os.IsNotExist(err) {
		t.Error("archive1.zip was not created")
	}
	if _, err := os.Stat(filepath.Join(tmpDir, "archive2.zip")); os.IsNotExist(err) {
		t.Error("archive2.zip was not created")
	}
}

func TestBatchCompressParallel(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-batch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test directories
	dir1 := filepath.Join(tmpDir, "source1")
	dir2 := filepath.Join(tmpDir, "source2")
	dir3 := filepath.Join(tmpDir, "source3")
	os.Mkdir(dir1, 0755)
	os.Mkdir(dir2, 0755)
	os.Mkdir(dir3, 0755)

	// Create test files
	os.WriteFile(filepath.Join(dir1, "test1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(dir2, "test2.txt"), []byte("content2"), 0644)
	os.WriteFile(filepath.Join(dir3, "test3.txt"), []byte("content3"), 0644)

	// Create batch config
	configs := []*models.CompressConfig{
		{
			SourcePath:       dir1,
			OutputPath:       filepath.Join(tmpDir, "archive1.tar.gz"),
			ArchiveType:      models.TARGZ,
			CompressionLevel: 5,
		},
		{
			SourcePath:       dir2,
			OutputPath:       filepath.Join(tmpDir, "archive2.tar.gz"),
			ArchiveType:      models.TARGZ,
			CompressionLevel: 5,
		},
		{
			SourcePath:       dir3,
			OutputPath:       filepath.Join(tmpDir, "archive3.tar.gz"),
			ArchiveType:      models.TARGZ,
			CompressionLevel: 5,
		},
	}

	batchConfig := &BatchCompressConfig{
		Configs:    configs,
		Parallel:   true,
		MaxWorkers: 2,
	}

	errors := BatchCompress(batchConfig)

	// Check for errors
	for i, err := range errors {
		if err != nil {
			t.Errorf("Batch compress failed for config %d: %v", i, err)
		}
	}

	// Verify archives were created
	for i := 1; i <= 3; i++ {
		archivePath := filepath.Join(tmpDir, "archive"+string(rune('0'+i))+".tar.gz")
		if _, err := os.Stat(archivePath); os.IsNotExist(err) {
			t.Errorf("archive%d.tar.gz was not created", i)
		}
	}
}

func TestBatchExtract(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-batch-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test archives first
	dir1 := filepath.Join(tmpDir, "source1")
	dir2 := filepath.Join(tmpDir, "source2")
	os.Mkdir(dir1, 0755)
	os.Mkdir(dir2, 0755)

	os.WriteFile(filepath.Join(dir1, "test1.txt"), []byte("content1"), 0644)
	os.WriteFile(filepath.Join(dir2, "test2.txt"), []byte("content2"), 0644)

	archive1 := filepath.Join(tmpDir, "archive1.zip")
	archive2 := filepath.Join(tmpDir, "archive2.zip")

	Compress(&models.CompressConfig{
		SourcePath:       dir1,
		OutputPath:       archive1,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})
	Compress(&models.CompressConfig{
		SourcePath:       dir2,
		OutputPath:       archive2,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})

	// Now batch extract
	dest1 := filepath.Join(tmpDir, "dest1")
	dest2 := filepath.Join(tmpDir, "dest2")

	configs := []*models.ExtractConfig{
		{
			ArchivePath:   archive1,
			DestPath:      dest1,
			ArchiveType:   models.ZIP,
			OverwriteAll:  true,
			PreservePerms: true,
		},
		{
			ArchivePath:   archive2,
			DestPath:      dest2,
			ArchiveType:   models.ZIP,
			OverwriteAll:  true,
			PreservePerms: true,
		},
	}

	batchConfig := &BatchExtractConfig{
		Configs:    configs,
		Parallel:   false,
		MaxWorkers: 2,
	}

	errors := BatchExtract(batchConfig)

	// Check for errors
	for i, err := range errors {
		if err != nil {
			t.Errorf("Batch extract failed for config %d: %v", i, err)
		}
	}

	// Verify files were extracted
	if _, err := os.Stat(filepath.Join(dest1, "test1.txt")); os.IsNotExist(err) {
		t.Error("test1.txt was not extracted")
	}
	if _, err := os.Stat(filepath.Join(dest2, "test2.txt")); os.IsNotExist(err) {
		t.Error("test2.txt was not extracted")
	}
}

func TestConvertArchive(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "zipprine-convert-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create source directory and file
	sourceDir := filepath.Join(tmpDir, "source")
	os.Mkdir(sourceDir, 0755)
	os.WriteFile(filepath.Join(sourceDir, "test.txt"), []byte("test content"), 0644)

	// Create ZIP archive
	zipPath := filepath.Join(tmpDir, "test.zip")
	err = Compress(&models.CompressConfig{
		SourcePath:       sourceDir,
		OutputPath:       zipPath,
		ArchiveType:      models.ZIP,
		CompressionLevel: 5,
	})
	if err != nil {
		t.Fatalf("Failed to create ZIP: %v", err)
	}

	// Convert to TAR.GZ
	targzPath := filepath.Join(tmpDir, "test.tar.gz")
	err = ConvertArchive(zipPath, targzPath, models.ZIP, models.TARGZ)
	if err != nil {
		t.Fatalf("Failed to convert archive: %v", err)
	}

	// Verify TAR.GZ was created
	if _, err := os.Stat(targzPath); os.IsNotExist(err) {
		t.Error("Converted archive was not created")
	}

	// Extract and verify contents
	destDir := filepath.Join(tmpDir, "dest")
	err = Extract(&models.ExtractConfig{
		ArchivePath:   targzPath,
		DestPath:      destDir,
		ArchiveType:   models.TARGZ,
		OverwriteAll:  true,
		PreservePerms: true,
	})
	if err != nil {
		t.Fatalf("Failed to extract converted archive: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(filepath.Join(destDir, "test.txt")); os.IsNotExist(err) {
		t.Error("File was not found in converted archive")
	}
}

func BenchmarkBatchCompressParallel(b *testing.B) {
	tmpDir, _ := os.MkdirTemp("", "zipprine-bench-*")
	defer os.RemoveAll(tmpDir)

	// Create test directories
	configs := make([]*models.CompressConfig, 4)
	for i := 0; i < 4; i++ {
		dir := filepath.Join(tmpDir, "source"+string(rune('0'+i)))
		os.Mkdir(dir, 0755)
		os.WriteFile(filepath.Join(dir, "test.txt"), []byte("content"), 0644)

		configs[i] = &models.CompressConfig{
			SourcePath:       dir,
			OutputPath:       filepath.Join(tmpDir, "archive"+string(rune('0'+i))+".zip"),
			ArchiveType:      models.ZIP,
			CompressionLevel: 5,
		}
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		batchConfig := &BatchCompressConfig{
			Configs:    configs,
			Parallel:   true,
			MaxWorkers: 2,
		}
		BatchCompress(batchConfig)
	}
}
