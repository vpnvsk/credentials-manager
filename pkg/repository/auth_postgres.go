package repository

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vpnvsk/p_s/internal/models"
)

type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}

}

func (r *AuthPostgres) CreateUser(user models.User) (uuid.UUID, error) {
	var id uuid.UUID
	query := fmt.Sprintf("INSERT INTO %s (username, password_hash) values ($1, $2) RETURNING id",
		userTable)
	row := r.db.QueryRow(query, user.Username, user.Password)
	if err := row.Scan(&id); err != nil {
		return uuid.Nil, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (models.User, error) {
	const op = "repository.GetUser"
	var user models.User
	query := fmt.Sprintf("SELECT id, username FROM %s WHERE username=$1 AND password_hash=$2",
		userTable)
	err := r.db.Get(&user, query, username, password)
	if err != nil {
		e := fmt.Errorf("%s:%w", op, err)
		return user, e
	}
	return user, nil
}
