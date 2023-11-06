package models

import "github.com/google/uuid"

type User struct {
	Id uuid.UUID `json:"-" db:"id"`

	Username       string `json:"username" binding:"required"`
	Password       string `json:"password" binding:"required"`
}
