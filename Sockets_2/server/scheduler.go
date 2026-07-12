package main

import (
	"fmt"
	"strconv"
	"strings"
)

// GetTerminvorschlaege ist der zentrale Einstiegspunkt für eine Terminanfrage.
// Parameter: kuerzelListe - getrimmte Liste von Mitarbeiterkürzeln.
// Rückgabe: Formatierte Zeilen ("Tag: Montag  Zeit: 15:30 - 16:45").
func GetTerminvorschlaege(kuerzelListe []string) ([]string, error) {
	// 1. Kürzel in IDs auflösen
	var personIDs []int
	for _, kuerzel := range kuerzelListe {
		id, ok := lookupPersonID(kuerzel)
		if !ok {
			return nil, fmt.Errorf("unbekanntes Kürzel: %q", kuerzel)
		}
		personIDs = append(personIDs, id)
	}

	// 2. Wertepaar-Request bauen (z.B. "1,25#")
	var reqBuilder strings.Builder
	for i, id := range personIDs {
		reqBuilder.WriteString(strconv.Itoa(id))
		if i < len(personIDs)-1 {
			reqBuilder.WriteString(",")
		}
	}
	reqBuilder.WriteString("#")

	// 3. Eigene Wertepaar-Implementierung aufrufen!
	rawResponse := ComputeWertepaare(reqBuilder.String())

	// 4. Antwort parsen und in menschenlesbaren Text umwandeln
	return parseAndFormatResponse(rawResponse), nil
}

// lookupPersonID durchsucht alle Mitarbeiter-Gruppen nach einem Kürzel.
func lookupPersonID(kuerzel string) (int, bool) {
	for _, group := range globalSchedule.StaffGroups {
		for _, person := range group.Persons {
			if person.Details.Kuerzel == kuerzel {
				return person.ID, true
			}
		}
	}
	return 0, false
}

// parseAndFormatResponse zerlegt "(0,5)(0,6)#" in lesbaren Text.
func parseAndFormatResponse(response string) []string {
	var lines []string
	cleanResp := strings.TrimSuffix(response, "#")
	if cleanResp == "" {
		return lines
	}

	// Trennt "(0,5)(0,6)" in "0,5", "0,6" auf
	cleanResp = strings.TrimPrefix(cleanResp, "(")
	cleanResp = strings.TrimSuffix(cleanResp, ")")
	pairs := strings.Split(cleanResp, ")(")

	for _, pair := range pairs {
		parts := strings.Split(pair, ",")
		if len(parts) == 2 {
			dayIndex, _ := strconv.Atoi(parts[0])
			timeIndex, _ := strconv.Atoi(parts[1])

			dayName, start, end := resolveIndicesToText(dayIndex, timeIndex)

			// Exaktes Format laut Aufgabenstellung-Beispiel (Doppeltes Leerzeichen vor Zeit)
			lines = append(lines, fmt.Sprintf("Tag: %s  Zeit: %s - %s", dayName, start, end))
		}
	}
	return lines
}

// resolveIndicesToText übersetzt die 0-basierten Indizes zurück in Namen.
func resolveIndicesToText(dayIndex, timeIndex int) (string, string, string) {
	dayName, start, end := "Unbekannt", "?", "?"

	// XML-IDs sind Index + 1
	targetDayID := dayIndex + 1
	targetTimeID := timeIndex + 1

	for _, d := range globalSchedule.Days {
		if d.ID == targetDayID {
			// Workaround für exakte Übereinstimmung mit Aufgabenblatt
			if d.Name == "Sonnabend" {
				dayName = "Samstag"
			} else {
				dayName = d.Name
			}
			break
		}
	}
	for _, t := range globalSchedule.Times {
		if t.ID == targetTimeID {
			start = t.Start
			end = t.End
			break
		}
	}
	return dayName, start, end
}
