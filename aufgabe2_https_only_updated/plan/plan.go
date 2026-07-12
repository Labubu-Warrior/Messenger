package plan

import (
	"encoding/xml"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"aufgabe2/shared"
)

const (
	MIN_EVENT_DURATION_BLOCKS = 1
	XML_INDEX_OFFSET          = 1
	TIME_LABEL_PARTS          = 2
)

// Plan enthält die für die Anwendung aufbereiteten Stundenplandaten.
type Plan struct {
	Days           []shared.Day
	TimeBlocks     []shared.TimeBlock
	Employees      []shared.Employee
	EmployeeByID   map[int]shared.Employee
	EmployeeByCode map[string]shared.Employee
	Events         []Event
}

// Event ist eine normalisierte Veranstaltung aus dem XML.
type Event struct {
	ID          int
	Name        string
	Additional  string
	Hours       int
	Terms       []Term
	EmployeeIDs []int
}

// Term enthält Tag und Startzeitblock einer Veranstaltung.
type Term struct {
	DayID  int
	TimeID int
}

// Parse wandelt ein vollständiges splan-XML in ein Plan-Objekt um.
func Parse(data []byte) (*Plan, error) {
	var raw xmlPlan
	if err := xml.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("ungültiges Stundenplan-XML: %w", err)
	}
	if len(raw.Days) == 0 || len(raw.Times) == 0 || len(raw.EmployeeGroups) == 0 || len(raw.Events) == 0 {
		return nil, fmt.Errorf("Stundenplan-XML enthält nicht alle benötigten Bereiche")
	}

	result := &Plan{
		EmployeeByID:   make(map[int]shared.Employee),
		EmployeeByCode: make(map[string]shared.Employee),
	}

	for _, rawDay := range raw.Days {
		result.Days = append(result.Days, shared.Day{
			ID:        rawDay.ID,
			ShortName: strings.TrimSpace(rawDay.ShortName),
			LongName:  strings.TrimSpace(rawDay.LongName),
		})
	}

	for _, rawTime := range raw.Times {
		from, to, err := splitTimeLabel(rawTime.Label)
		if err != nil {
			return nil, fmt.Errorf("Zeitblock %d: %w", rawTime.ID, err)
		}
		result.TimeBlocks = append(result.TimeBlocks, shared.TimeBlock{
			ID:    rawTime.ID,
			Label: strings.TrimSpace(rawTime.Label),
			From:  from,
			To:    to,
		})
	}

	for _, group := range raw.EmployeeGroups {
		for _, person := range group.Persons {
			code := strings.TrimSpace(person.Code)
			if code == "" {
				continue
			}
			employee := shared.Employee{
				ID:        person.ID,
				Code:      code,
				FirstName: strings.TrimSpace(person.FirstName),
				LastName:  strings.TrimSpace(person.LastName),
			}
			result.Employees = append(result.Employees, employee)
			result.EmployeeByID[employee.ID] = employee
			result.EmployeeByCode[strings.ToLower(employee.Code)] = employee
		}
	}

	eventsWithTeachers := 0
	eventsWithTerms := 0

	for _, rawEvent := range raw.Events {
		event := Event{
			ID:          rawEvent.ID,
			Name:        strings.TrimSpace(rawEvent.Name),
			Additional:  strings.TrimSpace(rawEvent.Additional),
			Hours:       rawEvent.Hours,
			EmployeeIDs: uniqueInts(rawEvent.Teachers),
		}
		for _, rawTerm := range rawEvent.Terms {
			event.Terms = append(event.Terms, Term{
				DayID:  rawTerm.DayID,
				TimeID: rawTerm.TimeID,
			})
		}
		if len(event.EmployeeIDs) > 0 {
			eventsWithTeachers++
		}
		if len(event.Terms) > 0 {
			eventsWithTerms++
		}
		result.Events = append(result.Events, event)
	}

	if len(result.Employees) == 0 {
		return nil, fmt.Errorf("Stundenplan-XML enthält keine Mitarbeiterkürzel")
	}
	if eventsWithTeachers == 0 {
		return nil, fmt.Errorf("Stundenplan-XML enthält keine Dozentenzuordnungen")
	}
	if eventsWithTerms == 0 {
		return nil, fmt.Errorf("Stundenplan-XML enthält keine Veranstaltungstermine")
	}

	sort.Slice(result.Days, func(i, j int) bool {
		return result.Days[i].ID < result.Days[j].ID
	})
	sort.Slice(result.TimeBlocks, func(i, j int) bool {
		return result.TimeBlocks[i].ID < result.TimeBlocks[j].ID
	})
	sort.Slice(result.Employees, func(i, j int) bool {
		return strings.ToLower(result.Employees[i].Code) < strings.ToLower(result.Employees[j].Code)
	})

	return result, nil
}

func splitTimeLabel(label string) (string, string, error) {
	parts := strings.Split(label, "-")
	if len(parts) != TIME_LABEL_PARTS {
		return "", "", fmt.Errorf("ungültige Zeitangabe %q", label)
	}
	return strings.TrimSpace(parts[0]), strings.TrimSpace(parts[1]), nil
}

func uniqueInts(values []int) []int {
	seen := make(map[int]struct{}, len(values))
	result := make([]int, 0, len(values))
	for _, value := range values {
		if _, exists := seen[value]; exists {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	return result
}

// Codes liefert alle Mitarbeiterkürzel sortiert zurück.
func (p *Plan) Codes() []string {
	codes := make([]string, 0, len(p.Employees))
	for _, employee := range p.Employees {
		codes = append(codes, employee.Code)
	}
	return codes
}

// IDsForCodes prüft Kürzel und liefert die zugehörigen Mitarbeiter-IDs.
func (p *Plan) IDsForCodes(codes []string) ([]int, error) {
	if len(codes) == 0 {
		return nil, fmt.Errorf("mindestens ein Mitarbeiterkürzel muss ausgewählt werden")
	}

	seen := make(map[int]struct{})
	ids := make([]int, 0, len(codes))
	for _, code := range codes {
		cleanCode := strings.ToLower(strings.TrimSpace(code))
		if cleanCode == "" {
			return nil, fmt.Errorf("leeres Mitarbeiterkürzel ist nicht erlaubt")
		}
		employee, exists := p.EmployeeByCode[cleanCode]
		if !exists {
			return nil, fmt.Errorf("unbekanntes Mitarbeiterkürzel %q", code)
		}
		if _, duplicate := seen[employee.ID]; duplicate {
			continue
		}
		seen[employee.ID] = struct{}{}
		ids = append(ids, employee.ID)
	}
	return ids, nil
}

// ValidateIDs prüft Mitarbeiter-IDs aus dem Wertepaar-Protokoll.
func (p *Plan) ValidateIDs(ids []int) error {
	if len(ids) == 0 {
		return fmt.Errorf("mindestens eine Mitarbeiter-ID ist erforderlich")
	}
	for _, id := range ids {
		if _, exists := p.EmployeeByID[id]; !exists {
			return fmt.Errorf("unbekannte Mitarbeiter-ID %d", id)
		}
	}
	return nil
}

// FreePairs berechnet alle gemeinsamen freien Zeitblöcke der Mitarbeiter.
// Q1- und Q2-Veranstaltungen werden nicht gefiltert und damit wie Langläufer behandelt.
func (p *Plan) FreePairs(ids []int) ([]shared.ValuePair, error) {
	if err := p.ValidateIDs(ids); err != nil {
		return nil, err
	}
	if len(p.TimeBlocks) == 0 {
		return nil, fmt.Errorf("keine Zeitblöcke im Stundenplan vorhanden")
	}

	selectedEmployees := make(map[int]struct{}, len(ids))
	for _, id := range ids {
		selectedEmployees[id] = struct{}{}
	}

	busySlots := make(map[[TIME_LABEL_PARTS]int]struct{})
	maxTimeID := p.TimeBlocks[len(p.TimeBlocks)-1].ID
	for _, event := range p.Events {
		if !intersects(event.EmployeeIDs, selectedEmployees) {
			continue
		}
		duration := event.Hours
		if duration < MIN_EVENT_DURATION_BLOCKS {
			duration = MIN_EVENT_DURATION_BLOCKS
		}
		for _, term := range event.Terms {
			for offset := 0; offset < duration; offset++ {
				timeID := term.TimeID + offset
				if timeID > maxTimeID {
					break
				}
				busySlots[[TIME_LABEL_PARTS]int{term.DayID, timeID}] = struct{}{}
			}
		}
	}

	pairs := make([]shared.ValuePair, 0)
	for _, day := range p.Days {
		for _, block := range p.TimeBlocks {
			if _, occupied := busySlots[[TIME_LABEL_PARTS]int{day.ID, block.ID}]; occupied {
				continue
			}
			pairs = append(pairs, shared.ValuePair{
				DayIndex:  day.ID - XML_INDEX_OFFSET,
				TimeIndex: block.ID - XML_INDEX_OFFSET,
			})
		}
	}
	return pairs, nil
}

func intersects(eventIDs []int, selected map[int]struct{}) bool {
	for _, id := range eventIDs {
		if _, exists := selected[id]; exists {
			return true
		}
	}
	return false
}

// SuggestionForPair wandelt ein Index-Wertepaar in einen lesbaren Termin um.
func (p *Plan) SuggestionForPair(pair shared.ValuePair) (shared.Suggestion, error) {
	dayID := pair.DayIndex + XML_INDEX_OFFSET
	timeID := pair.TimeIndex + XML_INDEX_OFFSET

	var selectedDay *shared.Day
	var selectedTime *shared.TimeBlock
	for i := range p.Days {
		if p.Days[i].ID == dayID {
			selectedDay = &p.Days[i]
			break
		}
	}
	for i := range p.TimeBlocks {
		if p.TimeBlocks[i].ID == timeID {
			selectedTime = &p.TimeBlocks[i]
			break
		}
	}
	if selectedDay == nil || selectedTime == nil {
		return shared.Suggestion{}, fmt.Errorf(
			"ungültiges Wertepaar (%s,%s)",
			strconv.Itoa(pair.DayIndex),
			strconv.Itoa(pair.TimeIndex),
		)
	}

	return shared.Suggestion{
		DayID:  pair.DayIndex,
		TimeID: pair.TimeIndex,
		Day:    selectedDay.LongName,
		From:   selectedTime.From,
		To:     selectedTime.To,
	}, nil
}
