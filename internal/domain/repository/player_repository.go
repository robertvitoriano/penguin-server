package repository

import "github.com/robertvitoriano/penguin-server/internal/domain/entities"

type PlayerRepository interface {
	FindByID(id string) (*entities.Player, error)
	RemoveByID(id string) (*entities.Player, error)
	FindByUsername(id string) (*entities.Player, error)
	Save(player *entities.Player) error
	List() ([]*entities.Player, error)
}
