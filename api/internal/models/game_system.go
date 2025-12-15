package models

import "encoding/xml"

// GameSystem represents the root game system file
type GameSystem struct {
	XMLName          xml.Name         `xml:"gameSystem"`
	ID               string           `xml:"id,attr"`
	Name             string           `xml:"name,attr"`
	Revision         string           `xml:"revision,attr"`
	BattleScribeVersion string        `xml:"battleScribeVersion,attr"`
	Type             string           `xml:"type,attr"`
	Publications     []Publication    `xml:"publications>publication"`
	CostTypes        []CostType       `xml:"costTypes>costType"`
	ProfileTypes     []ProfileType    `xml:"profileTypes>profileType"`
	CategoryEntries  []CategoryEntry  `xml:"categoryEntries>categoryEntry"`
}

// Publication represents a source publication
type Publication struct {
	XMLName        xml.Name `xml:"publication"`
	ID             string   `xml:"id,attr"`
	Name           string   `xml:"name,attr"`
	ShortName      string   `xml:"shortName,attr"`
	Publisher      string   `xml:"publisher,attr"`
	PublicationDate string  `xml:"publicationDate,attr"`
}

// CostType represents a type of cost (points, Crusade points, etc.)
type CostType struct {
	XMLName         xml.Name `xml:"costType"`
	ID              string   `xml:"id,attr"`
	Name            string   `xml:"name,attr"`
	DefaultCostLimit string  `xml:"defaultCostLimit,attr"`
	Hidden         string   `xml:"hidden,attr"`
	Comment        string   `xml:"comment"`
}

// ProfileType defines the structure for different profile types
type ProfileType struct {
	XMLName            xml.Name           `xml:"profileType"`
	ID                 string             `xml:"id,attr"`
	Name               string             `xml:"name,attr"`
	CharacteristicTypes []CharacteristicType `xml:"characteristicTypes>characteristicType"`
}

// CharacteristicType defines a characteristic within a profile type
type CharacteristicType struct {
	XMLName xml.Name `xml:"characteristicType"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
}

// CategoryEntry represents a category (Character, Infantry, Vehicle, etc.)
type CategoryEntry struct {
	XMLName    xml.Name   `xml:"categoryEntry"`
	ID         string     `xml:"id,attr"`
	Name       string     `xml:"name,attr"`
	Hidden     string     `xml:"hidden,attr"`
	Constraints []Constraint `xml:"constraints>constraint"`
	Modifiers  []Modifier   `xml:"modifiers>modifier"`
	Comment    string     `xml:"comment"`
}

