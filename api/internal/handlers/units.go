package handlers

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"grimoire-api/internal/service"
	"grimoire-api/pkg/response"
)

// UnitHandler handles unit-related HTTP requests
type UnitHandler struct {
	service *service.UnitService
}

// NewUnitHandler creates a new unit handler
func NewUnitHandler(unitService *service.UnitService) *UnitHandler {
	return &UnitHandler{service: unitService}
}

// GetUnit handles GET /api/v1/units/:id
func (h *UnitHandler) GetUnit(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "unit ID is required")
		return
	}

	unit, err := h.service.GetUnit(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, unit)
}

// ListUnits handles GET /api/v1/units
func (h *UnitHandler) ListUnits(c *gin.Context) {
	faction := c.Query("faction")
	category := c.Query("category")
	search := c.Query("search")

	limit := 100
	if limitStr := c.Query("limit"); limitStr != "" {
		if l, err := strconv.Atoi(limitStr); err == nil && l > 0 && l <= 1000 {
			limit = l
		}
	}

	offset := 0
	if offsetStr := c.Query("offset"); offsetStr != "" {
		if o, err := strconv.Atoi(offsetStr); err == nil && o >= 0 {
			offset = o
		}
	}

	units, total, err := h.service.ListUnits(faction, category, search, limit, offset)
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Paginated(c, units, total, limit, offset)
}

// GetUnitWeapons handles GET /api/v1/units/:id/weapons
func (h *UnitHandler) GetUnitWeapons(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "unit ID is required")
		return
	}

	weapons, err := h.service.GetUnitWeapons(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, weapons)
}


