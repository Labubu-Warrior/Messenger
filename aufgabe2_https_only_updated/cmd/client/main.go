package main

import (
	"flag"
	"fmt"
	"os"

	"aufgabe2/client"
	"aufgabe2/protocol"
)

const (
	FLAG_SERVER_DESCRIPTION = "Adresse des Hauptservers"
	ERROR_START_CLIENT      = "Client konnte nicht gestartet werden: %v\n"
)

func main() {
	address := flag.String(
		"server",
		protocol.MAIN_SERVER_HOST+":"+protocol.MAIN_SERVER_PORT,
		FLAG_SERVER_DESCRIPTION,
	)
	flag.Parse()

	networkClient, err := client.Connect(*address)
	if err != nil {
		fmt.Fprintf(os.Stderr, ERROR_START_CLIENT, err)
		os.Exit(1)
	}
	client.RunGUI(networkClient)
}
