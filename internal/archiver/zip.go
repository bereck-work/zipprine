package archiver

import (
	"archive/zip"
	"compress/flate"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"zipprine/internal/models"
	"zipprine/pkg/fileutil"
)

func createZip(config *models.CompressConfig) error {
	outFile, err := os.Create(config.OutputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	zipWriter := zip.NewWriter(outFile)
	defer zipWriter.Close()

	// Set compression level
	if config.CompressionLevel > 0 {
		zipWriter.RegisterCompressor(zip.Deflate, func(out io.Writer) (io.WriteCloser, error) {
			return flate.NewWriter(out, config.CompressionLevel)
		})
	}

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

		if info.IsDir() {
			return nil
		}

		fmt.Printf("  → %s\n", relPath)

		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = relPath
		header.Method = zip.Deflate

		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}

		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(writer, file)
		return err
	})
}

func extractZip(config *models.ExtractConfig) error {
	r, err := zip.OpenReader(config.ArchivePath)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		destPath := filepath.Join(config.DestPath, f.Name)

		if f.FileInfo().IsDir() {
			os.MkdirAll(destPath, os.ModePerm)
			continue
		}

		if !config.OverwriteAll {
			if _, err := os.Stat(destPath); err == nil {
				fmt.Printf("  ⚠️  Skipping: %s (already exists)\n", f.Name)
				continue
			}
		}

		fmt.Printf("  → Extracting: %s\n", f.Name)

		if err := os.MkdirAll(filepath.Dir(destPath), os.ModePerm); err != nil {
			return err
		}

		outFile, err := os.Create(destPath)
		if err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			outFile.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}

		if config.PreservePerms {
			os.Chmod(destPath, f.Mode())
		}
	}

	return nil
}

func analyzeZip(path string) (*models.ArchiveInfo, error) {
	r, err := zip.OpenReader(path)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	info := &models.ArchiveInfo{
		Type:  models.ZIP,
		Files: []models.FileInfo{},
	}

	file, _ := os.Open(path)
	defer file.Close()
	fileStat, _ := file.Stat()
	info.CompressedSize = fileStat.Size()

	hash := sha256.New()
	io.Copy(hash, file)
	info.Checksum = fmt.Sprintf("%x", hash.Sum(nil))

	for _, f := range r.File {
		info.FileCount++
		info.TotalSize += int64(f.UncompressedSize64)

		if len(info.Files) < 100 {
			info.Files = append(info.Files, models.FileInfo{
				Name:    f.Name,
				Size:    int64(f.UncompressedSize64),
				IsDir:   f.FileInfo().IsDir(),
				ModTime: f.Modified.Format("2006-01-02 15:04:05"),
			})
		}
	}

	if info.TotalSize > 0 {
		info.CompressionRatio = (1 - float64(info.CompressedSize)/float64(info.TotalSize)) * 100
	}

	return info, nil
}