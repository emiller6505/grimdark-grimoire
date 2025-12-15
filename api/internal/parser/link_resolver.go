package parser

import (
	"fmt"

	"grimoire-api/internal/models"
)

// LinkResolver resolves entryLinks to their actual selectionEntries
type LinkResolver struct {
	parser *Parser
}

// NewLinkResolver creates a new link resolver
func NewLinkResolver(parser *Parser) *LinkResolver {
	return &LinkResolver{parser: parser}
}

// ResolveEntryLink resolves an entryLink to its selectionEntry
func (lr *LinkResolver) ResolveEntryLink(entryLink *models.EntryLink, catalogueID string) (*models.SelectionEntry, error) {
	targetID := entryLink.TargetID
	if targetID == "" {
		return nil, fmt.Errorf("entryLink has no targetId")
	}

	// First, check if catalogueID is a library and search there first
	if library, exists := lr.parser.GetLibrary(catalogueID); exists {
		if entry := findEntryInCatalogue(library, targetID); entry != nil {
			return entry, nil
		}
	}

	// Then check if the target is in a linked library from a catalogue
	if catalogue, exists := lr.parser.GetCatalogue(catalogueID); exists {
		// Check catalogueLinks for libraries
		for _, catLink := range catalogue.CatalogueLinks {
			if catLink.ImportRootEntries == "true" {
				library, libExists := lr.parser.GetLibrary(catLink.TargetID)
				if libExists {
					if entry := findEntryInCatalogue(library, targetID); entry != nil {
						return entry, nil
					}
				}
			}
		}
	}

	// Search all libraries
	for _, library := range lr.parser.GetAllLibraries() {
		if entry := findEntryInCatalogue(library, targetID); entry != nil {
			return entry, nil
		}
	}

	// Search all catalogues
	for _, cat := range lr.parser.GetAllCatalogues() {
		if entry := findEntryInCatalogue(cat, targetID); entry != nil {
			return entry, nil
		}
	}

	return nil, fmt.Errorf("selectionEntry with id %s not found", targetID)
}

// ResolveCatalogueLinks resolves all catalogueLinks for a catalogue
func (lr *LinkResolver) ResolveCatalogueLinks(catalogue *models.Catalogue) []*models.Catalogue {
	var linked []*models.Catalogue
	for _, catLink := range catalogue.CatalogueLinks {
		if library, exists := lr.parser.GetLibrary(catLink.TargetID); exists {
			linked = append(linked, library)
		} else if cat, exists := lr.parser.GetCatalogue(catLink.TargetID); exists {
			linked = append(linked, cat)
		}
	}
	return linked
}

// MergeEntryLinkWithSelectionEntry merges entryLink overrides with the resolved selectionEntry
func (lr *LinkResolver) MergeEntryLinkWithSelectionEntry(entryLink *models.EntryLink, selectionEntry *models.SelectionEntry) *models.SelectionEntry {
	// Create a copy of the selectionEntry
	merged := *selectionEntry

	// Override name if entryLink has one
	if entryLink.Name != "" {
		merged.Name = entryLink.Name
	}

	// Merge costs (entryLink costs override selectionEntry costs)
	if len(entryLink.Costs) > 0 {
		costMap := make(map[string]*models.Cost)
		for i := range merged.Costs {
			costMap[merged.Costs[i].TypeID] = &merged.Costs[i]
		}
		for _, cost := range entryLink.Costs {
			costMap[cost.TypeID] = &cost
		}
		merged.Costs = make([]models.Cost, 0, len(costMap))
		for _, cost := range costMap {
			merged.Costs = append(merged.Costs, *cost)
		}
	}

	// Merge categoryLinks
	if len(entryLink.CategoryLinks) > 0 {
		catMap := make(map[string]*models.CategoryLink)
		for i := range merged.CategoryLinks {
			catMap[merged.CategoryLinks[i].TargetID] = &merged.CategoryLinks[i]
		}
		for _, catLink := range entryLink.CategoryLinks {
			catMap[catLink.TargetID] = &catLink
		}
		merged.CategoryLinks = make([]models.CategoryLink, 0, len(catMap))
		for _, catLink := range catMap {
			merged.CategoryLinks = append(merged.CategoryLinks, *catLink)
		}
	}

	// Merge constraints
	if len(entryLink.Constraints) > 0 {
		merged.Constraints = append(merged.Constraints, entryLink.Constraints...)
	}

	// Merge modifiers
	// Modifiers from entryLink are appended to modifiers from selectionEntry
	// This allows entryLinks to add additional modifiers while preserving original ones
	if len(entryLink.Modifiers) > 0 {
		merged.Modifiers = append(merged.Modifiers, entryLink.Modifiers...)
	}

	return &merged
}

