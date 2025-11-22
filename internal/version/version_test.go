package version

import (
	"strings"
	"testing"
)

func TestVersion(t *testing.T) {
	v := Version()
	
	if v == "" {
		t.Error("Version() returned empty string")
	}

	parts := strings.Split(v, ".")
	if len(parts) != 3 {
		t.Errorf("Version() = %q; expected format X.Y.Z", v)
	}
	expected := "1.0.3"
	if v != expected {
		t.Errorf("Version() = %q; want %q", v, expected)
	}
}

func TestFullVersion(t *testing.T) {
	fv := FullVersion()
	
	if fv == "" {
		t.Error("FullVersion() returned empty string")
	}

	if !strings.Contains(fv, "Zipprine") {
		t.Errorf("FullVersion() = %q; expected to contain 'Zipprine'", fv)
	}

	if !strings.Contains(fv, Version()) {
		t.Errorf("FullVersion() = %q; expected to contain version %q", fv, Version())
	}

	expected := "Zipprine v1.0.3"
	if fv != expected {
		t.Errorf("FullVersion() = %q; want %q", fv, expected)
	}
}

func TestVersionConstants(t *testing.T) {
	if Major != 1 {
		t.Errorf("Major = %d; want 1", Major)
	}
	if Minor != 0 {
		t.Errorf("Minor = %d; want 0", Minor)
	}
	if Patch != 3 {
		t.Errorf("Patch = %d; want 3", Patch)
	}
}

func TestVersionFormat(t *testing.T) {
	v := Version()
	
	if strings.Contains(v, " ") {
		t.Errorf("Version() contains spaces: %q", v)
	}

	if strings.HasPrefix(v, "v") {
		t.Errorf("Version() should not have 'v' prefix: %q", v)
	}

	for i, c := range v {
		if c != '.' && (c < '0' || c > '9') {
			t.Errorf("Version() contains non-numeric character at position %d: %q", i, v)
		}
	}
}

func TestFullVersionFormat(t *testing.T) {
	fv := FullVersion()
	
	expectedPrefix := "Zipprine v"
	if !strings.HasPrefix(fv, expectedPrefix) {
		t.Errorf("FullVersion() should start with %q, got %q", expectedPrefix, fv)
	}

	if !strings.HasSuffix(fv, Version()) {
		t.Errorf("FullVersion() should end with version %q, got %q", Version(), fv)
	}
}
