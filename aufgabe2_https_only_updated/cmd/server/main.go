package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"aufgabe2/plan"
	"aufgabe2/protocol"
	"aufgabe2/server"
	"aufgabe2/valuepair"
)

const (
	FLAG_XML_URL_DESCRIPTION   = "URL der aktuellen Stundenplan-XML-Datei"
	FLAG_PORT_DESCRIPTION      = "TCP-Port des Hauptservers"
	FLAG_VALUEPAIR_DESCRIPTION = "Adresse des eigenen Wertepaar-Servers"
	ERROR_START_MAIN_SERVER    = "Hauptserver konnte nicht gestartet werden: %v\n"
	ERROR_STOP_MAIN_SERVER     = "Hauptserver beendet: %v\n"
	INFO_MAIN_SERVER_RUNNING   = "Hauptserver läuft auf Port %s.\n"
)

func main() {
	xmlURL := flag.String("xml-url", protocol.SCHEDULE_XML_URL, FLAG_XML_URL_DESCRIPTION)
	port := flag.String("port", protocol.MAIN_SERVER_PORT, FLAG_PORT_DESCRIPTION)
	valuePairAddress := flag.String(
		"valuepair",
		protocol.VALUEPAIR_SERVER_HOST+":"+protocol.VALUEPAIR_SERVER_PORT,
		FLAG_VALUEPAIR_DESCRIPTION,
	)
	flag.Parse()

	provider := plan.NewProvider(*xmlURL)
	service := &server.TerminService{
		PlanProvider:    provider,
		ValuePairClient: valuepair.NewClient(*valuePairAddress),
	}

	listener, err := net.Listen(protocol.NETWORK_TYPE_TCP, ":"+*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, ERROR_START_MAIN_SERVER, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf(INFO_MAIN_SERVER_RUNNING, *port)
	if err := server.New(service).Serve(listener); err != nil {
		fmt.Fprintf(os.Stderr, ERROR_STOP_MAIN_SERVER, err)
	}
}
