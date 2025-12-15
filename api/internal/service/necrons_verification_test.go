package service

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"grimoire-api/internal/cache"
	"grimoire-api/internal/parser"
)

// TestAllNecronUnits verifies that all Necron units match their XML data
func TestAllNecronUnits(t *testing.T) {
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

	service := NewUnitService(p, resolver, transformer, cache)

	// Find Necrons catalogue
	var necronsCatalogueID string
	for id, cat := range p.GetAllCatalogues() {
		if cat.Name == "Xenos - Necrons" {
			necronsCatalogueID = id
			break
		}
	}

	if necronsCatalogueID == "" {
		t.Fatal("Necrons catalogue not found")
	}

	necronsCatalogue, exists := p.GetCatalogue(necronsCatalogueID)
	if !exists {
		t.Fatal("Necrons catalogue not found")
	}

	// Filter out UI/metadata entries
	skipNames := map[string]bool{
		"Show/Hide Options": true,
		"Order of Battle":  true,
		"Detachment":       true,
	}

	// Get all unit entryLinks
	unitCount := 0
	successCount := 0
	errorCount := 0
	errors := []string{}

	for _, entryLink := range necronsCatalogue.EntryLinks {
		if entryLink.Type != "selectionEntry" {
			continue
		}

		// Skip UI/metadata entries
		if skipNames[entryLink.Name] {
			continue
		}

		unitCount++
		unitID := entryLink.ID
		unitName := entryLink.Name

		// Get unit via API
		unit, err := service.GetUnit(unitID)
		if err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("%s (ID: %s): Failed to get unit: %v", unitName, unitID, err))
			continue
		}

		if unit == nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("%s (ID: %s): Unit is nil", unitName, unitID))
			continue
		}

		// Deep validation: Compare XML data with API response
		hasErrors := false
		unitErrors := []string{}

		// Get the resolved entry from XML
		resolvedEntry, err := resolver.ResolveEntryLink(&entryLink, necronsCatalogueID)
		if err != nil {
			errorCount++
			errors = append(errors, fmt.Sprintf("%s (ID: %s): Failed to resolve entryLink: %v", unitName, unitID, err))
			continue
		}

		// Merge entryLink overrides
		mergedEntry := resolver.MergeEntryLinkWithSelectionEntry(&entryLink, resolvedEntry)

		// Validate ID
		if unit.ID != unitID {
			hasErrors = true
			unitErrors = append(unitErrors, fmt.Sprintf("ID mismatch: expected %s, got %s", unitID, unit.ID))
		}

		// Validate name
		if unit.Name == "" {
			hasErrors = true
			unitErrors = append(unitErrors, "Name is empty")
		} else if unit.Name != mergedEntry.Name {
			hasErrors = true
			unitErrors = append(unitErrors, fmt.Sprintf("Name mismatch: XML='%s', API='%s'", mergedEntry.Name, unit.Name))
		}

		// Validate profiles structure
		if unit.Profiles == nil {
			hasErrors = true
			unitErrors = append(unitErrors, "Profiles is nil")
		} else {
			// Check if Unit profile exists in XML (using same logic as transformer)
			hasUnitProfile := countUnitProfilesInEntry(mergedEntry) > 0
			// Unit profile might not exist for some unit types (like upgrades), so only warn if it's a unit type
			if mergedEntry.Type == "unit" && !hasUnitProfile && unit.Profiles.Unit == nil {
				hasErrors = true
				unitErrors = append(unitErrors, "Unit profile missing in XML but should exist for unit type")
			} else if hasUnitProfile && unit.Profiles.Unit == nil {
				hasErrors = true
				unitErrors = append(unitErrors, "Unit profile exists in XML but missing in API")
			}

			// Count abilities in XML (using same logic as transformer - only top-level profiles)
			// Note: The transformer only extracts abilities from the profiles passed to transformProfiles,
			// which are the direct profiles of the entry, not nested entries
			xmlAbilityCount := 0
			for _, profile := range mergedEntry.Profiles {
				if profile.TypeName == "Abilities" {
					xmlAbilityCount++
				}
			}
			apiAbilityCount := len(unit.Profiles.Abilities)
			if xmlAbilityCount != apiAbilityCount {
				hasErrors = true
				unitErrors = append(unitErrors, fmt.Sprintf("Ability count mismatch: XML=%d, API=%d (XML abilities: %v)", xmlAbilityCount, apiAbilityCount, getAbilityNames(mergedEntry)))
			}
		}

		// Validate costs structure
		if unit.Costs == nil {
			hasErrors = true
			unitErrors = append(unitErrors, "Costs is nil (should at least be empty map)")
		} else {
			// Check if costs exist in XML
			hasCostsInXML := len(mergedEntry.Costs) > 0
			hasCostsInAPI := len(unit.Costs) > 0
			// Some units might not have costs at the entryLink level but have them in the resolved entry
			if !hasCostsInXML {
				// Check resolved entry
				hasCostsInXML = len(resolvedEntry.Costs) > 0
			}
			// Costs might be empty for some unit types, so this is informational
			if hasCostsInXML && !hasCostsInAPI {
				hasErrors = true
				unitErrors = append(unitErrors, "Costs exist in XML but missing in API")
			}
		}

		// Validate weapons structure
		if unit.Weapons == nil {
			hasErrors = true
			unitErrors = append(unitErrors, "Weapons is nil")
		} else {
			// Count weapons in XML (this is approximate - we check if any weapon profiles exist)
			xmlRangedCount := 0
			xmlMeleeCount := 0
			countWeaponsInEntry(mergedEntry, &xmlRangedCount, &xmlMeleeCount, resolver, necronsCatalogueID)
			
			apiRangedCount := len(unit.Weapons.Ranged)
			apiMeleeCount := len(unit.Weapons.Melee)
			
			// Only report if there's a significant mismatch (allowing for some flexibility)
			if xmlRangedCount > 0 && apiRangedCount == 0 {
				hasErrors = true
				unitErrors = append(unitErrors, fmt.Sprintf("Ranged weapons missing: XML has %d, API has %d", xmlRangedCount, apiRangedCount))
			}
			if xmlMeleeCount > 0 && apiMeleeCount == 0 {
				hasErrors = true
				unitErrors = append(unitErrors, fmt.Sprintf("Melee weapons missing: XML has %d, API has %d", xmlMeleeCount, apiMeleeCount))
			}
		}

		// Validate categories
		if len(mergedEntry.CategoryLinks) > 0 && len(unit.Categories) == 0 {
			hasErrors = true
			unitErrors = append(unitErrors, "Categories missing in API response")
		}

		// Validate rules
		if len(mergedEntry.InfoLinks) > 0 && len(unit.Rules) == 0 {
			hasErrors = true
			unitErrors = append(unitErrors, "Rules missing in API response")
		}

		if hasErrors {
			errorCount++
			errors = append(errors, fmt.Sprintf("%s (ID: %s): %v", unitName, unitID, unitErrors))
		} else {
			successCount++
		}

		// Log first 10 units for detailed inspection
		if unitCount <= 10 {
			jsonData, _ := json.MarshalIndent(unit, "", "  ")
			filename := fmt.Sprintf("/tmp/necron_unit_%d_%s.json", unitCount, sanitizeFilename(unitName))
			os.WriteFile(filename, jsonData, 0644)
			t.Logf("Unit %d: %s (ID: %s) - written to %s", unitCount, unitName, unitID, filename)
		}
	}

	// Write summary
	summary := fmt.Sprintf(`
# Necron Units Verification Summary

Total Units: %d
Successful: %d
Errors: %d

## Errors:
%s
`, unitCount, successCount, errorCount, formatErrors(errors))

	os.WriteFile("/tmp/necrons_verification_summary.md", []byte(summary), 0644)
	t.Logf("Verification summary written to /tmp/necrons_verification_summary.md")

	if errorCount > 0 {
		t.Errorf("Found %d units with errors out of %d total units", errorCount, unitCount)
		t.Logf("First 20 errors:\n%s", formatErrors(errors[:min(20, len(errors))]))
	}
}

