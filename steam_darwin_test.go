//go:build darwin
// +build darwin

package steamutils

import (
	"os"
	"path/filepath"
	"testing"
)

// TestGetSteamPathDarwin tests Steam path detection on macOS
func TestGetSteamPathDarwin(t *testing.T) {
	// Save original HOME
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	// Create temporary test environment
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create standard macOS Steam directory
	steamPath := filepath.Join(tempDir, "Library", "Application Support", "Steam")
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

// TestGetSteamPathDarwinLegacy tests legacy lowercase steam path
func TestGetSteamPathDarwinLegacy(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create legacy lowercase steam directory
	steamPath := filepath.Join(tempDir, "Library", "Application Support", "steam")
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

// TestGetSteamPathDarwinEnvVar tests STEAM_PATH environment variable on macOS
func TestGetSteamPathDarwinEnvVar(t *testing.T) {
	originalHome := os.Getenv("HOME")
	originalSteamPath := os.Getenv("STEAM_PATH")
	defer func() {
		os.Setenv("HOME", originalHome)
		os.Setenv("STEAM_PATH", originalSteamPath)
	}()

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create custom Steam directory
	customPath := filepath.Join(tempDir, "CustomSteam")
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

// TestGetSteamPathDarwinNotFound tests when Steam is not installed
func TestGetSteamPathDarwinNotFound(t *testing.T) {
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

// TestGetAutoLoggedInSteamUsernameDarwin tests username detection on macOS
func TestGetAutoLoggedInSteamUsernameDarwin(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create loginusers.vdf
	configDir := filepath.Join(tempDir, "Library", "Application Support", "Steam", "config")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	loginUsersContent := `"users"
{
	"76561198096000762"
	{
		"AccountName"		"macuser"
		"PersonaName"		"Mac User"
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

	if username != "macuser" {
		t.Errorf("Expected username 'macuser', got '%s'", username)
	}
}

// TestGetAutoLoggedInSteamUsernameDarwinMultipleUsers tests selecting most recent user
func TestGetAutoLoggedInSteamUsernameDarwinMultipleUsers(t *testing.T) {
	originalHome := os.Getenv("HOME")
	defer os.Setenv("HOME", originalHome)

	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)

	// Create loginusers.vdf with multiple users
	configDir := filepath.Join(tempDir, "Library", "Application Support", "Steam", "config")
	err := os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config directory: %v", err)
	}

	loginUsersContent := `"users"
{
	"76561198096000762"
	{
		"AccountName"		"olduser"
		"PersonaName"		"Old User"
		"Timestamp"		"1700000000"
	}
	"76561198096000999"
	{
		"AccountName"		"recentuser"
		"PersonaName"		"Recent User"
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

	if username != "recentuser" {
		t.Errorf("Expected most recent username 'recentuser', got '%s'", username)
	}
}

// TestPathSeparatorDarwin tests correct path separator
func TestPathSeparatorDarwin(t *testing.T) {
	sep := pathSeparator()
	if sep != "/" {
		t.Errorf("Expected path separator '/', got '%s'", sep)
	}
}
