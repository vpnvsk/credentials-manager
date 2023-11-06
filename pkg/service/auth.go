package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
	"github.com/vpnvsk/p_s/internal/lib"
	"github.com/vpnvsk/p_s/internal/models"
	"github.com/vpnvsk/p_s/pkg/repository"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
}

type AuthService struct {
	repo   repository.Authorization
	config ServiceConfig
}

func NewAuthService(repo repository.Authorization, cfg ServiceConfig) *AuthService {
	return &AuthService{
		repo:   repo,
		config: cfg,
	}
}

func (s *AuthService) CreateUser(user models.User) (uuid.UUID, error) {
	password := lib.GeneratePasswordHash(user.Password, s.config.salt)
	user.Password = password
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateToken(username, password string) (string, error) {
	const op = "service.generateToken"
	user, err := s.repo.GetUser(username, lib.GeneratePasswordHash(password, s.config.salt))
	if err != nil {
		return "", fmt.Errorf("%s:%w", op, err)
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{ExpiresAt: time.Now().Add(s.config.tokenTTL).Unix(),
			IssuedAt: time.Now().Unix()},
		user.Id,
	})
	return token.SignedString([]byte(s.config.signingKey))

}

func (s *AuthService) ParseToken(accessToken string) (uuid.UUID, error) {
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(s.config.signingKey), nil
	})
	if err != nil {
		return uuid.Nil, err
	}
	claims, ok := token.Claims.(*tokenClaims)
	if !ok {
		return uuid.Nil, errors.New("invalid type of token")
	}
	return claims.UserId, nil
}
