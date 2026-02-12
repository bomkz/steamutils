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

// InstalledDepot represents a Steam depot (DLC or content package) installed for an app
type InstalledDepot struct {
	DepotID  string
	Manifest string
	Size     int64
	DLCAppID string // Only present if this depot is DLC
}

// InstalledApp represents a Steam application with its installation details
type InstalledApp struct {
	AppID           string
	Name            string
	InstallDir      string
	FullPath        string
	BuildID         string
	SizeOnDisk      int64
	LastUpdated     int64
	LastPlayed      int64
	LibraryPath     string
	InstalledDepots []InstalledDepot
}

// The code in this file was made by an LLM, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
