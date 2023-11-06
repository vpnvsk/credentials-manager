package models

import (
	"errors"
	"github.com/google/uuid"
)

type PS struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Userlogin   string    `json:"userlogin" db:"userlogin" binding:"required"`
	Password    string    `json:"password" binding:"required"`
	Description string    `json:"description" db:"description"`
}
type PSList struct {
	Id          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title" binding:"required"`
	Description string    `json:"description" db:"description"`
}
type PSItem struct {
	Userlogin     string `json:"userlogin" db:"userlogin" binding:"required"`
	Password_Hash string `json:"password" binding:"required"`
}
type PSItemUpdate struct {
	Title       *string `json:"title"`
	Userlogin   *string `json:"userlogin"`
	Password    *string `json:"password"`
	Description *string `json:"description"`
}
type UserPS struct {
	Id     uuid.UUID
	UserId uuid.UUID
	PSId   uuid.UUID
}

func (i PSItemUpdate) Validate() error {
	if i.Title == nil && i.Description == nil && i.Userlogin == nil && i.Password == nil {
		return errors.New("update structure has no values")
	}
	return nil
}
