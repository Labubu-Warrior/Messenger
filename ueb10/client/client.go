package client

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"ueb10/protocol"
)

// ChatClient kapselt die Verbindung und Ein-/Ausgabe eines Chat-Clients.
type ChatClient struct {
	Name   string
	Conn   net.Conn
	Reader *bufio.Reader
}

// NewChatClient stellt eine TCP-Verbindung zum ChatServer her.
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

// Register sendet einen gewünschten Namen an den Server und wertet die Antwort aus.
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

// ReadLoop liest Nachrichten vom Server und beendet sich bei GOODBYE oder Verbindungsende.
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

func (c *ChatClient) showClientList(text string) {
	list := strings.TrimPrefix(text, protocol.MESSAGE_CLIENTS_PREFIX)

	if strings.TrimSpace(list) == "" {
		fmt.Println(protocol.CLIENT_TEXT_NO_CLIENTS)
	} else {
		fmt.Println(protocol.CLIENT_TEXT_CLIENT_LIST, list)
	}
}

// SendInput verarbeitet eine Benutzereingabe und sendet den passenden Befehl an den Server.
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

func (c *ChatClient) sendMessageCommand(parts []string) {
	if len(parts) < protocol.MIN_MESSAGE_PARTS {
		fmt.Println(protocol.CLIENT_TEXT_MESSAGE_USAGE)
	} else {
		message := strings.Join(parts[1:], " ")
		fmt.Fprintf(c.Conn, "%s %s\n", protocol.COMMAND_MESSAGE, message)
	}
}

func (c *ChatClient) sendPrivateCommand(parts []string) {
	if len(parts) < protocol.MIN_PRIVATE_PARTS {
		fmt.Println(protocol.CLIENT_TEXT_PRIVATE_USAGE)
	} else {
		name := parts[1]
		message := strings.Join(parts[2:], " ")
		fmt.Fprintf(c.Conn, "%s %s %s\n", protocol.COMMAND_PRIVATE, name, message)
	}
}

// WriteLoop liest Benutzereingaben von der Konsole und sendet sie an den Server.
func (c *ChatClient) WriteLoop() {
	stdin := bufio.NewReader(os.Stdin)

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
