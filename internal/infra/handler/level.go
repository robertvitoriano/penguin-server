package handler

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
	"github.com/robertvitoriano/penguin-server/internal/domain/usecase"
	"github.com/robertvitoriano/penguin-server/internal/infra/repository/mysql"
	"gorm.io/gorm"
)

type LevelHandler struct {
	levelUseCase usecase.LevelUseCase
}

func NewLevelHandler() *LevelHandler {
	return &LevelHandler{
		levelUseCase: *usecase.NewLevelUseCase(),
	}
}

type LoadLevelResponse struct {
	Enemies []entities.Enemy `json:"enemies"`
	Items   []entities.Item  `json:"items"`
}

func (l *LevelHandler) LoadLevel(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var request struct {
		LevelName string `json:"level_name"`
	}

	enemyRepository := mysql.NewEnemiesRepository(db)
	itemsRepository := mysql.NewItemsRepository(db)

	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		log.Panic("level name is required")
	}

	tileMap := l.levelUseCase.NewTileMap(fmt.Sprintf("assets/maps/%v.json", request.LevelName))

	mapEntitiesWaitGroup := sync.WaitGroup{}

	mapEntitiesWaitGroup.Add(len(tileMap.Enemies) + len(tileMap.Items))

	responseItemsChan := make(chan entities.Item)
	responseEnemiesChan := make(chan entities.Enemy)

	for _, enemy := range tileMap.Enemies {
		go func(enemy entities.Enemy) {
			defer mapEntitiesWaitGroup.Done()

			query := mysql.EnemyQuery{
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

		}(enemy)
	}

	for _, item := range tileMap.Items {
		go func(item entities.Item) {
			defer mapEntitiesWaitGroup.Done()

			query := mysql.ItemQuery{
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

		}(item)

	}

	go func() {
		mapEntitiesWaitGroup.Wait()
		close(responseItemsChan)
		close(responseEnemiesChan)
	}()

	var responseItems []entities.Item
	var responseEnemies []entities.Enemy

	var responseWaitGroup sync.WaitGroup

	responseWaitGroup.Add(2)

	go func() {
		defer responseWaitGroup.Done()
		for item := range responseItemsChan {
			responseItems = append(responseItems, item)
		}
	}()

	go func() {
		defer responseWaitGroup.Done()
		for enemy := range responseEnemiesChan {
			responseEnemies = append(responseEnemies, enemy)
		}
	}()

	responseWaitGroup.Wait()

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
