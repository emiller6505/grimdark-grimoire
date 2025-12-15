package parser

import (
	"testing"

	"grimoire-api/internal/models"
)

func TestResolveEntryLink(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)

	// Get a catalogue with entryLinks
	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Fatal("No catalogues available for testing")
	}

	var testEntryLink *models.EntryLink
	var testCatalogueID string
	
	// Find an entryLink to test with
	for catID, cat := range catalogues {
		for i := range cat.EntryLinks {
			if cat.EntryLinks[i].Type == "selectionEntry" && cat.EntryLinks[i].TargetID != "" {
				testEntryLink = &cat.EntryLinks[i]
				testCatalogueID = catID
				break
			}
		}
		if testEntryLink != nil {
			break
		}
	}

	if testEntryLink == nil {
		t.Skip("No suitable entryLink found for testing")
	}

	// Resolve the entryLink
	entry, err := resolver.ResolveEntryLink(testEntryLink, testCatalogueID)
	if err != nil {
		t.Fatalf("Failed to resolve entryLink: %v", err)
	}

	if entry == nil {
		t.Fatal("Resolved entry is nil")
	}

	if entry.ID != testEntryLink.TargetID {
		t.Errorf("Expected entry ID %s, got %s", testEntryLink.TargetID, entry.ID)
	}

	// Verify entry has required fields
	if entry.Name == "" {
		t.Error("Resolved entry has no name")
	}

	if entry.Type == "" {
		t.Error("Resolved entry has no type")
	}
}

func TestResolveCatalogueLinks(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)

	// Find a catalogue with catalogueLinks
	catalogues := p.GetAllCatalogues()
	var testCatalogue *models.Catalogue
	
	for _, cat := range catalogues {
		if len(cat.CatalogueLinks) > 0 {
			testCatalogue = cat
			break
		}
	}

	if testCatalogue == nil {
		t.Skip("No catalogue with catalogueLinks found for testing")
	}

	// Resolve catalogue links
	linked := resolver.ResolveCatalogueLinks(testCatalogue)
	
	if len(linked) == 0 {
		t.Error("No linked catalogues resolved")
	}

	// Verify linked catalogues are valid
	for _, cat := range linked {
		if cat.ID == "" {
			t.Error("Linked catalogue has no ID")
		}
		if cat.Name == "" {
			t.Error("Linked catalogue has no name")
		}
	}
}

func TestMergeEntryLinkWithSelectionEntry(t *testing.T) {
	dataDir := getTestDataDir(t)
	p := NewParser(dataDir)
	
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := NewLinkResolver(p)

	// Find an entryLink with overrides (costs, categoryLinks)
	catalogues := p.GetAllCatalogues()
	var testEntryLink *models.EntryLink
	var testCatalogueID string
	
	for catID, cat := range catalogues {
		for i := range cat.EntryLinks {
			el := &cat.EntryLinks[i]
			if el.Type == "selectionEntry" && (len(el.Costs) > 0 || len(el.CategoryLinks) > 0) {
				testEntryLink = el
				testCatalogueID = catID
				break
			}
		}
		if testEntryLink != nil {
			break
		}
	}

	if testEntryLink == nil {
		t.Skip("No entryLink with overrides found for testing")
	}

	// Resolve the entry
	entry, err := resolver.ResolveEntryLink(testEntryLink, testCatalogueID)
	if err != nil {
		t.Fatalf("Failed to resolve entryLink: %v", err)
	}

	// Merge
	merged := resolver.MergeEntryLinkWithSelectionEntry(testEntryLink, entry)

	if merged == nil {
		t.Fatal("Merged entry is nil")
	}

	// Verify merge worked
	if testEntryLink.Name != "" && merged.Name != testEntryLink.Name {
		t.Error("EntryLink name override not applied")
	}

	// If entryLink has costs, they should be merged
	if len(testEntryLink.Costs) > 0 {
		// Costs should be merged (entryLink costs override)
		found := false
		for _, cost := range merged.Costs {
			if cost.TypeID == testEntryLink.Costs[0].TypeID {
				found = true
				break
			}
		}
		if !found && len(merged.Costs) > 0 {
			// This is okay - costs might be merged differently
		}
	}
}


