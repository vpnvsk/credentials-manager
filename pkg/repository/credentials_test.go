package repository

import (
	"errors"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/vpnvsk/p_s/internal/models"
	sqlmock "github.com/zhashkevych/go-sqlxmock"
	"testing"
)

func TestPSPostgres_CreatePS(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewCredentialsPostgres(db)

	knownUUID, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	if err != nil {
		t.Fatalf("error generating known UUID: %v", err)
	}

	type args struct {
		userId uuid.UUID
		ps     models.Credentials
	}
	type mockBehavior func(args args, id uuid.UUID)

	tests := []struct {
		name    string
		mock    mockBehavior
		input   args
		want    uuid.UUID
		wantErr bool
	}{
		{
			name: "Ok",
			input: args{
				userId: knownUUID,
				ps: models.Credentials{
					Title:       "test title",
					Userlogin:   "test login",
					Password:    "test_password",
					Description: "test description",
				},
			},
			want: knownUUID,
			mock: func(args args, id uuid.UUID) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO ps_item").
					WithArgs(args.ps.Title, args.ps.Userlogin, args.ps.Password, args.ps.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_item").WithArgs(args.userId, id).
					WillReturnResult(sqlmock.NewResult(1, 1))

				mock.ExpectCommit()
			},
		},
		{
			name: "Empty Fields",
			input: args{
				userId: knownUUID,
				ps: models.Credentials{
					Title:       "",
					Userlogin:   "test login",
					Password:    "test_password",
					Description: "description",
				},
			},
			mock: func(args args, id uuid.UUID) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id).RowError(0, errors.New("insert error"))
				mock.ExpectQuery("INSERT INTO ps_item").
					WithArgs(args.ps.Title, args.ps.Userlogin, args.ps.Password, args.ps.Description).WillReturnRows(rows)

				mock.ExpectRollback()
			},
			wantErr: true,
		},
		{
			name: "Failed 2nd Insert",
			input: args{
				userId: knownUUID,
				ps: models.Credentials{
					Title:       "test title",
					Userlogin:   "test login",
					Password:    "test_password",
					Description: "test description",
				},
			},
			mock: func(args args, id uuid.UUID) {
				mock.ExpectBegin()

				rows := sqlmock.NewRows([]string{"id"}).AddRow(id)
				mock.ExpectQuery("INSERT INTO ps_item").
					WithArgs(args.ps.Title, args.ps.Userlogin, args.ps.Password, args.ps.Description).WillReturnRows(rows)

				mock.ExpectExec("INSERT INTO users_item").WithArgs(args.userId, id).
					WillReturnError(errors.New("insert error"))

				mock.ExpectRollback()
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock(tt.input, tt.want)

			got, err := r.CreateCredentials(tt.input.userId, tt.input.ps)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestPSPostgres_GetAllPS(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	r := NewCredentialsPostgres(db)

	id1, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	if err != nil {
		t.Fatalf("error generating known UUID: %v", err)
	}
	id2, err := uuid.Parse("00000000-0000-0000-0000-000000000002")
	if err != nil {
		t.Fatalf("error generating known UUID: %v", err)
	}

	id3, err := uuid.Parse("00000000-0000-0000-0000-000000000003")

	if err != nil {
		t.Fatalf("error generating known UUID: %v", err)
	}

	type args struct {
		userId uuid.UUID
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    []models.CredentialsList
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"}).
					AddRow(id1, "title1", "description1").
					AddRow(id2, "title2", "description2").
					AddRow(id3, "title3", "description3")

				mock.ExpectQuery("SELECT (.+) FROM ps_item ps JOIN users_item ui ON (.+) WHERE (.+)").
					WithArgs(id1).WillReturnRows(rows)

			},
			input: args{
				userId: id1,
			},
			want: []models.CredentialsList{
				{id1, "title1", "description1"},
				{id2, "title2", "description2"},
				{id3, "title3", "description3"},
			},
		},
		{
			name: "No records",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "title", "description"})
				mock.ExpectQuery("SELECT (.+) FROM ps_item ps JOIN users_item ui ON (.+) WHERE (.+)").
					WithArgs(id1).WillReturnRows(rows)

			},
			input: args{
				userId: id1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			got, err := r.GetAllCredentials(tt.input.userId)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}

}
func TestPSPostgres_GetPSByID(t *testing.T) {
	db, mock, err := sqlmock.Newx()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()
	id1, err := uuid.Parse("00000000-0000-0000-0000-000000000001")
	iderr, err := uuid.Parse("00000000-0000-0000-0000-000000000002")

	if err != nil {
		t.Fatalf("error generating known UUID: %v", err)
	}

	r := NewCredentialsPostgres(db)
	type args struct {
		userID uuid.UUID
		psID   uuid.UUID
	}
	tests := []struct {
		name    string
		mock    func()
		input   args
		want    models.CredentialsItemGet
		wantErr bool
	}{
		{
			name: "Ok",
			mock: func() {
				rows := sqlmock.NewRows([]string{"cduserlogin", "password_hash"}).
					AddRow("title1", "description1")
				mock.ExpectQuery("SELECT (.+) FROM ps_item ps JOIN users_item ui on (.+) WHERE (.+)").
					WithArgs(id1, id1).WillReturnRows(rows)
			},
			input: args{
				userID: id1,
				psID:   id1,
			},
			want: models.CredentialsItemGet{"title1", "description1"},
		},
		{
			name: "Not found",
			mock: func() {
				rows := sqlmock.NewRows([]string{"id", "userlogin", "password_hash"}).
					AddRow(id1, "title1", "description1")
				mock.ExpectQuery("SELECT (.+) FROM ps_item ps JOIN users_item ui on (.+) WHERE (.+)").
					WithArgs(id1, iderr).WillReturnRows(rows)
			},
			input: args{
				userID: id1,
				psID:   iderr,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.mock()
			get, err := r.GetCredentialsByID(tt.input.userID, tt.input.psID)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, get)
			}
		})
		assert.NoError(t, mock.ExpectationsWereMet())

	}
}
