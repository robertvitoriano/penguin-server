package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/models"
	"github.com/robertvitoriano/penguin-server/repositories"
)

type PlayerCreationResponse struct {
	Player models.Player `json:"player"`
	Token  string        `json:"token"`
}

func GetPlayers(w http.ResponseWriter, r *http.Request) {
	players := repositories.GetPlayers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func GetPlayer(w http.ResponseWriter, r *http.Request) {
	Players := repositories.GetPlayers()

	params := mux.Vars(r)
	for _, player := range Players {
		if player.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(player)
			return
		}
	}
	http.Error(w, "Player not found", http.StatusNotFound)
}

func CreatePlayer(responseWriter http.ResponseWriter, request *http.Request, websocketConnection *websocket.Conn) {
	responseWriter.Header().Set("Content-Type", "application/json")

	var newPlayer models.Player

	err := json.NewDecoder(request.Body).Decode(&newPlayer)
	if err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	min, max := 10, 255
	r := rand.Intn(max-min+1) + min
	g := rand.Intn(max-min+1) + min
	b := rand.Intn(max-min+1) + min
	a := rand.Intn(100)
	newColor := fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, float64(a)/100)

	newPlayer.Color = newColor
	newPlayer.ID = uuid.New().String()

	repositories.CreatePlayer(newPlayer)

	if websocketConnection != nil {
		err = websocketConnection.WriteMessage(websocket.TextMessage, []byte("User created"))
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}

	var (
		secretKey   string
		jwtToken    *jwt.Token
		signedToken string
	)

	secretKey = "hello"
	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       newPlayer.ID,
			"username": newPlayer.Username,
		})

	signedToken, err = jwtToken.SignedString([]byte(secretKey))

	if err != nil {
		log.Fatalf("Error signing token: %v", err)
	}

	response := PlayerCreationResponse{
		Player: newPlayer,
		Token:  signedToken,
	}

	responseWriter.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
	}
}
