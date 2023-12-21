package service

import (
	"errors"
	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId uuid.UUID `json:"user_id"`
	AppId  int32     `json:"app_id"`
}

type AuthService struct {
	config Config
}

func NewAuthService(cfg Config) *AuthService {
	return &AuthService{
		config: cfg,
	}
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

	if claims.AppId != s.config.appId {
		return uuid.Nil, errors.New("invalid app id")
	}
	return claims.UserId, nil
}
