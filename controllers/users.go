package controllers

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/robertvitoriano/penguin-server/models"
)

var users = []models.User{}

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

func CreateUser(responseWriter http.ResponseWriter, request *http.Request) {

	min, max := 10, 255
	var r = rand.Intn(max-min+1) + min
	var g = rand.Intn(max-min+1) + min
	var b = rand.Intn(max-min+1) + min
	var a = rand.Intn(100)

	newColor := fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, float64(a)/100)

	var newUser models.User
	json.NewDecoder(request.Body).Decode(&newUser)

	newUser.Color = newColor
	newUUID := uuid.New().String()
	newUser.ID = newUUID
	users = append(users, newUser)

	responseWriter.Header().Set("Content-Type", "application/json")

	json.NewEncoder(responseWriter).Encode(newUser)
}
