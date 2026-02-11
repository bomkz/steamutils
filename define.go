package steamutils

import "github.com/iancoleman/orderedmap"

type SteamReader struct {
	libraryVdfPath    string
	steamPath         string
	libraryVdfMap     *orderedmap.OrderedMap
	SteamReaderConfig SteamReaderConfig
}

type SteamReaderConfig struct {
	LibraryVdfPathFinder  func(steamPath string) (string, error)
	SteamPathFinder       func() (string, error)
	customSteamPathFinder bool
	CustomLibraryVdfPath  string
	CustomSteamPath       string

	// Fixes Steam Path returning completely lowercase. Only affects if using the default SteamPathFinder
	FormatSteamPath bool
}

// The code in this file was made by ChatGPT, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
