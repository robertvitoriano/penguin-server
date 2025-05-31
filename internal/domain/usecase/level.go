package usecase

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
)

type Animation struct {
	Duration int `json:"duration"`
	TileId   int `json:"tileid"`
}

type RawTile struct {
	Animation []Animation `json:"animation"`
	Id        int         `json:"id"`
}
type TileSet struct {
	Columns     int       `json:"columns"`
	FirstGID    int       `json:"firstgid"`
	Image       string    `json:"image"`
	ImageHeight float64   `json:"imageheight"`
	ImageWidth  float64   `json:"imagewidth"`
	Margin      int       `json:"margin"`
	Name        string    `json:"name"`
	Spacing     int       `json:"spacing"`
	TileCount   int       `json:"tilecount"`
	TileHeight  float64   `json:"tileheight"`
	TileWidth   float64   `json:"tilewidth"`
	Tiles       []RawTile `json:"tiles"`
}
type PolylinePoint struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type ObjectProperty struct {
	Name  string `json:"name"`
	Type  string `json:"type"`
	Value string `json:"value"`
}

type Object struct {
	Id         int              `json:"id"`
	Height     float64          `json:"height"`
	Width      float64          `json:"width"`
	X          float64          `json:"x"`
	Y          float64          `json:"y"`
	Visible    bool             `json:"visible"`
	Name       string           `json:"name"`
	Rotation   float64          `json:"rotation"`
	Properties []ObjectProperty `json:"properties"`
	Polyline   []PolylinePoint  `json:"polyline"`
	Ellipse    bool             `json:"ellipse"`
	Type       string           `json:"type"`
}

type Layer struct {
	Data      []int    `json:"data"`
	Height    float64  `json:"height"`
	Width     float64  `json:"width"`
	Id        int      `json:"id"`
	Name      string   `json:"name"`
	Opacity   int      `json:"opacity"`
	Type      string   `json:"type"`
	Visible   bool     `json:"visible"`
	X         float64  `json:"x"`
	Y         float64  `json:"y"`
	DrawOrder string   `json:"draworder"`
	Objects   []Object `json:"objects"`
	Color     string   `json:"color"`
}

type LevelData struct {
	Width            float64   `json:"width"`
	Height           float64   `json:"height"`
	Version          string    `json:"version"`
	Type             string    `json:"type"`
	CompressionLevel int       `json:"compressionlevel"`
	Infinite         bool      `json:"infinite"`
	NextLayerId      int       `json:"nextlayerid"`
	NextObjectId     int       `json:"nextobjectid"`
	Orientation      string    `json:"orientation"`
	RenderOrder      string    `json:"renderorder"`
	TiledVersion     string    `json:"tiledversion"`
	TileHeight       float64   `json:"tileheight"`
	TileWidth        float64   `json:"tilewidth"`
	Layers           []Layer   `json:"layers"`
	TileSets         []TileSet `json:"tilesets"`
}

type LevelUseCase struct {
}

func NewLevelUseCase() *LevelUseCase {
	return &LevelUseCase{}
}

type TileMap struct {
	Enemies []entities.Enemy
	Items   []entities.Item
}

func (l *LevelUseCase) NewTileMap(path string) *TileMap {
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
						tileMap.Enemies = append(tileMap.Enemies, entities.Enemy{
							ID:   &id,
							Name: enemy.Name,
							Position: &entities.Position{
								X: &enemy.X,
								Y: &enemy.Y,
							},
							Size: &entities.Size{
								Height: &enemy.Height,
								Width:  &enemy.Width,
							},
						})
					}
				}
				break
			}
		case "items":
			{
				for _, item := range layer.Objects {

					var propertyIndex *int

					for index, property := range item.Properties {
						if property.Name == "id" {
							propertyIndex = &index
						}
						if property.Name == "type" {
							propertyIndex = &index
						}
					}
					if propertyIndex != nil {
						id, err := strconv.Atoi(item.Properties[*propertyIndex].Value)

						if err != nil {
							log.Fatalf("Failed to convert id to int: %v", err)
						}
						tileMap.Items = append(tileMap.Items, entities.Item{
							Type: item.Properties[*propertyIndex].Value,
							Name: item.Name,
							ID:   &id,
							Position: &entities.Position{
								X: &item.X,
								Y: &item.Y,
							},
							Size: &entities.Size{
								Height: &item.Height,
								Width:  &item.Width,
							},
						})
					}
				}
			}
		}

	}
	return tileMap
}
