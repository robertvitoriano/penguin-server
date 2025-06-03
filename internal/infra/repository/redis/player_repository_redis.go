package redis

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"github.com/robertvitoriano/penguin-server/internal/domain/entities"
)

var Players = []*entities.Player{}

type PlayerRedisRepository struct {
	client *redis.Client
}

func NewPlayerRepository(client *redis.Client) *PlayerRedisRepository {
	return &PlayerRedisRepository{
		client: client,
	}
}

var ctx = context.Background()

func (p *PlayerRedisRepository) Save(newPlayer *entities.Player) error {
	data, err := json.Marshal(newPlayer)

	if err != nil {
		return err
	}
	return p.client.Set(ctx, "player:"+newPlayer.ID, data, 0).Err()
}

func (p *PlayerRedisRepository) List() ([]*entities.Player, error) {
	var players []*entities.Player
	iter := p.client.Scan(ctx, 0, "player:*", 0).Iterator()
	for iter.Next(ctx) {
		val, err := p.client.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}
		var player entities.Player
		if err := json.Unmarshal([]byte(val), &player); err != nil {
			continue
		}
		players = append(players, &player)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}
	return players, nil
}

func (p *PlayerRedisRepository) RemoveByID(id string) (*entities.Player, error) {
	playerToRemove, err := p.FindByID(id)

	if err != nil {
		return nil, err
	}

	if err := p.client.Del(ctx, "player:"+id).Err(); err != nil {
		return nil, err
	}

	return playerToRemove, nil
}

func (p *PlayerRedisRepository) FindByUsername(username string) (*entities.Player, error) {
	iter := p.client.Scan(ctx, 0, "player:*", 0).Iterator()

	for iter.Next(ctx) {
		val, err := p.client.Get(ctx, iter.Val()).Result()
		if err != nil {
			continue
		}

		var player entities.Player
		if err := json.Unmarshal([]byte(val), &player); err != nil {
			continue
		}
		if player.Username == username {
			return &player, nil
		}
	}
	return nil, fmt.Errorf("player not found")
}
func (p *PlayerRedisRepository) FindByID(id string) (*entities.Player, error) {
	playerRawData, err := p.client.Get(ctx, "player:"+id).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("player not found")
	} else if err != nil {
		return nil, err
	}
	var player entities.Player
	if err := json.Unmarshal([]byte(playerRawData), &player); err != nil {
		return nil, err
	}
	return &player, nil
}
