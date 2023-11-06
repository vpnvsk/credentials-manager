package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vpnvsk/p_s/internal/models"
)

type PSPostgres struct {
	db *sqlx.DB
}

func NewPSPostgres(db *sqlx.DB) *PSPostgres {
	return &PSPostgres{db: db}
}

func (r *PSPostgres) CreatePS(userId uuid.UUID, ps models.PS) (uuid.UUID, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	var psId uuid.UUID
	createPsItemQuery := fmt.Sprintf(
		"INSERT INTO %s (title, userlogin, password_hash, description) values ($1, $2, $3, $4) RETURNING id",
		psTable)

	row := tx.QueryRow(createPsItemQuery, ps.Title, ps.Userlogin, ps.Password, ps.Description)
	err = row.Scan(&psId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	createUserPsItemsQuery := fmt.Sprintf("INSERT INTO %s (user_id, list_id) values ($1, $2)", userPsTable)
	_, err = tx.Exec(createUserPsItemsQuery, userId, psId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return psId, tx.Commit()
}
func (r *PSPostgres) GetAllPS(userId uuid.UUID) ([]models.PSList, error) {
	var list []models.PSList
	query := fmt.Sprintf(
		"SELECT ps.id, ps.title, ps.description FROM %s ps JOIN %s ui ON ps.id=ui.list_id WHERE ui.user_id=$1",
		psTable, userPsTable)
	err := r.db.Select(&list, query, userId)
	return list, err
}

func (r *PSPostgres) GetPSByID(userID, psID uuid.UUID) (models.PSItem, error) {
	var m models.PSItem
	query := fmt.Sprintf(
		"SELECT ps.userlogin, ps.password_hash FROM %s ps JOIN %s ui on ps.id=ui.list_id WHERE ui.user_id=$1 and ps.id=$2",
		psTable, userPsTable)
	err := r.db.Get(&m, query, userID, psID)
	return m, err
}
func (r *PSPostgres) DeletePS(userId, psId uuid.UUID) error {
	query := fmt.Sprintf(
		"DELETE FROM %s ps USING %s ul WHERE ps.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		psTable, userPsTable)
	_, err := r.db.Exec(query, userId, psId)
	return err
}
func (r *PSPostgres) UpdatePS(userId, psId uuid.UUID, input models.PSItemUpdate) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}

	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	if input.Userlogin != nil {
		setValues = append(setValues, fmt.Sprintf("userlogin=$%d", argId))
		args = append(args, *input.Userlogin)
		argId++
	}
	if input.Password != nil {
		setValues = append(setValues, fmt.Sprintf("password_hash=$%d", argId))
		args = append(args, *input.Password)
		argId++
	}

	// title=$1
	// description=$1
	// title=$1, description=$2
	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE %s ps SET %s FROM %s ul WHERE ps.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		psTable, setQuery, userPsTable, argId, argId+1)
	args = append(args, psId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
