package steamutils

import (
	"strings"
	"testing"

	"github.com/iancoleman/orderedmap"
)

// TestUnmarshalSimple tests unmarshaling a simple VDF string
func TestUnmarshalSimple(t *testing.T) {
	vdfData := `"AppState"
{
	"appid"		"667970"
	"name"		"VTOL VR"
}`

	result, err := Unmarshal([]byte(vdfData))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	appStateVal, exists := result.Get("AppState")
	if !exists {
		t.Error("AppState key not found")
	}

	appState, ok := appStateVal.(*orderedmap.OrderedMap)
	if !ok {
		t.Error("AppState is not an OrderedMap")
	}

	appIDVal, exists := appState.Get("appid")
	if !exists {
		t.Error("appid key not found")
	}

	if appIDVal != "667970" {
		t.Errorf("Expected appid 667970, got %s", appIDVal)
	}

	nameVal, exists := appState.Get("name")
	if !exists {
		t.Error("name key not found")
	}

	if nameVal != "VTOL VR" {
		t.Errorf("Expected name 'VTOL VR', got %s", nameVal)
	}
}

// TestUnmarshalNested tests unmarshaling nested VDF structures
func TestUnmarshalNested(t *testing.T) {
	vdfData := `"AppState"
{
	"appid"		"667970"
	"InstalledDepots"
	{
		"667971"
		{
			"manifest"		"809774009354886606"
			"size"		"3488450479"
		}
	}
}`

	result, err := Unmarshal([]byte(vdfData))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	appStateVal, _ := result.Get("AppState")
	appState := appStateVal.(*orderedmap.OrderedMap)

	depotsVal, exists := appState.Get("InstalledDepots")
	if !exists {
		t.Error("InstalledDepots key not found")
	}

	depots, ok := depotsVal.(*orderedmap.OrderedMap)
	if !ok {
		t.Error("InstalledDepots is not an OrderedMap")
	}

	depotVal, exists := depots.Get("667971")
	if !exists {
		t.Error("Depot 667971 not found")
	}

	depot, ok := depotVal.(*orderedmap.OrderedMap)
	if !ok {
		t.Error("Depot is not an OrderedMap")
	}

	manifestVal, exists := depot.Get("manifest")
	if !exists {
		t.Error("manifest key not found")
	}

	if manifestVal != "809774009354886606" {
		t.Errorf("Expected manifest 809774009354886606, got %s", manifestVal)
	}
}

// TestUnmarshalMultipleDepots tests unmarshaling multiple depots
func TestUnmarshalMultipleDepots(t *testing.T) {
	vdfData := `"AppState"
{
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
	}
}`

	result, err := Unmarshal([]byte(vdfData))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	appStateVal, _ := result.Get("AppState")
	appState := appStateVal.(*orderedmap.OrderedMap)

	depotsVal, _ := appState.Get("InstalledDepots")
	depots := depotsVal.(*orderedmap.OrderedMap)

	// Check depot count
	if len(depots.Keys()) != 3 {
		t.Errorf("Expected 3 depots, got %d", len(depots.Keys()))
	}

	// Verify each depot
	depotIDs := []string{"667971", "1770481", "2141700"}
	for _, depotID := range depotIDs {
		_, exists := depots.Get(depotID)
		if !exists {
			t.Errorf("Depot %s not found", depotID)
		}
	}
}

// TestUnmarshalWhitespace tests handling of various whitespace
func TestUnmarshalWhitespace(t *testing.T) {
	vdfData := `  "AppState"  
	{
		"appid"			"667970"
		"name"		"VTOL VR"  
	}  `

	result, err := Unmarshal([]byte(vdfData))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	appStateVal, exists := result.Get("AppState")
	if !exists {
		t.Error("AppState key not found")
	}

	appState := appStateVal.(*orderedmap.OrderedMap)
	appIDVal, _ := appState.Get("appid")

	if appIDVal != "667970" {
		t.Errorf("Expected appid 667970, got %s", appIDVal)
	}
}

// TestUnmarshalEmpty tests unmarshaling empty VDF
func TestUnmarshalEmpty(t *testing.T) {
	vdfData := `"AppState"
{
}`

	result, err := Unmarshal([]byte(vdfData))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	appStateVal, exists := result.Get("AppState")
	if !exists {
		t.Error("AppState key not found")
	}

	appState := appStateVal.(*orderedmap.OrderedMap)
	if len(appState.Keys()) != 0 {
		t.Errorf("Expected empty AppState, got %d keys", len(appState.Keys()))
	}
}

// TestUnmarshalInvalidNoQuotes tests invalid VDF without quotes
func TestUnmarshalInvalidNoQuotes(t *testing.T) {
	vdfData := `AppState
{
	appid		667970
}`

	_, err := Unmarshal([]byte(vdfData))
	if err == nil {
		t.Error("Expected error for VDF without quotes, got nil")
	}
}

// TestUnmarshalInvalidUnterminatedString tests unterminated string
func TestUnmarshalInvalidUnterminatedString(t *testing.T) {
	vdfData := `"AppState
{
	"appid"		"667970"
}`

	_, err := Unmarshal([]byte(vdfData))
	if err == nil {
		t.Error("Expected error for unterminated string, got nil")
	}
}

// TestUnmarshalInvalidMissingBrace tests missing closing brace
func TestUnmarshalInvalidMissingBrace(t *testing.T) {
	vdfData := `"AppState"
{
	"appid"		"667970"
`

	result, err := Unmarshal([]byte(vdfData))
	// Should not error, just reach end of input
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result == nil {
		t.Error("Result should not be nil")
	}
}

// TestMarshalSimple tests marshaling a simple structure
func TestMarshalSimple(t *testing.T) {
	root := orderedmap.New()
	appState := orderedmap.New()
	appState.Set("appid", "667970")
	appState.Set("name", "VTOL VR")
	root.Set("AppState", appState)

	data, err := Marshal(root)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	result := string(data)

	// Check that output contains expected values
	if !strings.Contains(result, "AppState") {
		t.Error("Marshal output missing AppState")
	}

	if !strings.Contains(result, "667970") {
		t.Error("Marshal output missing appid value")
	}

	if !strings.Contains(result, "VTOL VR") {
		t.Error("Marshal output missing name value")
	}
}

// TestMarshalNested tests marshaling nested structures
func TestMarshalNested(t *testing.T) {
	root := orderedmap.New()
	appState := orderedmap.New()

	depots := orderedmap.New()
	depot1 := orderedmap.New()
	depot1.Set("manifest", "809774009354886606")
	depot1.Set("size", "3488450479")
	depots.Set("667971", depot1)

	appState.Set("appid", "667970")
	appState.Set("InstalledDepots", depots)
	root.Set("AppState", appState)

	data, err := Marshal(root)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	result := string(data)

	// Check for nested structure
	if !strings.Contains(result, "InstalledDepots") {
		t.Error("Marshal output missing InstalledDepots")
	}

	if !strings.Contains(result, "667971") {
		t.Error("Marshal output missing depot ID")
	}

	if !strings.Contains(result, "809774009354886606") {
		t.Error("Marshal output missing manifest")
	}
}

// TestUnmarshalMarshalRoundtrip tests roundtrip conversion
func TestUnmarshalMarshalRoundtrip(t *testing.T) {
	originalVDF := `"AppState"
{
	"appid"		"667970"
	"name"		"VTOL VR"
	"buildid"		"20275350"
}`

	// Unmarshal
	unmarshaled, err := Unmarshal([]byte(originalVDF))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	// Marshal
	marshaled, err := Marshal(unmarshaled)
	if err != nil {
		t.Fatalf("Marshal failed: %v", err)
	}

	// Unmarshal again
	remarshaled, err := Unmarshal(marshaled)
	if err != nil {
		t.Fatalf("Second unmarshal failed: %v", err)
	}

	// Verify data integrity
	appStateVal, exists := remarshaled.Get("AppState")
	if !exists {
		t.Error("AppState not found after roundtrip")
	}

	appState := appStateVal.(*orderedmap.OrderedMap)

	if val, _ := appState.Get("appid"); val != "667970" {
		t.Error("appid value changed after roundtrip")
	}

	if val, _ := appState.Get("name"); val != "VTOL VR" {
		t.Error("name value changed after roundtrip")
	}

	if val, _ := appState.Get("buildid"); val != "20275350" {
		t.Error("buildid value changed after roundtrip")
	}
}

// TestSkipWhitespace tests the skipWhitespace function
func TestSkipWhitespace(t *testing.T) {
	testCases := []struct {
		input    string
		start    int
		expected int
	}{
		{"   abc", 0, 3},
		{"abc", 0, 0},
		{"\t\n  abc", 0, 4},
		{"abc   def", 3, 6},
	}

	for _, tc := range testCases {
		result := skipWhitespace(tc.input, tc.start)
		if result != tc.expected {
			t.Errorf("For input '%s' at position %d, expected %d, got %d",
				tc.input, tc.start, tc.expected, result)
		}
	}
}

// TestParseString tests the parseString function
func TestParseString(t *testing.T) {
	testCases := []struct {
		input         string
		start         int
		expectedStr   string
		expectedIndex int
		expectError   bool
	}{
		{`"hello"`, 0, "hello", 7, false},
		{`"VTOL VR"`, 0, "VTOL VR", 9, false},
		{`"667970"`, 0, "667970", 8, false},
		{`abc"test"`, 3, "test", 9, false},
		{`"`, 0, "", 0, true},   // Unterminated
		{`abc`, 0, "", 0, true}, // No starting quote
	}

	for _, tc := range testCases {
		str, idx, err := parseString(tc.input, tc.start)

		if tc.expectError {
			if err == nil {
				t.Errorf("Expected error for input '%s', got nil", tc.input)
			}
			continue
		}

		if err != nil {
			t.Errorf("Unexpected error for input '%s': %v", tc.input, err)
			continue
		}

		if str != tc.expectedStr {
			t.Errorf("For input '%s', expected string '%s', got '%s'",
				tc.input, tc.expectedStr, str)
		}

		if idx != tc.expectedIndex {
			t.Errorf("For input '%s', expected index %d, got %d",
				tc.input, tc.expectedIndex, idx)
		}
	}
}

// TestUnmarshalLibraryFolders tests unmarshaling libraryfolders.vdf format
func TestUnmarshalLibraryFolders(t *testing.T) {
	vdfData := `"libraryfolders"
{
	"0"
	{
		"path"		"C:\\Program Files (x86)\\Steam"
		"label"		""
		"contentid"		"1234567890"
		"totalsize"		"0"
		"apps"
		{
			"667970"		"3606612198"
			"440"		"1000000"
		}
	}
	"1"
	{
		"path"		"D:\\SteamLibrary"
		"apps"
		{
			"730"		"5000000"
		}
	}
}`

	result, err := Unmarshal([]byte(vdfData))
	if err != nil {
		t.Fatalf("Unmarshal failed: %v", err)
	}

	libFoldersVal, exists := result.Get("libraryfolders")
	if !exists {
		t.Error("libraryfolders key not found")
	}

	libFolders := libFoldersVal.(*orderedmap.OrderedMap)

	// Check library 0
	lib0Val, exists := libFolders.Get("0")
	if !exists {
		t.Error("Library 0 not found")
	}

	lib0 := lib0Val.(*orderedmap.OrderedMap)
	pathVal, _ := lib0.Get("path")
	if pathVal != "C:\\Program Files (x86)\\Steam" {
		t.Errorf("Unexpected path for library 0: %s", pathVal)
	}

	// Check apps in library 0
	appsVal, exists := lib0.Get("apps")
	if !exists {
		t.Error("apps not found in library 0")
	}

	apps := appsVal.(*orderedmap.OrderedMap)
	app667970Val, exists := apps.Get("667970")
	if !exists {
		t.Error("App 667970 not found in library 0")
	}

	if app667970Val != "3606612198" {
		t.Errorf("Expected app size 3606612198, got %s", app667970Val)
	}

	// Check library 1
	lib1Val, exists := libFolders.Get("1")
	if !exists {
		t.Error("Library 1 not found")
	}

	lib1 := lib1Val.(*orderedmap.OrderedMap)
	path1Val, _ := lib1.Get("path")
	if path1Val != "D:\\SteamLibrary" {
		t.Errorf("Unexpected path for library 1: %s", path1Val)
	}
}

// BenchmarkUnmarshal benchmarks unmarshaling
func BenchmarkUnmarshal(b *testing.B) {
	vdfData := []byte(`"AppState"
{
	"appid"		"667970"
	"name"		"VTOL VR"
	"buildid"		"20275350"
	"InstalledDepots"
	{
		"667971"
		{
			"manifest"		"809774009354886606"
			"size"		"3488450479"
		}
	}
}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Unmarshal(vdfData)
		if err != nil {
			b.Fatalf("Unmarshal failed: %v", err)
		}
	}
}

// BenchmarkMarshal benchmarks marshaling
func BenchmarkMarshal(b *testing.B) {
	root := orderedmap.New()
	appState := orderedmap.New()
	appState.Set("appid", "667970")
	appState.Set("name", "VTOL VR")
	appState.Set("buildid", "20275350")

	depots := orderedmap.New()
	depot := orderedmap.New()
	depot.Set("manifest", "809774009354886606")
	depot.Set("size", "3488450479")
	depots.Set("667971", depot)

	appState.Set("InstalledDepots", depots)
	root.Set("AppState", appState)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Marshal(root)
		if err != nil {
			b.Fatalf("Marshal failed: %v", err)
		}
	}
}

// The code in this file was made by an LLM, use in production is highly discouraged as unexpected results may occur. The code in this file is not vetted for stability or edge cases.
