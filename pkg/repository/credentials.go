package repository

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"strings"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vpnvsk/p_s/internal/models"
)

type CredentialsPostgres struct {
	db *sqlx.DB
}

func NewCredentialsPostgres(db *sqlx.DB) *CredentialsPostgres {
	return &CredentialsPostgres{db: db}
}

func (r *CredentialsPostgres) CreateCredentials(userId uuid.UUID, cr models.Credentials) (uuid.UUID, error) {

	tx, err := r.db.Begin()
	if err != nil {
		return uuid.Nil, err
	}

	var psId uuid.UUID
	createPsItemQuery := fmt.Sprintf(
		"INSERT INTO %s (title, userlogin, password_hash, description) values ($1, $2, $3, $4) RETURNING id",
		credentialsTable)

	row := tx.QueryRow(createPsItemQuery, cr.Title, cr.Userlogin, cr.Password, cr.Description)
	err = row.Scan(&psId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	createUserPsItemsQuery := fmt.Sprintf(
		"INSERT INTO %s (user_id, list_id) values ($1, $2)",
		userCredentialsTable,
	)
	_, err = tx.Exec(createUserPsItemsQuery, userId, psId)
	if err != nil {
		tx.Rollback()
		return uuid.Nil, err
	}

	return psId, tx.Commit()
}
func (r *CredentialsPostgres) GetAllCredentials(userId uuid.UUID) ([]models.CredentialsList, error) {
	var list []models.CredentialsList
	query := fmt.Sprintf(
		"SELECT cr.id, cr.title, cr.description FROM %s cr JOIN %s ui ON cr.id=ui.list_id WHERE ui.user_id=$1",
		credentialsTable,
		userCredentialsTable,
	)
	err := r.db.Select(&list, query, userId)
	return list, err
}

func (r *CredentialsPostgres) GetCredentialsByID(userID, credentialsID uuid.UUID) (models.CredentialsItemGet, error) {
	var m models.CredentialsItemGet
	query := fmt.Sprintf(
		"SELECT cr.userlogin, cr.password_hash FROM %s cr JOIN %s ui on cr.id=ui.list_id WHERE ui.user_id=$1 and cr.id=$2",
		credentialsTable, userCredentialsTable)
	err := r.db.Get(&m, query, userID, credentialsID)
	return m, err
}
func (r *CredentialsPostgres) DeleteCredentials(userId, credentialsId uuid.UUID) error {
	query := fmt.Sprintf(
		"DELETE FROM %s cr USING %s ul WHERE cr.id = ul.list_id AND ul.user_id=$1 AND ul.list_id=$2",
		credentialsTable, userCredentialsTable)
	_, err := r.db.Exec(query, userId, credentialsId)
	return err
}
func (r *CredentialsPostgres) UpdateCredentials(userId, credentialsId uuid.UUID, input models.CredentialsItemUpdate) error {
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

	query := fmt.Sprintf("UPDATE %s cr SET %s FROM %s ul WHERE cr.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		credentialsTable, setQuery, userCredentialsTable, argId, argId+1)
	args = append(args, credentialsId, userId)

	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)

	_, err := r.db.Exec(query, args...)
	return err
}
