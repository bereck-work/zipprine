package archiver

import (
	"zipprine/internal/models"
)

func Compress(config *models.CompressConfig) error {
	switch config.ArchiveType {
	case models.ZIP:
		return createZip(config)
	case models.TARGZ:
		return createTarGz(config)
	case models.TAR:
		return createTar(config)
	case models.GZIP:
		return createGzip(config)
	default:
		return nil
	}
}

func Extract(config *models.ExtractConfig) error {
	switch config.ArchiveType {
	case models.ZIP:
		return extractZip(config)
	case models.TARGZ:
		return extractTarGz(config)
	case models.TAR:
		return extractTar(config)
	case models.GZIP:
		return extractGzip(config)
	default:
		return nil
	}
}