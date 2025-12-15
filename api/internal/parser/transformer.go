package parser

import (
	"strconv"
	"strings"

	"grimoire-api/internal/models"
)

// Transformer converts XML models to JSON-friendly response models
type Transformer struct {
	resolver *LinkResolver
}

// NewTransformer creates a new transformer
func NewTransformer(resolver *LinkResolver) *Transformer {
	return &Transformer{resolver: resolver}
}

// TransformUnit transforms a SelectionEntry to a UnitResponse
func (t *Transformer) TransformUnit(entry *models.SelectionEntry, catalogueID string) *models.UnitResponse {
	response := &models.UnitResponse{
		ID:   entry.ID,
		Name: entry.Name,
		Type: entry.Type,
	}

	// Set publication info
	if entry.PublicationID != "" {
		pubInfo := &models.PublicationInfo{
			ID:   entry.PublicationID,
			Page: entry.Page,
		}
		
		// Try to resolve full publication details from catalogue
		if catalogueID != "" {
			if cat, exists := t.resolver.parser.GetCatalogue(catalogueID); exists {
				for _, pub := range cat.Publications {
					if pub.ID == entry.PublicationID {
						pubInfo.Name = pub.Name
						pubInfo.ShortName = pub.ShortName
						pubInfo.PublicationDate = pub.PublicationDate
						break
					}
				}
			} else if lib, exists := t.resolver.parser.GetLibrary(catalogueID); exists {
				for _, pub := range lib.Publications {
					if pub.ID == entry.PublicationID {
						pubInfo.Name = pub.Name
						pubInfo.ShortName = pub.ShortName
						pubInfo.PublicationDate = pub.PublicationDate
						break
					}
				}
			}
		}
		
		response.Publication = pubInfo
	}

	// Transform profiles (also check nested entries for Unit profiles)
	response.Profiles = t.transformProfiles(entry.Profiles, entry, catalogueID)

	// Transform weapons (pass catalogueID for proper entryLink resolution)
	response.Weapons = t.transformWeaponsWithCatalogue(entry, catalogueID)

	// Transform categories
	response.Categories = t.transformCategories(entry.CategoryLinks)

	// Transform rules
	response.Rules = t.transformRules(entry.InfoLinks)

	// Transform costs
	response.Costs = t.TransformCosts(entry.Costs)

	// Transform constraints
	response.Constraints = t.transformConstraints(entry.Constraints)

	// Set catalogue info
	if catalogueID != "" {
		// Try regular catalogue first
		if cat, exists := t.resolver.parser.GetCatalogue(catalogueID); exists {
			response.Catalogue = &models.CatalogueInfo{
				ID:       cat.ID,
				Name:     cat.Name,
				Revision: cat.Revision,
				Library:  cat.Library == "true",
			}
		} else if lib, exists := t.resolver.parser.GetLibrary(catalogueID); exists {
			// If not found, try library
			response.Catalogue = &models.CatalogueInfo{
				ID:       lib.ID,
				Name:     lib.Name,
				Revision: lib.Revision,
				Library:  true,
			}
		}
	}

	// Extract faction from categories
	for _, catLink := range entry.CategoryLinks {
		if strings.HasPrefix(catLink.Name, "Faction:") {
			response.Faction = &models.FactionInfo{
				ID:   catLink.TargetID,
				Name: catLink.Name,
			}
			break
		}
	}

	return response
}

// transformProfiles transforms profiles to UnitProfiles
// Also checks nested selectionEntries for profiles (some units have profiles in child entries)
func (t *Transformer) transformProfiles(profiles []models.Profile, entry *models.SelectionEntry, catalogueID string) *models.UnitProfiles {
	result := &models.UnitProfiles{
		Abilities: make([]models.AbilityProfile, 0),
	}

	// First, check direct profiles
	for _, profile := range profiles {
		switch profile.TypeName {
		case "Unit":
			if result.Unit == nil { // Only set if not already found in nested entries
				result.Unit = t.transformUnitProfile(profile)
			}
		case "Abilities":
			result.Abilities = append(result.Abilities, t.transformAbilityProfile(profile))
		case "Transport":
			result.Transport = t.transformTransportProfile(profile)
		}
	}

	// Also check nested selectionEntries for Unit profiles (e.g., Wraithguard)
	if result.Unit == nil {
		for i := range entry.SelectionEntries {
			subEntry := &entry.SelectionEntries[i]
			for _, profile := range subEntry.Profiles {
				if profile.TypeName == "Unit" {
					result.Unit = t.transformUnitProfile(profile)
					break // Found it, no need to continue
				}
			}
			// If still not found, check deeper nested entries
			if result.Unit == nil {
				for j := range subEntry.SelectionEntries {
					deepEntry := &subEntry.SelectionEntries[j]
					for _, profile := range deepEntry.Profiles {
						if profile.TypeName == "Unit" {
							result.Unit = t.transformUnitProfile(profile)
							break
						}
					}
					if result.Unit != nil {
						break
					}
				}
			}
			if result.Unit != nil {
				break
			}
		}
	}

	// Also check selectionEntryGroups for Unit profiles (e.g., Boyz, Gretchin, Grot Tanks)
	// Use a recursive helper to search all possible nested structures
	if result.Unit == nil {
		result.Unit = t.findUnitProfileInGroups(entry.SelectionEntryGroups, catalogueID)
	}
	
	// Also check selectionEntryGroups for Unit profiles (e.g., Boyz, Gretchin) - keep original code as fallback
	if result.Unit == nil {
		for i := range entry.SelectionEntryGroups {
			group := &entry.SelectionEntryGroups[i]
			// Check SelectionEntries in the group
			for j := range group.SelectionEntries {
				subEntry := &group.SelectionEntries[j]
				for _, profile := range subEntry.Profiles {
					if profile.TypeName == "Unit" {
						result.Unit = t.transformUnitProfile(profile)
						break
					}
				}
				// Also check entryLinks within selectionEntries (e.g., Gretchin unit profile)
				if result.Unit == nil {
					for k := range subEntry.EntryLinks {
						entryLink := &subEntry.EntryLinks[k]
						if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
							resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
							if err == nil {
								// Check if resolved entry has unit profile
								for _, profile := range resolvedEntry.Profiles {
									if profile.TypeName == "Unit" {
										result.Unit = t.transformUnitProfile(profile)
										break
									}
								}
								// Also check nested selectionEntries in resolved entry
								if result.Unit == nil {
									for l := range resolvedEntry.SelectionEntries {
										deepEntry := &resolvedEntry.SelectionEntries[l]
										for _, profile := range deepEntry.Profiles {
											if profile.TypeName == "Unit" {
												result.Unit = t.transformUnitProfile(profile)
												break
											}
										}
										if result.Unit != nil {
											break
										}
									}
								}
							}
						}
						if result.Unit != nil {
							break
						}
					}
				}
				// Check nested selectionEntryGroups within selectionEntries (e.g., Grot Tanks - "Wargear" group)
				if result.Unit == nil {
					for k := range subEntry.SelectionEntryGroups {
						nestedGroupInEntry := &subEntry.SelectionEntryGroups[k]
						// Check entryLinks in nested groups within selectionEntries
						for l := range nestedGroupInEntry.EntryLinks {
							entryLink := &nestedGroupInEntry.EntryLinks[l]
							if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
								resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
								if err == nil {
									for _, profile := range resolvedEntry.Profiles {
										if profile.TypeName == "Unit" {
											result.Unit = t.transformUnitProfile(profile)
											break
										}
									}
									// Also check nested selectionEntries in resolved entry
									if result.Unit == nil {
										for m := range resolvedEntry.SelectionEntries {
											deepEntry := &resolvedEntry.SelectionEntries[m]
											for _, profile := range deepEntry.Profiles {
												if profile.TypeName == "Unit" {
													result.Unit = t.transformUnitProfile(profile)
													break
												}
											}
											if result.Unit != nil {
												break
											}
										}
									}
								}
							}
							if result.Unit != nil {
								break
							}
						}
						if result.Unit != nil {
							break
						}
					}
				}
				if result.Unit != nil {
					break
				}
			}
			if result.Unit != nil {
				break
			}
			// Check EntryLinks in the group (e.g., Gretchin unit profile is in an entryLink)
			for j := range group.EntryLinks {
				entryLink := &group.EntryLinks[j]
				if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
					resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
					if err == nil {
						// Check if resolved entry has unit profile
						for _, profile := range resolvedEntry.Profiles {
							if profile.TypeName == "Unit" {
								result.Unit = t.transformUnitProfile(profile)
								break
							}
						}
						// Also check nested selectionEntries in resolved entry
						if result.Unit == nil {
							for k := range resolvedEntry.SelectionEntries {
								subEntry := &resolvedEntry.SelectionEntries[k]
								for _, profile := range subEntry.Profiles {
									if profile.TypeName == "Unit" {
										result.Unit = t.transformUnitProfile(profile)
										break
									}
								}
								if result.Unit != nil {
									break
								}
							}
						}
					}
				}
				if result.Unit != nil {
					break
				}
			}
			if result.Unit != nil {
				break
			}
			// Also check nested groups
			for j := range group.SelectionEntryGroups {
				nestedGroup := &group.SelectionEntryGroups[j]
				for k := range nestedGroup.SelectionEntries {
					subEntry := &nestedGroup.SelectionEntries[k]
					for _, profile := range subEntry.Profiles {
						if profile.TypeName == "Unit" {
							result.Unit = t.transformUnitProfile(profile)
							break
						}
					}
					// Check entryLinks within nested group selectionEntries (e.g., Grot Tanks)
					if result.Unit == nil {
						for l := range subEntry.EntryLinks {
							entryLink := &subEntry.EntryLinks[l]
							if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
								resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
								if err == nil {
									for _, profile := range resolvedEntry.Profiles {
										if profile.TypeName == "Unit" {
											result.Unit = t.transformUnitProfile(profile)
											break
										}
									}
								}
							}
							if result.Unit != nil {
								break
							}
						}
					}
					if result.Unit != nil {
						break
					}
				}
				// Check entryLinks directly in nested groups (e.g., Grot Tanks - "Ramshackle hull" is here)
				if result.Unit == nil {
					for k := range nestedGroup.EntryLinks {
						entryLink := &nestedGroup.EntryLinks[k]
						if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
							resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
							if err == nil {
								for _, profile := range resolvedEntry.Profiles {
									if profile.TypeName == "Unit" {
										result.Unit = t.transformUnitProfile(profile)
										break
									}
								}
								// Also check nested selectionEntries in resolved entry
								if result.Unit == nil {
									for l := range resolvedEntry.SelectionEntries {
										deepEntry := &resolvedEntry.SelectionEntries[l]
										for _, profile := range deepEntry.Profiles {
											if profile.TypeName == "Unit" {
												result.Unit = t.transformUnitProfile(profile)
												break
											}
										}
										if result.Unit != nil {
											break
										}
									}
								}
							}
						}
						if result.Unit != nil {
							break
						}
					}
				}
				// Also recurse into deeper nested groups within nested groups
				if result.Unit == nil {
					for k := range nestedGroup.SelectionEntryGroups {
						deepNestedGroup := &nestedGroup.SelectionEntryGroups[k]
						for l := range deepNestedGroup.EntryLinks {
							entryLink := &deepNestedGroup.EntryLinks[l]
							if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
								resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
								if err == nil {
									for _, profile := range resolvedEntry.Profiles {
										if profile.TypeName == "Unit" {
											result.Unit = t.transformUnitProfile(profile)
											break
										}
									}
								}
							}
							if result.Unit != nil {
								break
							}
						}
						if result.Unit != nil {
							break
						}
					}
				}
				if result.Unit != nil {
					break
				}
			}
			if result.Unit != nil {
				break
			}
		}
	}

	// Also check infoLinks for Unit profiles (e.g., Grot Tanks [Legends])
	// Some units reference Unit profiles via infoLinks with type="profile"
	// Check infoLinks in the entry itself and recursively in nested entries
	if result.Unit == nil && catalogueID != "" {
		result.Unit = t.findUnitProfileInInfoLinks(entry, catalogueID)
	}

	return result
}

// findUnitProfileInInfoLinks recursively searches infoLinks for Unit profiles
func (t *Transformer) findUnitProfileInInfoLinks(entry *models.SelectionEntry, catalogueID string) *models.UnitProfile {
	// Check infoLinks in this entry
	for _, infoLink := range entry.InfoLinks {
		if infoLink.Type == "profile" && infoLink.TargetID != "" {
			// Try to resolve the profile from sharedProfiles
			if profile, found := t.resolver.parser.GetProfile(infoLink.TargetID, catalogueID); found {
				if profile.TypeName == "Unit" {
					return t.transformUnitProfile(*profile)
				}
			}
		}
	}

	// Check nested selectionEntries
	for i := range entry.SelectionEntries {
		if unitProfile := t.findUnitProfileInInfoLinks(&entry.SelectionEntries[i], catalogueID); unitProfile != nil {
			return unitProfile
		}
	}

	// Check selectionEntryGroups
	for i := range entry.SelectionEntryGroups {
		if unitProfile := t.findUnitProfileInInfoLinksInGroup(&entry.SelectionEntryGroups[i], catalogueID); unitProfile != nil {
			return unitProfile
		}
	}

	return nil
}

// findUnitProfileInInfoLinksInGroup recursively searches infoLinks in a group for Unit profiles
func (t *Transformer) findUnitProfileInInfoLinksInGroup(group *models.SelectionEntryGroup, catalogueID string) *models.UnitProfile {
	// Check SelectionEntries in the group
	for i := range group.SelectionEntries {
		if unitProfile := t.findUnitProfileInInfoLinks(&group.SelectionEntries[i], catalogueID); unitProfile != nil {
			return unitProfile
		}
	}

	// Recursively check nested groups
	for i := range group.SelectionEntryGroups {
		if unitProfile := t.findUnitProfileInInfoLinksInGroup(&group.SelectionEntryGroups[i], catalogueID); unitProfile != nil {
			return unitProfile
		}
	}

	return nil
}

// findUnitProfileInGroups recursively searches selectionEntryGroups for Unit profiles
func (t *Transformer) findUnitProfileInGroups(groups []models.SelectionEntryGroup, catalogueID string) *models.UnitProfile {
	for i := range groups {
		group := &groups[i]
		
		// Check SelectionEntries in the group
		for j := range group.SelectionEntries {
			subEntry := &group.SelectionEntries[j]
			// Check profiles directly in selectionEntry
			for _, profile := range subEntry.Profiles {
				if profile.TypeName == "Unit" {
					return t.transformUnitProfile(profile)
				}
			}
			// Check entryLinks in selectionEntry
			for k := range subEntry.EntryLinks {
				entryLink := &subEntry.EntryLinks[k]
				if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
					if resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID); err == nil {
						for _, profile := range resolvedEntry.Profiles {
							if profile.TypeName == "Unit" {
								return t.transformUnitProfile(profile)
							}
						}
						// Check nested selectionEntries in resolved entry
						for l := range resolvedEntry.SelectionEntries {
							deepEntry := &resolvedEntry.SelectionEntries[l]
							for _, profile := range deepEntry.Profiles {
								if profile.TypeName == "Unit" {
									return t.transformUnitProfile(profile)
								}
							}
						}
					}
				}
			}
			// Recursively check nested groups in selectionEntry
			if unitProfile := t.findUnitProfileInGroups(subEntry.SelectionEntryGroups, catalogueID); unitProfile != nil {
				return unitProfile
			}
		}
		
		// Check EntryLinks directly in the group
		for j := range group.EntryLinks {
			entryLink := &group.EntryLinks[j]
			if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
				if resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID); err == nil {
					for _, profile := range resolvedEntry.Profiles {
						if profile.TypeName == "Unit" {
							return t.transformUnitProfile(profile)
						}
					}
					// Check nested selectionEntries in resolved entry
					for k := range resolvedEntry.SelectionEntries {
						subEntry := &resolvedEntry.SelectionEntries[k]
						for _, profile := range subEntry.Profiles {
							if profile.TypeName == "Unit" {
								return t.transformUnitProfile(profile)
							}
						}
					}
				}
			}
		}
		
		// Recursively check nested groups
		if unitProfile := t.findUnitProfileInGroups(group.SelectionEntryGroups, catalogueID); unitProfile != nil {
			return unitProfile
		}
	}
	return nil
}

// transformUnitProfile transforms a Unit profile
func (t *Transformer) transformUnitProfile(profile models.Profile) *models.UnitProfile {
	unit := &models.UnitProfile{}
	charMap := make(map[string]string)

	for _, char := range profile.Characteristics {
		charMap[char.Name] = char.Value
	}

	unit.Movement = charMap["M"]
	unit.Toughness = parseInt(charMap["T"])
	unit.Save = charMap["SV"]
	unit.Wounds = parseInt(charMap["W"])
	unit.Leadership = charMap["LD"]
	unit.ObjectiveControl = parseInt(charMap["OC"])

	return unit
}

// transformAbilityProfile transforms an Ability profile
func (t *Transformer) transformAbilityProfile(profile models.Profile) models.AbilityProfile {
	ability := models.AbilityProfile{
		Name: profile.Name,
	}

	for _, char := range profile.Characteristics {
		if char.Name == "Description" {
			ability.Description = char.Value
		}
	}

	return ability
}

// transformTransportProfile transforms a Transport profile
func (t *Transformer) transformTransportProfile(profile models.Profile) *models.TransportProfile {
	transport := &models.TransportProfile{}

	for _, char := range profile.Characteristics {
		if char.Name == "Capacity" {
			transport.Capacity = char.Value
		}
	}

	return transport
}

// transformWeapons extracts weapons from a selectionEntry (backwards compatibility)
func (t *Transformer) transformWeapons(entry *models.SelectionEntry) *models.WeaponSet {
	return t.transformWeaponsWithCatalogue(entry, "")
}

// transformWeaponsWithCatalogue extracts weapons from a selectionEntry with catalogue context
func (t *Transformer) transformWeaponsWithCatalogue(entry *models.SelectionEntry, catalogueID string) *models.WeaponSet {
	weapons := &models.WeaponSet{
		Ranged: make([]models.RangedWeapon, 0),
		Melee:  make([]models.MeleeWeapon, 0),
	}

	// Search through selectionEntries for weapons
	t.extractWeaponsRecursive(entry, weapons, catalogueID)

	return weapons
}

// extractWeaponsRecursive recursively extracts weapons from selectionEntries
func (t *Transformer) extractWeaponsRecursive(entry *models.SelectionEntry, weapons *models.WeaponSet, catalogueID string) {
	// Check direct selectionEntries
	for i := range entry.SelectionEntries {
		subEntry := &entry.SelectionEntries[i]
		
		// Check if this entry has weapon profiles
		for _, profile := range subEntry.Profiles {
			if profile.TypeName == "Ranged Weapons" {
				weapons.Ranged = append(weapons.Ranged, t.transformRangedWeapon(profile))
			} else if profile.TypeName == "Melee Weapons" {
				weapons.Melee = append(weapons.Melee, t.transformMeleeWeapon(profile))
			}
		}

		// Recurse into nested entries
		t.extractWeaponsRecursive(subEntry, weapons, catalogueID)
	}

	// Check EntryLinks directly in the entry
	for i := range entry.EntryLinks {
		entryLink := &entry.EntryLinks[i]
		if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
			resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
			if err == nil {
				// Extract weapons from resolved entry
				t.extractWeaponsFromEntry(resolvedEntry, weapons, catalogueID)
			}
		}
	}

	// Check selectionEntryGroups
	for i := range entry.SelectionEntryGroups {
		group := &entry.SelectionEntryGroups[i]
		
		// Check SelectionEntries in the group
		for j := range group.SelectionEntries {
			subEntry := &group.SelectionEntries[j]
			for _, profile := range subEntry.Profiles {
				if profile.TypeName == "Ranged Weapons" {
					weapons.Ranged = append(weapons.Ranged, t.transformRangedWeapon(profile))
				} else if profile.TypeName == "Melee Weapons" {
					weapons.Melee = append(weapons.Melee, t.transformMeleeWeapon(profile))
				}
			}
			t.extractWeaponsRecursive(subEntry, weapons, catalogueID)
		}
		
		// Check EntryLinks in the group (this is the missing piece!)
		for j := range group.EntryLinks {
			entryLink := &group.EntryLinks[j]
			if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
				resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
				if err == nil {
					// Extract weapons from resolved entry
					t.extractWeaponsFromEntry(resolvedEntry, weapons, catalogueID)
				}
			}
		}
		
		// Recurse into nested selectionEntryGroups
		for j := range group.SelectionEntryGroups {
			t.extractWeaponsFromGroup(&group.SelectionEntryGroups[j], weapons, catalogueID)
		}
	}
}

// extractWeaponsFromEntry extracts weapons from a selectionEntry (used for resolved entryLinks)
func (t *Transformer) extractWeaponsFromEntry(entry *models.SelectionEntry, weapons *models.WeaponSet, catalogueID string) {
	// Check if this entry has weapon profiles
	for _, profile := range entry.Profiles {
		if profile.TypeName == "Ranged Weapons" {
			weapons.Ranged = append(weapons.Ranged, t.transformRangedWeapon(profile))
		} else if profile.TypeName == "Melee Weapons" {
			weapons.Melee = append(weapons.Melee, t.transformMeleeWeapon(profile))
		}
	}
	
	// Recurse into nested structures
	t.extractWeaponsRecursive(entry, weapons, catalogueID)
}

// extractWeaponsFromGroup extracts weapons from a selectionEntryGroup
func (t *Transformer) extractWeaponsFromGroup(group *models.SelectionEntryGroup, weapons *models.WeaponSet, catalogueID string) {
	// Check SelectionEntries in the group
	for i := range group.SelectionEntries {
		subEntry := &group.SelectionEntries[i]
		for _, profile := range subEntry.Profiles {
			if profile.TypeName == "Ranged Weapons" {
				weapons.Ranged = append(weapons.Ranged, t.transformRangedWeapon(profile))
			} else if profile.TypeName == "Melee Weapons" {
				weapons.Melee = append(weapons.Melee, t.transformMeleeWeapon(profile))
			}
		}
		t.extractWeaponsRecursive(subEntry, weapons, catalogueID)
	}
	
	// Check EntryLinks in the group
	for i := range group.EntryLinks {
		entryLink := &group.EntryLinks[i]
		if entryLink.Type == "selectionEntry" || entryLink.Type == "upgrade" {
			resolvedEntry, err := t.resolver.ResolveEntryLink(entryLink, catalogueID)
			if err == nil {
				t.extractWeaponsFromEntry(resolvedEntry, weapons, catalogueID)
			}
		}
	}
	
	// Recurse into nested groups
	for i := range group.SelectionEntryGroups {
		t.extractWeaponsFromGroup(&group.SelectionEntryGroups[i], weapons, catalogueID)
	}
}

// transformRangedWeapon transforms a Ranged Weapons profile
func (t *Transformer) transformRangedWeapon(profile models.Profile) models.RangedWeapon {
	weapon := models.RangedWeapon{
		Name: profile.Name,
	}

	charMap := make(map[string]string)
	for _, char := range profile.Characteristics {
		charMap[char.Name] = char.Value
	}

	weapon.Range = charMap["Range"]
	weapon.Attacks = charMap["A"]
	weapon.BallisticSkill = charMap["BS"]
	weapon.Strength = charMap["S"]
	weapon.ArmorPenetration = charMap["AP"]
	weapon.Damage = charMap["D"]
	
	if keywords := charMap["Keywords"]; keywords != "" {
		weapon.Keywords = parseKeywords(keywords)
	}

	return weapon
}

// transformMeleeWeapon transforms a Melee Weapons profile
func (t *Transformer) transformMeleeWeapon(profile models.Profile) models.MeleeWeapon {
	weapon := models.MeleeWeapon{
		Name: profile.Name,
	}

	charMap := make(map[string]string)
	for _, char := range profile.Characteristics {
		charMap[char.Name] = char.Value
	}

	weapon.Range = charMap["Range"]
	weapon.Attacks = charMap["A"]
	weapon.WeaponSkill = charMap["WS"]
	weapon.Strength = charMap["S"]
	weapon.ArmorPenetration = charMap["AP"]
	weapon.Damage = charMap["D"]
	
	if keywords := charMap["Keywords"]; keywords != "" {
		weapon.Keywords = parseKeywords(keywords)
	}

	return weapon
}

// transformCategories transforms categoryLinks
func (t *Transformer) transformCategories(categoryLinks []models.CategoryLink) []models.CategoryInfo {
	categories := make([]models.CategoryInfo, 0, len(categoryLinks))
	for _, catLink := range categoryLinks {
		categories = append(categories, models.CategoryInfo{
			ID:      catLink.TargetID,
			Name:    catLink.Name,
			Primary: catLink.Primary == "true",
		})
	}
	return categories
}

// transformRules transforms infoLinks
func (t *Transformer) transformRules(infoLinks []models.InfoLink) []models.RuleInfo {
	rules := make([]models.RuleInfo, 0, len(infoLinks))
	for _, infoLink := range infoLinks {
		rules = append(rules, models.RuleInfo{
			ID:   infoLink.TargetID,
			Name: infoLink.Name,
		})
	}
	return rules
}

// TransformCosts transforms costs to a map (public method for use in services)
func (t *Transformer) TransformCosts(costs []models.Cost) map[string]int {
	costMap := make(map[string]int)
	for _, cost := range costs {
		if cost.Value == "" {
			continue
		}
		if value, err := strconv.Atoi(cost.Value); err == nil {
			costMap[cost.Name] = value
		}
		// Note: Silently ignore parse errors - some costs might be non-numeric (e.g., "Variable")
	}
	return costMap
}

// transformConstraints extracts constraint information
func (t *Transformer) transformConstraints(constraints []models.Constraint) *models.UnitConstraints {
	result := &models.UnitConstraints{}
	
	for _, constraint := range constraints {
		value := parseInt(constraint.Value)
		if value < 0 {
			continue // -1 means unlimited
		}

		switch {
		case constraint.Type == "max" && constraint.Scope == "roster":
			result.MaxPerRoster = value
		case constraint.Type == "min" && constraint.Scope == "roster":
			result.MinPerRoster = value
		case constraint.Type == "max" && constraint.Scope == "force":
			result.MaxPerForce = value
		case constraint.Type == "min" && constraint.Scope == "force":
			result.MinPerForce = value
		}
	}

	if result.MaxPerRoster == 0 && result.MinPerRoster == 0 && 
	   result.MaxPerForce == 0 && result.MinPerForce == 0 {
		return nil
	}

	return result
}

// TransformCatalogue transforms a Catalogue to CatalogueResponse
func (t *Transformer) TransformCatalogue(catalogue *models.Catalogue) *models.CatalogueResponse {
	response := &models.CatalogueResponse{
		ID:           catalogue.ID,
		Name:         catalogue.Name,
		Revision:     catalogue.Revision,
		Library:      catalogue.Library == "true",
		GameSystemID: catalogue.GameSystemID,
	}

	// Transform linked catalogues
	linked := t.resolver.ResolveCatalogueLinks(catalogue)
	response.LinkedCatalogues = make([]models.CatalogueInfo, 0, len(linked))
	for _, cat := range linked {
		response.LinkedCatalogues = append(response.LinkedCatalogues, models.CatalogueInfo{
			ID:       cat.ID,
			Name:     cat.Name,
			Revision: cat.Revision,
			Library:  cat.Library == "true",
		})
	}

	// Transform publications
	response.Publications = make([]models.PublicationInfo, 0, len(catalogue.Publications))
	for _, pub := range catalogue.Publications {
		response.Publications = append(response.Publications, models.PublicationInfo{
			ID:              pub.ID,
			Name:            pub.Name,
			ShortName:       pub.ShortName,
			PublicationDate: pub.PublicationDate,
		})
	}

	// Transform entryLinks to unit summaries
	response.Units = make([]models.UnitSummary, 0, len(catalogue.EntryLinks))
	for _, entryLink := range catalogue.EntryLinks {
		summary := models.UnitSummary{
			ID:       entryLink.ID,
			Name:     entryLink.Name,
			TargetID: entryLink.TargetID,
			Type:     entryLink.Type,
			Costs:    t.TransformCosts(entryLink.Costs),
		}
		response.Units = append(response.Units, summary)
	}

	return response
}

// Helper functions

// parseInt safely parses an integer string, returning 0 on error
// Note: Returns 0 for both "0" and parse errors - caller should validate if needed
func parseInt(s string) int {
	if s == "" {
		return 0
	}
	val, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}
	return val
}

func parseKeywords(keywords string) []string {
	if keywords == "" || keywords == "-" {
		return []string{}
	}
	parts := strings.Split(keywords, ",")
	result := make([]string, 0, len(parts))
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

