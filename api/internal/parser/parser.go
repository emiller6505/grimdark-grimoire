package parser

import (
	"encoding/xml"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"grimoire-api/internal/models"
)

// Parser handles parsing of BattleScribe XML files
type Parser struct {
	dataDir      string
	gameSystem   *models.GameSystem
	catalogues   map[string]*models.Catalogue
	libraries    map[string]*models.Catalogue
	mu           sync.RWMutex
}

// NewParser creates a new parser instance
func NewParser(dataDir string) *Parser {
	return &Parser{
		dataDir:    dataDir,
		catalogues: make(map[string]*models.Catalogue),
		libraries:  make(map[string]*models.Catalogue),
	}
}

// LoadGameSystem loads and parses the game system file
func (p *Parser) LoadGameSystem() error {
	gstFile := filepath.Join(p.dataDir, "Warhammer 40,000.gst")
	
	data, err := os.ReadFile(gstFile)
	if err != nil {
		return fmt.Errorf("failed to read game system file: %w", err)
	}

	var gameSystem models.GameSystem
	if err := xml.Unmarshal(data, &gameSystem); err != nil {
		return fmt.Errorf("failed to parse game system file: %w", err)
	}

	p.mu.Lock()
	p.gameSystem = &gameSystem
	p.mu.Unlock()

	log.Printf("Loaded game system: %s (revision %s)", gameSystem.Name, gameSystem.Revision)
	return nil
}

// LoadAllCatalogues loads all catalogue files from the data directory
func (p *Parser) LoadAllCatalogues() error {
	return filepath.WalkDir(p.dataDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			return nil
		}

		if !strings.HasSuffix(strings.ToLower(path), ".cat") {
			return nil
		}

		return p.LoadCatalogue(path)
	})
}

// LoadCatalogue loads and parses a single catalogue file
func (p *Parser) LoadCatalogue(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read catalogue file %s: %w", filePath, err)
	}

	var catalogue models.Catalogue
	if err := xml.Unmarshal(data, &catalogue); err != nil {
		return fmt.Errorf("failed to parse catalogue file %s: %w", filePath, err)
	}

	p.mu.Lock()
	if catalogue.Library == "true" {
		p.libraries[catalogue.ID] = &catalogue
		log.Printf("Loaded library: %s (revision %s)", catalogue.Name, catalogue.Revision)
	} else {
		p.catalogues[catalogue.ID] = &catalogue
		log.Printf("Loaded catalogue: %s (revision %s)", catalogue.Name, catalogue.Revision)
	}
	p.mu.Unlock()

	return nil
}

// GetGameSystem returns the loaded game system
func (p *Parser) GetGameSystem() *models.GameSystem {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.gameSystem
}

// GetCatalogue returns a catalogue by ID
func (p *Parser) GetCatalogue(id string) (*models.Catalogue, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	catalogue, exists := p.catalogues[id]
	return catalogue, exists
}

// GetLibrary returns a library by ID
func (p *Parser) GetLibrary(id string) (*models.Catalogue, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	library, exists := p.libraries[id]
	return library, exists
}

// GetAllCatalogues returns all catalogues
func (p *Parser) GetAllCatalogues() map[string]*models.Catalogue {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make(map[string]*models.Catalogue)
	for k, v := range p.catalogues {
		result[k] = v
	}
	return result
}

// GetAllLibraries returns all libraries
func (p *Parser) GetAllLibraries() map[string]*models.Catalogue {
	p.mu.RLock()
	defer p.mu.RUnlock()
	result := make(map[string]*models.Catalogue)
	for k, v := range p.libraries {
		result[k] = v
	}
	return result
}

// GetProfile returns a profile from sharedProfiles by ID
func (p *Parser) GetProfile(profileID string, catalogueID string) (*models.Profile, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	
	// Try regular catalogue first
	if cat, exists := p.catalogues[catalogueID]; exists {
		for i := range cat.SharedProfiles {
			if cat.SharedProfiles[i].ID == profileID {
				return &cat.SharedProfiles[i], true
			}
		}
	}
	
	// Try library
	if lib, exists := p.libraries[catalogueID]; exists {
		for i := range lib.SharedProfiles {
			if lib.SharedProfiles[i].ID == profileID {
				return &lib.SharedProfiles[i], true
			}
		}
	}
	
	return nil, false
}

// FindSelectionEntryByID finds a selectionEntry by ID across all loaded catalogues and libraries
func (p *Parser) FindSelectionEntryByID(id string) (*models.SelectionEntry, string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	// Search in libraries first (most common location)
	for libID, library := range p.libraries {
		if entry := findEntryInCatalogue(library, id); entry != nil {
			return entry, libID, true
		}
	}

	// Search in catalogues
	for catID, catalogue := range p.catalogues {
		if entry := findEntryInCatalogue(catalogue, id); entry != nil {
			return entry, catID, true
		}
	}

	return nil, "", false
}

// FindEntryLinkByID finds an entryLink by ID across all catalogues
func (p *Parser) FindEntryLinkByID(id string) (*models.EntryLink, string, bool) {
	p.mu.RLock()
	defer p.mu.RUnlock()

	for catID, catalogue := range p.catalogues {
		for i := range catalogue.EntryLinks {
			if catalogue.EntryLinks[i].ID == id {
				return &catalogue.EntryLinks[i], catID, true
			}
		}
	}

	return nil, "", false
}

// findEntryInCatalogue searches for a selectionEntry in a catalogue's sharedSelectionEntries
func findEntryInCatalogue(catalogue *models.Catalogue, id string) *models.SelectionEntry {
	for i := range catalogue.SharedSelectionEntries {
		if catalogue.SharedSelectionEntries[i].ID == id {
			return &catalogue.SharedSelectionEntries[i]
		}
		// Also search nested entries
		if entry := findEntryRecursive(&catalogue.SharedSelectionEntries[i], id); entry != nil {
			return entry
		}
	}

	// Search in entry groups
	for i := range catalogue.SharedSelectionEntryGroups {
		if entry := findEntryInGroup(&catalogue.SharedSelectionEntryGroups[i], id); entry != nil {
			return entry
		}
	}

	return nil
}

// findEntryRecursive searches recursively through nested selectionEntries
func findEntryRecursive(entry *models.SelectionEntry, id string) *models.SelectionEntry {
	for i := range entry.SelectionEntries {
		if entry.SelectionEntries[i].ID == id {
			return &entry.SelectionEntries[i]
		}
		if found := findEntryRecursive(&entry.SelectionEntries[i], id); found != nil {
			return found
		}
	}
	return nil
}

// findEntryInGroup searches for an entry in a selectionEntryGroup
func findEntryInGroup(group *models.SelectionEntryGroup, id string) *models.SelectionEntry {
	for i := range group.SelectionEntries {
		if group.SelectionEntries[i].ID == id {
			return &group.SelectionEntries[i]
		}
		if found := findEntryRecursive(&group.SelectionEntries[i], id); found != nil {
			return found
		}
	}

	for i := range group.SelectionEntryGroups {
		if found := findEntryInGroup(&group.SelectionEntryGroups[i], id); found != nil {
			return found
		}
	}

	return nil
}

