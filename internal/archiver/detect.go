package archiver

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"zipprine/internal/models"
)

func DetectArchiveType(path string) (models.ArchiveType, error) {
	// First, try by extension
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".zip":
		return models.ZIP, nil
	case ".gz":
		if strings.HasSuffix(strings.ToLower(path), ".tar.gz") {
			return models.TARGZ, nil
		}
		return models.GZIP, nil
	case ".tar":
		return models.TAR, nil
	case ".tgz":
		return models.TARGZ, nil
	case ".rar":
		return models.RAR, nil
	}

	// Try by magic bytes
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	header := make([]byte, 512)
	n, err := file.Read(header)
	if err != nil && err != io.EOF {
		return "", err
	}
	header = header[:n]

	// ZIP magic: PK (0x504B)
	if len(header) >= 2 && header[0] == 0x50 && header[1] == 0x4B {
		return models.ZIP, nil
	}

	// GZIP magic: 0x1F 0x8B
	if len(header) >= 2 && header[0] == 0x1F && header[1] == 0x8B {
		// Check if it's a tar.gz by trying to decompress and check for tar header
		file.Seek(0, 0)
		gzReader, err := gzip.NewReader(file)
		if err == nil {
			defer gzReader.Close()
			tarHeader := make([]byte, 512)
			if n, _ := gzReader.Read(tarHeader); n >= 257 {
				// TAR magic: "ustar" at offset 257
				if bytes.Equal(tarHeader[257:262], []byte("ustar")) {
					return models.TARGZ, nil
				}
			}
		}
		return models.GZIP, nil
	}

	// TAR magic: "ustar" at offset 257
	if len(header) >= 262 && bytes.Equal(header[257:262], []byte("ustar")) {
		return models.TAR, nil
	}

	// RAR magic: Rar! (0x52 0x61 0x72 0x21)
	if len(header) >= 4 && header[0] == 0x52 && header[1] == 0x61 && header[2] == 0x72 && header[3] == 0x21 {
		return models.RAR, nil
	}

	return models.AUTO, nil
}

func Analyze(path string) (*models.ArchiveInfo, error) {
	archiveType, err := DetectArchiveType(path)
	if err != nil {
		return nil, err
	}

	switch archiveType {
	case models.ZIP:
		return analyzeZip(path)
	case models.TARGZ:
		return analyzeTar(path, true)
	case models.TAR:
		return analyzeTar(path, false)
	case models.RAR:
		return analyzeRar(path)
	case models.GZIP:
		// For GZIP, provide basic file info
		file, err := os.Open(path)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		fileStat, _ := file.Stat()
		return &models.ArchiveInfo{
			Type:           models.GZIP,
			CompressedSize: fileStat.Size(),
			FileCount:      1,
			Files:          []models.FileInfo{},
		}, nil
	default:
		return nil, nil
	}
}