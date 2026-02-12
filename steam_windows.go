//go:build windows
// +build windows

package steamutils

import (
	"strings"

	"golang.org/x/sys/windows/registry"
)

// GetSteamPath finds Steam's installation path from Windows Registry
func GetSteamPath() (string, error) {
	root := registry.CURRENT_USER
	keyPath := `Software\Valve\Steam`

	SteamPath, err := readStringValueWithDefault(root, keyPath, "SteamPath", "")
	if err != nil {
		return "", err
	}

	SteamPath = strings.ReplaceAll(SteamPath, "/", "\\")
	return SteamPath, nil
}

// GetAutoLoggedInSteamUsername returns Steam username that has autologin enabled (Windows only)
func GetAutoLoggedInSteamUsername() (string, error) {
	root := registry.CURRENT_USER
	keyPath := `Software\Valve\Steam`

	UserName, err := readStringValueWithDefault(root, keyPath, "AutoLoginUser", "")
	return UserName, err
}

// readStringValueWithDefault reads a string value from the Windows Registry with a default value
func readStringValueWithDefault(root registry.Key, keyPath, valueName, defaultValue string) (string, error) {
	k, err := registry.OpenKey(root, keyPath, registry.QUERY_VALUE)
	if err != nil {
		return defaultValue, nil // Return the default value if the key or value doesn't exist
	}
	defer k.Close()

	value, _, err := k.GetStringValue(valueName)
	if err != nil {
		return defaultValue, nil // Return the default value if the value doesn't exist
	}

	return value, nil
}

// pathSeparator returns the OS-specific path separator
func pathSeparator() string {
	return "\\"
}
