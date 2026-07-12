package main

import "encoding/xml"

// ScheduleData repräsentiert die Wurzel der Stundenplan-XML.
type ScheduleData struct {
	XMLName         xml.Name        `xml:"stundenplan"`
	StaffGroups     []StaffGroup    `xml:"mitarbeiter>typ"`
	Days            []Day           `xml:"tag"`
	Times           []Time          `xml:"zeit"`
	Veranstaltungen []Veranstaltung `xml:"veranstaltung"`
}

// StaffGroup bündelt Mitarbeiter (z.B. Dozenten, Lehrbeauftragte).
type StaffGroup struct {
	ID      int      `xml:"id"`
	Persons []Person `xml:"person"`
}

// Person speichert die Mitarbeiterdaten (aus dem Pascal-Snippet adaptiert).
type Person struct {
	ID      int           `xml:"id"`
	Details PersonDetails `xml:"bezeichnung"`
}

// PersonDetails enthält das Kürzel (z.B. ahr, wol).
type PersonDetails struct {
	Kuerzel string `xml:"kuerzel"`
}

// Day speichert die Wochentage (XML-ID ist meist 1..7).
type Day struct {
	ID   int    `xml:"id"`
	Name string `xml:"name"`
}

// Time speichert die Zeitblöcke (XML-ID ist meist 1..N).
type Time struct {
	ID    int    `xml:"id"`
	Start string `xml:"start"`
	End   string `xml:"ende"`
}

// Veranstaltung verknüpft eine Person mit einem Zeitblock an einem Tag.
type Veranstaltung struct {
	PersonID int    `xml:"person_id"`
	DayID    int    `xml:"tag_id"`
	TimeID   int    `xml:"zeit_id"`
	Quartal  string `xml:"quartal"`
}
