package service

import (
	"github.com/google/uuid"
	"github.com/vpnvsk/p_s/internal/models"
	"github.com/vpnvsk/p_s/pkg/repository"
	"time"
)

type Authorization interface {
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
type Config struct {
	salt       string
	signingKey string
	tokenTTL   time.Duration
	encryptKey []byte
	appId      int32
}

func (s *Config) SetFields(salt, signingKey, encryptKey string, tokenTTL time.Duration, appId int32) {
	s.salt = salt
	s.signingKey = signingKey
	s.tokenTTL = tokenTTL
	s.encryptKey = []byte(encryptKey)
	s.appId = appId
}

func NewService(r *repository.Repository, cfg Config) *Service {
	return &Service{
		Authorization: NewAuthService(cfg),
		Credentials:   NewCredentialsService(r.Credentials, cfg),
	}
}
