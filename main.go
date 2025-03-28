package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/controllers"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

var ws Websocket

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/users/{id}", controllers.GetUser).Methods("GET")
	router.HandleFunc("/users", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreateUser(w, r, ws.connection)
	}).Methods("POST")
	router.HandleFunc("/users", controllers.GetUsers).Methods("GET")

	router.HandleFunc("/", ws.serveWebsocket).Methods("GET")

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", router))
}
