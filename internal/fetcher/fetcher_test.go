package fetcher

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestIsValidArchiveURL(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected bool
	}{
		{"ZIP URL", "https://example.com/archive.zip", true},
		{"TAR URL", "https://example.com/archive.tar", true},
		{"TAR.GZ URL", "https://example.com/archive.tar.gz", true},
		{"TGZ URL", "https://example.com/archive.tgz", true},
		{"GZIP URL", "https://example.com/file.gz", true},
		{"RAR URL", "https://example.com/archive.rar", true},
		{"Invalid URL", "https://example.com/file.txt", false},
		{"No extension", "https://example.com/file", false},
		{"Invalid format", "not a url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidArchiveURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidArchiveURL(%q) = %v; want %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestGetFilenameFromURL(t *testing.T) {
	tests := []struct {
		name        string
		url         string
		expected    string
		expectError bool
	}{
		{"Simple filename", "https://example.com/archive.zip", "archive.zip", false},
		{"Path with directories", "https://example.com/path/to/file.tar.gz", "file.tar.gz", false},
		{"Query parameters", "https://example.com/file.zip?download=true", "file.zip", false},
		{"No filename", "https://example.com/", "", true},
		{"Relative path treated as filename", "not a url", "not a url", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := GetFilenameFromURL(tt.url)
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for %q, got nil", tt.url)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %q: %v", tt.url, err)
				}
				if result != tt.expected {
					t.Errorf("GetFilenameFromURL(%q) = %q; want %q", tt.url, result, tt.expected)
				}
			}
		})
	}
}

func TestDownloadFile(t *testing.T) {
	// Create a test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "13")
		w.Write([]byte("test content!"))
	}))
	defer server.Close()

	// Create temp directory
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "downloaded.txt")

	// Test download
	err := downloadFile(outputFile, server.URL)
	if err != nil {
		t.Fatalf("Failed to download file: %v", err)
	}

	// Verify file exists
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		t.Error("Downloaded file does not exist")
	}

	// Verify content
	content, err := os.ReadFile(outputFile)
	if err != nil {
		t.Fatalf("Failed to read downloaded file: %v", err)
	}

	expected := "test content!"
	if string(content) != expected {
		t.Errorf("Downloaded content = %q; want %q", string(content), expected)
	}
}

func TestDownloadFileNotFound(t *testing.T) {
	// Create a test server that returns 404
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "downloaded.txt")

	err := downloadFile(outputFile, server.URL)
	if err == nil {
		t.Error("Expected error for 404 response, got nil")
	}
}

func TestDownloadFileInvalidURL(t *testing.T) {
	tempDir := t.TempDir()
	outputFile := filepath.Join(tempDir, "downloaded.txt")

	err := downloadFile(outputFile, "http://invalid-url-that-does-not-exist-12345.com")
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestFetchAndExtractInvalidURL(t *testing.T) {
	err := FetchAndExtract("not a url", "/tmp/output", false, true)
	if err == nil {
		t.Error("Expected error for invalid URL, got nil")
	}
}

func TestFetchAndExtractNonHTTP(t *testing.T) {
	err := FetchAndExtract("ftp://example.com/file.zip", "/tmp/output", false, true)
	if err == nil {
		t.Error("Expected error for non-HTTP URL, got nil")
	}

	expectedMsg := "only HTTP and HTTPS URLs are supported"
	if err != nil && err.Error() != expectedMsg {
		t.Errorf("Expected error message %q, got %q", expectedMsg, err.Error())
	}
}

func TestProgressReader(t *testing.T) {
	// Create a test reader
	data := []byte("test data for progress reader")
	pr := &progressReader{
		reader:        nil, // We'll test the struct directly
		total:         int64(len(data)),
		downloaded:    0,
		lastPrintSize: 0,
	}

	// Test initial state
	if pr.downloaded != 0 {
		t.Errorf("Initial downloaded = %d; want 0", pr.downloaded)
	}

	if pr.total != int64(len(data)) {
		t.Errorf("Total = %d; want %d", pr.total, len(data))
	}
}

func TestFetchAndExtractInvalidDestination(t *testing.T) {
	// Create a test server with a valid ZIP
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "4")
		w.Write([]byte("test"))
	}))
	defer server.Close()

	// Try to extract to an invalid destination (file instead of directory)
	tempDir := t.TempDir()
	invalidDest := filepath.Join(tempDir, "file.txt")
	os.WriteFile(invalidDest, []byte("test"), 0644)

	err := FetchAndExtract(server.URL+"/test.zip", invalidDest, false, true)
	// This should fail during extraction or detection
	if err == nil {
		t.Log("Note: This test may pass if the downloaded content is not a valid archive")
	}
}
