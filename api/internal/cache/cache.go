package cache

import (
	"sync"

	"grimoire-api/internal/models"
)

// Cache provides in-memory caching for parsed data
type Cache struct {
	units       map[string]*models.UnitResponse
	catalogues  map[string]*models.CatalogueResponse
	mu          sync.RWMutex
}

// NewCache creates a new cache instance
func NewCache() *Cache {
	return &Cache{
		units:      make(map[string]*models.UnitResponse),
		catalogues: make(map[string]*models.CatalogueResponse),
	}
}

// SetUnit caches a unit response
func (c *Cache) SetUnit(id string, unit *models.UnitResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.units[id] = unit
}

// GetUnit retrieves a cached unit response
func (c *Cache) GetUnit(id string) (*models.UnitResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	unit, exists := c.units[id]
	return unit, exists
}

// SetCatalogue caches a catalogue response
func (c *Cache) SetCatalogue(id string, catalogue *models.CatalogueResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.catalogues[id] = catalogue
}

// GetCatalogue retrieves a cached catalogue response
func (c *Cache) GetCatalogue(id string) (*models.CatalogueResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	catalogue, exists := c.catalogues[id]
	return catalogue, exists
}

// Clear clears all cached data
func (c *Cache) Clear() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.units = make(map[string]*models.UnitResponse)
	c.catalogues = make(map[string]*models.CatalogueResponse)
}

// ClearUnit removes a specific unit from cache
func (c *Cache) ClearUnit(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.units, id)
}

// ClearCatalogue removes a specific catalogue from cache
func (c *Cache) ClearCatalogue(id string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.catalogues, id)
}

