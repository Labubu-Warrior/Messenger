package main

import (
	client "Merd_Chat/ClientLogik"
	"Merd_Chat/protocol"
	"fmt"
)

func main() {
	fmt.Println("Verbinde mit Chat-Server...")
	c, err := client.Connect(protocol.SERVER_HOST, protocol.SERVER_PORT)
	if err != nil {
		fmt.Println("Fehler:", err)
		return
	}

	c.Run()
}
