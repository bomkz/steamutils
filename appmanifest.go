package steamutils

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/orderedmap"
)

// GetAllInstalledApps returns a list of all installed Steam applications
// by reading appmanifest_<appid>.acf files from all Steam library folders
func (steamreader *SteamReader) GetAllInstalledApps() ([]InstalledApp, error) {
	var installedApps []InstalledApp

	libFoldersVal, exists := steamreader.libraryVdfMap.Get("libraryfolders")
	if !exists {
		return nil, fmt.Errorf("libraryfolders key not found in the VDF data")
	}

	libFolders, ok := libFoldersVal.(*orderedmap.OrderedMap)
	if !ok {
		return nil, fmt.Errorf("libraryfolders is not of the expected type")
	}

	// Iterate through each library folder
	for _, libKey := range libFolders.Keys() {
		libraryVal, exists := libFolders.Get(libKey)
		if !exists {
			continue
		}

		library, ok := libraryVal.(*orderedmap.OrderedMap)
		if !ok {
			continue
		}

		// Get the library path
		pathVal, exists := library.Get("path")
		if !exists {
			continue
		}

		libraryPath, ok := pathVal.(string)
		if !ok {
			continue
		}

		libraryPath = strings.ReplaceAll(libraryPath, "\\\\", "\\")

		// Get the apps from this library
		appsVal, exists := library.Get("apps")
		if !exists {
			continue
		}

		apps, ok := appsVal.(*orderedmap.OrderedMap)
		if !ok {
			continue
		}

		// Read each app's manifest file
		for _, appID := range apps.Keys() {
			app, err := readAppManifest(libraryPath, appID)
			if err != nil {
				// Skip apps that can't be read
				continue
			}
			app.LibraryPath = libraryPath
			installedApps = append(installedApps, app)
		}
	}

	return installedApps, nil
}

// GetInstalledAppByID retrieves details for a specific installed app by AppID
func (steamreader *SteamReader) GetInstalledAppByID(appID string) (*InstalledApp, error) {
	libraryPath, err := steamreader.FindAppIDPath(appID)
	if err != nil {
		return nil, err
	}

	app, err := readAppManifest(libraryPath, appID)
	if err != nil {
		return nil, err
	}

	app.LibraryPath = libraryPath
	return &app, nil
}

// readAppManifest reads and parses an appmanifest_<appid>.acf file
func readAppManifest(libraryPath, appID string) (InstalledApp, error) {
	manifestPath := filepath.Join(libraryPath, "steamapps", fmt.Sprintf("appmanifest_%s.acf", appID))

	data, err := os.ReadFile(manifestPath)
	if err != nil {
		return InstalledApp{}, fmt.Errorf("failed to read manifest file: %w", err)
	}

	manifestMap, err := Unmarshal(data)
	if err != nil {
		return InstalledApp{}, fmt.Errorf("failed to parse manifest file: %w", err)
	}

	appStateVal, exists := manifestMap.Get("AppState")
	if !exists {
		return InstalledApp{}, fmt.Errorf("AppState not found in manifest")
	}

	appState, ok := appStateVal.(*orderedmap.OrderedMap)
	if !ok {
		return InstalledApp{}, fmt.Errorf("AppState is not of expected type")
	}

	app := InstalledApp{
		AppID: appID,
	}

	// Extract app details from the manifest
	if nameVal, exists := appState.Get("name"); exists {
		if name, ok := nameVal.(string); ok {
			app.Name = name
		}
	}

	if installDirVal, exists := appState.Get("installdir"); exists {
		if installDir, ok := installDirVal.(string); ok {
			app.InstallDir = installDir
			// Construct full path
			app.FullPath = filepath.Join(libraryPath, "steamapps", "common", installDir)
		}
	}

	if buildIDVal, exists := appState.Get("buildid"); exists {
		if buildID, ok := buildIDVal.(string); ok {
			app.BuildID = buildID
		}
	}

	if sizeVal, exists := appState.Get("SizeOnDisk"); exists {
		if sizeStr, ok := sizeVal.(string); ok {
			fmt.Sscanf(sizeStr, "%d", &app.SizeOnDisk)
		}
	}

	if lastUpdatedVal, exists := appState.Get("LastUpdated"); exists {
		if lastUpdatedStr, ok := lastUpdatedVal.(string); ok {
			fmt.Sscanf(lastUpdatedStr, "%d", &app.LastUpdated)
		}
	}

	if lastPlayedVal, exists := appState.Get("LastPlayed"); exists {
		if lastPlayedStr, ok := lastPlayedVal.(string); ok {
			fmt.Sscanf(lastPlayedStr, "%d", &app.LastPlayed)
		}
	}

	// Parse InstalledDepots
	if installedDepotsVal, exists := appState.Get("InstalledDepots"); exists {
		if installedDepots, ok := installedDepotsVal.(*orderedmap.OrderedMap); ok {
			for _, depotID := range installedDepots.Keys() {
				depotVal, exists := installedDepots.Get(depotID)
				if !exists {
					continue
				}

				depotMap, ok := depotVal.(*orderedmap.OrderedMap)
				if !ok {
					continue
				}

				depot := InstalledDepot{
					DepotID: depotID,
				}

				// Get manifest ID
				if manifestVal, exists := depotMap.Get("manifest"); exists {
					if manifest, ok := manifestVal.(string); ok {
						depot.Manifest = manifest
					}
				}

				// Get size
				if sizeVal, exists := depotMap.Get("size"); exists {
					if sizeStr, ok := sizeVal.(string); ok {
						fmt.Sscanf(sizeStr, "%d", &depot.Size)
					}
				}

				// Get DLC App ID if this is DLC
				if dlcAppIDVal, exists := depotMap.Get("dlcappid"); exists {
					if dlcAppID, ok := dlcAppIDVal.(string); ok {
						depot.DLCAppID = dlcAppID
					}
				}

				app.InstalledDepots = append(app.InstalledDepots, depot)
			}
		}
	}

	return app, nil
}
