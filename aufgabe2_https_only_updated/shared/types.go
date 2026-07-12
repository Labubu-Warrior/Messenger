// Package shared enthält Datentypen, die von Client und Server gemeinsam benutzt werden.
package shared

// ClientRequest ist eine Anfrage des GUI-Clients an den Hauptserver.
type ClientRequest struct {
	Type  string   `json:"type"`
	Codes []string `json:"codes,omitempty"`
}

// ClientResponse ist eine Antwort des Hauptservers an den GUI-Client.
type ClientResponse struct {
	Type        string       `json:"type"`
	Codes       []string     `json:"codes,omitempty"`
	Suggestions []Suggestion `json:"suggestions,omitempty"`
	Message     string       `json:"message,omitempty"`
}

// Suggestion beschreibt einen lesbaren Terminvorschlag.
type Suggestion struct {
	DayID  int    `json:"day_id"`
	TimeID int    `json:"time_id"`
	Day    string `json:"day"`
	From   string `json:"from"`
	To     string `json:"to"`
}

// ValuePair enthält nullbasierte Tages- und Zeitindizes.
type ValuePair struct {
	DayIndex  int
	TimeIndex int
}

// Day enthält einen im XML definierten Wochentag.
type Day struct {
	ID        int
	ShortName string
	LongName  string
}

// TimeBlock enthält einen im XML definierten Zeitblock.
type TimeBlock struct {
	ID    int
	Label string
	From  string
	To    string
}

// Employee enthält die für die Aufgabe relevanten Mitarbeiterdaten.
type Employee struct {
	ID        int
	Code      string
	FirstName string
	LastName  string
}
