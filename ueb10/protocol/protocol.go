package protocol

// Verbindungseinstellungen
const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "4000"
)

// Befehle, die der Client an den Server sendet
const (
	COMMAND_MESSAGE = "MSG"
	COMMAND_PRIVATE = "PRIVMSG"
	COMMAND_LIST    = "LIST"
	COMMAND_QUIT    = "QUIT"
)

// Befehle, die der Benutzer im Client eingibt
const (
	USER_COMMAND_MESSAGE = "msg"
	USER_COMMAND_PRIVATE = "privmsg"
	USER_COMMAND_LIST    = "list"
	USER_COMMAND_QUIT    = "quit"
)

// Antworten des Servers während der Anmeldung
const (
	RESPONSE_NAME_OK      = "NAME_OK"
	RESPONSE_NAME_TAKEN   = "NAME_TAKEN"
	RESPONSE_NAME_INVALID = "NAME_INVALID"
)

// Servernachrichten
const (
	MESSAGE_CLIENTS_PREFIX = "CLIENTS:"
	MESSAGE_GOODBYE        = "GOODBYE"
)

// Mindestanzahl von Eingabeteilen für Befehle
const (
	MIN_MESSAGE_PARTS = 2
	MIN_PRIVATE_PARTS = 3
)

// Fehlermeldungen vom Server an Clients
const (
	ERROR_RECEIVER_NOT_FOUND = "ERROR: Empfänger nicht gefunden"
	ERROR_MESSAGE_MISSING    = "ERROR: Nachricht fehlt"
	ERROR_UNKNOWN_COMMAND    = "ERROR: Unbekannter Befehl"
	ERROR_PRIVATE_USAGE      = "ERROR: Verwendung: PRIVMSG <Name> <Text>"
)

// Fehlermeldungen für Registrierung im Client
const (
	ERROR_NAME_TAKEN   = "Name schon vergeben"
	ERROR_NAME_INVALID = "Name ungültig: nur Buchstaben, Zahlen, _ und - erlaubt"
)

// Textvorlagen für Chatnachrichten
const (
	TEXT_SYSTEM_JOIN  = "[System]: %s hat den Chat betreten"
	TEXT_SYSTEM_LEAVE = "[System]: %s hat den Chat verlassen"
	TEXT_PUBLIC_MSG   = "[%s]: %s"
	TEXT_PRIVATE_MSG  = "[Privat von %s]: %s"
)

// Server-Konsolenausgabe
const (
	SERVER_TEXT_RUNNING             = "Server läuft auf Port %s\n"
	SERVER_TEXT_ACCEPT_ERROR        = "Accept-Fehler: %v\n"
	SERVER_TEXT_CLIENT_CONNECTED    = "[%s] verbunden\n"
	SERVER_TEXT_CLIENT_DISCONNECTED = "[%s] getrennt\n"
	SERVER_TEXT_SEND_FAILED         = "Nachricht an %s konnte nicht gesendet werden: %v\n"
	SERVER_ERROR_SEND_BUFFER_FULL   = "Sendepuffer von %s ist voll"
)

// Texte für die Client-Ausgabe
const (
	CLIENT_TEXT_CONNECTION_CLOSED = "[Chat] Verbindung beendet."
	CLIENT_TEXT_ASK_NAME          = "[Chat] Bitte Namen eingeben."
	CLIENT_TEXT_ALLOWED_NAME      = "[Chat] Erlaubt sind Buchstaben, Zahlen, _ und -"
	CLIENT_TEXT_INPUT_ERROR       = "[Chat] Eingabefehler."
	CLIENT_TEXT_COMMANDS          = "[Chat] Befehle: msg <Text>, privmsg <Name> <Text>, list, quit"
	CLIENT_TEXT_NO_CLIENTS        = "[Chat] Bereits angemeldete Clients: keine"
	CLIENT_TEXT_CLIENT_LIST       = "[Chat] Bereits angemeldete Clients:"
	CLIENT_TEXT_MESSAGE_USAGE     = "[Chat] Verwendung: msg <Text>"
	CLIENT_TEXT_PRIVATE_USAGE     = "[Chat] Verwendung: privmsg <Name> <Text>"
	CLIENT_TEXT_WELCOME           = "[Chat] Willkommen, %s!\n"
	CLIENT_TEXT_CONNECTION_ERROR  = "[Chat] Verbindungsfehler: %v\n"
	CLIENT_TEXT_NAME_PROMPT       = "Name: "
	CLIENT_TEXT_INPUT_PROMPT      = "> "
)
