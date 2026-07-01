package server

import (
	"fmt"
	"regexp"
	"sort"
	"strings"
	"ueb10/protocol"
)

const VALID_NAME_PATTERN = `^[A-Za-z0-9_-]+$`

var validNamePattern = regexp.MustCompile(VALID_NAME_PATTERN)

// Methode von ChatServer.
// Verarbeitet einen Registrierungsantrag aus dem register-Kanal.
// Prüft den gewünschten Namen auf Gültigkeit und Eindeutigkeit.
// Liefert das Ergebnis über req.Resp zurück und legt erfolgreiche Namen zuerst in pendingNames ab.
func (s *ChatServer) handleRegister(req RegisterRequest) {
	name := strings.TrimSpace(req.Name)
	_, activeExists := s.clients[name]
	_, pendingExists := s.pendingNames[name]

	// Normale Hilfsfunktion.
	// Prüft, ob ein Benutzername die erlaubte Länge und das erlaubte Zeichenmuster erfüllt.
	// Liefert true, wenn der Name für die Registrierung gültig ist.
	if !isValidClientName(name) {
		req.Resp <- RegisterResponse{
			OK:     false,
			Reason: protocol.RESPONSE_NAME_INVALID,
			Client: nil,
		}
	} else if activeExists || pendingExists {
		req.Resp <- RegisterResponse{
			OK:     false,
			Reason: protocol.RESPONSE_NAME_TAKEN,
			Client: nil,
		}
	} else {
		client := NewConnectedClient(name, req.Conn)
		s.pendingNames[name] = struct{}{} //

		req.Resp <- RegisterResponse{
			OK:     true,
			Reason: protocol.RESPONSE_NAME_OK,
			Client: client,
		}
	}
}

func isValidClientName(name string) bool {
	return len(name) <= protocol.MAX_NAME_LENGTH && validNamePattern.MatchString(name) //Мы используем:validNamePattern.MatchString(name) соответствует ли name разрешённому шаблону?
}

// Prüft, ob ein Benutzername die erlaubte Länge und das erlaubte Zeichenmuster erfüllt.
// Liefert true, wenn der Name für die Registrierung gültig ist.
func (s *ChatServer) handleJoin(client *ConnectedClient) {
	delete(s.pendingNames, client.Name)
	s.clients[client.Name] = client

	names := s.getClientNames(client.Name)

	s.sendOrLog(client, protocol.MESSAGE_CLIENTS_PREFIX+strings.Join(names, ","))

	systemText := fmt.Sprintf(protocol.TEXT_SYSTEM_JOIN, client.Name)
	s.sendToAllExcept(client.Name, systemText)

	fmt.Printf(protocol.SERVER_TEXT_CLIENT_CONNECTED, client.Name)
}

// Methode von ChatServer.
// Entfernt einen reservierten Namen wieder aus pendingNames.
// Wird benutzt, wenn die Registrierung nach der Namensreservierung doch nicht abgeschlossen werden konnte.
func (s *ChatServer) handleCancelPending(name string) {
	delete(s.pendingNames, name)
}

// Methode von ChatServer.
// Entfernt einen aktiven Client aus dem Chat.
// Schließt seinen Sendekanal und informiert die übrigen Clients über das Verlassen.
func (s *ChatServer) handleLeave(client *ConnectedClient) {
	_, exists := s.clients[client.Name]

	if exists {
		delete(s.clients, client.Name)
		close(client.Send)

		systemText := fmt.Sprintf(protocol.TEXT_SYSTEM_LEAVE, client.Name)
		s.sendToAllExcept(client.Name, systemText)

		fmt.Printf(protocol.SERVER_TEXT_CLIENT_DISCONNECTED, client.Name)
	}
}

// Methode von ChatServer.
// Verarbeitet eine öffentliche Chatnachricht.
// Formatiert die Nachricht und sendet sie an alle Clients außer dem Sender.
func (s *ChatServer) handleBroadcast(message BroadcastMessage) {
	out := fmt.Sprintf(protocol.TEXT_PUBLIC_MSG, message.From, message.Text)
	s.sendToAllExcept(message.From, out)
}

// Methode von ChatServer.
// Verarbeitet eine private Nachricht zwischen zwei Clients.
// Sendet die Nachricht nur an den Empfänger oder meldet dem Sender,
// wenn der Empfänger nicht gefunden wurde.
func (s *ChatServer) handlePrivateMessage(pm PrivateMessage) {
	receiver, exists := s.clients[pm.To]

	if exists {
		out := fmt.Sprintf(protocol.TEXT_PRIVATE_MSG, pm.From, pm.Text)
		s.sendOrLog(receiver, out)
	} else {
		sender, ok := s.clients[pm.From]

		if ok {
			s.sendOrLog(sender, protocol.ERROR_RECEIVER_NOT_FOUND)
		}
	}
}

// Methode von ChatServer.
// Beantwortet eine LIST-Anfrage eines Clients.
// Erstellt die aktuelle Clientliste und sendet sie nur an den anfragenden Client.
func (s *ChatServer) handleListRequest(client *ConnectedClient) {
	names := s.getClientNames("")
	s.sendOrLog(client, protocol.MESSAGE_CLIENTS_PREFIX+strings.Join(names, ","))
}

// Methode von ChatServer.
// Erstellt eine sortierte Liste der aktuell angemeldeten Clientnamen.
// excludeName kann genutzt werden, um einen bestimmten Client aus der Liste auszuschließen.
func (s *ChatServer) getClientNames(excludeName string) []string {
	names := make([]string, 0, len(s.clients))

	for name := range s.clients {
		if excludeName == "" || name != excludeName {
			names = append(names, name)
		}
	}

	sort.Strings(names)

	return names
}

// Methode von ChatServer.
// Sendet eine Nachricht an alle aktiven Clients außer excludeName.
// Wird für Systemmeldungen und öffentliche Chatnachrichten benutzt.
func (s *ChatServer) sendToAllExcept(excludeName string, message string) {
	for name, client := range s.clients {
		if name != excludeName {
			s.sendOrLog(client, message)
		}
	}
}

// Methode von ChatServer.
// Übergibt eine Nachricht an den Sendekanal eines einzelnen Clients.
// Falls das nicht möglich ist, wird der Fehler nur in der Serverkonsole protokolliert.
func (s *ChatServer) sendOrLog(client *ConnectedClient, message string) { // diese funktion hat text bekommt von conecteed_client go
	err := client.SendMessage(message)

	if err != nil {
		fmt.Printf(protocol.SERVER_TEXT_SEND_FAILED, client.Name, err) // а если отправка не удалась — пишет ошибку в лог сервера.
	}
}

/*handleBroadcast делает готовый текст
→ sendToAllExcept выбирает получателей
→ sendOrLog передаёт текст конкретному клиенту     ^
→ SendMessage кладёт текст в client.Send
→ WriteLoop отправляет клиенту
*/
