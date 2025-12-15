package main

import (
	"log"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/handlers"
	"grimoire-api/internal/parser"
	"grimoire-api/internal/service"
)

func main() {
	// Get data directory from environment or use default
	dataDir := os.Getenv("DATA_DIR")
	if dataDir == "" {
		dataDir = "../wh40k-10e"
	}

	log.Printf("Initializing parser with data directory: %s", dataDir)

	// Initialize parser
	p := parser.NewParser(dataDir)

	// Load game system
	if err := p.LoadGameSystem(); err != nil {
		log.Fatalf("Failed to load game system: %v", err)
	}

	// Load all catalogues
	if err := p.LoadAllCatalogues(); err != nil {
		log.Fatalf("Failed to load catalogues: %v", err)
	}

	log.Printf("Loaded %d catalogues and %d libraries", len(p.GetAllCatalogues()), len(p.GetAllLibraries()))

	// Initialize components
	cache := cache.NewCache()
	resolver := parser.NewLinkResolver(p)
	transformer := parser.NewTransformer(resolver)

	// Initialize services
	unitService := service.NewUnitService(p, resolver, transformer, cache)
	catalogueService := service.NewCatalogueService(p, resolver, transformer, cache)

	// Initialize handlers
	unitHandler := handlers.NewUnitHandler(unitService)
	catalogueHandler := handlers.NewCatalogueHandler(catalogueService)
	factionHandler := handlers.NewFactionHandler(unitService, catalogueService)
	searchHandler := handlers.NewSearchHandler(unitService)
	gameSystemHandler := handlers.NewGameSystemHandler(p)

	// Setup Gin router
	if os.Getenv("GIN_MODE") == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// CORS middleware
	config := cors.DefaultConfig()
	config.AllowAllOrigins = true
	config.AllowMethods = []string{"GET", "OPTIONS"}
	router.Use(cors.New(config))

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
		})
	})

	// API routes
	v1 := router.Group("/api/v1")
	{
		// Game system
		v1.GET("/game-system", gameSystemHandler.GetGameSystem)

		// Catalogues
		v1.GET("/catalogues", catalogueHandler.ListCatalogues)
		v1.GET("/catalogues/:id", catalogueHandler.GetCatalogue)
		v1.GET("/catalogues/:id/units", catalogueHandler.GetCatalogueUnits)

		// Units
		v1.GET("/units", unitHandler.ListUnits)
		v1.GET("/units/:id", unitHandler.GetUnit)
		v1.GET("/units/:id/weapons", unitHandler.GetUnitWeapons)

		// Factions
		v1.GET("/factions", factionHandler.ListFactions)
		v1.GET("/factions/:name/units", factionHandler.GetFactionUnits)

		// Search
		v1.GET("/search", searchHandler.Search)
	}

	// Root endpoint
	router.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"name":    "Warhammer 40K 10th Edition API",
			"version": "1.0.0",
			"endpoints": gin.H{
				"game-system": "/api/v1/game-system",
				"catalogues":  "/api/v1/catalogues",
				"units":       "/api/v1/units",
				"factions":    "/api/v1/factions",
				"search":      "/api/v1/search",
			},
		})
	})

	// Get port from environment or use default
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Starting server on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


