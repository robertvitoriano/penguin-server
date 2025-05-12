package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/robertvitoriano/penguin-server/internal/models"
	"github.com/robertvitoriano/penguin-server/internal/repositories/mysqlrepositories"
	"github.com/robertvitoriano/penguin-server/internal/tiled"
	"gorm.io/gorm"
)

type LoadLevelResponse struct {
	Enemies []models.Enemy
	Items   []models.Item
}

func LoadLevel(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var request struct {
		LevelName string `json:"level_name"`
	}

	enemyRepository := mysqlrepositories.NewEnemiesRepository(db)
	itemsRepository := mysqlrepositories.NewItemsRepository(db)

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		log.Panic("level name is required")
	}

	tileMap := tiled.NewTileMap(fmt.Sprintf("assets/maps/%v.json", request.LevelName))

	mapEntitiesWaitGroup := sync.WaitGroup{}

	mapEntitiesWaitGroup.Add(len(tileMap.Enemies) + len(tileMap.Items))

	responseItemsChan := make(chan models.Item, len(tileMap.Items))
	responseEnemiesChan := make(chan models.Enemy, len(tileMap.Enemies))

	for _, enemy := range tileMap.Enemies {
		go func(enemy models.Enemy) {
			defer mapEntitiesWaitGroup.Done()

			query := mysqlrepositories.EnemyQuery{
				ID: *enemy.ID,
			}
			enemyFound, err := enemyRepository.FindEnemy(query)

			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			if enemyFound != nil {
				responseEnemiesChan <- *enemyFound
				return
			}
			err = enemyRepository.CreateEnemy(&enemy)

			if err != nil {
				log.Printf("Error: %v", err)
				return
			}

			responseEnemiesChan <- enemy

			fmt.Printf("Enemy %v added!\n", *enemy.ID)

		}(enemy)
	}

	for _, item := range tileMap.Items {
		go func(item models.Item) {
			defer mapEntitiesWaitGroup.Done()

			query := mysqlrepositories.ItemQuery{
				ID: *item.ID,
			}

			itemFound, err := itemsRepository.FindItem(query)

			if err != nil {
				log.Printf("Error: %v", err)
				return
			}
			if itemFound != nil {
				responseItemsChan <- *itemFound
				return
			}

			err = itemsRepository.CreateItem(&item)

			if err != nil {
				log.Printf("Error: %v", err)
				return
			}

			responseItemsChan <- item

			fmt.Printf("Item %v created!\n", *item.ID)
		}(item)

	}

	go func() {
		mapEntitiesWaitGroup.Wait()
		close(responseItemsChan)
		close(responseEnemiesChan)
	}()

	var responseItems []models.Item
	var responseEnemies []models.Enemy

	for item := range responseItemsChan {
		responseItems = append(responseItems, item)
	}

	for enemy := range responseEnemiesChan {
		responseEnemies = append(responseEnemies, enemy)
	}

	jsonResponse, err := json.Marshal(LoadLevelResponse{
		Enemies: responseEnemies,
		Items:   responseItems,
	})

	if err != nil {
		http.Error(w, "Error parsing load level response to JSON", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonResponse)

}
