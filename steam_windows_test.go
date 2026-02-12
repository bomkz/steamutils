//go:build windows
// +build windows

package steamutils

import (
	"testing"
)

// TestPathSeparatorWindows tests correct path separator on Windows
func TestPathSeparatorWindows(t *testing.T) {
	sep := pathSeparator()
	if sep != "\\" {
		t.Errorf("Expected path separator '\\', got '%s'", sep)
	}
}

// Note: Full GetSteamPath() testing on Windows requires registry access
// which is difficult to mock in unit tests. Integration tests should be
// run on actual Windows systems with Steam installed.

// TestGetSteamPathWindowsIntegration is a helper for manual integration testing
// Run with: go test -v -tags=integration
func TestGetSteamPathWindowsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	path, err := GetSteamPath()
	if err != nil {
		t.Logf("Steam not found (this is OK if Steam isn't installed): %v", err)
		return
	}

	t.Logf("Found Steam at: %s", path)

	// Basic validation
	if path == "" {
		t.Error("Steam path should not be empty")
	}
}

// TestGetAutoLoggedInSteamUsernameWindowsIntegration tests username detection
func TestGetAutoLoggedInSteamUsernameWindowsIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	username, err := GetAutoLoggedInSteamUsername()
	if err != nil {
		t.Logf("Auto-login user not found (this is OK if not configured): %v", err)
		return
	}

	t.Logf("Found auto-login user: %s", username)
}

// The code in this file was made by an LLM, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
