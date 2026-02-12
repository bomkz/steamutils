//go:build linux
// +build linux

package steamutils

import (
	"fmt"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// GetSteamPath finds Steam's installation path on Linux
// Checks multiple common locations in order of likelihood
func GetSteamPath() (string, error) {
	// Get current user
	currentUser, err := user.Current()
	if err != nil {
		return "", fmt.Errorf("failed to get current user: %w", err)
	}

	// List of common Steam installation paths on Linux (in order of priority)
	steamPaths := []string{
		// Standard Steam installation
		filepath.Join(currentUser.HomeDir, ".steam", "steam"),
		filepath.Join(currentUser.HomeDir, ".local", "share", "Steam"),

		// Flatpak Steam
		filepath.Join(currentUser.HomeDir, ".var", "app", "com.valvesoftware.Steam", ".steam", "steam"),
		filepath.Join(currentUser.HomeDir, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam"),

		// Snap Steam
		filepath.Join(currentUser.HomeDir, "snap", "steam", "common", ".steam", "steam"),
		filepath.Join(currentUser.HomeDir, "snap", "steam", "common", ".local", "share", "Steam"),

		// Legacy paths
		filepath.Join(currentUser.HomeDir, ".steam", "debian-installation"),
		filepath.Join(currentUser.HomeDir, ".steam", "root"),

		// System-wide installation (uncommon but possible)
		"/usr/share/steam",
		"/usr/local/share/steam",

		// Check if STEAM_PATH environment variable is set
		os.Getenv("STEAM_PATH"),
	}

	// Try each path
	for _, path := range steamPaths {
		if path == "" {
			continue
		}

		// Follow symlinks
		resolvedPath, err := filepath.EvalSymlinks(path)
		if err != nil {
			// Try original path if symlink resolution fails
			resolvedPath = path
		}

		// Check if directory exists and contains Steam
		if _, err := os.Stat(filepath.Join(resolvedPath, "steamapps")); err == nil {
			return resolvedPath, nil
		}

		// Some installations might have steamapps as a symlink
		if _, err := os.Stat(resolvedPath); err == nil {
			// Check if this looks like a Steam directory
			entries, err := os.ReadDir(resolvedPath)
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
					return resolvedPath, nil
				}
			}
		}
	}

	// If no standard path found, try to find via registry file (Linux has a registry.vdf)
	registryPath := filepath.Join(currentUser.HomeDir, ".steam", "registry.vdf")
	if steamPath, err := findSteamPathFromRegistry(registryPath); err == nil && steamPath != "" {
		return steamPath, nil
	}

	return "", fmt.Errorf("Steam installation not found. Searched paths: %v", steamPaths)
}

// findSteamPathFromRegistry attempts to read Steam path from registry.vdf on Linux
func findSteamPathFromRegistry(registryPath string) (string, error) {
	data, err := os.ReadFile(registryPath)
	if err != nil {
		return "", err
	}

	// Parse the VDF file to find SteamPath or InstallPath
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Look for InstallPath or SteamPath entries
		if strings.Contains(line, "InstallPath") || strings.Contains(line, "SteamPath") {
			// Extract the path between quotes
			parts := strings.Split(line, "\"")
			if len(parts) >= 4 {
				path := parts[3]
				if _, err := os.Stat(path); err == nil {
					return path, nil
				}
			}
		}
	}

	return "", fmt.Errorf("Steam path not found in registry.vdf")
}

// GetAutoLoggedInSteamUsername returns Steam username on Linux
// This is less reliable than Windows but attempts to find it
func GetAutoLoggedInSteamUsername() (string, error) {
	currentUser, err := user.Current()
	if err != nil {
		return "", err
	}

	// Try to read from loginusers.vdf
	loginUsersPath := filepath.Join(currentUser.HomeDir, ".local", "share", "Steam", "config", "loginusers.vdf")

	// Try alternate paths
	alternatePaths := []string{
		loginUsersPath,
		filepath.Join(currentUser.HomeDir, ".steam", "steam", "config", "loginusers.vdf"),
		filepath.Join(currentUser.HomeDir, ".var", "app", "com.valvesoftware.Steam", ".local", "share", "Steam", "config", "loginusers.vdf"),
	}

	for _, path := range alternatePaths {
		username, err := readAutoLoginFromVDF(path)
		if err == nil && username != "" {
			return username, nil
		}
	}

	return "", fmt.Errorf("could not determine auto-login username")
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

// The code in this file was made by an LLM, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
