package models

// This file contains JSON response models for weapons
// The XML models are in profile.go

// Weapon represents a weapon (ranged or melee) in JSON responses
type Weapon struct {
	Name            string   `json:"name"`
	Range           string   `json:"range"`
	Attacks         string   `json:"attacks"`
	BallisticSkill  string   `json:"ballisticSkill,omitempty"`
	WeaponSkill     string   `json:"weaponSkill,omitempty"`
	Strength        string   `json:"strength"`
	ArmorPenetration string  `json:"armorPenetration"`
	Damage          string   `json:"damage"`
	Keywords        []string `json:"keywords"`
	Type            string   `json:"type"` // "ranged" or "melee"
}

// RangedWeapon represents a ranged weapon profile
type RangedWeapon struct {
	Name            string   `json:"name"`
	Range           string   `json:"range"`
	Attacks         string   `json:"attacks"`
	BallisticSkill  string   `json:"ballisticSkill"`
	Strength        string   `json:"strength"`
	ArmorPenetration string  `json:"armorPenetration"`
	Damage          string   `json:"damage"`
	Keywords        []string `json:"keywords"`
}

// MeleeWeapon represents a melee weapon profile
type MeleeWeapon struct {
	Name            string   `json:"name"`
	Range           string   `json:"range"`
	Attacks         string   `json:"attacks"`
	WeaponSkill     string   `json:"weaponSkill"`
	Strength        string   `json:"strength"`
	ArmorPenetration string  `json:"armorPenetration"`
	Damage          string   `json:"damage"`
	Keywords        []string `json:"keywords"`
}


