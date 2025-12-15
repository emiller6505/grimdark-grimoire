package service

import (
	"fmt"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/models"
	"grimoire-api/internal/parser"
)

// CatalogueService handles catalogue-related business logic
type CatalogueService struct {
	parser     *parser.Parser
	resolver   *parser.LinkResolver
	transformer *parser.Transformer
	cache      *cache.Cache
}

// NewCatalogueService creates a new catalogue service
func NewCatalogueService(p *parser.Parser, r *parser.LinkResolver, t *parser.Transformer, c *cache.Cache) *CatalogueService {
	return &CatalogueService{
		parser:     p,
		resolver:   r,
		transformer: t,
		cache:      c,
	}
}

// GetCatalogue retrieves a catalogue by ID
func (s *CatalogueService) GetCatalogue(id string) (*models.CatalogueResponse, error) {
	// Check cache first
	if catalogue, exists := s.cache.GetCatalogue(id); exists {
		return catalogue, nil
	}

	// Get catalogue from parser
	catalogue, exists := s.parser.GetCatalogue(id)
	if !exists {
		return nil, fmt.Errorf("catalogue not found: %s", id)
	}

	// Transform to response
	response := s.transformer.TransformCatalogue(catalogue)

	// Cache the result
	s.cache.SetCatalogue(id, response)

	return response, nil
}

// ListCatalogues lists all catalogues
func (s *CatalogueService) ListCatalogues() ([]models.CatalogueInfo, error) {
	catalogues := s.parser.GetAllCatalogues()
	result := make([]models.CatalogueInfo, 0, len(catalogues))

	for _, cat := range catalogues {
		result = append(result, models.CatalogueInfo{
			ID:       cat.ID,
			Name:     cat.Name,
			Revision: cat.Revision,
			Library:  cat.Library == "true",
		})
	}

	return result, nil
}

// GetCatalogueUnits retrieves all units in a catalogue
func (s *CatalogueService) GetCatalogueUnits(id string) ([]models.UnitSummary, error) {
	catalogue, exists := s.parser.GetCatalogue(id)
	if !exists {
		return nil, fmt.Errorf("catalogue not found: %s", id)
	}

	units := make([]models.UnitSummary, 0, len(catalogue.EntryLinks))
	for _, entryLink := range catalogue.EntryLinks {
		if entryLink.Type == "selectionEntry" {
			entry, err := s.resolver.ResolveEntryLink(&entryLink, id)
			if err != nil {
				continue
			}

			summary := models.UnitSummary{
				ID:       entryLink.ID,
				Name:     entry.Name,
				TargetID: entryLink.TargetID,
				Type:     entry.Type,
				Costs:    s.transformer.TransformCosts(entry.Costs),
			}
			units = append(units, summary)
		}
	}

	return units, nil
}

