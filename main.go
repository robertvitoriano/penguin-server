package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"

	"github.com/gorilla/mux"
)

type User struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

var users = []User{
	{ID: "1", Name: "Penguin 1"},
	{ID: "2", Name: "Penguin 2"},
}

func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}
func getUser(w http.ResponseWriter, r *http.Request) {
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

func createUser(responseWriter http.ResponseWriter, request *http.Request) {

	min, max := 10, 255
	var r = rand.Intn(max-min+1) + min
	var g = rand.Intn(max-min+1) + min
	var b = rand.Intn(max-min+1) + min
	var a = rand.Intn(100)

	newColor := fmt.Sprintf("rgba(%d, %d, %d, %.2f)", r, g, b, float64(a)/100)

	var newUser User
	json.NewDecoder(request.Body).Decode(&newUser)

	newUser.Color = newColor
	users = append(users, newUser)

	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(newUser)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/users/{id}", getUser).Methods("GET")
	r.HandleFunc("/users", createUser).Methods("POST")
	r.HandleFunc("/users", getUsers).Methods("GET")

	fmt.Println("Server running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
