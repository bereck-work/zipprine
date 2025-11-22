package version

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	
	// Check that version is not empty
	if v == "" {
		t.Error("Version() returned empty string")
	}

	// Check format (should be X.Y.Z)
	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		t.Errorf("Version() = %q; expected format X.Y.Z", v)
	}

	// Verify it matches the constants
	expected := "1.0.0"
	if v != expected {
		t.Errorf("Version() = %q; want %q", v, expected)
	}
}

func TestFullVersion(t *testing.T) {
	fv := FullVersion()
	
	// Check that full version is not empty
	if fv == "" {
		t.Error("FullVersion() returned empty string")
	}

	// Check that it contains "Zipprine"
	if !strings.Contains(fv, "Zipprine") {
		t.Errorf("FullVersion() = %q; expected to contain 'Zipprine'", fv)
	}

	// Check that it contains the version
	if !strings.Contains(fv, Version()) {
		t.Errorf("FullVersion() = %q; expected to contain version %q", fv, Version())
	}

	// Verify exact format
	expected := "Zipprine v1.0.0"
	if fv != expected {
		t.Errorf("FullVersion() = %q; want %q", fv, expected)
	}
}

func TestVersionConstants(t *testing.T) {
	// Test that constants are set correctly
	if Major != 1 {
		t.Errorf("Major = %d; want 1", Major)
	}
	if Minor != 0 {
		t.Errorf("Minor = %d; want 0", Minor)
	}
	if Patch != 0 {
		t.Errorf("Patch = %d; want 0", Patch)
	}
}

func TestVersionFormat(t *testing.T) {
	// Test that version follows semantic versioning
	v := Version()
	
	// Should not contain spaces
	if strings.Contains(v, " ") {
		t.Errorf("Version() contains spaces: %q", v)
	}

	// Should not contain 'v' prefix
	if strings.HasPrefix(v, "v") {
		t.Errorf("Version() should not have 'v' prefix: %q", v)
	}

	// Should be numeric with dots
	for i, c := range v {
		if c != '.' && (c < '0' || c > '9') {
			t.Errorf("Version() contains non-numeric character at position %d: %q", i, v)
		}
	}
}

func TestFullVersionFormat(t *testing.T) {
	fv := FullVersion()
	
	// Should start with "Zipprine v"
	expectedPrefix := "Zipprine v"
	if !strings.HasPrefix(fv, expectedPrefix) {
		t.Errorf("FullVersion() should start with %q, got %q", expectedPrefix, fv)
	}

	// Should end with version number
	if !strings.HasSuffix(fv, Version()) {
		t.Errorf("FullVersion() should end with version %q, got %q", Version(), fv)
	}
}
