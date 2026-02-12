//go:build linux
// +build linux

package steamutils

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetSteamPathLinux tests Steam path detection on Linux
func TestGetSteamPathLinux(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary test environment
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create standard Steam directory
	steamPath := filepath.Join(tempDir, ".local", "share", "Steam")
	steamappsPath := filepath.Join(steamPath, "steamapps")
	err := os.MkdirAll(steamappsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test Steam directory: %v", err)
	}

	// Test detection
	detectedPath, err := GetSteamPath()
	if err != nil {
		t.Fatalf("GetSteamPath failed: %v", err)
	}

	if detectedPath != steamPath {
		t.Errorf("Expected Steam path %s, got %s", steamPath, detectedPath)
	}
}

// TestGetSteamPathLinuxLegacy tests legacy Steam path detection
func TestGetSteamPathLinuxLegacy(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create legacy Steam directory
	steamPath := filepath.Join(tempDir, ".steam", "steam")
	steamappsPath := filepath.Join(steamPath, "steamapps")
	err := os.MkdirAll(steamappsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test Steam directory: %v", err)
	}

	detectedPath, err := GetSteamPath()
	if err != nil {
		t.Fatalf("GetSteamPath failed: %v", err)
	}

	if detectedPath != steamPath {
		t.Errorf("Expected Steam path %s, got %s", steamPath, detectedPath)
	}
}

// TestGetSteamPathLinuxFlatpak tests Flatpak Steam detection
func TestGetSteamPathLinuxFlatpak(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create Flatpak Steam directory
	steamPath := filepath.Join(tempDir, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam")
	steamappsPath := filepath.Join(steamPath, "steamapps")
	err := os.MkdirAll(steamappsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test Steam directory: %v", err)
	}

	detectedPath, err := GetSteamPath()
	if err != nil {
		t.Fatalf("GetSteamPath failed: %v", err)
	}

	if detectedPath != steamPath {
		t.Errorf("Expected Steam path %s, got %s", steamPath, detectedPath)
	}
}

// TestGetSteamPathLinuxSnap tests Snap Steam detection
func TestGetSteamPathLinuxSnap(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create Snap Steam directory
	steamPath := filepath.Join(tempDir, "snap", "steam", "common", ".local", "share", "Steam")
	steamappsPath := filepath.Join(steamPath, "steamapps")
	err := os.MkdirAll(steamappsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test Steam directory: %v", err)
	}

	detectedPath, err := GetSteamPath()
	if err != nil {
		t.Fatalf("GetSteamPath failed: %v", err)
	}

	if detectedPath != steamPath {
		t.Errorf("Expected Steam path %s, got %s", steamPath, detectedPath)
	}
}

// TestGetSteamPathLinuxEnvVar tests STEAM_PATH environment variable
func TestGetSteamPathLinuxEnvVar(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalSteamPath := os.Getenv("STEAM_PATH")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("STEAM_PATH", originalSteamPath)
	}()

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create custom Steam directory
	customPath := filepath.Join(tempDir, "custom", "steam")
	steamappsPath := filepath.Join(customPath, "steamapps")
	err := os.MkdirAll(steamappsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create test Steam directory: %v", err)
	}

	// Set environment variable
	os.Setenv("STEAM_PATH", customPath)

	detectedPath, err := GetSteamPath()
	if err != nil {
		t.Fatalf("GetSteamPath failed: %v", err)
	}

	if detectedPath != customPath {
		t.Errorf("Expected Steam path %s, got %s", customPath, detectedPath)
	}
}

// TestGetSteamPathLinuxSymlink tests symlink resolution
func TestGetSteamPathLinuxSymlink(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create actual Steam directory
	actualPath := filepath.Join(tempDir, "actual", "steam")
	steamappsPath := filepath.Join(actualPath, "steamapps")
	err := os.MkdirAll(steamappsPath, 0755)
	if err != nil {
		t.Fatalf("Failed to create actual Steam directory: %v", err)
	}

	// Create symlink
	symlinkDir := filepath.Join(tempDir, ".steam")
	err = os.MkdirAll(symlinkDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create symlink directory: %v", err)
	}

	symlinkPath := filepath.Join(symlinkDir, "steam")
	err = os.Symlink(actualPath, symlinkPath)
	if err != nil {
		t.Fatalf("Failed to create symlink: %v", err)
	}

	detectedPath, err := GetSteamPath()
	if err != nil {
		t.Fatalf("GetSteamPath failed: %v", err)
	}

	// Should resolve to actual path
	if detectedPath != actualPath {
		t.Errorf("Expected resolved path %s, got %s", actualPath, detectedPath)
	}
}

// TestGetSteamPathLinuxNotFound tests when Steam is not installed
func TestGetSteamPathLinuxNotFound(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalSteamPath := os.Getenv("STEAM_PATH")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("STEAM_PATH", originalSteamPath)
	}()

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	os.Setenv("STEAM_PATH", "")

	_, err := GetSteamPath()
	if err == nil {
		t.Error("Expected error when Steam not found, got nil")
	}
}

// TestGetAutoLoggedInSteamUsernameLinux tests username detection on Linux
func TestGetAutoLoggedInSteamUsernameLinux(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create loginusers.vdf
	configDir := filepath.Join(tempDir, ".local", "share", "Steam", "config")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	loginUsersContent := `"users"
{
	"76561198096000762"
	{
		"AccountName"		"testuser"
		"PersonaName"		"Test User"
		"RememberPassword"		"1"
		"MostRecent"		"1"
		"Timestamp"		"1766280889"
	}
}`

	loginUsersPath := filepath.Join(configDir, "loginusers.vdf")
	err = os.WriteFile(loginUsersPath, []byte(loginUsersContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create loginusers.vdf: %v", err)
	}

	username, err := GetAutoLoggedInSteamUsername()
	if err != nil {
		t.Fatalf("GetAutoLoggedInSteamUsername failed: %v", err)
	}

	if username != "testuser" {
		t.Errorf("Expected username 'testuser', got '%s'", username)
	}
}

// TestPathSeparatorLinux tests correct path separator
func TestPathSeparatorLinux(t *testing.T) {
	sep := pathSeparator()
	if sep != "/" {
		t.Errorf("Expected path separator '/', got '%s'", sep)
	}
}
