package plan

import "encoding/xml"

type xmlPlan struct {
	XMLName        xml.Name          `xml:"splan"`
	Days           []xmlDay          `xml:"tage>tag"`
	Times          []xmlTime         `xml:"zeiten>zeit"`
	EmployeeGroups []xmlEmployeeType `xml:"mitarbeiter>typ"`
	Events         []xmlEvent        `xml:"veranstaltungen>veranstaltung"`
}

type xmlDay struct {
	ID        int    `xml:"id"`
	ShortName string `xml:"bezeichnung>kurz"`
	LongName  string `xml:"bezeichnung>lang"`
}

type xmlTime struct {
	ID    int    `xml:"id"`
	Label string `xml:"bezeichnung"`
}

type xmlEmployeeType struct {
	ID      int           `xml:"id"`
	Label   string        `xml:"bezeichnung"`
	Persons []xmlEmployee `xml:"person"`
}

type xmlEmployee struct {
	ID        int    `xml:"id"`
	Code      string `xml:"bezeichnung>kuerzel"`
	FirstName string `xml:"bezeichnung>vorname"`
	LastName  string `xml:"bezeichnung>nachname"`
}

type xmlEvent struct {
	ID         int       `xml:"id"`
	Name       string    `xml:"bezeichnung"`
	Additional string    `xml:"bez_zusatz"`
	Hours      int       `xml:"stunden"`
	Terms      []xmlTerm `xml:"termine>termin"`
	Teachers   []int     `xml:"mitarbeiter>dozent>id"`
}

type xmlTerm struct {
	DayID  int `xml:"tag"`
	TimeID int `xml:"zeit"`
}
