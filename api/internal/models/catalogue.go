package models

import "encoding/xml"

// Catalogue represents a BattleScribe catalogue file
type Catalogue struct {
	XMLName                xml.Name              `xml:"catalogue"`
	ID                     string                `xml:"id,attr"`
	Name                   string                `xml:"name,attr"`
	Revision               string                `xml:"revision,attr"`
	BattleScribeVersion    string                `xml:"battleScribeVersion,attr"`
	Library                string                `xml:"library,attr"`
	GameSystemID           string                `xml:"gameSystemId,attr"`
	GameSystemRevision     string                `xml:"gameSystemRevision,attr"`
	Type                   string                `xml:"type,attr"`
	Publications           []Publication         `xml:"publications>publication"`
	CategoryEntries        []CategoryEntry       `xml:"categoryEntries>categoryEntry"`
	SharedSelectionEntries []SelectionEntry      `xml:"sharedSelectionEntries>selectionEntry"`
	SharedSelectionEntryGroups []SelectionEntryGroup `xml:"sharedSelectionEntryGroups>selectionEntryGroup"`
	SharedProfiles         []Profile             `xml:"sharedProfiles>profile"`
	EntryLinks             []EntryLink           `xml:"entryLinks>entryLink"`
	CatalogueLinks         []CatalogueLink       `xml:"catalogueLinks>catalogueLink"`
}

// CatalogueLink links to another catalogue (typically a library)
type CatalogueLink struct {
	XMLName          xml.Name `xml:"catalogueLink"`
	ID               string   `xml:"id,attr"`
	Name             string   `xml:"name,attr"`
	TargetID         string   `xml:"targetId,attr"`
	Type             string   `xml:"type,attr"`
	ImportRootEntries string  `xml:"importRootEntries,attr"`
}

// EntryLink references a selectionEntry from another catalogue
type EntryLink struct {
	XMLName         xml.Name         `xml:"entryLink"`
	ID              string           `xml:"id,attr"`
	Name            string           `xml:"name,attr"`
	Hidden          string           `xml:"hidden,attr"`
	Collective      string           `xml:"collective,attr"`
	Import          string           `xml:"import,attr"`
	TargetID        string           `xml:"targetId,attr"`
	Type            string           `xml:"type,attr"`
	CategoryLinks   []CategoryLink   `xml:"categoryLinks>categoryLink"`
	Costs           []Cost           `xml:"costs>cost"`
	Constraints     []Constraint     `xml:"constraints>constraint"`
	Modifiers       []Modifier       `xml:"modifiers>modifier"`
	EntryLinks      []EntryLink      `xml:"entryLinks>entryLink"`
}

// CategoryLink associates an entry with a category
type CategoryLink struct {
	XMLName xml.Name `xml:"categoryLink"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	TargetID string  `xml:"targetId,attr"`
	Primary string   `xml:"primary,attr"`
}

