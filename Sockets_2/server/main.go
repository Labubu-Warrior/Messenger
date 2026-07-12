package main

import (
	"bufio"
	"log"
	"net"
	"strings"
)

const ServerPort = ":8080"

func main() {
	LoadSchedule()

	listener, err := net.Listen("tcp", ServerPort)
	if err != nil {
		log.Fatalf("Server konnte nicht starten: %v", err)
	}
	log.Printf("Server lauscht auf Port %s...", ServerPort)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Verbindungsfehler:", err)
			continue
		}
		// Goroutine erfüllt Anforderung: "verwaltet viele unabhängige Clients"
		go handleClientConnection(conn)
	}
}

func handleClientConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// 1. Anforderung: Liste der Kürzel an Client senden
	var kuerzelList []string
	for _, group := range globalSchedule.StaffGroups {
		for _, p := range group.Persons {
			if p.Details.Kuerzel != "" {
				kuerzelList = append(kuerzelList, p.Details.Kuerzel)
			}
		}
	}
	conn.Write([]byte("INIT:" + strings.Join(kuerzelList, ",") + "\n"))

	// 2. Anfragen verarbeiten
	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			return // Client disconnected
		}
		message = strings.TrimSpace(message)

		if strings.HasPrefix(message, "REQ:") {
			reqStr := strings.TrimPrefix(message, "REQ:")
			kuerzelArray := strings.Split(reqStr, ",")

			// Logik komplett ausgelagert an scheduler.go
			lines, err := GetTerminvorschlaege(kuerzelArray)

			if err != nil {
				// Anforderung: Fehlermeldung bei fehlerhaften Anfragen
				conn.Write([]byte("ERROR:" + err.Error() + "\n"))
			} else if len(lines) == 0 {
				conn.Write([]byte("RES:Keine gemeinsamen Termine gefunden.\n"))
			} else {
				// Zeilen mit "|" trennen für die TCP-Übertragung
				conn.Write([]byte("RES:" + strings.Join(lines, "|") + "\n"))
			}
		}
	}
}
