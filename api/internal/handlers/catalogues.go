package handlers

import (
	"github.com/gin-gonic/gin"
	"grimoire-api/internal/service"
	"grimoire-api/pkg/response"
)

// CatalogueHandler handles catalogue-related HTTP requests
type CatalogueHandler struct {
	service *service.CatalogueService
}

// NewCatalogueHandler creates a new catalogue handler
func NewCatalogueHandler(catalogueService *service.CatalogueService) *CatalogueHandler {
	return &CatalogueHandler{service: catalogueService}
}

// GetCatalogue handles GET /api/v1/catalogues/:id
func (h *CatalogueHandler) GetCatalogue(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "catalogue ID is required")
		return
	}

	catalogue, err := h.service.GetCatalogue(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, catalogue)
}

// ListCatalogues handles GET /api/v1/catalogues
func (h *CatalogueHandler) ListCatalogues(c *gin.Context) {
	catalogues, err := h.service.ListCatalogues()
	if err != nil {
		response.InternalServerError(c, err.Error())
		return
	}

	response.Success(c, catalogues)
}

// GetCatalogueUnits handles GET /api/v1/catalogues/:id/units
func (h *CatalogueHandler) GetCatalogueUnits(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		response.BadRequest(c, "catalogue ID is required")
		return
	}

	units, err := h.service.GetCatalogueUnits(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.Success(c, units)
}


