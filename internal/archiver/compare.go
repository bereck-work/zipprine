package archiver

import (
	"fmt"
	"sort"

	"zipprine/internal/models"
)

// ComparisonResult holds the result of comparing two archives
type ComparisonResult struct {
	OnlyInFirst  []string
	OnlyInSecond []string
	InBoth       []string
	Different    []DifferentFile
	Summary      string
}

// DifferentFile represents a file that exists in both archives but differs
type DifferentFile struct {
	Name      string
	Size1     int64
	Size2     int64
	ModTime1  string
	ModTime2  string
}

// CompareArchives compares two archives and returns differences
func CompareArchives(path1, path2 string, type1, type2 models.ArchiveType) (*ComparisonResult, error) {
	// Analyze both archives
	info1, err := AnalyzeArchive(path1, type1)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze first archive: %w", err)
	}

	info2, err := AnalyzeArchive(path2, type2)
	if err != nil {
		return nil, fmt.Errorf("failed to analyze second archive: %w", err)
	}

	// Create maps for quick lookup
	files1 := make(map[string]models.FileInfo)
	files2 := make(map[string]models.FileInfo)

	for _, f := range info1.Files {
		files1[f.Name] = f
	}

	for _, f := range info2.Files {
		files2[f.Name] = f
	}

	result := &ComparisonResult{
		OnlyInFirst:  []string{},
		OnlyInSecond: []string{},
		InBoth:       []string{},
		Different:    []DifferentFile{},
	}

	// Find files only in first archive
	for name := range files1 {
		if _, exists := files2[name]; !exists {
			result.OnlyInFirst = append(result.OnlyInFirst, name)
		}
	}

	// Find files only in second archive and files in both
	for name, file2 := range files2 {
		if file1, exists := files1[name]; !exists {
			result.OnlyInSecond = append(result.OnlyInSecond, name)
		} else {
			result.InBoth = append(result.InBoth, name)
			
			// Check if files are different
			if file1.Size != file2.Size || file1.ModTime != file2.ModTime {
				result.Different = append(result.Different, DifferentFile{
					Name:     name,
					Size1:    file1.Size,
					Size2:    file2.Size,
					ModTime1: file1.ModTime,
					ModTime2: file2.ModTime,
				})
			}
		}
	}

	// Sort results for consistent output
	sort.Strings(result.OnlyInFirst)
	sort.Strings(result.OnlyInSecond)
	sort.Strings(result.InBoth)

	// Generate summary
	result.Summary = fmt.Sprintf(
		"Comparison Summary:\n"+
			"  Files in both: %d\n"+
			"  Only in first: %d\n"+
			"  Only in second: %d\n"+
			"  Different: %d",
		len(result.InBoth),
		len(result.OnlyInFirst),
		len(result.OnlyInSecond),
		len(result.Different),
	)

	return result, nil
}

// AnalyzeArchive analyzes an archive and returns information about it
func AnalyzeArchive(path string, archiveType models.ArchiveType) (*models.ArchiveInfo, error) {
	switch archiveType {
	case models.ZIP:
		return analyzeZip(path)
	case models.TARGZ:
		return analyzeTar(path, true)
	case models.TAR:
		return analyzeTar(path, false)
	default:
		return nil, fmt.Errorf("unsupported archive type: %s", archiveType)
	}
}

// GetArchiveStats returns quick statistics about an archive
func GetArchiveStats(path string, archiveType models.ArchiveType) (string, error) {
	info, err := AnalyzeArchive(path, archiveType)
	if err != nil {
		return "", err
	}

	compressionRatio := 0.0
	if info.TotalSize > 0 {
		compressionRatio = (1.0 - float64(info.CompressedSize)/float64(info.TotalSize)) * 100
	}

	stats := fmt.Sprintf(
		"Archive Statistics:\n"+
			"  Type: %s\n"+
			"  Files: %d\n"+
			"  Total Size: %d bytes\n"+
			"  Compressed Size: %d bytes\n"+
			"  Compression Ratio: %.2f%%\n"+
			"  Checksum: %s",
		info.Type,
		info.FileCount,
		info.TotalSize,
		info.CompressedSize,
		compressionRatio,
		info.Checksum,
	)

	return stats, nil
}
