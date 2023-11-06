package service

import (
	"github.com/google/uuid"
	"github.com/vpnvsk/p_s/internal/lib"
	"github.com/vpnvsk/p_s/internal/models"
	"github.com/vpnvsk/p_s/pkg/repository"
)

type PSService struct {
	repo repository.PS
	cfg  ServiceConfig
}

func NewPSService(repo repository.PS, cfg ServiceConfig) *PSService {
	return &PSService{
		repo: repo,
		cfg:  cfg,
	}
}
func (s *PSService) CreatePS(userId uuid.UUID, ps models.PS) (uuid.UUID, error) {
	password, err := lib.EncryptPassword(ps.Password, s.cfg.encryptKey)
	if err != nil {
		return uuid.Nil, err
	}
	ps.Password = password
	return s.repo.CreatePS(userId, ps)
}

func (s *PSService) GetAllPS(userId uuid.UUID) ([]models.PSList, error) {
	return s.repo.GetAllPS(userId)
}
func (s *PSService) GetPSByID(userId, psId uuid.UUID) (models.PSItem, error) {
	model, err := s.repo.GetPSByID(userId, psId)
	if err != nil {
		return model, err
	}
	pass := model.Password_Hash
	model.Password_Hash, err = lib.DecryptPassword(pass, s.cfg.encryptKey)
	return model, err
}
func (s *PSService) DeletePS(userId, psId uuid.UUID) error {
	return s.repo.DeletePS(userId, psId)
}
func (s *PSService) UpdatePS(userId, psId uuid.UUID, input models.PSItemUpdate) error {
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
	return s.repo.UpdatePS(userId, psId, input)
}
