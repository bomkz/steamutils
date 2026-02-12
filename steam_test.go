package steamutils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/iancoleman/orderedmap"
)

// TestSteamReaderConfig tests the SteamReaderConfig structure
func TestSteamReaderConfig(t *testing.T) {
	config := SteamReaderConfig{
		CustomSteamPath:      "C:\\Program Files (x86)\\Steam",
		CustomLibraryVdfPath: "C:\\Program Files (x86)\\Steam\\steamapps\\libraryfolders.vdf",
		FormatSteamPath:      true,
	}

	if config.CustomSteamPath == "" {
		t.Error("CustomSteamPath should not be empty")
	}

	if config.CustomLibraryVdfPath == "" {
		t.Error("CustomLibraryVdfPath should not be empty")
	}

	if !config.FormatSteamPath {
		t.Error("FormatSteamPath should be true")
	}
}

// TestCheckDefaultLibraryPath tests the checkDefaultLibraryPath function
func TestCheckDefaultLibraryPath(t *testing.T) {
	tempDir := t.TempDir()
	steamappsDir := filepath.Join(tempDir, "steamapps")
	err := os.MkdirAll(steamappsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create steamapps directory: %v", err)
	}

	libraryVdfPath := filepath.Join(steamappsDir, "libraryfolders.vdf")
	err = os.WriteFile(libraryVdfPath, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create libraryfolders.vdf: %v", err)
	}

	result, err := checkDefaultLibraryPath(tempDir)
	if err != nil {
		t.Errorf("checkDefaultLibraryPath failed: %v", err)
	}

	expectedPath := tempDir + "\\steamapps\\libraryfolders.vdf"
	if result != expectedPath {
		t.Errorf("Expected path %s, got %s", expectedPath, result)
	}
}

// TestCheckDefaultLibraryPathNotFound tests when libraryfolders.vdf doesn't exist
func TestCheckDefaultLibraryPathNotFound(t *testing.T) {
	tempDir := t.TempDir()

	_, err := checkDefaultLibraryPath(tempDir)
	if err == nil {
		t.Error("Expected error when libraryfolders.vdf doesn't exist, got nil")
	}
}

// TestFindAppIDPath tests finding the library path for an app
func TestFindAppIDPath(t *testing.T) {
	// Create mock library structure

	library0 := orderedmap.New()
	library0.Set("path", "D:\\SteamLibrary")

	apps := orderedmap.New()
	apps.Set("667970", "3606612198")
	apps.Set("440", "1000000")

	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	// Test finding existing app
	path, err := steamReader.FindAppIDPath("667970")
	if err != nil {
		t.Errorf("FindAppIDPath failed: %v", err)
	}

	if path != "D:\\SteamLibrary" {
		t.Errorf("Expected path 'D:\\SteamLibrary', got %s", path)
	}
}

// TestFindAppIDPathNotFound tests finding non-existent app
func TestFindAppIDPathNotFound(t *testing.T) {
	library0 := orderedmap.New()
	library0.Set("path", "D:\\SteamLibrary")
	apps := orderedmap.New()
	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	_, err := steamReader.FindAppIDPath("999999")
	if err == nil {
		t.Error("Expected error for non-existent app, got nil")
	}
}

// TestFindAppIDPathMultipleLibraries tests searching across multiple libraries
func TestFindAppIDPathMultipleLibraries(t *testing.T) {

	// Library 0
	library0 := orderedmap.New()
	library0.Set("path", "C:\\Program Files (x86)\\Steam")
	apps0 := orderedmap.New()
	apps0.Set("440", "1000000")
	library0.Set("apps", apps0)

	// Library 1
	library1 := orderedmap.New()
	library1.Set("path", "D:\\SteamLibrary")
	apps1 := orderedmap.New()
	apps1.Set("667970", "3606612198")
	library1.Set("apps", apps1)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)
	libFolders.Set("1", library1)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	// Find app in second library
	path, err := steamReader.FindAppIDPath("667970")
	if err != nil {
		t.Errorf("FindAppIDPath failed: %v", err)
	}

	if path != "D:\\SteamLibrary" {
		t.Errorf("Expected path 'D:\\SteamLibrary', got %s", path)
	}

	// Find app in first library
	path, err = steamReader.FindAppIDPath("440")
	if err != nil {
		t.Errorf("FindAppIDPath failed: %v", err)
	}

	if path != "C:\\Program Files (x86)\\Steam" {
		t.Errorf("Expected path 'C:\\Program Files (x86)\\Steam', got %s", path)
	}
}

// TestFindAppIDBuildID tests finding build ID for an app
func TestFindAppIDBuildID(t *testing.T) {
	tempDir := t.TempDir()
	steamappsDir := filepath.Join(tempDir, "steamapps")
	err := os.MkdirAll(steamappsDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create steamapps directory: %v", err)
	}

	// Create test manifest
	manifestContent := `"AppState"
{
	"appid"		"667970"
	"buildid"		"20275350"
	"name"		"VTOL VR"
}`

	manifestPath := filepath.Join(steamappsDir, "appmanifest_667970.acf")
	err = os.WriteFile(manifestPath, []byte(manifestContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write manifest file: %v", err)
	}

	// Create mock library structure
	library0 := orderedmap.New()
	library0.Set("path", tempDir)
	apps := orderedmap.New()
	apps.Set("667970", "3606612198")
	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	buildID, err := steamReader.FindAppIDBuildID("667970")
	if err != nil {
		t.Fatalf("FindAppIDBuildID failed: %v", err)
	}

	if buildID != "20275350" {
		t.Errorf("Expected buildID 20275350, got %s", buildID)
	}
}

// TestFindAppIDBuildIDNotFound tests when app or manifest doesn't exist
func TestFindAppIDBuildIDNotFound(t *testing.T) {
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

	_, err := steamReader.FindAppIDBuildID("999999")
	if err == nil {
		t.Error("Expected error for non-existent app, got nil")
	}
}

// TestGetLibraryVdfMap tests getting the library VDF map
func TestGetLibraryVdfMap(t *testing.T) {
	testMap := orderedmap.New()
	testMap.Set("test", "value")

	steamReader := SteamReader{
		libraryVdfMap: testMap,
	}

	result := steamReader.GetLibraryVdfMap()
	if result == nil {
		t.Error("GetLibraryVdfMap returned nil")
	}

	val, exists := result.Get("test")
	if !exists {
		t.Error("Expected 'test' key to exist")
	}

	if val != "value" {
		t.Errorf("Expected value 'value', got %s", val)
	}
}

// TestGetSteamPath tests getting the Steam path
func TestGetSteamPath(t *testing.T) {
	steamReader := SteamReader{
		steamPath: "c:\\program files (x86)\\steam",
		SteamReaderConfig: SteamReaderConfig{
			customSteamPathFinder: false,
		},
	}

	path := steamReader.GetSteamPath()

	// Should format the path properly
	if path == "" {
		t.Error("GetSteamPath returned empty string")
	}

	// Should start with uppercase drive letter
	if len(path) > 0 && path[0] < 'A' || path[0] > 'Z' {
		t.Errorf("Expected path to start with uppercase letter, got %c", path[0])
	}
}

// TestGetSteamPathCustom tests getting custom Steam path
func TestGetSteamPathCustom(t *testing.T) {
	customPath := "D:\\Games\\Steam"
	steamReader := SteamReader{
		steamPath: customPath,
		SteamReaderConfig: SteamReaderConfig{
			customSteamPathFinder: true,
		},
	}

	path := steamReader.GetSteamPath()
	if path != customPath {
		t.Errorf("Expected path %s, got %s", customPath, path)
	}
}

// TestGetLibraryVdfPath tests getting library VDF path
func TestGetLibraryVdfPath(t *testing.T) {
	testPath := "C:\\Program Files (x86)\\Steam\\steamapps\\libraryfolders.vdf"
	steamReader := SteamReader{
		libraryVdfPath: testPath,
	}

	path := steamReader.GetLibraryVdfPath()
	if path != testPath {
		t.Errorf("Expected path %s, got %s", testPath, path)
	}
}

// TestSteamReaderStruct tests the SteamReader structure
func TestSteamReaderStruct(t *testing.T) {
	libraryMap := orderedmap.New()
	libraryMap.Set("test", "data")

	reader := SteamReader{
		libraryVdfPath: "C:\\Steam\\steamapps\\libraryfolders.vdf",
		steamPath:      "C:\\Steam",
		libraryVdfMap:  libraryMap,
		SteamReaderConfig: SteamReaderConfig{
			FormatSteamPath: true,
		},
	}

	if reader.libraryVdfPath == "" {
		t.Error("libraryVdfPath should not be empty")
	}

	if reader.steamPath == "" {
		t.Error("steamPath should not be empty")
	}

	if reader.libraryVdfMap == nil {
		t.Error("libraryVdfMap should not be nil")
	}

	if !reader.SteamReaderConfig.FormatSteamPath {
		t.Error("FormatSteamPath should be true")
	}
}

// TestFindAppIDPathDoubleBackslash tests handling of double backslashes
func TestFindAppIDPathDoubleBackslash(t *testing.T) {

	library0 := orderedmap.New()
	library0.Set("path", "D:\\\\SteamLibrary") // Double backslash

	apps := orderedmap.New()
	apps.Set("667970", "3606612198")

	library0.Set("apps", apps)

	libFolders := orderedmap.New()
	libFolders.Set("0", library0)

	root := orderedmap.New()
	root.Set("libraryfolders", libFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	path, err := steamReader.FindAppIDPath("667970")
	if err != nil {
		t.Errorf("FindAppIDPath failed: %v", err)
	}

	// Should replace double backslashes with single
	if path != "D:\\SteamLibrary" {
		t.Errorf("Expected path 'D:\\SteamLibrary', got %s", path)
	}
}

// TestGetSteamPathFormatting tests path formatting with different cases
func TestGetSteamPathFormatting(t *testing.T) {
	testCases := []struct {
		input    string
		expected string
	}{
		{
			input:    "c:\\program files (x86)\\steam",
			expected: "C:\\Program Files (x86)\\Steam\\",
		},
		{
			input:    "d:\\games\\steam",
			expected: "D:\\games\\Steam\\",
		},
	}

	for _, tc := range testCases {
		steamReader := SteamReader{
			steamPath: tc.input,
			SteamReaderConfig: SteamReaderConfig{
				customSteamPathFinder: false,
			},
		}

		result := steamReader.GetSteamPath()
		if result != tc.expected {
			t.Errorf("For input %s, expected %s, got %s", tc.input, tc.expected, result)
		}
	}
}

// BenchmarkFindAppIDPath benchmarks finding app path
func BenchmarkFindAppIDPath(b *testing.B) {
	libraryFolders := orderedmap.New()

	for i := 0; i < 5; i++ {
		library := orderedmap.New()
		library.Set("path", "D:\\SteamLibrary")

		apps := orderedmap.New()
		for j := 0; j < 100; j++ {
			apps.Set(string(rune(i*100+j)), "1000000")
		}
		apps.Set("667970", "3606612198")

		library.Set("apps", apps)
		libraryFolders.Set(string(rune(i)), library)
	}

	root := orderedmap.New()
	root.Set("libraryfolders", libraryFolders)

	steamReader := SteamReader{
		libraryVdfMap: root,
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = steamReader.FindAppIDPath("667970")
	}
}
