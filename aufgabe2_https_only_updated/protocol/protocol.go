// Package protocol enthält alle benannten Konstanten des Netzwerkprotokolls.
package protocol

import "time"

const (
	MAIN_SERVER_HOST       = "localhost"
	MAIN_SERVER_PORT       = "4000"
	VALUEPAIR_SERVER_HOST  = "localhost"
	VALUEPAIR_SERVER_PORT  = "16667"
	SCHEDULE_XML_URL       = "https://intern.fh-wedel.de/~tho/api/splan.php"
	NETWORK_TYPE_TCP       = "tcp"
	JSON_LINE_DELIMITER    = '\n'
	VALUEPAIR_END_MARKER   = '#'
	VALUEPAIR_SEPARATOR    = ','
	VALUEPAIR_ERROR_PREFIX = "ERROR:"
)

const (
	REQUEST_SUGGESTIONS = "SUGGEST"
	REQUEST_QUIT        = "QUIT"

	RESPONSE_WELCOME     = "WELCOME"
	RESPONSE_SUGGESTIONS = "SUGGESTIONS"
	RESPONSE_ERROR       = "ERROR"
	RESPONSE_GOODBYE     = "GOODBYE"
)

const (
	MAX_CLIENT_LINE_BYTES     = 64 * 1024
	MAX_VALUEPAIR_LINE_BYTES  = 16 * 1024
	MAX_XML_SIZE_BYTES        = 8 * 1024 * 1024
	NETWORK_READ_BUFFER_BYTES = 4 * 1024
)

const (
	HTTP_REQUEST_TIMEOUT      = 30 * time.Second
	CLIENT_CONNECT_TIMEOUT    = 10 * time.Second
	CLIENT_READ_TIMEOUT       = 30 * time.Minute
	VALUEPAIR_REQUEST_TIMEOUT = 20 * time.Second
)

const (
	MESSAGE_CONNECTION_CLOSED          = "Verbindung beendet"
	MESSAGE_INVALID_REQUEST_FORMAT     = "ungültiges Anfrageformat"
	MESSAGE_UNKNOWN_COMMAND            = "unbekannter Befehl"
	MESSAGE_REQUEST_TOO_LONG           = "Anfrage ist zu lang"
	MESSAGE_VALUEPAIR_REQUEST_TOO_LONG = "Wertepaar-Anfrage ist zu lang"
	MESSAGE_XML_UNAVAILABLE            = "aktuelle Stundenplandaten konnten nicht geladen werden"
)
