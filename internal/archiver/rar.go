package archiver

import (
	"fmt"
	"os"
	"path/filepath"

	"zipprine/internal/models"

	"github.com/nwaples/rardecode"
)

// extractRar extracts a RAR archive
func extractRar(config *models.ExtractConfig) error {
	file, err := os.Open(config.ArchivePath)
	if err != nil {
		return fmt.Errorf("failed to open RAR file: %w", err)
	}
	defer file.Close()

	reader, err := rardecode.NewReader(file, "")
	if err != nil {
		return fmt.Errorf("failed to create RAR reader: %w", err)
	}

	if err := os.MkdirAll(config.DestPath, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	for {
		header, err := reader.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return fmt.Errorf("failed to read RAR entry: %w", err)
		}

		if header.IsDir {
			targetPath := filepath.Join(config.DestPath, header.Name)
			if err := os.MkdirAll(targetPath, 0755); err != nil {
				return fmt.Errorf("failed to create directory: %w", err)
			}
			continue
		}

		targetPath := filepath.Join(config.DestPath, header.Name)

		if _, err := os.Stat(targetPath); err == nil && !config.OverwriteAll {
			fmt.Printf("Skipping existing file: %s\n", header.Name)
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
			return fmt.Errorf("failed to create parent directory: %w", err)
		}

		outFile, err := os.Create(targetPath)
		if err != nil {
			return fmt.Errorf("failed to create file %s: %w", header.Name, err)
		}

		if _, err := outFile.ReadFrom(reader); err != nil {
			outFile.Close()
			return fmt.Errorf("failed to write file %s: %w", header.Name, err)
		}
		outFile.Close()

		// Set permissions if requested
		if config.PreservePerms {
			if err := os.Chmod(targetPath, header.Mode()); err != nil {
				fmt.Printf("Warning: failed to set permissions for %s: %v\n", header.Name, err)
			}
		}

		fmt.Printf("Extracted: %s\n", header.Name)
	}

	return nil
}

// analyzeRar analyzes a RAR archive and returns information about it
func analyzeRar(path string) (*models.ArchiveInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open RAR file: %w", err)
	}
	defer file.Close()

	fileStat, _ := file.Stat()

	reader, err := rardecode.NewReader(file, "")
	if err != nil {
		return nil, fmt.Errorf("failed to create RAR reader: %w", err)
	}

	info := &models.ArchiveInfo{
		Type:           models.RAR,
		CompressedSize: fileStat.Size(),
		Files:          []models.FileInfo{},
	}

	for {
		header, err := reader.Next()
		if err != nil {
			if err.Error() == "EOF" {
				break
			}
			return nil, fmt.Errorf("failed to read RAR entry: %w", err)
		}

		if !header.IsDir {
			info.FileCount++
			info.TotalSize += header.UnPackedSize
			info.Files = append(info.Files, models.FileInfo{
				Name:    header.Name,
				Size:    header.UnPackedSize,
				IsDir:   header.IsDir,
				ModTime: header.ModificationTime.Format("2006-01-02 15:04:05"),
			})
		}
	}

	if info.TotalSize > 0 {
		info.CompressionRatio = float64(info.CompressedSize) / float64(info.TotalSize)
	}

	return info, nil
}

// Note: RAR compression is proprietary and requires a license.
// This implementation only supports extraction using the rardecode library.
// For compression, users should use WinRAR or other licensed tools.
func createRar(config *models.CompressConfig) error {
	return fmt.Errorf("RAR compression is not supported due to proprietary format. Please use ZIP, TAR, or TAR.GZ for compression")
}
