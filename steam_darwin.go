//go:build darwin
// +build darwin

package steamutils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetSteamPath finds Steam's installation path on macOS
func GetSteamPath() (string, error) {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	// List of common Steam installation paths on macOS
	steamPaths := []string{
		// Standard Steam installation
		filepath.Join(currentUser.HomeDir, "Library", "Application Support", "Steam"),

		// Legacy path
		filepath.Join(currentUser.HomeDir, "Library", "Application Support", "steam"),

		// Alternative installation via Homebrew Cask
		"/Applications/Steam.app/Contents/MacOS/Steam",

		// Check if STEAM_PATH environment variable is set
		os.Getenv("STEAM_PATH"),
	}

	// Try each path
	for _, path := range steamPaths {
		if path == "" {
			continue
		}

		// For .app bundle, we need to go to the data directory
		if strings.HasSuffix(path, "Steam.app/Contents/MacOS/Steam") {
			// Steam data is in ~/Library/Application Support/Steam
			path = filepath.Join(currentUser.HomeDir, "Library", "Application Support", "Steam")
		}

		// Check if directory exists and contains Steam
		if _, err := os.Stat(filepath.Join(path, "steamapps")); err == nil {
			return path, nil
		}

		// Check if this looks like a Steam directory
		if _, err := os.Stat(path); err == nil {
			entries, err := os.ReadDir(path)
			if err == nil {
				hasSteamFiles := false
				for _, entry := range entries {
					name := entry.Name()
					if name == "steamapps" || name == "userdata" || name == "config" {
						hasSteamFiles = true
						break
					}
				}
				if hasSteamFiles {
					return path, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Steam installation not found. Searched paths: %v", steamPaths)
}

// GetAutoLoggedInSteamUsername returns Steam username on macOS
func GetAutoLoggedInSteamUsername() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// Path to loginusers.vdf on macOS
	loginUsersPath := filepath.Join(currentUser.HomeDir, "Library", "Application Support", "Steam", "config", "loginusers.vdf")

	username, err := readAutoLoginFromVDF(loginUsersPath)
	if err != nil {
		return "", fmt.Errorf("could not determine auto-login username: %w", err)
	}

	return username, nil
}

// readAutoLoginFromVDF reads the auto-login username from loginusers.vdf
func readAutoLoginFromVDF(vdfPath string) (string, error) {
	data, err := os.ReadFile(vdfPath)
	if err != nil {
		return "", err
	}

	lines := strings.Split(string(data), "\n")
	var currentUser string
	mostRecentTimestamp := int64(0)

	for i, line := range lines {
		line = strings.TrimSpace(line)

		// Look for Steam ID entries (they're the user identifiers)
		if strings.HasPrefix(line, "\"7656") { // Steam IDs start with 765
			parts := strings.Split(line, "\"")
			if len(parts) >= 2 {
				// Look ahead for AccountName and Timestamp
				for j := i + 1; j < len(lines) && j < i+20; j++ {
					innerLine := strings.TrimSpace(lines[j])

					if strings.Contains(innerLine, "AccountName") {
						nameParts := strings.Split(innerLine, "\"")
						if len(nameParts) >= 4 {
							accountName := nameParts[3]

							// Check if this is the most recent login
							for k := i + 1; k < len(lines) && k < i+20; k++ {
								timestampLine := strings.TrimSpace(lines[k])
								if strings.Contains(timestampLine, "Timestamp") {
									tsParts := strings.Split(timestampLine, "\"")
									if len(tsParts) >= 4 {
										var timestamp int64
										fmt.Sscanf(tsParts[3], "%d", &timestamp)
										if timestamp > mostRecentTimestamp {
											mostRecentTimestamp = timestamp
											currentUser = accountName
										}
									}
									break
								}
							}
						}
						break
					}
				}
			}
		}
	}

	if currentUser != "" {
		return currentUser, nil
	}

	return "", fmt.Errorf("no user found in loginusers.vdf")
}

// pathSeparator returns the OS-specific path separator
func pathSeparator() string {
	return "/"
}
