package service

import (
	"fmt"
	"strings"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/models"
	"grimoire-api/internal/parser"
)

// UnitService handles unit-related business logic
type UnitService struct {
	parser     *parser.Parser
	resolver   *parser.LinkResolver
	transformer *parser.Transformer
	cache      *cache.Cache
}

// NewUnitService creates a new unit service
func NewUnitService(p *parser.Parser, r *parser.LinkResolver, t *parser.Transformer, c *cache.Cache) *UnitService {
	return &UnitService{
		parser:     p,
		resolver:   r,
		transformer: t,
		cache:      c,
	}
}

// GetUnit retrieves a unit by ID
// The ID can be either an entryLink ID (from a catalogue) or a selectionEntry ID (from a library)
func (s *UnitService) GetUnit(id string) (*models.UnitResponse, error) {
	// Check cache first
	if unit, exists := s.cache.GetUnit(id); exists {
		return unit, nil
	}

	var entry *models.SelectionEntry
	var catalogueID string
	var found bool
	var foundViaEntryLink bool

	// First, try to find as an entryLink (most common case for API consumers)
	if entryLink, catID, linkFound := s.parser.FindEntryLinkByID(id); linkFound {
		// Resolve the entryLink to its selectionEntry
		resolvedEntry, err := s.resolver.ResolveEntryLink(entryLink, catID)
		if err == nil {
			entry = resolvedEntry
			catalogueID = catID
			found = true
			foundViaEntryLink = true
			// Merge entryLink overrides with resolved entry
			entry = s.resolver.MergeEntryLinkWithSelectionEntry(entryLink, entry)
		}
	}

	// If not found as entryLink, try as selectionEntry ID directly
	if !found {
		entry, catalogueID, found = s.parser.FindSelectionEntryByID(id)
	}

	if !found {
		return nil, fmt.Errorf("unit not found: %s (tried as entryLink and selectionEntry)", id)
	}

	// Transform to response
	unit := s.transformer.TransformUnit(entry, catalogueID)

	// If we found via entryLink, preserve the original entryLink ID
	// (the transformer uses the selectionEntry ID by default)
	if foundViaEntryLink {
		unit.ID = id
	}

	// Cache the result
	s.cache.SetUnit(id, unit)

	return unit, nil
}

// ListUnits lists all units with optional filters
func (s *UnitService) ListUnits(faction, category, search string, limit, offset int) ([]models.UnitSummary, int, error) {
	var allUnits []models.UnitSummary

	// Collect units from all catalogues
	for _, catalogue := range s.parser.GetAllCatalogues() {
		for _, entryLink := range catalogue.EntryLinks {
			if entryLink.Type == "selectionEntry" {
				// Try to resolve the entry
				entry, err := s.resolver.ResolveEntryLink(&entryLink, catalogue.ID)
				if err != nil {
					continue
				}

				// Apply filters
				if faction != "" {
					hasFaction := false
					for _, catLink := range entry.CategoryLinks {
						if strings.Contains(catLink.Name, faction) {
							hasFaction = true
							break
						}
					}
					if !hasFaction {
						continue
					}
				}

				if category != "" {
					hasCategory := false
					for _, catLink := range entry.CategoryLinks {
						if strings.Contains(catLink.Name, category) {
							hasCategory = true
							break
						}
					}
					if !hasCategory {
						continue
					}
				}

				if search != "" {
					if !strings.Contains(strings.ToLower(entry.Name), strings.ToLower(search)) {
						continue
					}
				}

				summary := models.UnitSummary{
					ID:       entryLink.ID,
					Name:     entry.Name,
					TargetID: entryLink.TargetID,
					Type:     entry.Type,
					Costs:    s.transformer.TransformCosts(entry.Costs),
				}
				allUnits = append(allUnits, summary)
			}
		}
	}

	total := len(allUnits)

	// Apply pagination
	start := offset
	if start > total {
		start = total
	}
	end := start + limit
	if end > total {
		end = total
	}

	if start < total {
		return allUnits[start:end], total, nil
	}

	return []models.UnitSummary{}, total, nil
}

// SearchUnits searches for units by name
func (s *UnitService) SearchUnits(query string, limit int) ([]models.SearchResult, error) {
	query = strings.ToLower(query)
	var results []models.SearchResult

	// Search through all catalogues
	for _, catalogue := range s.parser.GetAllCatalogues() {
		for _, entryLink := range catalogue.EntryLinks {
			if entryLink.Type == "selectionEntry" {
				entry, err := s.resolver.ResolveEntryLink(&entryLink, catalogue.ID)
				if err != nil {
					continue
				}

				if strings.Contains(strings.ToLower(entry.Name), query) {
					results = append(results, models.SearchResult{
						Type: "unit",
						ID:   entryLink.ID,
						Name: entry.Name,
					})

					if len(results) >= limit {
						return results, nil
					}
				}
			}
		}
	}

	return results, nil
}

// GetUnitWeapons retrieves weapons for a unit
func (s *UnitService) GetUnitWeapons(id string) (*models.WeaponSet, error) {
	unit, err := s.GetUnit(id)
	if err != nil {
		return nil, err
	}

	return unit.Weapons, nil
}

