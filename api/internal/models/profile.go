package models

import "encoding/xml"

// Profile represents a unit, weapon, ability, or transport profile
type Profile struct {
	XMLName        xml.Name         `xml:"profile"`
	ID             string           `xml:"id,attr"`
	Name           string           `xml:"name,attr"`
	TypeID         string           `xml:"typeId,attr"`
	TypeName       string           `xml:"typeName,attr"`
	Hidden         string           `xml:"hidden,attr"`
	PublicationID  string           `xml:"publicationId,attr"`
	Page           string           `xml:"page,attr"`
	Characteristics []Characteristic `xml:"characteristics>characteristic"`
	Modifiers      []Modifier       `xml:"modifiers>modifier"`
}

// Characteristic represents a single characteristic value
type Characteristic struct {
	XMLName xml.Name `xml:"characteristic"`
	Name    string   `xml:"name,attr"`
	TypeID  string   `xml:"typeId,attr"`
	Value   string   `xml:",chardata"`
}

