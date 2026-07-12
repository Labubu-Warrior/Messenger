package main

import (
	"encoding/xml"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

// =====================================================================
// UMSCHALTER: true = Zuhause (lokale XML), false = Uni (Live-Server)
// =====================================================================
const isLocalTesting = true

const ScheduleURL = "https://intern.fh-wedel.de/~tho/api/splan.php"

var globalSchedule ScheduleData

// LoadSchedule lädt und parst die XML-Datei.
func LoadSchedule() {
	var body []byte
	var err error

	if isLocalTesting {
		log.Println("Lade lokale Dummy-XML-Daten (Offline-Modus)...")
		file, errOpen := os.Open("splan.xml")
		if errOpen != nil {
			log.Fatalf("Fehler: Lokale splan.xml nicht gefunden: %v", errOpen)
		}
		defer file.Close()
		body, err = io.ReadAll(file)
	} else {
		log.Println("Lade XML-Stundenplan live vom FH-Server...")
		resp, errHttp := http.Get(ScheduleURL)
		if errHttp != nil {
			log.Fatalf("Netzwerkfehler: %v", errHttp)
		}
		defer resp.Body.Close()
		body, err = io.ReadAll(resp.Body)
	}

	if err != nil {
		log.Fatalf("Fehler beim Lesen: %v", err)
	}

	xmlString := string(body)
	if !strings.Contains(xmlString, "<stundenplan>") {
		xmlString = "<stundenplan>" + xmlString + "</stundenplan>"
	}

	err = xml.Unmarshal([]byte(xmlString), &globalSchedule)
	if err != nil {
		log.Fatalf("XML Parsing-Fehler: %v", err)
	}
	log.Println("Stundenplan erfolgreich geladen und geparst.")
}
