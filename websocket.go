package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	gameEvents "github.com/robertvitoriano/penguin-server/enums"
)

type Websocket struct {
	connection *websocket.Conn
}

type Message struct {
	Content string `json:"content"`
	Token   string `json:token`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (ws *Websocket) serveWebsocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	ws.connection = conn

	if err != nil {
		log.Println(err)
		return
	}

	fmt.Println("trying to connect...")

	for {
		messageType, data, err := conn.ReadMessage()

		if err != nil {
			fmt.Println("There was an error")
		}

		var message Message

		err = json.Unmarshal(data, &message)

		if err != nil {
			fmt.Println("Error parsing json")
		}

		if message.Content == "close" {
			break
		}

		if message.Content == string(gameEvents.START_GAME) {
			// claims, err := auth.ParseToken(message.Token)

			if err != nil {
				fmt.Println("Error parsing token")
			}
			conn.WriteMessage(messageType, []byte("Hello client"))
		}

	}

	conn.Close()

}
