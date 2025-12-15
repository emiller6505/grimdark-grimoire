package models

// Response models for JSON API responses

// UnitResponse represents a unit in API responses
type UnitResponse struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"`
	Publication *PublicationInfo       `json:"publication,omitempty"`
	Profiles    *UnitProfiles          `json:"profiles"`
	Weapons     *WeaponSet             `json:"weapons"`
	Categories  []CategoryInfo         `json:"categories"`
	Rules       []RuleInfo             `json:"rules"`
	Costs       map[string]int         `json:"costs"`
	TieredCosts *TieredCosts           `json:"tieredCosts,omitempty"`
	Constraints *UnitConstraints       `json:"constraints,omitempty"`
	Faction     *FactionInfo          `json:"faction,omitempty"`
	Catalogue   *CatalogueInfo        `json:"catalogue,omitempty"`
}

// UnitProfiles contains all profile types for a unit
type UnitProfiles struct {
	Unit      *UnitProfile      `json:"unit,omitempty"`
	Abilities []AbilityProfile  `json:"abilities,omitempty"`
	Transport *TransportProfile `json:"transport,omitempty"`
}

// UnitProfile represents unit statistics
type UnitProfile struct {
	Movement         string `json:"movement"`
	Toughness        int    `json:"toughness"`
	Save             string `json:"save"`
	Wounds           int    `json:"wounds"`
	Leadership       string `json:"leadership"`
	ObjectiveControl int    `json:"objectiveControl"`
}

// AbilityProfile represents an ability
type AbilityProfile struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

// TransportProfile represents transport capacity
type TransportProfile struct {
	Capacity string `json:"capacity"`
}

// WeaponSet contains all weapons for a unit
type WeaponSet struct {
	Ranged []RangedWeapon `json:"ranged,omitempty"`
	Melee  []MeleeWeapon  `json:"melee,omitempty"`
}

// CategoryInfo represents a category association
type CategoryInfo struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Primary bool   `json:"primary"`
}

// RuleInfo represents a game rule reference
type RuleInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// TieredCosts represents costs that vary based on model count
type TieredCosts struct {
	BaseCost int          `json:"baseCost"` // Base cost (minimum model count)
	Tiers    []CostTier   `json:"tiers"`    // Cost tiers based on model count
}

// CostTier represents a cost tier
type CostTier struct {
	MinModels int `json:"minModels"` // Minimum number of models for this tier
	Cost      int `json:"cost"`      // Cost at this tier
}

// UnitConstraints represents unit selection constraints
type UnitConstraints struct {
	MaxPerRoster int `json:"maxPerRoster,omitempty"`
	MinPerRoster int `json:"minPerRoster,omitempty"`
	MaxPerForce  int `json:"maxPerForce,omitempty"`
	MinPerForce  int `json:"minPerForce,omitempty"`
}

// PublicationInfo represents publication metadata
type PublicationInfo struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	ShortName      string `json:"shortName,omitempty"`
	PublicationDate string `json:"publicationDate,omitempty"`
	Page           string `json:"page,omitempty"`
}

// FactionInfo represents faction information
type FactionInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// CatalogueInfo represents catalogue information
type CatalogueInfo struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Revision string `json:"revision"`
	Library  bool   `json:"library"`
}

// CatalogueResponse represents a catalogue in API responses
type CatalogueResponse struct {
	ID                string            `json:"id"`
	Name              string            `json:"name"`
	Revision          string            `json:"revision"`
	Library           bool              `json:"library"`
	GameSystemID      string            `json:"gameSystemId"`
	LinkedCatalogues  []CatalogueInfo   `json:"linkedCatalogues"`
	Units             []UnitSummary     `json:"units"`
	Publications      []PublicationInfo `json:"publications"`
}

// UnitSummary represents a summary of a unit (for lists)
type UnitSummary struct {
	ID          string       `json:"id"`
	Name        string       `json:"name"`
	TargetID    string       `json:"targetId,omitempty"`
	Costs       map[string]int `json:"costs"`
	TieredCosts *TieredCosts `json:"tieredCosts,omitempty"`
	Type        string       `json:"type,omitempty"`
}

// GameSystemResponse represents game system information
type GameSystemResponse struct {
	ID              string          `json:"id"`
	Name            string          `json:"name"`
	Revision        string          `json:"revision"`
	BattleScribeVersion string      `json:"battleScribeVersion"`
	ProfileTypes    []ProfileTypeInfo `json:"profileTypes"`
	Categories      []CategoryInfo  `json:"categories"`
	CostTypes       []CostTypeInfo  `json:"costTypes"`
}

// ProfileTypeInfo represents a profile type
type ProfileTypeInfo struct {
	ID                string   `json:"id"`
	Name              string   `json:"name"`
	Characteristics   []string `json:"characteristics"`
}

// CostTypeInfo represents a cost type
type CostTypeInfo struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	DefaultCostLimit string `json:"defaultCostLimit"`
	Hidden          bool   `json:"hidden"`
}

// FactionResponse represents faction information
type FactionResponse struct {
	Name     string         `json:"name"`
	Catalogues []CatalogueInfo `json:"catalogues"`
	Units    []UnitSummary  `json:"units"`
}

// SearchResponse represents search results
type SearchResponse struct {
	Query   string         `json:"query"`
	Type    string         `json:"type,omitempty"`
	Results []SearchResult `json:"results"`
	Total   int            `json:"total"`
}

// SearchResult represents a single search result
type SearchResult struct {
	Type     string `json:"type"` // "unit", "weapon", "ability"
	ID       string `json:"id"`
	Name     string `json:"name"`
	Summary  string `json:"summary,omitempty"`
}


