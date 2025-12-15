package parser

import (
	"testing"

	"grimoire-api/internal/models"
)

func TestTransformUnit(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)
	transformer := NewTransformer(resolver)

	// Get a known unit (Asurmen)
	unitID := "828d-840a-9a67-9074"
	entry, catalogueID, found := p.FindSelectionEntryByID(unitID)
	if !found {
		t.Fatalf("Unit %s not found", unitID)
	}

	// Transform
	unit := transformer.TransformUnit(entry, catalogueID)

	if unit == nil {
		t.Fatal("Transformed unit is nil")
	}

	// Verify basic fields
	if unit.ID != unitID {
		t.Errorf("Expected ID %s, got %s", unitID, unit.ID)
	}

	if unit.Name != "Asurmen" {
		t.Errorf("Expected name 'Asurmen', got '%s'", unit.Name)
	}

	if unit.Type != "model" {
		t.Errorf("Expected type 'model', got '%s'", unit.Type)
	}

	// Verify profiles
	if unit.Profiles == nil {
		t.Fatal("Unit profiles is nil")
	}

	if unit.Profiles.Unit == nil {
		t.Error("Unit profile is nil")
	} else {
		// Verify unit profile fields
		if unit.Profiles.Unit.Movement == "" {
			t.Error("Unit movement is empty")
		}
		if unit.Profiles.Unit.Toughness == 0 {
			t.Error("Unit toughness is 0")
		}
		if unit.Profiles.Unit.Wounds == 0 {
			t.Error("Unit wounds is 0")
		}
	}

	if len(unit.Profiles.Abilities) == 0 {
		t.Error("Unit has no abilities")
	}

	// Verify weapons
	if unit.Weapons == nil {
		t.Fatal("Unit weapons is nil")
	}

	// Asurmen should have both ranged and melee weapons
	if len(unit.Weapons.Ranged) == 0 && len(unit.Weapons.Melee) == 0 {
		t.Error("Unit has no weapons")
	}

	// Verify costs
	if unit.Costs == nil {
		t.Fatal("Unit costs is nil")
	}

	if pts, exists := unit.Costs["pts"]; !exists {
		t.Error("Unit has no points cost")
	} else if pts == 0 {
		t.Error("Unit points cost is 0")
	}

	// Verify categories
	if len(unit.Categories) == 0 {
		t.Error("Unit has no categories")
	}

	// Verify faction
	if unit.Faction == nil {
		t.Error("Unit has no faction")
	} else {
		if unit.Faction.Name == "" {
			t.Error("Faction name is empty")
		}
	}
}

func TestTransformUnitProfile(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)
	transformer := NewTransformer(resolver)

	// Get a unit with a Unit profile
	unitID := "828d-840a-9a67-9074"
	entry, _, found := p.FindSelectionEntryByID(unitID)
	if !found {
		t.Fatalf("Unit %s not found", unitID)
	}

	// Find Unit profile
	var unitProfile *models.Profile
	for i := range entry.Profiles {
		if entry.Profiles[i].TypeName == "Unit" {
			unitProfile = &entry.Profiles[i]
			break
		}
	}

	if unitProfile == nil {
		t.Fatal("Unit profile not found")
	}

	transformed := transformer.transformUnitProfile(*unitProfile)

	if transformed == nil {
		t.Fatal("Transformed unit profile is nil")
	}

	// Verify all fields are populated
	if transformed.Movement == "" {
		t.Error("Movement is empty")
	}
	if transformed.Toughness == 0 {
		t.Error("Toughness is 0")
	}
	if transformed.Save == "" {
		t.Error("Save is empty")
	}
	if transformed.Wounds == 0 {
		t.Error("Wounds is 0")
	}
	if transformed.Leadership == "" {
		t.Error("Leadership is empty")
	}
	if transformed.ObjectiveControl == 0 {
		t.Error("ObjectiveControl is 0")
	}
}

func TestTransformWeapons(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)
	transformer := NewTransformer(resolver)

	// Get a unit with weapons (Asurmen)
	unitID := "828d-840a-9a67-9074"
	entry, _, found := p.FindSelectionEntryByID(unitID)
	if !found {
		t.Fatalf("Unit %s not found", unitID)
	}

	weapons := transformer.transformWeapons(entry)

	if weapons == nil {
		t.Fatal("Weapons is nil")
	}

	// Asurmen should have ranged weapons (The Bloody Twins)
	if len(weapons.Ranged) == 0 {
		t.Error("Unit has no ranged weapons")
	} else {
		ranged := weapons.Ranged[0]
		if ranged.Name == "" {
			t.Error("Ranged weapon has no name")
		}
		if ranged.Range == "" {
			t.Error("Ranged weapon has no range")
		}
		if ranged.Attacks == "" {
			t.Error("Ranged weapon has no attacks")
		}
		if ranged.Strength == "" {
			t.Error("Ranged weapon has no strength")
		}
	}

	// Asurmen should have melee weapons (The Sword of Asur)
	if len(weapons.Melee) == 0 {
		t.Error("Unit has no melee weapons")
	} else {
		melee := weapons.Melee[0]
		if melee.Name == "" {
			t.Error("Melee weapon has no name")
		}
		if melee.Range == "" {
			t.Error("Melee weapon has no range")
		}
		if melee.Attacks == "" {
			t.Error("Melee weapon has no attacks")
		}
		if melee.Strength == "" {
			t.Error("Melee weapon has no strength")
		}
	}
}

func TestTransformWeaponsWithEntryLinks(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)
	transformer := NewTransformer(resolver)

	// Get Warlock - a unit with weapons in entryLinks within selectionEntryGroups
	unitID := "291f-2885-4512-f0cd"
	entry, catID, found := p.FindSelectionEntryByID(unitID)
	if !found {
		t.Fatalf("Unit %s not found", unitID)
	}

	weapons := transformer.transformWeaponsWithCatalogue(entry, catID)

	if weapons == nil {
		t.Fatal("Weapons is nil")
	}

	// Warlock should have ranged weapons (Destructor and Shuriken Pistol)
	if len(weapons.Ranged) == 0 {
		t.Error("Warlock has no ranged weapons - entryLinks in selectionEntryGroups not being resolved")
	} else {
		// Check for Destructor
		foundDestructor := false
		for _, w := range weapons.Ranged {
			if w.Name == "Destructor" {
				foundDestructor = true
				if w.Range == "" {
					t.Error("Destructor has no range")
				}
				if w.Attacks == "" {
					t.Error("Destructor has no attacks")
				}
				break
			}
		}
		if !foundDestructor {
			t.Error("Destructor weapon not found in Warlock's weapons")
		}

		// Check for Shuriken Pistol
		foundShurikenPistol := false
		for _, w := range weapons.Ranged {
			if w.Name == "Shuriken Pistol" {
				foundShurikenPistol = true
				if w.Range == "" {
					t.Error("Shuriken Pistol has no range")
				}
				break
			}
		}
		if !foundShurikenPistol {
			t.Error("Shuriken Pistol weapon not found in Warlock's weapons")
		}
	}
}

func TestTransformCosts(t *testing.T) {
	resolver := &LinkResolver{} // Dummy resolver for transformer
	transformer := NewTransformer(resolver)

	costs := []models.Cost{
		{Name: "pts", TypeID: "test-id", Value: "125"},
		{Name: "Crusade Points", TypeID: "test-id-2", Value: "0"},
		{Name: "Invalid", TypeID: "test-id-3", Value: "not-a-number"},
		{Name: "Empty", TypeID: "test-id-4", Value: ""},
	}

	result := transformer.TransformCosts(costs)

	if result == nil {
		t.Fatal("TransformCosts returned nil")
	}

	// Verify valid costs
	if pts, exists := result["pts"]; !exists {
		t.Error("Points cost not found")
	} else if pts != 125 {
		t.Errorf("Expected points 125, got %d", pts)
	}

	if cp, exists := result["Crusade Points"]; !exists {
		t.Error("Crusade Points cost not found")
	} else if cp != 0 {
		t.Errorf("Expected Crusade Points 0, got %d", cp)
	}

	// Invalid costs should be silently ignored
	if _, exists := result["Invalid"]; exists {
		t.Error("Invalid cost should not be in result")
	}

	// Empty costs should be silently ignored
	if _, exists := result["Empty"]; exists {
		t.Error("Empty cost should not be in result")
	}
}

func TestTransformCatalogue(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)
	transformer := NewTransformer(resolver)

	// Get a catalogue
	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Fatal("No catalogues available")
	}

	var testCatalogue *models.Catalogue
	for _, cat := range catalogues {
		if len(cat.EntryLinks) > 0 {
			testCatalogue = cat
			break
		}
	}

	if testCatalogue == nil {
		t.Skip("No catalogue with entryLinks found")
	}

	response := transformer.TransformCatalogue(testCatalogue)

	if response == nil {
		t.Fatal("Transformed catalogue is nil")
	}

	// Verify basic fields
	if response.ID != testCatalogue.ID {
		t.Errorf("Expected ID %s, got %s", testCatalogue.ID, response.ID)
	}

	if response.Name != testCatalogue.Name {
		t.Errorf("Expected name %s, got %s", testCatalogue.Name, response.Name)
	}

	if response.Revision != testCatalogue.Revision {
		t.Errorf("Expected revision %s, got %s", testCatalogue.Revision, response.Revision)
	}

	// Verify units
	if len(response.Units) == 0 {
		t.Error("Catalogue has no units")
	}

	// Verify linked catalogues
	if len(testCatalogue.CatalogueLinks) > 0 && len(response.LinkedCatalogues) == 0 {
		t.Error("Linked catalogues not resolved")
	}
}


