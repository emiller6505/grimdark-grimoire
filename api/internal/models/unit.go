package models

import "encoding/xml"

// SelectionEntry represents a unit, upgrade, or other selectable item
type SelectionEntry struct {
	XMLName              xml.Name              `xml:"selectionEntry"`
	ID                   string                `xml:"id,attr"`
	Name                 string                `xml:"name,attr"`
	Hidden               string                `xml:"hidden,attr"`
	Collective           string                `xml:"collective,attr"`
	Import               string                `xml:"import,attr"`
	Type                 string                `xml:"type,attr"`
	PublicationID        string                `xml:"publicationId,attr"`
	Page                 string                `xml:"page,attr"`
	SortIndex            string                `xml:"sortIndex,attr"`
	Profiles             []Profile             `xml:"profiles>profile"`
	InfoLinks            []InfoLink            `xml:"infoLinks>infoLink"`
	CategoryLinks        []CategoryLink        `xml:"categoryLinks>categoryLink"`
	SelectionEntries     []SelectionEntry      `xml:"selectionEntries>selectionEntry"`
	SelectionEntryGroups []SelectionEntryGroup `xml:"selectionEntryGroups>selectionEntryGroup"`
	EntryLinks           []EntryLink           `xml:"entryLinks>entryLink"`
	Costs                []Cost                `xml:"costs>cost"`
	Constraints          []Constraint          `xml:"constraints>constraint"`
	Modifiers            []Modifier            `xml:"modifiers>modifier"`
	ModifierGroups       []ModifierGroup       `xml:"modifierGroups>modifierGroup"`
	Comment              string                `xml:"comment"`
}

// SelectionEntryGroup groups related selection entries
type SelectionEntryGroup struct {
	XMLName              xml.Name              `xml:"selectionEntryGroup"`
	ID                   string                `xml:"id,attr"`
	Name                 string                `xml:"name,attr"`
	Hidden               string                `xml:"hidden,attr"`
	Collapsible          string                `xml:"collapsible,attr"`
	Flatten              string                `xml:"flatten,attr"`
	SortIndex            string                `xml:"sortIndex,attr"`
	SelectionEntries     []SelectionEntry      `xml:"selectionEntries>selectionEntry"`
	SelectionEntryGroups []SelectionEntryGroup `xml:"selectionEntryGroups>selectionEntryGroup"`
	EntryLinks           []EntryLink           `xml:"entryLinks>entryLink"`
	Constraints          []Constraint          `xml:"constraints>constraint"`
	Modifiers            []Modifier            `xml:"modifiers>modifier"`
}

// Constraint enforces game rules
type Constraint struct {
	XMLName              xml.Name `xml:"constraint"`
	ID                   string   `xml:"id,attr"`
	Type                 string   `xml:"type,attr"`
	Value                string   `xml:"value,attr"`
	Field                string   `xml:"field,attr"`
	Scope                string   `xml:"scope,attr"`
	Shared               string   `xml:"shared,attr"`
	IncludeChildSelections string `xml:"includeChildSelections,attr"`
	IncludeChildForces    string  `xml:"includeChildForces,attr"`
	PercentValue         string   `xml:"percentValue,attr"`
	Negative             string   `xml:"negative,attr"`
}

// Modifier applies conditional changes
type Modifier struct {
	XMLName         xml.Name       `xml:"modifier"`
	ID              string         `xml:"id,attr"`
	Type            string         `xml:"type,attr"`
	Value           string         `xml:"value,attr"`
	Field           string         `xml:"field,attr"`
	Scope           string         `xml:"scope,attr"`
	Affects         string         `xml:"affects,attr"`
	Join            string         `xml:"join,attr"`
	Conditions      []Condition    `xml:"conditions>condition"`
	ConditionGroups []ConditionGroup `xml:"conditionGroups>conditionGroup"`
	Repeats         []Repeat       `xml:"repeats>repeat"`
}

// ModifierGroup groups related modifiers
type ModifierGroup struct {
	XMLName         xml.Name       `xml:"modifierGroup"`
	ID              string         `xml:"id,attr"`
	Type            string         `xml:"type,attr"`
	Modifiers       []Modifier     `xml:"modifiers>modifier"`
	ConditionGroups []ConditionGroup `xml:"conditionGroups>conditionGroup"`
	Comment         string         `xml:"comment"`
}

// Condition defines a condition for modifiers
type Condition struct {
	XMLName              xml.Name `xml:"condition"`
	ID                   string   `xml:"id,attr"`
	Type                 string   `xml:"type,attr"`
	Value                string   `xml:"value,attr"`
	Field                string   `xml:"field,attr"`
	Scope                string   `xml:"scope,attr"`
	ChildID              string   `xml:"childId,attr"`
	Shared               string   `xml:"shared,attr"`
	IncludeChildSelections string `xml:"includeChildSelections,attr"`
	IncludeChildForces    string  `xml:"includeChildForces,attr"`
	PercentValue         string   `xml:"percentValue,attr"`
}

// ConditionGroup groups conditions with AND/OR logic
type ConditionGroup struct {
	XMLName         xml.Name       `xml:"conditionGroup"`
	ID              string         `xml:"id,attr"`
	Type            string         `xml:"type,attr"`
	Conditions      []Condition    `xml:"conditions>condition"`
	ConditionGroups []ConditionGroup `xml:"conditionGroups>conditionGroup"`
}

// Repeat defines a repeat pattern for modifiers
type Repeat struct {
	XMLName xml.Name `xml:"repeat"`
	Value   string   `xml:"value,attr"`
	Repeats string   `xml:"repeats,attr"`
	Field   string   `xml:"field,attr"`
	RoundUp string   `xml:"roundUp,attr"`
	IncludeChildSelections string `xml:"includeChildSelections,attr"`
}

// Cost represents a point cost or other resource cost
type Cost struct {
	XMLName xml.Name `xml:"cost"`
	Name    string   `xml:"name,attr"`
	TypeID  string   `xml:"typeId,attr"`
	Value   string   `xml:"value,attr"`
}

// InfoLink references a game rule
type InfoLink struct {
	XMLName xml.Name `xml:"infoLink"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	TargetID string  `xml:"targetId,attr"`
	Type    string   `xml:"type,attr"`
	Hidden  string   `xml:"hidden,attr"`
}

