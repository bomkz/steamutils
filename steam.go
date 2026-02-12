package steamutils

import (
	"fmt"
	"os"
	"strings"

	"github.com/iancoleman/orderedmap"
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
	// Use Stat instead of opening the file to avoid leaking file handles
	_, err = os.Stat(steamPath + "\\steamapps\\libraryfolders.vdf")
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
		if len(pathStrings) == 0 {
			return steamreader.steamPath
		}
		rebuiltPathString := strings.ToUpper(pathStrings[0]) + "\\"

		// Handle "Program Files (x86)" case insensitively
		if len(pathStrings) > 1 && strings.ToLower(pathStrings[1]) == "program files (x86)" {
			rebuiltPathString += "Program Files (x86)\\"

			for x, y := range pathStrings {
				if x == 0 || x == 1 {
					continue
				}
				if strings.ToLower(y) == "steam" {
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
				if strings.ToLower(y) == "steam" {
					rebuiltPathString += "Steam\\"
				} else {
					rebuiltPathString += y + "\\"
				}
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

// The code in this file was made by an LLM, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
