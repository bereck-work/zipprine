package models

type ArchiveType string

const (
	ZIP    ArchiveType = "ZIP"
	TARGZ  ArchiveType = "TAR.GZ"
	TAR    ArchiveType = "TAR"
	GZIP   ArchiveType = "GZIP"
	RAR    ArchiveType = "RAR"
	AUTO   ArchiveType = "AUTO"
)

type CompressConfig struct {
	SourcePath      string
	OutputPath      string
	ArchiveType     ArchiveType
	ExcludePaths    []string
	IncludePaths    []string
	VerifyIntegrity bool
	CompressionLevel int
}

type ExtractConfig struct {
	ArchivePath   string
	DestPath      string
	ArchiveType   ArchiveType
	OverwriteAll  bool
	PreservePerms bool
}

type ArchiveInfo struct {
	Type            ArchiveType
	FileCount       int
	TotalSize       int64
	CompressedSize  int64
	CompressionRatio float64
	Files           []FileInfo
	Checksum        string
}

type FileInfo struct {
	Name string
	Size int64
	IsDir bool
	ModTime string
}