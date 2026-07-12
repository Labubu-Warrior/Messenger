package server

import (
	"fmt"

	"aufgabe2/plan"
	"aufgabe2/shared"
	"aufgabe2/valuepair"
)

// TerminService enthält die komplette fachliche Logik des Hauptservers.
type TerminService struct {
	PlanProvider    *plan.Provider
	ValuePairClient *valuepair.Client
}

// Codes lädt das aktuelle XML und liefert alle Mitarbeiterkürzel.
func (s *TerminService) Codes() ([]string, error) {
	currentPlan, err := s.PlanProvider.Load()
	if err != nil {
		return nil, err
	}
	return currentPlan.Codes(), nil
}

// Suggestions lädt aktuelle XML-Daten, validiert Kürzel und übersetzt Wertepaare.
func (s *TerminService) Suggestions(codes []string) ([]shared.Suggestion, error) {
	currentPlan, err := s.PlanProvider.Load()
	if err != nil {
		return nil, err
	}

	ids, err := currentPlan.IDsForCodes(codes)
	if err != nil {
		return nil, err
	}
	pairs, err := s.ValuePairClient.RequestPairs(ids)
	if err != nil {
		return nil, err
	}

	suggestions := make([]shared.Suggestion, 0, len(pairs))
	for _, pair := range pairs {
		suggestion, err := currentPlan.SuggestionForPair(pair)
		if err != nil {
			return nil, fmt.Errorf("Wertepaar-Server lieferte ungültige Daten: %w", err)
		}
		suggestions = append(suggestions, suggestion)
	}
	return suggestions, nil
}
