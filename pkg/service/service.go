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
type PS interface {
	CreatePS(userId uuid.UUID, ps models.PS) (uuid.UUID, error)
	GetAllPS(userId uuid.UUID) ([]models.PSList, error)
	GetPSByID(userId, psId uuid.UUID) (models.PSItem, error)
	DeletePS(userId, psId uuid.UUID) error
	UpdatePS(userId, psId uuid.UUID, input models.PSItemUpdate) error
}
type Service struct {
	Authorization
	PS
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
		PS:            NewPSService(r.PS, cfg),
	}
}
