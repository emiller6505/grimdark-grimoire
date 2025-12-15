package service

import (
	"testing"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/parser"
)

func TestGetCatalogue(t *testing.T) {
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

	service := NewCatalogueService(p, resolver, transformer, cache)

	// Get a catalogue
	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Fatal("No catalogues available")
	}

	var testCatalogueID string
	for id := range catalogues {
		testCatalogueID = id
		break
	}

	catalogue, err := service.GetCatalogue(testCatalogueID)
	if err != nil {
		t.Fatalf("Failed to get catalogue: %v", err)
	}

	if catalogue == nil {
		t.Fatal("Catalogue is nil")
	}

	if catalogue.ID != testCatalogueID {
		t.Errorf("Expected ID %s, got %s", testCatalogueID, catalogue.ID)
	}

	if catalogue.Name == "" {
		t.Error("Catalogue name is empty")
	}

	// Test non-existent catalogue
	_, err = service.GetCatalogue("nonexistent-id")
	if err == nil {
		t.Error("Expected error for non-existent catalogue")
	}
}

func TestListCatalogues(t *testing.T) {
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

	service := NewCatalogueService(p, resolver, transformer, cache)

	catalogues, err := service.ListCatalogues()
	if err != nil {
		t.Fatalf("Failed to list catalogues: %v", err)
	}

	if len(catalogues) == 0 {
		t.Error("No catalogues returned")
	}

	// Verify catalogue structure
	for _, cat := range catalogues {
		if cat.ID == "" {
			t.Error("Catalogue ID is empty")
		}
		if cat.Name == "" {
			t.Error("Catalogue name is empty")
		}
		if cat.Revision == "" {
			t.Error("Catalogue revision is empty")
		}
	}
}

func TestGetCatalogueUnits(t *testing.T) {
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

	service := NewCatalogueService(p, resolver, transformer, cache)

	// Find a catalogue with units
	catalogues := p.GetAllCatalogues()
	var testCatalogueID string
	for id, cat := range catalogues {
		if len(cat.EntryLinks) > 0 {
			testCatalogueID = id
			break
		}
	}

	if testCatalogueID == "" {
		t.Skip("No catalogue with units found")
	}

	units, err := service.GetCatalogueUnits(testCatalogueID)
	if err != nil {
		t.Fatalf("Failed to get catalogue units: %v", err)
	}

	if len(units) == 0 {
		t.Error("Catalogue has no units")
	}

	// Verify unit structure
	for _, unit := range units {
		if unit.ID == "" {
			t.Error("Unit ID is empty")
		}
		if unit.Name == "" {
			t.Error("Unit name is empty")
		}
	}
}


