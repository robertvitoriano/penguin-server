package controllers

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/robertvitoriano/penguin-server/models"
)

var users = []models.User{}

type UserCreationResponse struct {
	User  models.User `json:"user"`
	Token string      `json:"token"`
}

func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func GetUser(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for _, user := range users {
		if user.ID == params["id"] {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
			return
		}
	}
	http.Error(w, "User not found", http.StatusNotFound)
}

func CreateUser(responseWriter http.ResponseWriter, request *http.Request, websocketConnection *websocket.Conn) {
	responseWriter.Header().Set("Content-Type", "application/json")

	var newUser models.User
	err := json.NewDecoder(request.Body).Decode(&newUser)
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

	newUser.Color = newColor
	newUser.ID = uuid.New().String()

	users = append(users, newUser)

	if websocketConnection != nil {
		err = websocketConnection.WriteMessage(websocket.TextMessage, []byte("User created"))
		if err != nil {
			log.Println("WebSocket write error:", err)
		}
	}

	response := UserCreationResponse{
		User:  newUser,
		Token: "DD",
	}

	responseWriter.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(responseWriter).Encode(response); err != nil {
		log.Println("Error encoding response:", err)
	}
}
