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
type Credentials interface {
	CreateCredentials(userId uuid.UUID, ps models.Credentials) (uuid.UUID, error)
	GetAllCredentials(userId uuid.UUID) ([]models.CredentialsList, error)
	GetCredentialsByID(userId, psId uuid.UUID) (models.CredentialsItemGet, error)
	DeleteCredentials(userId, psId uuid.UUID) error
	UpdateCredentials(userId, psId uuid.UUID, input models.CredentialsItemUpdate) error
}
type Repository struct {
	Authorization
	Credentials
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		Credentials:   NewCredentialsPostgres(db),
	}
}
