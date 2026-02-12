package steamutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/iancoleman/orderedmap"
)

// TestInstalledDepotStruct tests the InstalledDepot structure
func TestInstalledDepotStruct(t *testing.T) {
	depot := InstalledDepot{
		DepotID:  "667971",
		Manifest: "809774009354886606",
		Size:     3488450479,
		DLCAppID: "",
	}

	if depot.DepotID != "667971" {
		t.Errorf("Expected DepotID 667971, got %s", depot.DepotID)
	}

	if depot.Manifest != "809774009354886606" {
		t.Errorf("Expected Manifest 809774009354886606, got %s", depot.Manifest)
	}

	if depot.Size != 3488450479 {
		t.Errorf("Expected Size 3488450479, got %d", depot.Size)
	}

	if depot.DLCAppID != "" {
		t.Errorf("Expected empty DLCAppID for base game depot, got %s", depot.DLCAppID)
	}
}

// TestInstalledDepotDLC tests DLC depot structure
func TestInstalledDepotDLC(t *testing.T) {
	depot := InstalledDepot{
		DepotID:  "1770481",
		Manifest: "4625979481897414804",
		Size:     57120063,
		DLCAppID: "1770480",
	}

	if depot.DepotID != "1770481" {
		t.Errorf("Expected DepotID 1770481, got %s", depot.DepotID)
	}

	if depot.Manifest != "4625979481897414804" {
		t.Errorf("Expected Manifest 4625979481897414804, got %s", depot.Manifest)
	}

	if depot.Size != 57120063 {
		t.Errorf("Expected Size 57120063, got %d", depot.Size)
	}

	if depot.DLCAppID == "" {
		t.Error("Expected DLCAppID to be populated for DLC depot")
	}

	if depot.DLCAppID != "1770480" {
		t.Errorf("Expected DLCAppID 1770480, got %s", depot.DLCAppID)
	}
}

// TestInstalledAppStruct tests the InstalledApp structure
func TestInstalledAppStruct(t *testing.T) {
	app := InstalledApp{
		AppID:       "667970",
		Name:        "VTOL VR",
		InstallDir:  "VTOL VR",
		FullPath:    "D:\\SteamLibrary\\steamapps\\common\\VTOL VR",
		BuildID:     "20275350",
		SizeOnDisk:  3606612198,
		LastUpdated: 1766280889,
		LastPlayed:  0,
		LibraryPath: "D:\\SteamLibrary",
		InstalledDepots: []InstalledDepot{
			{
				DepotID:  "667971",
				Manifest: "809774009354886606",
				Size:     3488450479,
				DLCAppID: "",
			},
		},
	}

	if app.AppID != "667970" {
		t.Errorf("Expected AppID 667970, got %s", app.AppID)
	}

	if app.Name != "VTOL VR" {
		t.Errorf("Expected Name 'VTOL VR', got %s", app.Name)
	}

	if app.InstallDir != "VTOL VR" {
		t.Errorf("Expected InstallDir 'VTOL VR', got %s", app.InstallDir)
	}

	if app.FullPath != "D:\\SteamLibrary\\steamapps\\common\\VTOL VR" {
		t.Errorf("Expected FullPath 'D:\\SteamLibrary\\steamapps\\common\\VTOL VR', got %s", app.FullPath)
	}

	if app.BuildID != "20275350" {
		t.Errorf("Expected BuildID '20275350', got %s", app.BuildID)
	}

	if app.SizeOnDisk != 3606612198 {
		t.Errorf("Expected SizeOnDisk 3606612198, got %d", app.SizeOnDisk)
	}

	if app.LastUpdated != 1766280889 {
		t.Errorf("Expected LastUpdated 1766280889, got %d", app.LastUpdated)
	}

	if app.LastPlayed != 0 {
		t.Errorf("Expected LastPlayed 0, got %d", app.LastPlayed)
	}

	if app.LibraryPath != "D:\\SteamLibrary" {
		t.Errorf("Expected LibraryPath 'D:\\SteamLibrary', got %s", app.LibraryPath)
	}

	if len(app.InstalledDepots) != 1 {
		t.Errorf("Expected 1 depot, got %d", len(app.InstalledDepots))
	}
}

// createTestManifest creates a test appmanifest file
func createTestManifest(t *testing.T, tempDir string, appID string) string {
	manifestContent := `"AppState"
{
	"appid"		"667970"
	"universe"		"1"
	"LauncherPath"		"C:\\Program Files (x86)\\Steam\\steam.exe"
	"name"		"VTOL VR"
	"StateFlags"		"4"
	"installdir"		"VTOL VR"
	"LastUpdated"		"1766280889"
	"LastPlayed"		"0"
	"SizeOnDisk"		"3606612198"
	"StagingSize"		"0"
	"buildid"		"20275350"
	"LastOwner"		"76561198096000762"
	"DownloadType"		"2"
	"UpdateResult"		"0"
	"BytesToDownload"		"0"
	"BytesDownloaded"		"0"
	"BytesToStage"		"0"
	"BytesStaged"		"0"
	"TargetBuildID"		"0"
	"AutoUpdateBehavior"		"0"
	"AllowOtherDownloadsWhileRunning"		"0"
	"ScheduledAutoUpdate"		"0"
	"FullValidateAfterNextUpdate"		"1"
	"InstalledDepots"
	{
		"667971"
		{
			"manifest"		"809774009354886606"
			"size"		"3488450479"
		}
		"1770481"
		{
			"manifest"		"4625979481897414804"
			"size"		"57120063"
			"dlcappid"		"1770480"
		}
		"2141700"
		{
			"manifest"		"1013815941531505868"
			"size"		"767090"
			"dlcappid"		"2141700"
		}
		"2531290"
		{
			"manifest"		"1746642230860805602"
			"size"		"60274566"
			"dlcappid"		"2531290"
		}
	}
	"UserConfig"
	{
		"language"		"english"
	}
	"MountedConfig"
	{
		"language"		"english"
	}
}`

	steamappsDir := filepath.Join(tempDir, "steamapps")
	err := os.MkdirAll(steamappsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create steamapps directory: %v", err)
	}

	manifestPath := filepath.Join(steamappsDir, "appmanifest_"+appID+".acf")
	err = os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write manifest file: %v", err)
	}

	return tempDir
}

// createTestManifestB creates a test appmanifest file
func createTestManifestB(b *testing.B, tempDir string, appID string) string {
	manifestContent := `"AppState"
{
	"appid"		"667970"
	"universe"		"1"
	"LauncherPath"		"C:\\Program Files (x86)\\Steam\\steam.exe"
	"name"		"VTOL VR"
	"StateFlags"		"4"
	"installdir"		"VTOL VR"
	"LastUpdated"		"1766280889"
	"LastPlayed"		"0"
	"SizeOnDisk"		"3606612198"
	"StagingSize"		"0"
	"buildid"		"20275350"
	"LastOwner"		"76561198096000762"
	"DownloadType"		"2"
	"UpdateResult"		"0"
	"BytesToDownload"		"0"
	"BytesDownloaded"		"0"
	"BytesToStage"		"0"
	"BytesStaged"		"0"
	"TargetBuildID"		"0"
	"AutoUpdateBehavior"		"0"
	"AllowOtherDownloadsWhileRunning"		"0"
	"ScheduledAutoUpdate"		"0"
	"FullValidateAfterNextUpdate"		"1"
	"InstalledDepots"
	{
		"667971"
		{
			"manifest"		"809774009354886606"
			"size"		"3488450479"
		}
		"1770481"
		{
			"manifest"		"4625979481897414804"
			"size"		"57120063"
			"dlcappid"		"1770480"
		}
		"2141700"
		{
			"manifest"		"1013815941531505868"
			"size"		"767090"
			"dlcappid"		"2141700"
		}
		"2531290"
		{
			"manifest"		"1746642230860805602"
			"size"		"60274566"
			"dlcappid"		"2531290"
		}
	}
	"UserConfig"
	{
		"language"		"english"
	}
	"MountedConfig"
	{
		"language"		"english"
	}
}`

	steamappsDir := filepath.Join(tempDir, "steamapps")
	err := os.MkdirAll(steamappsDir, 0755)
	if err != nil {
		b.Fatalf("Failed to create steamapps directory: %v", err)
	}

	manifestPath := filepath.Join(steamappsDir, "appmanifest_"+appID+".acf")
	err = os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	if err != nil {
		b.Fatalf("Failed to write manifest file: %v", err)
	}

	return tempDir
}

// TestReadAppManifest tests reading a manifest file
func TestReadAppManifest(t *testing.T) {
	tempDir := t.TempDir()
	libraryPath := createTestManifest(t, tempDir, "667970")

	app, err := readAppManifest(libraryPath, "667970")
	if err != nil {
		t.Fatalf("readAppManifest failed: %v", err)
	}

	// Test basic fields
	if app.AppID != "667970" {
		t.Errorf("Expected AppID 667970, got %s", app.AppID)
	}

	if app.Name != "VTOL VR" {
		t.Errorf("Expected Name 'VTOL VR', got %s", app.Name)
	}

	if app.InstallDir != "VTOL VR" {
		t.Errorf("Expected InstallDir 'VTOL VR', got %s", app.InstallDir)
	}

	if app.BuildID != "20275350" {
		t.Errorf("Expected BuildID 20275350, got %s", app.BuildID)
	}

	if app.SizeOnDisk != 3606612198 {
		t.Errorf("Expected SizeOnDisk 3606612198, got %d", app.SizeOnDisk)
	}

	if app.LastUpdated != 1766280889 {
		t.Errorf("Expected LastUpdated 1766280889, got %d", app.LastUpdated)
	}

	if app.LastPlayed != 0 {
		t.Errorf("Expected LastPlayed 0, got %d", app.LastPlayed)
	}

	// Test full path construction
	expectedPath := filepath.Join(libraryPath, "steamapps", "common", "VTOL VR")
	if app.FullPath != expectedPath {
		t.Errorf("Expected FullPath %s, got %s", expectedPath, app.FullPath)
	}
}

// TestReadAppManifestDepots tests depot parsing
func TestReadAppManifestDepots(t *testing.T) {
	tempDir := t.TempDir()
	libraryPath := createTestManifest(t, tempDir, "667970")

	app, err := readAppManifest(libraryPath, "667970")
	if err != nil {
		t.Fatalf("readAppManifest failed: %v", err)
	}

	// Test depot count
	expectedDepotCount := 4
	if len(app.InstalledDepots) != expectedDepotCount {
		t.Errorf("Expected %d depots, got %d", expectedDepotCount, len(app.InstalledDepots))
	}

	// Test base game depot (667971)
	foundBaseDepot := false
	for _, depot := range app.InstalledDepots {
		if depot.DepotID == "667971" {
			foundBaseDepot = true

			if depot.Manifest != "809774009354886606" {
				t.Errorf("Expected Manifest 809774009354886606, got %s", depot.Manifest)
			}

			if depot.Size != 3488450479 {
				t.Errorf("Expected Size 3488450479, got %d", depot.Size)
			}

			if depot.DLCAppID != "" {
				t.Errorf("Expected empty DLCAppID for base depot, got %s", depot.DLCAppID)
			}
			break
		}
	}

	if !foundBaseDepot {
		t.Error("Base game depot 667971 not found")
	}

	// Test DLC depot (1770481)
	foundDLCDepot := false
	for _, depot := range app.InstalledDepots {
		if depot.DepotID == "1770481" {
			foundDLCDepot = true

			if depot.Manifest != "4625979481897414804" {
				t.Errorf("Expected Manifest 4625979481897414804, got %s", depot.Manifest)
			}

			if depot.Size != 57120063 {
				t.Errorf("Expected Size 57120063, got %d", depot.Size)
			}

			if depot.DLCAppID != "1770480" {
				t.Errorf("Expected DLCAppID 1770480, got %s", depot.DLCAppID)
			}
			break
		}
	}

	if !foundDLCDepot {
		t.Error("DLC depot 1770481 not found")
	}

	// Count base vs DLC depots
	baseCount := 0
	dlcCount := 0
	for _, depot := range app.InstalledDepots {
		if depot.DLCAppID == "" {
			baseCount++
		} else {
			dlcCount++
		}
	}

	if baseCount != 1 {
		t.Errorf("Expected 1 base depot, got %d", baseCount)
	}

	if dlcCount != 3 {
		t.Errorf("Expected 3 DLC depots, got %d", dlcCount)
	}
}

// TestReadAppManifestNonExistent tests reading a non-existent manifest
func TestReadAppManifestNonExistent(t *testing.T) {
	tempDir := t.TempDir()

	_, err := readAppManifest(tempDir, "999999")
	if err == nil {
		t.Error("Expected error for non-existent manifest, got nil")
	}
}

// TestReadAppManifestCorrupted tests reading a corrupted manifest
func TestReadAppManifestCorrupted(t *testing.T) {
	tempDir := t.TempDir()

	steamappsDir := filepath.Join(tempDir, "steamapps")
	err := os.MkdirAll(steamappsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create steamapps directory: %v", err)
	}

	manifestPath := filepath.Join(steamappsDir, "appmanifest_667970.acf")
	err = os.WriteFile(manifestPath, []byte("corrupted data {{{"), 0644)
	if err != nil {
		t.Fatalf("Failed to write corrupted manifest: %v", err)
	}

	_, err = readAppManifest(tempDir, "667970")
	if err == nil {
		t.Error("Expected error for corrupted manifest, got nil")
	}
}

// createMockLibraryVdf creates a mock libraryfolders.vdf structure
func createMockLibraryVdf(_ *testing.T, libraryPath string) *orderedmap.OrderedMap {

	library0 := orderedmap.New()
	library0.Set("path", libraryPath)

	apps := orderedmap.New()
	apps.Set("667970", "3606612198")

	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	return root
}

// createMockLibraryVdfB creates a mock libraryfolders.vdf structure
func createMockLibraryVdfB(_ *testing.B, libraryPath string) *orderedmap.OrderedMap {

	library0 := orderedmap.New()
	library0.Set("path", libraryPath)

	apps := orderedmap.New()
	apps.Set("667970", "3606612198")

	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	return root
}

// TestGetAllInstalledApps tests retrieving all installed apps
func TestGetAllInstalledApps(t *testing.T) {
	tempDir := t.TempDir()
	libraryPath := createTestManifest(t, tempDir, "667970")

	steamReader := SteamReader{
		libraryVdfMap: createMockLibraryVdf(t, libraryPath),
	}

	apps, err := steamReader.GetAllInstalledApps()
	if err != nil {
		t.Fatalf("GetAllInstalledApps failed: %v", err)
	}

	if len(apps) != 1 {
		t.Errorf("Expected 1 app, got %d", len(apps))
	}

	if len(apps) > 0 {
		app := apps[0]
		if app.AppID != "667970" {
			t.Errorf("Expected AppID 667970, got %s", app.AppID)
		}

		if app.Name != "VTOL VR" {
			t.Errorf("Expected Name 'VTOL VR', got %s", app.Name)
		}

		if len(app.InstalledDepots) != 4 {
			t.Errorf("Expected 4 depots, got %d", len(app.InstalledDepots))
		}
	}
}

// TestGetAllInstalledAppsEmpty tests with no apps
func TestGetAllInstalledAppsEmpty(t *testing.T) {
	tempDir := t.TempDir()

	library0 := orderedmap.New()
	library0.Set("path", tempDir)
	apps := orderedmap.New()
	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	apps2, err := steamReader.GetAllInstalledApps()
	if err != nil {
		t.Fatalf("GetAllInstalledApps failed: %v", err)
	}

	if len(apps2) != 0 {
		t.Errorf("Expected 0 apps, got %d", len(apps2))
	}
}

// TestGetInstalledAppByID tests retrieving a specific app
func TestGetInstalledAppByID(t *testing.T) {
	tempDir := t.TempDir()
	libraryPath := createTestManifest(t, tempDir, "667970")

	steamReader := SteamReader{
		libraryVdfMap: createMockLibraryVdf(t, libraryPath),
	}

	app, err := steamReader.GetInstalledAppByID("667970")
	if err != nil {
		t.Fatalf("GetInstalledAppByID failed: %v", err)
	}

	if app.AppID != "667970" {
		t.Errorf("Expected AppID 667970, got %s", app.AppID)
	}

	if app.Name != "VTOL VR" {
		t.Errorf("Expected Name 'VTOL VR', got %s", app.Name)
	}

	if app.LibraryPath != libraryPath {
		t.Errorf("Expected LibraryPath %s, got %s", libraryPath, app.LibraryPath)
	}
}

// TestGetInstalledAppByIDNotFound tests retrieving non-existent app
func TestGetInstalledAppByIDNotFound(t *testing.T) {
	tempDir := t.TempDir()

	steamReader := SteamReader{
		libraryVdfMap: createMockLibraryVdf(t, tempDir),
	}

	_, err := steamReader.GetInstalledAppByID("999999")
	if err == nil {
		t.Error("Expected error for non-existent app, got nil")
	}
}

// TestDepotSizeCalculation tests depot size calculations
func TestDepotSizeCalculation(t *testing.T) {
	tempDir := t.TempDir()
	libraryPath := createTestManifest(t, tempDir, "667970")

	app, err := readAppManifest(libraryPath, "667970")
	if err != nil {
		t.Fatalf("readAppManifest failed: %v", err)
	}

	var totalDepotSize int64 = 0
	for _, depot := range app.InstalledDepots {
		totalDepotSize += depot.Size
	}

	// Total depot sizes from the test manifest
	expectedTotal := int64(3488450479 + 57120063 + 767090 + 60274566)
	if totalDepotSize != expectedTotal {
		t.Errorf("Expected total depot size %d, got %d", expectedTotal, totalDepotSize)
	}
}

// TestDLCIdentification tests DLC depot identification
func TestDLCIdentification(t *testing.T) {
	tempDir := t.TempDir()
	libraryPath := createTestManifest(t, tempDir, "667970")

	app, err := readAppManifest(libraryPath, "667970")
	if err != nil {
		t.Fatalf("readAppManifest failed: %v", err)
	}

	dlcDepots := make(map[string]bool)
	for _, depot := range app.InstalledDepots {
		if depot.DLCAppID != "" {
			dlcDepots[depot.DLCAppID] = true
		}
	}

	expectedDLCs := map[string]bool{
		"1770480": true,
		"2141700": true,
		"2531290": true,
	}

	if len(dlcDepots) != len(expectedDLCs) {
		t.Errorf("Expected %d unique DLCs, got %d", len(expectedDLCs), len(dlcDepots))
	}

	for dlcID := range expectedDLCs {
		if !dlcDepots[dlcID] {
			t.Errorf("Expected DLC %s not found", dlcID)
		}
	}
}

// BenchmarkReadAppManifest benchmarks manifest reading
func BenchmarkReadAppManifest(b *testing.B) {
	tempDir := b.TempDir()
	libraryPath := createTestManifestB(b, tempDir, "667970")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := readAppManifest(libraryPath, "667970")
		if err != nil {
			b.Fatalf("readAppManifest failed: %v", err)
		}
	}
}

// BenchmarkGetAllInstalledApps benchmarks getting all apps
func BenchmarkGetAllInstalledApps(b *testing.B) {
	tempDir := b.TempDir()
	libraryPath := createTestManifestB(b, tempDir, "667970")

	steamReader := SteamReader{
		libraryVdfMap: createMockLibraryVdfB(b, libraryPath),
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := steamReader.GetAllInstalledApps()
		if err != nil {
			b.Fatalf("GetAllInstalledApps failed: %v", err)
		}
	}
}
