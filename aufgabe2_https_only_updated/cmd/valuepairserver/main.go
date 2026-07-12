package main

import (
	"flag"
	"fmt"
	"net"
	"os"

	"aufgabe2/plan"
	"aufgabe2/protocol"
	"aufgabe2/valuepairserver"
)

const (
	FLAG_XML_URL_DESCRIPTION      = "URL der aktuellen Stundenplan-XML-Datei"
	FLAG_PORT_DESCRIPTION         = "TCP-Port des Wertepaar-Servers"
	ERROR_START_VALUEPAIR_SERVER  = "Wertepaar-Server konnte nicht gestartet werden: %v\n"
	ERROR_STOP_VALUEPAIR_SERVER   = "Wertepaar-Server beendet: %v\n"
	INFO_VALUEPAIR_SERVER_RUNNING = "Eigener Wertepaar-Server läuft auf Port %s.\n"
)

func main() {
	xmlURL := flag.String("xml-url", protocol.SCHEDULE_XML_URL, FLAG_XML_URL_DESCRIPTION)
	port := flag.String("port", protocol.VALUEPAIR_SERVER_PORT, FLAG_PORT_DESCRIPTION)
	flag.Parse()

	provider := plan.NewProvider(*xmlURL)
	listener, err := net.Listen(protocol.NETWORK_TYPE_TCP, ":"+*port)
	if err != nil {
		fmt.Fprintf(os.Stderr, ERROR_START_VALUEPAIR_SERVER, err)
		os.Exit(1)
	}
	defer listener.Close()

	fmt.Printf(INFO_VALUEPAIR_SERVER_RUNNING, *port)
	if err := valuepairserver.New(provider).Serve(listener); err != nil {
		fmt.Fprintf(os.Stderr, ERROR_STOP_VALUEPAIR_SERVER, err)
	}
}
