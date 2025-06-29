package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/robertvitoriano/penguin-server/internal/infra/database"
	"github.com/robertvitoriano/penguin-server/internal/infra/handler"
	"github.com/robertvitoriano/penguin-server/internal/infra/middlewares"
	"github.com/robertvitoriano/penguin-server/internal/infra/repository/mysql"
	"github.com/robertvitoriano/penguin-server/internal/infra/repository/redis"
	"github.com/rs/cors"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	mysqlDatabase := database.NewMysqlDabase()

	mysqlDatabase.Dsn = fmt.Sprintf(
		"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		os.Getenv("MYSQL_USER"),
		os.Getenv("MYSQL_PASSWORD"),
		os.Getenv("MYSQL_HOST"),
		os.Getenv("MYSQL_PORT"),
		os.Getenv("MYSQL_DATABASE"))

	mysqlDatabase.Db = &gorm.DB{}
	mysqlDatabase.DbType = "mysql"
	mysqlDatabase.Connect()

	redisDatabase := database.NewRedisDatabase(fmt.Sprintf("%v:6379", os.Getenv("REDIS_HOST")), "", 0)
	redisClient := redisDatabase.Connect()
	router := mux.NewRouter()
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	playerPersistencyRepository := mysql.NewPlayerRepository(mysqlDatabase.Db)
	playerLiveDataRepository := redis.NewPlayerRepository(redisClient)

	playerHandler := handler.NewPlayerHandler(playerPersistencyRepository, playerLiveDataRepository)
	levelHandler := handler.NewLevelHandler()
	ws := handler.NewWebsocket(playerLiveDataRepository, playerPersistencyRepository)

	router.HandleFunc("/players/{id}", playerHandler.GetPlayer).Methods("GET")
	router.HandleFunc("/players", func(w http.ResponseWriter, r *http.Request) {
		playerHandler.CreatePlayer(w, r, ws, mysqlDatabase.Db)
	}).Methods("POST")
	router.HandleFunc("/players", playerHandler.GetPlayers).Methods("GET")
	router.HandleFunc("/load-level", func(w http.ResponseWriter, r *http.Request) {
		levelHandler.LoadLevel(w, r, mysqlDatabase.Db)
	}).Methods("POST")

	router.HandleFunc("/ws", ws.ServeWebsocket).Methods("GET")

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"https://penguim-adventure.robertvitoriano.com", "http://localhost:8000", fmt.Sprintf("http://%v:8000", os.Getenv("COMPUTER_IP"))},

		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Accept-Encoding", "X-CSRF-Token", "Authorization"},
		AllowCredentials: true,
		Debug:            true,
	})

	handler := c.Handler(router)

	rateLimiter := middlewares.NewRateLimiter(30, 1*time.Minute, context.Background(), *redisClient)

	handler = middlewares.RateLimiterMiddleware(handler, *rateLimiter)

	fmt.Println("Server running on port 7777...")
	log.Fatal(http.ListenAndServe(":7777", handler))

	removedPlayer := make(chan string)
	go func(removedPlayer chan string) {
		for {

			liveDataPlayers, _ := playerLiveDataRepository.List()

			for _, player := range liveDataPlayers {
				removedPlayer <- player.ID
				if player.LastTimeOnline == nil || time.Since(*player.LastTimeOnline) >= 7*24*time.Second {
					playerLiveDataRepository.RemoveByID(player.ID)
				}
			}
			time.Sleep(time.Minute)
		}
	}(removedPlayer)
	for removedPlayerId := range removedPlayer {
		fmt.Printf("PLAYER REMOVED %v", removedPlayerId)

	}

}
