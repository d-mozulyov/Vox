//go:build !darwin
// +build !darwin

package indicator

// isRetina always returns false on non-macOS platforms
func isRetina() bool {
	return false
}
