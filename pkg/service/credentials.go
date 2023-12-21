package service

import (
	"github.com/google/uuid"
	"github.com/vpnvsk/p_s/internal/lib"
	"github.com/vpnvsk/p_s/internal/models"
	"github.com/vpnvsk/p_s/pkg/repository"
)

type CredentialsService struct {
	repo repository.Credentials
	cfg  Config
}

func NewCredentialsService(repo repository.Credentials, cfg Config) *CredentialsService {
	return &CredentialsService{
		repo: repo,
		cfg:  cfg,
	}
}
func (s *CredentialsService) CreateCredentials(userId uuid.UUID, ps models.Credentials) (uuid.UUID, error) {
	password, err := lib.EncryptPassword(ps.Password, s.cfg.encryptKey)
	if err != nil {
		return uuid.Nil, err
	}
	ps.Password = password
	return s.repo.CreateCredentials(userId, ps)
}

func (s *CredentialsService) GetAllCredentials(userId uuid.UUID) ([]models.CredentialsList, error) {
	return s.repo.GetAllCredentials(userId)
}
func (s *CredentialsService) GetCredentialsByID(userId, psId uuid.UUID) (models.CredentialsItemGet, error) {
	model, err := s.repo.GetCredentialsByID(userId, psId)
	if err != nil {
		return model, err
	}
	pass := model.Password_Hash
	model.Password_Hash, err = lib.DecryptPassword(pass, s.cfg.encryptKey)
	return model, err
}
func (s *CredentialsService) DeleteCredentials(userId, psId uuid.UUID) error {
	return s.repo.DeleteCredentials(userId, psId)
}
func (s *CredentialsService) UpdateCredentials(userId, psId uuid.UUID, input models.CredentialsItemUpdate) error {
	if err := input.Validate(); err != nil {
		return err
	}
	if input.Password != nil {
		password, err := lib.EncryptPassword(*input.Password, s.cfg.encryptKey)
		if err != nil {
			return err
		}
		input.Password = &password
	}
	return s.repo.UpdateCredentials(userId, psId, input)
}
