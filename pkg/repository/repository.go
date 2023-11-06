package repository

import (
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vpnvsk/p_s/internal/models"
)

type Authorization interface {
	CreateUser(user models.User) (uuid.UUID, error)
	GetUser(username, password string) (models.User, error)
}
type PS interface {
	CreatePS(userId uuid.UUID, ps models.PS) (uuid.UUID, error)
	GetAllPS(userId uuid.UUID) ([]models.PSList, error)
	GetPSByID(userId, psId uuid.UUID) (models.PSItem, error)
	DeletePS(userId, psId uuid.UUID) error
	UpdatePS(userId, psId uuid.UUID, input models.PSItemUpdate) error
}
type Repository struct {
	Authorization
	PS
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		PS:            NewPSPostgres(db),
	}
}
