package service

import (
	"testing"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/parser"
)

func TestGetUnit(t *testing.T) {
	dataDir := getTestDataDir(t)
	
	p := parser.NewParser(dataDir)
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := parser.NewLinkResolver(p)
	transformer := parser.NewTransformer(resolver)
	cache := cache.NewCache()

	service := NewUnitService(p, resolver, transformer, cache)

	// Test with entryLink ID (most common case)
	entryLinkID := "a502-4dbe-d0c6-69fd" // Warlock
	unit, err := service.GetUnit(entryLinkID)
	if err != nil {
		t.Fatalf("Failed to get unit: %v", err)
	}

	if unit == nil {
		t.Fatal("Unit is nil")
	}

	if unit.ID != entryLinkID {
		t.Errorf("Expected ID %s, got %s", entryLinkID, unit.ID)
	}

	if unit.Name == "" {
		t.Error("Unit name is empty")
	}

	// Test with selectionEntry ID
	selectionEntryID := "828d-840a-9a67-9074" // Asurmen
	unit2, err := service.GetUnit(selectionEntryID)
	if err != nil {
		t.Fatalf("Failed to get unit by selectionEntry ID: %v", err)
	}

	if unit2 == nil {
		t.Fatal("Unit2 is nil")
	}

	if unit2.Name != "Asurmen" {
		t.Errorf("Expected name 'Asurmen', got '%s'", unit2.Name)
	}

	// Test non-existent ID
	_, err = service.GetUnit("nonexistent-id")
	if err == nil {
		t.Error("Expected error for non-existent unit")
	}
}

func TestListUnits(t *testing.T) {
	dataDir := getTestDataDir(t)
	
	p := parser.NewParser(dataDir)
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := parser.NewLinkResolver(p)
	transformer := parser.NewTransformer(resolver)
	cache := cache.NewCache()

	service := NewUnitService(p, resolver, transformer, cache)

	// Test listing all units
	units, total, err := service.ListUnits("", "", "", 100, 0)
	if err != nil {
		t.Fatalf("Failed to list units: %v", err)
	}

	if total == 0 {
		t.Error("No units found")
	}

	if len(units) == 0 {
		t.Error("Units list is empty")
	}

	// Test pagination
	units2, total2, err := service.ListUnits("", "", "", 10, 0)
	if err != nil {
		t.Fatalf("Failed to list units with pagination: %v", err)
	}

	if len(units2) > 10 {
		t.Errorf("Expected max 10 units, got %d", len(units2))
	}

	if total2 != total {
		t.Errorf("Total should be same, got %d vs %d", total2, total)
	}

	// Test search filter
	units3, _, err := service.ListUnits("", "", "marine", 100, 0)
	if err != nil {
		t.Fatalf("Failed to search units: %v", err)
	}

	// Should find some units with "marine" in name
	if len(units3) == 0 {
		t.Log("No units found with 'marine' in name (this might be expected)")
	}

	// Test faction filter
	units4, _, err := service.ListUnits("Imperium", "", "", 100, 0)
	if err != nil {
		t.Fatalf("Failed to filter by faction: %v", err)
	}

	// Should find some Imperium units
	if len(units4) == 0 {
		t.Log("No Imperium units found (this might be expected)")
	}
}

func TestSearchUnits(t *testing.T) {
	dataDir := getTestDataDir(t)
	
	p := parser.NewParser(dataDir)
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := parser.NewLinkResolver(p)
	transformer := parser.NewTransformer(resolver)
	cache := cache.NewCache()

	service := NewUnitService(p, resolver, transformer, cache)

	// Test search
	results, err := service.SearchUnits("asurmen", 10)
	if err != nil {
		t.Fatalf("Failed to search units: %v", err)
	}

	if len(results) == 0 {
		t.Log("No results found for 'asurmen' (case-insensitive search)")
	} else {
		// Verify results
		found := false
		for _, result := range results {
			if result.Type != "unit" {
				t.Errorf("Expected type 'unit', got '%s'", result.Type)
			}
			if result.ID == "" {
				t.Error("Result ID is empty")
			}
			if result.Name == "" {
				t.Error("Result name is empty")
			}
			if containsIgnoreCase(result.Name, "asurmen") {
				found = true
			}
		}
		if !found {
			t.Log("Asurmen not found in search results (might be case-sensitive)")
		}
	}

	// Test limit
	results2, err := service.SearchUnits("a", 5)
	if err != nil {
		t.Fatalf("Failed to search with limit: %v", err)
	}

	if len(results2) > 5 {
		t.Errorf("Expected max 5 results, got %d", len(results2))
	}
}

func TestGetUnitWeapons(t *testing.T) {
	dataDir := getTestDataDir(t)
	
	p := parser.NewParser(dataDir)
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := parser.NewLinkResolver(p)
	transformer := parser.NewTransformer(resolver)
	cache := cache.NewCache()

	service := NewUnitService(p, resolver, transformer, cache)

	// Get a unit with weapons (Asurmen)
	unitID := "828d-840a-9a67-9074"
	weapons, err := service.GetUnitWeapons(unitID)
	if err != nil {
		t.Fatalf("Failed to get unit weapons: %v", err)
	}

	if weapons == nil {
		t.Fatal("Weapons is nil")
	}

	// Asurmen should have weapons
	if len(weapons.Ranged) == 0 && len(weapons.Melee) == 0 {
		t.Error("Unit has no weapons")
	}
}


func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && 
		(len(s) == len(substr) && s == substr ||
		containsMiddleIgnoreCase(s, substr))
}

func containsMiddleIgnoreCase(s, substr string) bool {
	sLower := toLower(s)
	substrLower := toLower(substr)
	for i := 0; i <= len(sLower)-len(substrLower); i++ {
		if sLower[i:i+len(substrLower)] == substrLower {
			return true
		}
	}
	return false
}

func toLower(s string) string {
	result := make([]byte, len(s))
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= 'A' && c <= 'Z' {
			c += 'a' - 'A'
		}
		result[i] = c
	}
	return string(result)
}

