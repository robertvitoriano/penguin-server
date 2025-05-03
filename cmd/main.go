package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/robertvitoriano/penguin-server/internal/controllers"
	"github.com/rs/cors"
)

func main() {
	ws := controllers.NewWebsocket()
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	router.HandleFunc("/players/{id}", controllers.GetPlayer).Methods("GET")
	router.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {
		controllers.CreatePlayer(w, r, ws)
	}).Methods("POST")
	router.HandleFunc("/players", controllers.GetPlayers).Methods("GET")
	router.HandleFunc("/ws", ws.ServeWebsocket).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://penguim-adventure.robertvitoriano.com", "http://localhost:8000", fmt.Sprintf("http://%v:8000", os.Getenv("COMPUTER_IP"))},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)

	fmt.Println("Server running on port 7777...")
	log.Fatal(http.ListenAndServe(":7777", handler))
}
