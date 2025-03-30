package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/auth"
	gameEvents "github.com/robertvitoriano/penguin-server/enums"
	"github.com/robertvitoriano/penguin-server/models"
	"github.com/robertvitoriano/penguin-server/repositories"
)

type Websocket struct {
	connection *websocket.Conn
}

type Message struct {
	Event    string          `json:"event"`
	Token    string          `json:"token"`
	Position models.Position `json:"position"`
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

		if message.Event == "close" {
			break
		}

		claims, err := auth.ParseToken(message.Token)

		if message.Event == string(gameEvents.START_GAME) {

			for _, player := range repositories.Players {
				if player.ID == claims["id"] {

					player.Position.X = message.Position.X
					player.Position.Y = message.Position.Y

					playersJSON, err := json.Marshal(repositories.Players)

					if err != nil {
						fmt.Println("Error conveting players to json")
					}

					conn.WriteMessage(messageType, playersJSON)

					break
				}
			}

			if err != nil {
				fmt.Println("Error parsing token")
			}
		}

	}

	conn.Close()

}
