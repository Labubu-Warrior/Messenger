// Package plan lädt und interpretiert die aktuelle XML-Datenbasis des Stundenplans.
package plan

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	"aufgabe2/protocol"
)

// Provider lädt bei jedem Aufruf die aktuelle XML-Datei über HTTPS.
type Provider struct {
	URL        string
	HTTPClient *http.Client
}

// NewProvider erstellt eine HTTPS-Datenquelle für den Stundenplan.
func NewProvider(url string) *Provider {
	return &Provider{
		URL: strings.TrimSpace(url),
		HTTPClient: &http.Client{
			Timeout: protocol.HTTP_REQUEST_TIMEOUT,
		},
	}
}

// Load lädt die aktuelle XML-Datei und wandelt sie in einen Plan um.
func (p *Provider) Load() (*Plan, error) {
	if p == nil || strings.TrimSpace(p.URL) == "" {
		return nil, fmt.Errorf("keine Stundenplan-URL angegeben")
	}
	if p.HTTPClient == nil {
		return nil, fmt.Errorf("kein HTTP-Client konfiguriert")
	}

	response, err := p.HTTPClient.Get(p.URL)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", protocol.MESSAGE_XML_UNAVAILABLE, err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Stundenplan-Server antwortet mit %s", response.Status)
	}

	limitedReader := io.LimitReader(response.Body, protocol.MAX_XML_SIZE_BYTES+1)
	data, readErr := io.ReadAll(limitedReader)
	if len(data) > protocol.MAX_XML_SIZE_BYTES {
		return nil, fmt.Errorf("Stundenplan-XML überschreitet die erlaubte Größe")
	}
	// Der FH-Server beendet die Übertragung gelegentlich mit einem ungenauen
	// Content-Length-Wert. Wenn das XML vollständig ist, kann Parse es trotzdem
	// sicher validieren. Andere Lesefehler werden abgelehnt.
	if readErr != nil && !errors.Is(readErr, io.ErrUnexpectedEOF) {
		return nil, fmt.Errorf("Stundenplan-XML konnte nicht gelesen werden: %w", readErr)
	}

	parsedPlan, err := Parse(data)
	if err != nil {
		return nil, err
	}
	return parsedPlan, nil
}
