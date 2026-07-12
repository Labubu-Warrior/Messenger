package main

import (
	"bufio"
	"log"
	"net"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

const ServerAddress = "localhost:8080"

func main() {
	// 1. Verbindung aufbauen
	conn, err := net.Dial("tcp", ServerAddress)
	if err != nil {
		log.Fatalf("Keine Verbindung zum Server: %v", err)
	}
	defer conn.Close()
	reader := bufio.NewReader(conn)

	// 2. Kürzel empfangen
	initMsg, _ := reader.ReadString('\n')
	initMsg = strings.TrimSpace(initMsg)
	var staffList []string
	if strings.HasPrefix(initMsg, "INIT:") {
		staffList = strings.Split(strings.TrimPrefix(initMsg, "INIT:"), ",")
	}

	// 3. GUI aufbauen (Erfüllt: "arbeitet mit grafischer Oberfläche")
	guiApp := app.New()
	window := guiApp.NewWindow("Termin-Finder FH Wedel")
	window.Resize(fyne.NewSize(600, 500))

	var selectedStaff []string

	// Erfüllt: "komfortable und fehlerfreie Auswahl" (CheckGroup)
	staffCheckGroup := widget.NewCheckGroup(staffList, func(selected []string) {
		selectedStaff = selected
	})

	// Erfüllt: "präsentiert Ergebnisse übersichtlich"
	resultDisplay := widget.NewLabel("")
	scrollArea := container.NewVScroll(resultDisplay) // Macht den Text scrollbar

	requestButton := widget.NewButton("Terminvorschläge anfragen", func() {
		if len(selectedStaff) == 0 {
			resultDisplay.SetText("Bitte wählen Sie mindestens einen Mitarbeiter aus.")
			return
		}

		// Erfüllt: "lediglich für Eingabe und Senden zuständig"
		requestString := "REQ:" + strings.Join(selectedStaff, ",") + "\n"
		conn.Write([]byte(requestString))

		// Antwort verarbeiten
		responseMsg, _ := reader.ReadString('\n')
		responseMsg = strings.TrimSpace(responseMsg)

		if strings.HasPrefix(responseMsg, "ERROR:") {
			resultDisplay.SetText("Fehler vom Server:\n" + strings.TrimPrefix(responseMsg, "ERROR:"))
		} else if strings.HasPrefix(responseMsg, "RES:") {
			content := strings.TrimPrefix(responseMsg, "RES:")
			// Protokoll-Trenner "|" in echte Zeilenumbrüche verwandeln
			formattedOutput := strings.ReplaceAll(content, "|", "\n")
			resultDisplay.SetText("Terminvorschläge:\n\n" + formattedOutput)
		}
	})

	topPanel := container.NewVBox(
		widget.NewLabelWithStyle("Gewünschte Teilnehmer:", fyne.TextAlignLeading, fyne.TextStyle{Bold: true}),
		staffCheckGroup,
		requestButton,
	)

	mainLayout := container.NewBorder(topPanel, nil, nil, nil, scrollArea)
	window.SetContent(mainLayout)
	window.ShowAndRun()
}
