package service

import (
	"github.com/google/uuid"
	"github.com/vpnvsk/p_s/internal/models"
	"github.com/vpnvsk/p_s/pkg/repository"
	"time"
)

type Authorization interface {
	CreateUser(u models.User) (uuid.UUID, error)
	GenerateToken(username, password string) (string, error)
	ParseToken(token string) (uuid.UUID, error)
}
type Credentials interface {
	CreateCredentials(userId uuid.UUID, ps models.Credentials) (uuid.UUID, error)
	GetAllCredentials(userId uuid.UUID) ([]models.CredentialsList, error)
	GetCredentialsByID(userId, psId uuid.UUID) (models.CredentialsItemGet, error)
	DeleteCredentials(userId, psId uuid.UUID) error
	UpdateCredentials(userId, psId uuid.UUID, input models.CredentialsItemUpdate) error
}
type Service struct {
	Authorization
	Credentials
}
type ServiceConfig struct {
	salt       string
	signingKey string
	tokenTTL   time.Duration
	encryptKey []byte
}

func (s *ServiceConfig) SetFields(salt, signingKey, encryptKey string, tokenTTL time.Duration) {
	s.salt = salt
	s.signingKey = signingKey
	s.tokenTTL = tokenTTL
	s.encryptKey = []byte(encryptKey)
}

func NewService(r *repository.Repository, cfg ServiceConfig) *Service {
	return &Service{
		Authorization: NewAuthService(r.Authorization, cfg),
		Credentials:   NewCredentialsService(r.Credentials, cfg),
	}
}
