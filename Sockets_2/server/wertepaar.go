package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ComputeWertepaare simuliert den externen FH-Server (Port 16667).
// Es nimmt einen String wie "1,25#" entgegen und liefert gemeinsame
// freie Termine als Index-Paare z.B. "(0,5)(0,6)#" zurück.
// Laut Aufgabe: Montag..Sonntag = 0..6, Erster..Letzter Block = 0..N
func ComputeWertepaare(request string) string {
	// 1. Request parsen (das '#' am Ende entfernen)
	req := strings.TrimSuffix(request, "#")
	if req == "" {
		return "#"
	}

	idStrings := strings.Split(req, ",")
	var requestedIDs []int
	for _, idStr := range idStrings {
		id, _ := strconv.Atoi(idStr)
		requestedIDs = append(requestedIDs, id)
	}

	// 2. Belegte Slots ermitteln
	// Map-Struktur: [DayID][TimeID]bool (Hier nutzen wir noch die echten XML-IDs)
	blockedSlots := make(map[int]map[int]bool)

	for _, v := range globalSchedule.Veranstaltungen {
		for _, reqID := range requestedIDs {
			if v.PersonID == reqID {
				if blockedSlots[v.DayID] == nil {
					blockedSlots[v.DayID] = make(map[int]bool)
				}
				// Q1/Q2 Kurzläufer werden vereinfachend als Langläufer gewertet
				blockedSlots[v.DayID][v.TimeID] = true
			}
		}
	}

	// 3. Freie Slots berechnen und Ergebnis-String bauen
	var resultBuilder strings.Builder

	// Umwandlung der XML-IDs (meist 1-basiert) in die geforderten 0-basierten Indizes
	for _, day := range globalSchedule.Days {
		// Index = 0..6 (Wir nehmen an, XML-ID 1 ist Montag, also ID-1 = 0)
		dayIndex := day.ID - 1

		for _, time := range globalSchedule.Times {
			timeIndex := time.ID - 1

			if !blockedSlots[day.ID][time.ID] {
				// Slot ist frei -> an den String anhängen, z.B. "(0,5)"
				resultBuilder.WriteString(fmt.Sprintf("(%d,%d)", dayIndex, timeIndex))
			}
		}
	}

	// Protokollkonform mit '#' abschließen
	resultBuilder.WriteString("#")
	return resultBuilder.String()
}
