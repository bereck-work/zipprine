package archiver

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"zipprine/internal/models"
	"zipprine/pkg/fileutil"
)

func createTar(config *models.CompressConfig) error {
	outFile, err := os.Create(config.OutputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	tarWriter := tar.NewWriter(outFile)
	defer tarWriter.Close()

	return addToTar(tarWriter, config)
}

func createTarGz(config *models.CompressConfig) error {
	outFile, err := os.Create(config.OutputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzWriter, err := gzip.NewWriterLevel(outFile, config.CompressionLevel)
	if err != nil {
		return err
	}
	defer gzWriter.Close()

	tarWriter := tar.NewWriter(gzWriter)
	defer tarWriter.Close()

	return addToTar(tarWriter, config)
}

func createGzip(config *models.CompressConfig) error {
	inFile, err := os.Open(config.SourcePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(config.OutputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	gzWriter, err := gzip.NewWriterLevel(outFile, config.CompressionLevel)
	if err != nil {
		return err
	}
	defer gzWriter.Close()

	_, err = io.Copy(gzWriter, inFile)
	return err
}

func addToTar(tarWriter *tar.Writer, config *models.CompressConfig) error {
	return filepath.Walk(config.SourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !fileutil.ShouldInclude(path, config.ExcludePaths, config.IncludePaths) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		relPath, err := filepath.Rel(config.SourcePath, path)
		if err != nil {
			return err
		}

		header, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		header.Name = relPath

		if err := tarWriter.WriteHeader(header); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		fmt.Printf("  → %s\n", relPath)

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(tarWriter, file)
		return err
	})
}

func extractTar(config *models.ExtractConfig) error {
	file, err := os.Open(config.ArchivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	tarReader := tar.NewReader(file)
	return extractFromTar(tarReader, config)
}

func extractTarGz(config *models.ExtractConfig) error {
	file, err := os.Open(config.ArchivePath)
	if err != nil {
		return err
	}
	defer file.Close()

	gzReader, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)
	return extractFromTar(tarReader, config)
}

func extractGzip(config *models.ExtractConfig) error {
	inFile, err := os.Open(config.ArchivePath)
	if err != nil {
		return err
	}
	defer inFile.Close()

	gzReader, err := gzip.NewReader(inFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	outPath := filepath.Join(config.DestPath, filepath.Base(config.ArchivePath))
	outPath = outPath[:len(outPath)-3] // Remove .gz extension

	outFile, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, gzReader)
	return err
}

func extractFromTar(tarReader *tar.Reader, config *models.ExtractConfig) error {
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		destPath := filepath.Join(config.DestPath, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			os.MkdirAll(destPath, os.ModePerm)
		case tar.TypeReg:
			if !config.OverwriteAll {
				if _, err := os.Stat(destPath); err == nil {
					fmt.Printf("  ⚠️  Skipping: %s\n", header.Name)
					continue
				}
			}

			fmt.Printf("  → Extracting: %s\n", header.Name)

			os.MkdirAll(filepath.Dir(destPath), os.ModePerm)

			outFile, err := os.Create(destPath)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tarReader); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

			if config.PreservePerms {
				os.Chmod(destPath, os.FileMode(header.Mode))
			}
		}
	}
	return nil
}

func analyzeTar(path string, isGzipped bool) (*models.ArchiveInfo, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	info := &models.ArchiveInfo{
		Type:  models.TAR,
		Files: []models.FileInfo{},
	}

	if isGzipped {
		info.Type = models.TARGZ
	}

	fileStat, _ := file.Stat()
	info.CompressedSize = fileStat.Size()

	hash := sha256.New()
	io.Copy(hash, file)
	info.Checksum = fmt.Sprintf("%x", hash.Sum(nil))

	// Reopen for tar reading
	file.Seek(0, 0)

	var tarReader *tar.Reader
	if isGzipped {
		gzReader, err := gzip.NewReader(file)
		if err != nil {
			return nil, err
		}
		defer gzReader.Close()
		tarReader = tar.NewReader(gzReader)
	} else {
		tarReader = tar.NewReader(file)
	}

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		info.FileCount++
		info.TotalSize += header.Size

		if len(info.Files) < 100 {
			info.Files = append(info.Files, models.FileInfo{
				Name:    header.Name,
				Size:    header.Size,
				IsDir:   header.Typeflag == tar.TypeDir,
				ModTime: header.ModTime.Format("2006-01-02 15:04:05"),
			})
		}
	}

	if info.TotalSize > 0 {
		info.CompressionRatio = (1 - float64(info.CompressedSize)/float64(info.TotalSize)) * 100
	}

	return info, nil
}