package archiver

import (
	"fmt"
	"sync"

	"zipprine/internal/models"
)

// BatchCompressConfig holds configuration for batch compression operations
type BatchCompressConfig struct {
	Configs     []*models.CompressConfig
	Parallel    bool
	MaxWorkers  int
	OnProgress  func(index int, total int, filename string)
	OnError     func(index int, filename string, err error)
	OnComplete  func(index int, filename string)
}

// BatchCompress compresses multiple sources in batch
func BatchCompress(batchConfig *BatchCompressConfig) []error {
	if batchConfig.MaxWorkers <= 0 {
		batchConfig.MaxWorkers = 4
	}

	errors := make([]error, len(batchConfig.Configs))
	
	if !batchConfig.Parallel {
		// Sequential processing
		for i, config := range batchConfig.Configs {
			if batchConfig.OnProgress != nil {
				batchConfig.OnProgress(i+1, len(batchConfig.Configs), config.OutputPath)
			}
			
			err := Compress(config)
			errors[i] = err
			
			if err != nil && batchConfig.OnError != nil {
				batchConfig.OnError(i, config.OutputPath, err)
			} else if batchConfig.OnComplete != nil {
				batchConfig.OnComplete(i, config.OutputPath)
			}
		}
		return errors
	}

	// Parallel processing with worker pool
	var wg sync.WaitGroup
	jobs := make(chan int, len(batchConfig.Configs))
	
	// Start workers
	for w := 0; w < batchConfig.MaxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				config := batchConfig.Configs[i]
				
				if batchConfig.OnProgress != nil {
					batchConfig.OnProgress(i+1, len(batchConfig.Configs), config.OutputPath)
				}
				
				err := Compress(config)
				errors[i] = err
				
				if err != nil && batchConfig.OnError != nil {
					batchConfig.OnError(i, config.OutputPath, err)
				} else if batchConfig.OnComplete != nil {
					batchConfig.OnComplete(i, config.OutputPath)
				}
			}
		}()
	}

	// Send jobs
	for i := range batchConfig.Configs {
		jobs <- i
	}
	close(jobs)
	
	wg.Wait()
	return errors
}

// BatchExtractConfig holds configuration for batch extraction operations
type BatchExtractConfig struct {
	Configs     []*models.ExtractConfig
	Parallel    bool
	MaxWorkers  int
	OnProgress  func(index int, total int, filename string)
	OnError     func(index int, filename string, err error)
	OnComplete  func(index int, filename string)
}

// BatchExtract extracts multiple archives in batch
func BatchExtract(batchConfig *BatchExtractConfig) []error {
	if batchConfig.MaxWorkers <= 0 {
		batchConfig.MaxWorkers = 4
	}

	errors := make([]error, len(batchConfig.Configs))
	
	if !batchConfig.Parallel {
		// Sequential processing
		for i, config := range batchConfig.Configs {
			if batchConfig.OnProgress != nil {
				batchConfig.OnProgress(i+1, len(batchConfig.Configs), config.ArchivePath)
			}
			
			err := Extract(config)
			errors[i] = err
			
			if err != nil && batchConfig.OnError != nil {
				batchConfig.OnError(i, config.ArchivePath, err)
			} else if batchConfig.OnComplete != nil {
				batchConfig.OnComplete(i, config.ArchivePath)
			}
		}
		return errors
	}

	// Parallel processing with worker pool
	var wg sync.WaitGroup
	jobs := make(chan int, len(batchConfig.Configs))
	
	// Start workers
	for w := 0; w < batchConfig.MaxWorkers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := range jobs {
				config := batchConfig.Configs[i]
				
				if batchConfig.OnProgress != nil {
					batchConfig.OnProgress(i+1, len(batchConfig.Configs), config.ArchivePath)
				}
				
				err := Extract(config)
				errors[i] = err
				
				if err != nil && batchConfig.OnError != nil {
					batchConfig.OnError(i, config.ArchivePath, err)
				} else if batchConfig.OnComplete != nil {
					batchConfig.OnComplete(i, config.ArchivePath)
				}
			}
		}()
	}

	// Send jobs
	for i := range batchConfig.Configs {
		jobs <- i
	}
	close(jobs)
	
	wg.Wait()
	return errors
}

// ConvertArchive converts an archive from one format to another
func ConvertArchive(sourcePath, destPath string, sourceType, destType models.ArchiveType) error {
	// Create temporary directory for extraction
	tmpDir := destPath + ".tmp"
	
	// Extract source archive
	extractConfig := &models.ExtractConfig{
		ArchivePath:   sourcePath,
		DestPath:      tmpDir,
		ArchiveType:   sourceType,
		OverwriteAll:  true,
		PreservePerms: true,
	}
	
	if err := Extract(extractConfig); err != nil {
		return fmt.Errorf("failed to extract source archive: %w", err)
	}
	
	// Compress to destination format
	compressConfig := &models.CompressConfig{
		SourcePath:       tmpDir,
		OutputPath:       destPath,
		ArchiveType:      destType,
		CompressionLevel: 5,
	}
	
	if err := Compress(compressConfig); err != nil {
		return fmt.Errorf("failed to create destination archive: %w", err)
	}
	
	return nil
}
