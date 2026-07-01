package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"strings"
	"ueb10/protocol"
)

type ChatClient struct {
	Name   string
	Conn   net.Conn
	Reader *bufio.Reader
}

// Erstellt einen ChatClient für die angegebene Serveradresse.
// Nutzt host und port zum Aufbau einer TCP-Verbindung über net.Conn.
// Initialisiert zusätzlich einen bufio.Reader für Serverantworten.
// Liefert den fertigen ChatClient oder einen Verbindungsfehler.
func NewChatClient(host string, port string) (*ChatClient, error) {
	conn, err := net.Dial("tcp", host+":"+port)

	if err != nil {
		return nil, err
	}

	return &ChatClient{
		Name:   "",
		Conn:   conn,
		Reader: bufio.NewReader(conn),
	}, nil
}

// Registriert den Client mit dem gewünschten Namen beim Server.
// Sendet den bereinigten Namen über c.Conn und liest die Antwort über c.Reader.
// Setzt c.Name nur bei erfolgreicher Registrierung.
// Liefert nil bei Erfolg oder einen passenden Registrierungsfehler.
func (c *ChatClient) Register(name string) error {
	name = strings.TrimSpace(name)

	_, err := fmt.Fprintf(c.Conn, "%s\n", name)

	if err != nil {
		return err
	}

	line, err := c.Reader.ReadString('\n')

	if err != nil {
		return err
	}

	response := strings.TrimSpace(line)

	if response == protocol.RESPONSE_NAME_OK {
		c.Name = name
		return nil
	}

	if response == protocol.RESPONSE_NAME_TAKEN {
		return errors.New(protocol.ERROR_NAME_TAKEN)
	}

	if response == protocol.RESPONSE_NAME_INVALID {
		return errors.New(protocol.ERROR_NAME_INVALID)
	}

	return fmt.Errorf("unerwartete Antwort: %s", response)
}

// Liest dauerhaft eingehende Servernachrichten über c.Reader.
// Gibt normale Nachrichten aus und signalisiert main über done,
// wenn die Verbindung endet oder der Server GOODBYE sendet.
// Die Funktion liefert keinen Rückgabewert.
func (c *ChatClient) ReadLoop(done chan bool) {
	for {
		line, err := c.Reader.ReadString('\n')

		if err != nil {
			fmt.Println()
			fmt.Println(protocol.CLIENT_TEXT_CONNECTION_CLOSED)
			done <- true
			return
		}

		text := strings.TrimSpace(line)

		c.showServerMessage(text)

		if text == protocol.MESSAGE_GOODBYE {
			done <- true
			return
		}
	}
}

// Methode von ChatClient
// Entscheidet, wie eine bereits bereinigte Servernachricht angezeigt wird.
// Client-Listen werden an showClientList weitergegeben,
// alle anderen Nachrichten werden direkt im Terminal ausgegeben.
func (c *ChatClient) showServerMessage(text string) {
	if text == "" {
		return
	}

	if strings.HasPrefix(text, protocol.MESSAGE_CLIENTS_PREFIX) {
		c.showClientList(text)
	} else {
		fmt.Println(text)
	}
}

// Methode von ChatClient
// Zeigt die vom Server gesendete Client-Liste an.
// Erwartet eine Nachricht mit dem Protokoll-Prefix CLIENTS:
// und entfernt diesen Prefix vor der Ausgabe.
// Liefert keinen Wert zurück und verändert keinen Client-Zustand.
func (c *ChatClient) showClientList(text string) {
	list := strings.TrimPrefix(text, protocol.MESSAGE_CLIENTS_PREFIX)

	if strings.TrimSpace(list) == "" {
		fmt.Println(protocol.CLIENT_TEXT_NO_CLIENTS)
	} else {
		fmt.Println(protocol.CLIENT_TEXT_CLIENT_LIST, list)
	}
}

// Methode von ChatClient
// Analysiert eine Benutzereingabe und übersetzt sie in ein Serverkommando.
// Nutzt c.Conn zum Senden an den Server.
// Bekannte Befehle werden gezielt behandelt, unbekannte Eingaben
// werden als normale Chatnachricht gesendet.
func (c *ChatClient) SendInput(text string) {
	text = strings.TrimSpace(text)

	if text == "" {
		return
	}

	parts := strings.Fields(text)
	command := strings.ToLower(parts[0])

	switch command {
	case protocol.USER_COMMAND_MESSAGE:
		c.sendMessageCommand(parts)
	case protocol.USER_COMMAND_PRIVATE:
		c.sendPrivateCommand(parts)
	case protocol.USER_COMMAND_LIST:
		fmt.Fprintf(c.Conn, "%s\n", protocol.COMMAND_LIST)
	default:
		fmt.Fprintf(c.Conn, "%s %s\n", protocol.COMMAND_MESSAGE, text)
	}
}

// Methode von ChatClient
// Verarbeitet den Benutzerbefehl msg.
// Erwartet die bereits zerlegte Eingabe in parts.
// Baut den Nachrichtentext aus allen Teilen nach dem Befehl zusammen
// und sendet ihn über c.Conn als öffentliche Nachricht an den Server.
func (c *ChatClient) sendMessageCommand(parts []string) {
	if len(parts) < protocol.MIN_MESSAGE_PARTS {
		fmt.Println(protocol.CLIENT_TEXT_MESSAGE_USAGE)
	} else {
		message := strings.Join(parts[1:], " ")
		fmt.Fprintf(c.Conn, "%s %s\n", protocol.COMMAND_MESSAGE, message)
	}
}

// Methode von ChatClient
// Verarbeitet den Benutzerbefehl privmsg.
// Erwartet in parts den Befehl, den Empfängernamen und den Nachrichtentext.
// Sendet die private Nachricht über c.Conn an den Server
// oder gibt bei fehlenden Angaben einen Nutzungshinweis aus.
func (c *ChatClient) sendPrivateCommand(parts []string) {
	if len(parts) < protocol.MIN_PRIVATE_PARTS {
		fmt.Println(protocol.CLIENT_TEXT_PRIVATE_USAGE)
	} else {
		name := parts[1]
		message := strings.Join(parts[2:], " ")
		fmt.Fprintf(c.Conn, "%s %s %s\n", protocol.COMMAND_PRIVATE, name, message)
	}
}

// Methode von ChatClient
// Liest dauerhaft Benutzereingaben aus stdin.
// Sendet quit direkt als QUIT-Kommando an den Server,
// alle anderen Eingaben werden an SendInput übergeben.
// Die Funktion endet bei quit oder bei einem Lesefehler.
func (c *ChatClient) WriteLoop(stdin *bufio.Reader) {
	for {
		fmt.Print(protocol.CLIENT_TEXT_INPUT_PROMPT)
		line, err := stdin.ReadString('\n')

		if err != nil {
			fmt.Fprintf(c.Conn, "%s\n", protocol.COMMAND_QUIT)
			return
		}

		text := strings.TrimSpace(line)
		parts := strings.Fields(text)

		if len(parts) > 0 && strings.ToLower(parts[0]) == protocol.USER_COMMAND_QUIT {
			fmt.Fprintf(c.Conn, "%s\n", protocol.COMMAND_QUIT)
			return
		}

		c.SendInput(line)
	}
}
