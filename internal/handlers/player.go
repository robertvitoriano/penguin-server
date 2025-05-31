package handlers

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
	"github.com/robertvitoriano/penguin-server/internal/models"
	"github.com/robertvitoriano/penguin-server/internal/repositories/mysql"
	"github.com/robertvitoriano/penguin-server/internal/repositories/redis"

	"gorm.io/gorm"
)

type PlayerCreationResponse struct {
	Player models.Player `json:"player"`
	Token  string        `json:"token"`
	Result string        `json:"result"`
}

type PlayerHandler struct {
}

func NewPlayerHandler() *PlayerHandler {
	return &PlayerHandler{}
}

func (p *PlayerHandler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players := redis.GetPlayers()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}

func (p *PlayerHandler) GetPlayer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	playerID := vars["id"]

	var playerFound *models.Player
	for _, player := range redis.GetPlayers() {
		if player.ID == playerID {
			playerFound = player
			break
		}
	}

	if playerFound == nil {
		http.Error(w, "Player not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	err := json.NewEncoder(w).Encode(playerFound)
	if err != nil {
		http.Error(w, "Failed to encode player", http.StatusInternalServerError)
	}
}

func (p *PlayerHandler) CreatePlayer(responseWriter http.ResponseWriter, request *http.Request, ws *Websocket, db *gorm.DB) {
	responseWriter.Header().Set("Content-Type", "application/json")

	var newPlayer models.Player

	err := json.NewDecoder(request.Body).Decode(&newPlayer)
	if err != nil {
		http.Error(responseWriter, "Invalid request body", http.StatusBadRequest)
		return
	}

	existingUser, err := redis.FindPlayerByUsername(newPlayer.Username)

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

	redis.CreatePlayer(&newPlayer)
	playerPersistencyRepository := mysql.NewPlayerRepository(db)
	playerPersistencyRepository.CreatePlayer(&newPlayer)
	ws.Broadcast([]byte(`{"message":"User created"}`))

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

func (p *PlayerHandler) GetPlayerMessages(responseWriter http.ResponseWriter, request *http.Request, ws *Websocket) {

}
