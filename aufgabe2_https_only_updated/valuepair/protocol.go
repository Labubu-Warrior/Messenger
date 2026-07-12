// Package valuepair implementiert das textbasierte Wertepaar-Protokoll.
package valuepair

import (
	"fmt"
	"strconv"
	"strings"

	"aufgabe2/protocol"
	"aufgabe2/shared"
)

const (
	MIN_VALID_EMPLOYEE_ID = 1
	VALUEPAIR_PARTS       = 2
)

// ParseIDRequest liest eine Anfrage im Format "1,25#".
func ParseIDRequest(line string) ([]int, error) {
	if !strings.HasSuffix(line, string(protocol.VALUEPAIR_END_MARKER)) {
		return nil, fmt.Errorf("Anfrage muss mit %c enden", protocol.VALUEPAIR_END_MARKER)
	}
	body := strings.TrimSpace(strings.TrimSuffix(line, string(protocol.VALUEPAIR_END_MARKER)))
	if body == "" {
		return nil, fmt.Errorf("keine Mitarbeiter-IDs angegeben")
	}

	parts := strings.Split(body, string(protocol.VALUEPAIR_SEPARATOR))
	ids := make([]int, 0, len(parts))
	seen := make(map[int]struct{})
	for _, part := range parts {
		cleanPart := strings.TrimSpace(part)
		if cleanPart == "" {
			return nil, fmt.Errorf("leere Mitarbeiter-ID")
		}
		id, err := strconv.Atoi(cleanPart)
		if err != nil || id < MIN_VALID_EMPLOYEE_ID {
			return nil, fmt.Errorf("ungültige Mitarbeiter-ID %q", cleanPart)
		}
		if _, duplicate := seen[id]; duplicate {
			continue
		}
		seen[id] = struct{}{}
		ids = append(ids, id)
	}
	return ids, nil
}

// FormatIDRequest erzeugt eine Anfrage im Format "1,25#".
func FormatIDRequest(ids []int) (string, error) {
	if len(ids) == 0 {
		return "", fmt.Errorf("keine Mitarbeiter-IDs angegeben")
	}
	parts := make([]string, 0, len(ids))
	for _, id := range ids {
		if id < MIN_VALID_EMPLOYEE_ID {
			return "", fmt.Errorf("ungültige Mitarbeiter-ID %d", id)
		}
		parts = append(parts, strconv.Itoa(id))
	}
	return strings.Join(parts, string(protocol.VALUEPAIR_SEPARATOR)) + string(protocol.VALUEPAIR_END_MARKER), nil
}

// FormatPairs erzeugt die Antwort "(0,5)(0,6)#".
func FormatPairs(pairs []shared.ValuePair) string {
	var builder strings.Builder
	for _, pair := range pairs {
		fmt.Fprintf(&builder, "(%d,%d)", pair.DayIndex, pair.TimeIndex)
	}
	builder.WriteByte(protocol.VALUEPAIR_END_MARKER)
	return builder.String()
}

// ParsePairs liest eine Wertepaar-Antwort.
func ParsePairs(line string) ([]shared.ValuePair, error) {
	if strings.HasPrefix(line, protocol.VALUEPAIR_ERROR_PREFIX) {
		message := strings.TrimSuffix(strings.TrimPrefix(line, protocol.VALUEPAIR_ERROR_PREFIX), string(protocol.VALUEPAIR_END_MARKER))
		return nil, fmt.Errorf("Wertepaar-Server: %s", strings.TrimSpace(message))
	}
	if !strings.HasSuffix(line, string(protocol.VALUEPAIR_END_MARKER)) {
		return nil, fmt.Errorf("Antwort muss mit %c enden", protocol.VALUEPAIR_END_MARKER)
	}

	body := strings.TrimSpace(strings.TrimSuffix(line, string(protocol.VALUEPAIR_END_MARKER)))
	if body == "" {
		return []shared.ValuePair{}, nil
	}

	pairs := make([]shared.ValuePair, 0)
	for len(body) > 0 {
		if body[0] != '(' {
			return nil, fmt.Errorf("ungültige Wertepaar-Antwort")
		}
		endIndex := strings.IndexByte(body, ')')
		if endIndex < 0 {
			return nil, fmt.Errorf("fehlende schließende Klammer")
		}
		parts := strings.Split(body[1:endIndex], string(protocol.VALUEPAIR_SEPARATOR))
		if len(parts) != VALUEPAIR_PARTS {
			return nil, fmt.Errorf("ungültiges Wertepaar %q", body[:endIndex+1])
		}
		dayIndex, dayErr := strconv.Atoi(strings.TrimSpace(parts[0]))
		timeIndex, timeErr := strconv.Atoi(strings.TrimSpace(parts[1]))
		if dayErr != nil || timeErr != nil || dayIndex < 0 || timeIndex < 0 {
			return nil, fmt.Errorf("ungültiges Wertepaar %q", body[:endIndex+1])
		}
		pairs = append(pairs, shared.ValuePair{DayIndex: dayIndex, TimeIndex: timeIndex})
		body = strings.TrimSpace(body[endIndex+1:])
	}
	return pairs, nil
}
