package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/parser"
	"grimoire-api/internal/service"
)

func setupTestRouter(t *testing.T) *gin.Engine {
	gin.SetMode(gin.TestMode)

	dataDir := getTestDataDir(t)
	
	p := parser.NewParser(dataDir)
	if err := p.LoadGameSystem(); err != nil {
		t.Fatalf("Failed to load game system: %v", err)
	}
	if err := p.LoadAllCatalogues(); err != nil {
		t.Fatalf("Failed to load catalogues: %v", err)
	}

	resolver := parser.NewLinkResolver(p)
	transformer := parser.NewTransformer(resolver)
	cache := cache.NewCache()

	unitService := service.NewUnitService(p, resolver, transformer, cache)
	catalogueService := service.NewCatalogueService(p, resolver, transformer, cache)

	unitHandler := NewUnitHandler(unitService)
	catalogueHandler := NewCatalogueHandler(catalogueService)
	factionHandler := NewFactionHandler(unitService, catalogueService)
	searchHandler := NewSearchHandler(unitService)
	gameSystemHandler := NewGameSystemHandler(p)

	router := gin.New()
	v1 := router.Group("/api/v1")
	{
		v1.GET("/game-system", gameSystemHandler.GetGameSystem)
		v1.GET("/catalogues", catalogueHandler.ListCatalogues)
		v1.GET("/catalogues/:id", catalogueHandler.GetCatalogue)
		v1.GET("/catalogues/:id/units", catalogueHandler.GetCatalogueUnits)
		v1.GET("/units", unitHandler.ListUnits)
		v1.GET("/units/:id", unitHandler.GetUnit)
		v1.GET("/units/:id/weapons", unitHandler.GetUnitWeapons)
		v1.GET("/factions", factionHandler.ListFactions)
		v1.GET("/factions/:name/units", factionHandler.GetFactionUnits)
		v1.GET("/search", searchHandler.Search)
	}

	return router
}

func TestGetUnitHandler(t *testing.T) {
	router := setupTestRouter(t)

	// Test with entryLink ID
	req := httptest.NewRequest("GET", "/api/v1/units/a502-4dbe-d0c6-69fd", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "Warlock")
}

func TestListUnitsHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/units?limit=10", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
}

func TestGetUnitWeaponsHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/units/828d-840a-9a67-9074/weapons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "ranged")
}

func TestGetCatalogueHandler(t *testing.T) {
	router := setupTestRouter(t)

	// Get a catalogue ID first
	dataDir := getTestDataDir(t)
	p := parser.NewParser(dataDir)
	p.LoadGameSystem()
	p.LoadAllCatalogues()
	
	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Skip("No catalogues available")
	}

	var testCatalogueID string
	for id := range catalogues {
		testCatalogueID = id
		break
	}

	req := httptest.NewRequest("GET", "/api/v1/catalogues/"+testCatalogueID, nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

func TestListCataloguesHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/catalogues", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
}

func TestSearchHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/search?q=asurmen", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "results")
}

func TestSearchHandlerMissingQuery(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/search", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

func TestGetGameSystemHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/game-system", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "profileTypes")
}

func TestListFactionsHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/factions", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
}

func TestGetCatalogueUnitsHandler(t *testing.T) {
	router := setupTestRouter(t)

	// Get a catalogue ID first
	dataDir := getTestDataDir(t)
	p := parser.NewParser(dataDir)
	p.LoadGameSystem()
	p.LoadAllCatalogues()
	
	catalogues := p.GetAllCatalogues()
	if len(catalogues) == 0 {
		t.Skip("No catalogues available")
	}

	var testCatalogueID string
	for id, cat := range catalogues {
		if len(cat.EntryLinks) > 0 {
			testCatalogueID = id
			break
		}
	}

	if testCatalogueID == "" {
		t.Skip("No catalogue with units found")
	}

	req := httptest.NewRequest("GET", "/api/v1/catalogues/"+testCatalogueID+"/units", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
}

func TestGetFactionUnitsHandler(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/factions/Imperium/units", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
}

func TestGetUnitHandlerNotFound(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/units/nonexistent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetCatalogueHandlerNotFound(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/catalogues/nonexistent-id", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetCatalogueUnitsHandlerNotFound(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/catalogues/nonexistent-id/units", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestGetUnitWeaponsHandlerNotFound(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/units/nonexistent-id/weapons", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestListUnitsHandlerPagination(t *testing.T) {
	router := setupTestRouter(t)

	// Test with offset
	req := httptest.NewRequest("GET", "/api/v1/units?limit=5&offset=0", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")
	assert.Contains(t, w.Body.String(), "total")
}

func TestListUnitsHandlerFilters(t *testing.T) {
	router := setupTestRouter(t)

	// Test with search filter
	req := httptest.NewRequest("GET", "/api/v1/units?search=marine", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "data")

	// Test with faction filter
	req2 := httptest.NewRequest("GET", "/api/v1/units?faction=Imperium", nil)
	w2 := httptest.NewRecorder()
	router.ServeHTTP(w2, req2)

	assert.Equal(t, http.StatusOK, w2.Code)
	assert.Contains(t, w2.Body.String(), "data")
}

func TestSearchHandlerWithLimit(t *testing.T) {
	router := setupTestRouter(t)

	req := httptest.NewRequest("GET", "/api/v1/search?q=a&limit=5", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Contains(t, w.Body.String(), "results")
	assert.Contains(t, w.Body.String(), "total")
}

func getTestDataDir(t *testing.T) string {
	dataDir := os.Getenv("TEST_DATA_DIR")
	if dataDir == "" {
		dataDir = "../../../wh40k-10e"
	}

	if _, err := os.Stat(dataDir); os.IsNotExist(err) {
		t.Skipf("Data directory not found: %s", dataDir)
	}

	return dataDir
}

