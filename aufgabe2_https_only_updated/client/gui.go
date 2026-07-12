package client

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"

	"aufgabe2/shared"
)

const (
	GUI_APP_ID                   = "de.fh-wedel.aufgabe2.terminvorschlaege"
	GUI_WINDOW_TITLE             = "Terminvorschläge aus dem Stundenplan"
	GUI_WINDOW_WIDTH             = 920
	GUI_WINDOW_HEIGHT            = 620
	GUI_SPLIT_OFFSET             = 0.38
	GUI_STATUS_CONNECTED         = "Mit dem Server verbunden."
	GUI_STATUS_SELECTION_CLEARED = "Auswahl wurde gelöscht."
	GUI_STATUS_CALCULATING       = "Der Server berechnet die Terminvorschläge …"
	GUI_SEARCH_PLACEHOLDER       = "Kürzel suchen …"
	GUI_RESULT_PLACEHOLDER       = "Hier erscheinen die Terminvorschläge."
	GUI_TITLE_EMPLOYEE_SELECTION = "Mitarbeiter auswählen"
	GUI_TITLE_RESULT             = "Ergebnis"
	GUI_BUTTON_CALCULATE         = "Terminvorschläge berechnen"
	GUI_BUTTON_CLEAR             = "Auswahl löschen"
	GUI_NO_SUGGESTIONS           = "Keine gemeinsamen freien Termine gefunden."
	GUI_STATUS_RESULT_FORMAT     = "%d Terminvorschläge gefunden."
	GUI_ERROR_PREFIX             = "Fehler: "
	GUI_TIME_LINE_FORMAT         = "  %s – %s\n"
	GUI_DAY_LINE_FORMAT          = "%s\n"
)

// RunGUI öffnet die grafische Client-Oberfläche.
func RunGUI(networkClient *Client) {
	application := app.NewWithID(GUI_APP_ID)
	window := application.NewWindow(GUI_WINDOW_TITLE)
	window.Resize(fyne.NewSize(GUI_WINDOW_WIDTH, GUI_WINDOW_HEIGHT))
	window.SetCloseIntercept(func() {
		_ = networkClient.Close()
		window.Close()
	})

	statusLabel := widget.NewLabel(GUI_STATUS_CONNECTED)
	statusLabel.Wrapping = fyne.TextWrapWord

	searchEntry := widget.NewEntry()
	searchEntry.SetPlaceHolder(GUI_SEARCH_PLACEHOLDER)

	resultEntry := widget.NewMultiLineEntry()
	resultEntry.SetPlaceHolder(GUI_RESULT_PLACEHOLDER)
	resultEntry.Wrapping = fyne.TextWrapWord
	resultEntry.Disable()

	selectedCodes := make(map[string]bool)
	listBox := container.NewVBox()

	filteredCodes := func() []string {
		query := strings.ToLower(strings.TrimSpace(searchEntry.Text))
		codes := append([]string(nil), networkClient.Codes...)
		sort.Strings(codes)
		if query == "" {
			return codes
		}
		filtered := make([]string, 0)
		for _, code := range codes {
			if strings.Contains(strings.ToLower(code), query) {
				filtered = append(filtered, code)
			}
		}
		return filtered
	}

	var rebuildList func()
	rebuildList = func() {
		listBox.RemoveAll()
		for _, code := range filteredCodes() {
			currentCode := code
			checkBox := widget.NewCheck(currentCode, func(checked bool) {
				selectedCodes[currentCode] = checked
			})
			checkBox.SetChecked(selectedCodes[currentCode])
			listBox.Add(checkBox)
		}
		listBox.Refresh()
	}
	searchEntry.OnChanged = func(string) {
		rebuildList()
	}
	rebuildList()

	requestButton := widget.NewButton(GUI_BUTTON_CALCULATE, nil)
	clearButton := widget.NewButton(GUI_BUTTON_CLEAR, func() {
		selectedCodes = make(map[string]bool)
		rebuildList()
		resultEntry.SetText("")
		statusLabel.SetText(GUI_STATUS_SELECTION_CLEARED)
	})

	requestButton.OnTapped = func() {
		codes := make([]string, 0)
		for _, code := range networkClient.Codes {
			if selectedCodes[code] {
				codes = append(codes, code)
			}
		}

		requestButton.Disable()
		statusLabel.SetText(GUI_STATUS_CALCULATING)
		go func() {
			suggestions, err := networkClient.RequestSuggestions(codes)
			fyne.Do(func() {
				requestButton.Enable()
				if err != nil {
					statusLabel.SetText(GUI_ERROR_PREFIX + err.Error())
					resultEntry.SetText("")
					return
				}
				statusLabel.SetText(fmt.Sprintf(GUI_STATUS_RESULT_FORMAT, len(suggestions)))
				resultEntry.SetText(formatSuggestions(suggestions))
			})
		}()
	}

	leftHeader := container.NewVBox(
		widget.NewLabelWithStyle(
			GUI_TITLE_EMPLOYEE_SELECTION,
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		searchEntry,
	)
	leftFooter := container.NewHBox(clearButton, layout.NewSpacer(), requestButton)
	leftPanel := container.NewBorder(
		leftHeader,
		leftFooter,
		nil,
		nil,
		container.NewVScroll(listBox),
	)
	rightPanel := container.NewBorder(
		widget.NewLabelWithStyle(
			GUI_TITLE_RESULT,
			fyne.TextAlignLeading,
			fyne.TextStyle{Bold: true},
		),
		statusLabel,
		nil,
		nil,
		resultEntry,
	)

	split := container.NewHSplit(leftPanel, rightPanel)
	split.Offset = GUI_SPLIT_OFFSET
	window.SetContent(container.NewPadded(split))
	window.ShowAndRun()
}

func formatSuggestions(items []shared.Suggestion) string {
	if len(items) == 0 {
		return GUI_NO_SUGGESTIONS
	}

	var builder strings.Builder
	currentDay := ""
	for _, item := range items {
		if item.Day != currentDay {
			if currentDay != "" {
				builder.WriteByte('\n')
			}
			currentDay = item.Day
			fmt.Fprintf(&builder, GUI_DAY_LINE_FORMAT, currentDay)
		}
		fmt.Fprintf(&builder, GUI_TIME_LINE_FORMAT, item.From, item.To)
	}
	return strings.TrimRight(builder.String(), "\n")
}
