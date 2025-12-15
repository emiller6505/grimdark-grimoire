package parser

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewParser(t *testing.T) {
	p := NewParser("testdata")
	if p == nil {
		t.Fatal("NewParser returned nil")
	}
	if p.dataDir != "testdata" {
		t.Errorf("Expected dataDir 'testdata', got '%s'", p.dataDir)
	}
}

func TestLoadGameSystem(t *testing.T) {
	// Use actual data directory if available, otherwise skip
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	p := NewParser(dataDir)
	err := p.LoadGameSystem()
	if err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}

	gameSystem := p.GetGameSystem()
	if gameSystem == nil {
		t.Fatal("Game system is nil after loading")
	}

	if gameSystem.Name == "" {
		t.Error("Game system name is empty")
	}

	if gameSystem.ID == "" {
		t.Error("Game system ID is empty")
	}

	if len(gameSystem.ProfileTypes) == 0 {
		t.Error("No profile types loaded")
	}

	// Verify expected profile types
	profileTypeNames := make(map[string]bool)
	for _, pt := range gameSystem.ProfileTypes {
		profileTypeNames[pt.Name] = true
	}

	expectedTypes := []string{"Unit", "Ranged Weapons", "Melee Weapons", "Abilities", "Transport"}
	for _, expected := range expectedTypes {
		if !profileTypeNames[expected] {
			t.Errorf("Missing expected profile type: %s", expected)
		}
	}
}

func TestLoadCatalogue(t *testing.T) {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	p := NewParser(dataDir)
	
	// Try to load a known catalogue file
	catalogueFile := filepath.Join(dataDir, "Aeldari - Craftworlds.cat")
	if _, err := os.Stat(catalogueFile); os.IsNotExist(err) {
		t.Skipf("Catalogue file not found: %s", catalogueFile)
	}

	err := p.LoadCatalogue(catalogueFile)
	if err != nil {
		t.Fatalf("Failed to load catalogue: %v", err)
	}

	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Error("No catalogues loaded")
	}
}

func TestFindSelectionEntryByID(t *testing.T) {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	p := NewParser(dataDir)
	
	// Load game system and catalogues
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	// Test with known unit ID (Asurmen from Aeldari Library)
	unitID := "828d-840a-9a67-9074"
	entry, catalogueID, found := p.FindSelectionEntryByID(unitID)
	
	if !found {
		t.Fatalf("Unit %s not found", unitID)
	}

	if entry == nil {
		t.Fatal("Entry is nil")
	}

	if entry.ID != unitID {
		t.Errorf("Expected ID %s, got %s", unitID, entry.ID)
	}

	if entry.Name != "Asurmen" {
		t.Errorf("Expected name 'Asurmen', got '%s'", entry.Name)
	}

	if catalogueID == "" {
		t.Error("Catalogue ID is empty")
	}

	// Verify entry has profiles
	if len(entry.Profiles) == 0 {
		t.Error("Entry has no profiles")
	}

	// Verify entry has costs
	if len(entry.Costs) == 0 {
		t.Error("Entry has no costs")
	}
}

func TestFindEntryLinkByID(t *testing.T) {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	// Test with known entryLink ID (Warlock from Craftworlds)
	entryLinkID := "a502-4dbe-d0c6-69fd"
	entryLink, catalogueID, found := p.FindEntryLinkByID(entryLinkID)
	
	if !found {
		t.Fatalf("EntryLink %s not found", entryLinkID)
	}

	if entryLink == nil {
		t.Fatal("EntryLink is nil")
	}

	if entryLink.ID != entryLinkID {
		t.Errorf("Expected ID %s, got %s", entryLinkID, entryLink.ID)
	}

	if entryLink.Name != "Warlock" {
		t.Errorf("Expected name 'Warlock', got '%s'", entryLink.Name)
	}

	if entryLink.TargetID == "" {
		t.Error("TargetID is empty")
	}

	if catalogueID == "" {
		t.Error("Catalogue ID is empty")
	}

	// Verify entryLink has costs (from the XML we saw)
	if len(entryLink.Costs) == 0 {
		t.Error("EntryLink has no costs")
	}
}

func TestGetAllCatalogues(t *testing.T) {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Error("No catalogues loaded")
	}

	// Verify we have some expected catalogues
	catalogueNames := make(map[string]bool)
	for _, cat := range catalogues {
		catalogueNames[cat.Name] = true
	}

	// Check for at least one known catalogue
	hasMajorFaction := false
	for name := range catalogueNames {
		if len(name) > 0 {
			hasMajorFaction = true
			break
		}
	}

	if !hasMajorFaction {
		t.Error("Expected to find at least one catalogue")
	}
}

func TestGetAllLibraries(t *testing.T) {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	libraries := p.GetAllLibraries()
	if len(libraries) == 0 {
		t.Error("No libraries loaded")
	}

	// Verify we have library files
	hasLibrary := false
	for _, lib := range libraries {
		if lib.Library == "true" {
			hasLibrary = true
			break
		}
	}

	if !hasLibrary {
		t.Error("Expected to find at least one library file")
	}
}


