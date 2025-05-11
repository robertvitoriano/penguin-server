package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/robertvitoriano/penguin-server/internal/tiled"
)

func LoadLevel(w http.ResponseWriter, r *http.Request) {
	var request struct {
		LevelName string `json:"level_name"`
	}

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		log.Panic("level name is required")
	}

	tileMap := tiled.NewTileMap(fmt.Sprintf("assets/maps/%v.json", request.LevelName))

	fmt.Println(tileMap)

}
