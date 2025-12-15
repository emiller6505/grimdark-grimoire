package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"grimoire-api/internal/service"
	"grimoire-api/pkg/response"
)

// SearchHandler handles search-related HTTP requests
type SearchHandler struct {
	unitService *service.UnitService
}

// NewSearchHandler creates a new search handler
func NewSearchHandler(unitService *service.UnitService) *SearchHandler {
	return &SearchHandler{unitService: unitService}
}

// Search handles GET /api/v1/search
func (h *SearchHandler) Search(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "search query is required")
		return
	}

	limit := 50
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 200 {
			limit = l
		}
	}

	results, err := h.unitService.SearchUnits(query, limit)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, map[string]interface{}{
		"query":   query,
		"results": results,
		"total":   len(results),
	})
}

