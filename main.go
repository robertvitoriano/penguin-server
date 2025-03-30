package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/robertvitoriano/penguin-server/controllers"
	"github.com/rs/cors"
)

var ws Websocket

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/players/{id}", controllers.GetPlayer).Methods("GET")
	router.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreatePlayer(w, r, ws.connection)
	}).Methods("POST")
	router.HandleFunc("/players", controllers.GetPlayers).Methods("GET")

	router.HandleFunc("/", ws.serveWebsocket).Methods("GET")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	fmt.Println("Server running on port 8080...")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)
	log.Fatal(http.ListenAndServe(":8080", handler))
}
