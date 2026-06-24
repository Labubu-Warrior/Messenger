package protocol

// Netzwerkeinstellungen
const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "6666"
)

// Steuerbefehle (Client -> Server)
const (
	CMD_PRIVATE = "/msg"  // Format: /msg Name Text
	CMD_LIST    = "/list" // Fordert die Liste der Online-Nutzer an
	CMD_QUIT    = "/quit" // Beendet den Client sauber
)

// Systemnachrichten und Prompts
const (
	PROMPT_NAME        = "[System] Bitte wähle einen eindeutigen Namen:"
	ERR_NAME_TAKEN     = "[System] Name ist bereits vergeben oder ungültig! Neuer Versuch:"
	MSG_WELCOME        = "[System] Willkommen, %s! (Befehle: /msg <Name> <Text>, /list, /quit)"
	MSG_JOIN           = "+++ %s hat den Chat betreten +++"
	MSG_LEAVE          = "--- %s hat den Chat verlassen ---"
	MSG_CLIENT_LIST    = "[System] Aktuell online: %s"
	ERR_USER_NOT_FOUND = "[System] Fehler: Benutzer '%s' nicht gefunden."
	ERR_PRIVATE_FORMAT = "[System] Fehler: Falsches Format. Nutzung: /msg <Name> <Text>"
)

// Endezeichen für den Socket-Stream
const ENDSIGN = '#'
