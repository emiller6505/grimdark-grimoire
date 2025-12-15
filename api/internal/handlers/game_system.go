package handlers

import (
	"github.com/gin-gonic/gin"
	"grimoire-api/internal/models"
	"grimoire-api/internal/parser"
	"grimoire-api/pkg/response"
)

// GameSystemHandler handles game system-related HTTP requests
type GameSystemHandler struct {
	parser *parser.Parser
}

// NewGameSystemHandler creates a new game system handler
func NewGameSystemHandler(p *parser.Parser) *GameSystemHandler {
	return &GameSystemHandler{parser: p}
}

// GetGameSystem handles GET /api/v1/game-system
func (h *GameSystemHandler) GetGameSystem(c *gin.Context) {
	gameSystem := h.parser.GetGameSystem()
	if gameSystem == nil {
		response.InternalServerError(c, "game system not loaded")
		return
	}

	// Transform to response
	responseData := &models.GameSystemResponse{
		ID:                gameSystem.ID,
		Name:              gameSystem.Name,
		Revision:          gameSystem.Revision,
		BattleScribeVersion: gameSystem.BattleScribeVersion,
		ProfileTypes:      make([]models.ProfileTypeInfo, 0),
		Categories:        make([]models.CategoryInfo, 0),
		CostTypes:         make([]models.CostTypeInfo, 0),
	}

	// Transform profile types
	for _, pt := range gameSystem.ProfileTypes {
		characteristics := make([]string, 0, len(pt.CharacteristicTypes))
		for _, ct := range pt.CharacteristicTypes {
			characteristics = append(characteristics, ct.Name)
		}
		responseData.ProfileTypes = append(responseData.ProfileTypes, models.ProfileTypeInfo{
			ID:              pt.ID,
			Name:            pt.Name,
			Characteristics: characteristics,
		})
	}

	// Transform categories
	for _, cat := range gameSystem.CategoryEntries {
		if cat.Hidden != "true" {
			responseData.Categories = append(responseData.Categories, models.CategoryInfo{
				ID:   cat.ID,
				Name: cat.Name,
			})
		}
	}

	// Transform cost types
	for _, ct := range gameSystem.CostTypes {
		responseData.CostTypes = append(responseData.CostTypes, models.CostTypeInfo{
			ID:              ct.ID,
			Name:            ct.Name,
			DefaultCostLimit: ct.DefaultCostLimit,
			Hidden:          ct.Hidden == "true",
		})
	}

	response.Success(c, responseData)
}

