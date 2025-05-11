package tiled

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/robertvitoriano/penguin-server/internal/models"
)

type TileMap struct {
	Enemies []models.Enemy
	Items   []models.Item
}

func NewTileMap(path string) *TileMap {
	jsonFile, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	defer jsonFile.Close()

	jsonFileData, err := io.ReadAll(jsonFile)
	if err != nil {
		log.Fatal(err)
	}

	var parsedData LevelData
	err = json.Unmarshal(jsonFileData, &parsedData)
	if err != nil {
		log.Fatal(err)
	}
	var tileMap = &TileMap{}

	for _, layer := range parsedData.Layers {
		switch layer.Name {
		case "enemies":
			{
				for _, enemy := range layer.Objects {

					var idPropertyIndex *int

					for index, property := range enemy.Properties {
						if property.Name == "id" {
							idPropertyIndex = &index
						}
					}

					if idPropertyIndex != nil {
						id, err := strconv.Atoi(enemy.Properties[*idPropertyIndex].Value)

						if err != nil {
							log.Fatalf("Failed to convert id to int: %v", err)
						}
						tileMap.Enemies = append(tileMap.Enemies, models.Enemy{
							ID:   &id,
							Name: enemy.Name,
							Position: &models.Position{
								X: &enemy.X,
								Y: &enemy.Y,
							},
						})
					}
				}
				break
			}
		case "items":
			{
				for _, item := range layer.Objects {

					var idPropertyIndex *int
					var typePropertyIndex *int

					for index, property := range item.Properties {
						if property.Name == "id" {
							idPropertyIndex = &index
						}
						if property.Name == "type" {
							typePropertyIndex = &index
						}
					}

					if idPropertyIndex != nil && typePropertyIndex != nil {
						id, err := strconv.Atoi(item.Properties[*idPropertyIndex].Value)

						if err != nil {
							log.Fatalf("Failed to convert id to int: %v", err)
						}
						tileMap.Items = append(tileMap.Items, models.Item{
							Type: item.Properties[*idPropertyIndex].Value,
							ID:   id,
							Position: &models.Position{
								X: &item.X,
								Y: &item.Y,
							},
						})
					}
				}
			}
		}

	}
	fmt.Println(tileMap)
	return tileMap
}

// func (t *TileMap) ChangeMap(path string) {
// 	t.path = path

// }
