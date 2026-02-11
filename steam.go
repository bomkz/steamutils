package steamutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
	"golang.org/x/sys/windows/registry"
)

func NewSteamReader(steamReaderConfig SteamReaderConfig) (steamreader SteamReader, err error) {

	if steamReaderConfig.LibraryVdfPathFinder == nil {
		steamReaderConfig.LibraryVdfPathFinder = checkDefaultLibraryPath
	}

	steamReaderConfig.customSteamPathFinder = true
	if steamReaderConfig.SteamPathFinder == nil {
		steamReaderConfig.SteamPathFinder = GetSteamPath
		steamReaderConfig.customSteamPathFinder = false
	}

	steamreader.SteamReaderConfig = steamReaderConfig

	if steamReaderConfig.CustomSteamPath == "" {
		steamreader.steamPath, err = steamreader.SteamReaderConfig.SteamPathFinder()
		if err != nil {
			return
		}
	}

	if steamReaderConfig.CustomLibraryVdfPath == "" {
		steamreader.libraryVdfPath, err = steamreader.SteamReaderConfig.LibraryVdfPathFinder(steamreader.steamPath)
		if err != nil {
			return
		}
	}

	libraryVdfByte, err := os.ReadFile(steamreader.libraryVdfPath)
	if err != nil {
		return
	}

	steamreader.libraryVdfMap, err = Unmarshal(libraryVdfByte)
	if err != nil {
		return
	}

	return
}

func (steamreader *SteamReader) FindAppIDBuildID(AppID string) (buildId string, err error) {

	dir, err := steamreader.FindAppIDPath(AppID)
	if err != nil {
		return
	}

	f, err := os.ReadFile(dir + "\\steamapps\\appmanifest_" + AppID + ".acf")
	if err != nil {
		return
	}

	acf, err := Unmarshal(f)
	if err != nil {
		return
	}

	if appStateRaw, exists := acf.Get("AppState"); exists {
		if appState, ok := appStateRaw.(*orderedmap.OrderedMap); ok {
			buildIdInt, found := appState.Get("buildid")
			if found {
				buildId = buildIdInt.(string)
			}

		}
	}
	return

}

func checkDefaultLibraryPath(steamPath string) (librarypath string, err error) {
	_, err = os.Open(steamPath + "\\steamapps\\libraryfolders.vdf")
	if err != nil {
		return
	}

	librarypath = steamPath + "\\steamapps\\libraryfolders.vdf"
	return
}

func (steamreader *SteamReader) GetLibraryVdfMap() *orderedmap.OrderedMap {
	return steamreader.libraryVdfMap
}
func (steamreader *SteamReader) GetSteamPath() string {
	if !steamreader.SteamReaderConfig.customSteamPathFinder {
		pathStrings := strings.Split(steamreader.steamPath, "\\")
		rebuiltPathString := strings.ToUpper(pathStrings[0]) + "\\"
		if pathStrings[1] == "program files (x86)" {
			rebuiltPathString += "Program Files (x86)\\"

			for x, y := range pathStrings {
				if x == 0 {
					continue
				} else if x == 1 {
					continue
				}
				if y == "steam" {
					rebuiltPathString += "Steam\\"
				} else {
					rebuiltPathString += y + "\\"
				}
			}
		} else {
			for x, y := range pathStrings {
				if x == 0 {
					continue
				}
				rebuiltPathString += y + "\\"

			}

		}
		return rebuiltPathString
	}
	return steamreader.steamPath
}

func (steamreader *SteamReader) GetLibraryVdfPath() string {
	return steamreader.libraryVdfPath
}

// Goes through all libraries in libraryfolders.vdf to find the path of the library containing the target appid
func (steamreader *SteamReader) FindAppIDPath(targetAppID string) (string, error) {
	libFoldersVal, exists := steamreader.libraryVdfMap.Get("libraryfolders")
	if !exists {
		return "", fmt.Errorf("libraryfolders key not found in the VDF data")
	}

	libFolders, ok := libFoldersVal.(*orderedmap.OrderedMap)
	if !ok {
		return "", fmt.Errorf("libraryfolders is not of the expected type")
	}

	for _, libKey := range libFolders.Keys() {
		libraryVal, exists := libFolders.Get(libKey)
		if !exists {
			continue
		}

		library, ok := libraryVal.(*orderedmap.OrderedMap)
		if !ok {
			continue
		}

		appsVal, exists := library.Get("apps")
		if !exists {
			continue
		}

		apps, ok := appsVal.(*orderedmap.OrderedMap)
		if !ok {
			continue
		}

		for _, appKey := range apps.Keys() {
			if appKey == targetAppID {
				pathVal, exists := library.Get("path")
				if !exists {
					return "", fmt.Errorf("path not found in library %s for app %s", libKey, targetAppID)
				}
				pathStr, ok := pathVal.(string)
				if !ok {
					return "", fmt.Errorf("library path for library %s is not a string", libKey)
				}

				pathStr = strings.ReplaceAll(pathStr, "\\\\", "\\")
				return pathStr, nil
			}
		}
	}

	return "", fmt.Errorf("app with appid %s not found in any library", targetAppID)
}

// Returns Steam username that has autologin enabled
func GetAutoLoggedInSteamUsername() (string, error) {
	root := registry.CURRENT_USER
	keyPath := `Software\Valve\Steam`

	UserName, err := readStringValueWithDefault(root, keyPath, "AutoLoginUser", "")

	return UserName, err
}

// Finds Steam's Path from Windows Registry
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

// ReadStringValueWithDefault reads a string value from the Windows Registry with a default value.
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

// The code in this file was made by ChatGPT, and me, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
