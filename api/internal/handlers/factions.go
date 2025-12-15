package handlers

import (
	"strings"

	"github.com/gin-gonic/gin"
	"grimoire-api/internal/service"
	"grimoire-api/pkg/response"
)

// FactionHandler handles faction-related HTTP requests
type FactionHandler struct {
	unitService      *service.UnitService
	catalogueService *service.CatalogueService
}

// NewFactionHandler creates a new faction handler
func NewFactionHandler(unitService *service.UnitService, catalogueService *service.CatalogueService) *FactionHandler {
	return &FactionHandler{
		unitService:      unitService,
		catalogueService: catalogueService,
	}
}

// ListFactions handles GET /api/v1/factions
func (h *FactionHandler) ListFactions(c *gin.Context) {
	catalogues, err := h.catalogueService.ListCatalogues()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	// Group catalogues by faction
	factionMap := make(map[string][]string)
	for _, cat := range catalogues {
		// Extract faction from catalogue name
		faction := extractFactionFromName(cat.Name)
		if faction != "" {
			factionMap[faction] = append(factionMap[faction], cat.Name)
		}
	}

	// Convert to response format
	factions := make([]map[string]interface{}, 0, len(factionMap))
	for faction, catalogues := range factionMap {
		factions = append(factions, map[string]interface{}{
			"name":       faction,
			"catalogues": catalogues,
		})
	}

	response.Success(c, factions)
}

// GetFactionUnits handles GET /api/v1/factions/:name/units
func (h *FactionHandler) GetFactionUnits(c *gin.Context) {
	factionName := c.Param("name")
	if factionName == "" {
		response.BadRequest(c, "faction name is required")
		return
	}

	units, _, err := h.unitService.ListUnits(factionName, "", "", 1000, 0)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, units)
}

// extractFactionFromName extracts faction name from catalogue name
func extractFactionFromName(name string) string {
	// Remove prefixes like "Imperium - ", "Chaos - ", "Aeldari - "
	parts := strings.Split(name, " - ")
	if len(parts) > 1 {
		return parts[0]
	}
	
	// For simple names like "Necrons.cat", "Orks.cat"
	if strings.Contains(name, " ") {
		return strings.Split(name, " ")[0]
	}
	
	return name
}


