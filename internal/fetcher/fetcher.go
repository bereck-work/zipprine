package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"zipprine/internal/archiver"
	"zipprine/internal/models"
)

// FetchAndExtract downloads an archive from a URL and extracts it to the destination path
func FetchAndExtract(archiveURL, destPath string, overwriteAll, preservePerms bool) error {

	parsedURL, err := url.Parse(archiveURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return fmt.Errorf("only HTTP and HTTPS URLs are supported")
	}

	filename := filepath.Base(parsedURL.Path)
	if filename == "" || filename == "." || filename == "/" {
		filename = "archive.tmp"
	}

	tempDir, err := os.MkdirTemp("", "zipprine-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	tempFile := filepath.Join(tempDir, filename)

	fmt.Printf("📥 Downloading from %s...\n", archiveURL)
	if err := downloadFile(tempFile, archiveURL); err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}

	fmt.Printf("✅ Download complete: %s\n", tempFile)

	archiveType, err := archiver.DetectArchiveType(tempFile)
	if err != nil {
		return fmt.Errorf("failed to detect archive type: %w", err)
	}

	if archiveType == models.AUTO {
		return fmt.Errorf("could not detect archive type from downloaded file")
	}

	fmt.Printf("📦 Detected archive type: %s\n", archiveType)

	fmt.Printf("📂 Extracting to %s...\n", destPath)
	extractConfig := &models.ExtractConfig{
		ArchivePath:   tempFile,
		DestPath:      destPath,
		ArchiveType:   archiveType,
		OverwriteAll:  overwriteAll,
		PreservePerms: preservePerms,
	}

	if err := archiver.Extract(extractConfig); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	fmt.Printf("✨ Extraction complete!\n")
	return nil
}

// downloadFile downloads a file from a URL to a local path with progress indication
func downloadFile(filepath, url string) error {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	contentLength := resp.ContentLength

	var reader io.Reader = resp.Body
	if contentLength > 0 {
		reader = &progressReader{
			reader:        resp.Body,
			total:         contentLength,
			downloaded:    0,
			lastPrintSize: 0,
		}
	}

	// Write the body to file
	_, err = io.Copy(out, reader)
	if err != nil {
		return err
	}

	fmt.Println()
	return nil
}

// progressReader wraps an io.Reader to provide download progress
type progressReader struct {
	reader        io.Reader
	total         int64
	downloaded    int64
	lastPrintSize int64
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.downloaded += int64(n)

	// Print progress every 1MB or at completion
	if pr.downloaded-pr.lastPrintSize > 1024*1024 || err == io.EOF {
		pr.lastPrintSize = pr.downloaded
		percentage := float64(pr.downloaded) / float64(pr.total) * 100
		fmt.Printf("\r📊 Progress: %.2f%% (%d/%d bytes)", percentage, pr.downloaded, pr.total)
	}

	return n, err
}

// GetFilenameFromURL extracts a filename from a URL
func GetFilenameFromURL(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	filename := filepath.Base(parsedURL.Path)
	if filename == "" || filename == "." || filename == "/" {
		return "", fmt.Errorf("could not extract filename from URL")
	}

	return filename, nil
}

// IsValidArchiveURL checks if a URL points to a supported archive format
func IsValidArchiveURL(urlStr string) bool {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return false
	}

	path := strings.ToLower(parsedURL.Path)
	validExtensions := []string{".zip", ".tar", ".tar.gz", ".tgz", ".gz", ".rar"}

	for _, ext := range validExtensions {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}

	return false
}
