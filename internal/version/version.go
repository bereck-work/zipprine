package version

import "fmt"

const (
	Major = 1
	Minor = 0
	Patch = 3
)

// Version returns the semantic version string
func Version() string {
	return fmt.Sprintf("%d.%d.%d", Major, Minor, Patch)
}

// FullVersion returns the version with app name
func FullVersion() string {
	return fmt.Sprintf("Zipprine v%s", Version())
}
