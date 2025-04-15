package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/robertvitoriano/penguin-server/models"
	"github.com/robertvitoriano/penguin-server/repositories"
)

type PlayerCreationResponse struct {
	Player models.Player `json:"player"`
	Token  string        `json:"token"`
	Result string        `json:"result"`
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

func CreatePlayer(responseWriter http.ResponseWriter, request *http.Request, ws *Websocket) {
	responseWriter.Header().Set("Content-Type", "application/json")

	var newPlayer models.Player

	err := json.NewDecoder(request.Body).Decode(&newPlayer)
	if err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, err := repositories.FindPlayerByUsername(newPlayer.Username)

	if err == nil {

		responseWriter.WriteHeader(http.StatusCreated)

		response := PlayerCreationResponse{
			Player: existingUser,
			Result: "User already exists",
		}

		if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
			log.Println("Error encoding response:", err)
		}
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

	repositories.CreatePlayer(&newPlayer)

	ws.Broadcast([]byte("User created"))

	var (
		secretKey   string
		jwtToken    *jwt.Token
		signedToken string
	)

	secretKey = os.Getenv("JWT_SECRET_KEY")
	jwtToken = jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"id":       newPlayer.ID,
			"username": newPlayer.Username,
			"color":    newPlayer.Color,
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

func GetPlayerMessages(responseWriter http.ResponseWriter, request *http.Request, ws *Websocket) {

}
