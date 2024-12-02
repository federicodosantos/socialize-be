package repository_test

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository"
	customerror "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func setup() (*sqlx.DB, sqlmock.Sqlmock, error) {
	mockDb, mock, err := sqlmock.New()
	if err != nil {
		return nil, nil, err
	}

	db := sqlx.NewDb(mockDb, "sqlmock")

	return db, mock, nil
}

func createUser() *model.User {
	now := time.Now()

	return &model.User{
		Name:  "Jamal",
		Email: "jamalunyu@gmail.com",
		Photo: sql.NullString{
			String: "https://nctojleerwtdfjlhafsq.supabase.co/storage/v1/object/public/photo_profile/pas%20foto%20windah.jpg",
			Valid:  true,
		},
		Password:  "rahasia123",
		CreatedAt: now,
		UpdatedAt: now,
	}
}

func TestCreateUser(t *testing.T) {
	type testCase struct {
		name          string
		setupMock     func(mock sqlmock.Sqlmock, user *model.User)
		user          *model.User
		expectedError error
	}

	testCases := []testCase{
		{
			name: "Success - CreateUser",
			setupMock: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
								SELECT COUNT(*) FROM users WHERE email = '%s'
								`, user.Email))).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectExec(regexp.QuoteMeta(fmt.Sprintf(`
								INSERT INTO users(name, email, password, created_at, updated_at)
  								VALUES ('%s', '%s', '%s', '%s', '%s')`,
					user.Name, user.Email, user.Password, util.ConvertTimeToString(user.CreatedAt), util.ConvertTimeToString(user.UpdatedAt)))).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			user:          createUser(),
			expectedError: nil,
		},
		{
			name: "Error email already exists - CreateUser",
			setupMock: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf(`
								SELECT COUNT(*) FROM users WHERE email = '%s'
								`, user.Email))).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(1))
			},
			user:          createUser(),
			expectedError: customerror.ErrEmailExist,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := setup()
			if err != nil {
				t.Fatalf("error to create sql mock and db due to %s", err.Error())
			}
			defer db.Close()

			tc.setupMock(mock, tc.user)

			u := repository.NewUserRepo(db)

			err = u.CreateUser(context.Background(), tc.user)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
			} else {
				assert.NoError(t, err)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}

func TestGetUserByEmail(t *testing.T) {
	type testCase struct {
		name          string
		email         string
		setupMock     func(mock sqlmock.Sqlmock, email string)
		expectedUser  *model.User
		expectedError error
	}

	user := createUser()

	testCases := []testCase{
		{
			name:  "Success get user by email",
			email: user.Email,
			setupMock: func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}).
						AddRow(user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt))
			},
			expectedUser:  user,
			expectedError: nil,
		},
		{
			name:  "Error user not found",
			email: user.Email,
			setupMock: func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta(fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email))).
					WillReturnRows(sqlmock.NewRows([]string{"id", "name", "email", "password", "created_at", "updated_at"}))
			},
			expectedUser:  nil,
			expectedError: customerror.ErrEmailNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := setup()
			if err != nil {
				t.Fatalf("error to create sql mock and db due to %s", err.Error())
			}
			defer db.Close()

			tc.setupMock(mock, tc.email)

			u := repository.NewUserRepo(db)

			user, err := u.GetUserByEmail(context.Background(), tc.email)

			if tc.expectedError != nil {
				assert.Error(t, err)
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tc.expectedUser.ID, user.ID)
				assert.Equal(t, tc.expectedUser.Name, user.Name)
				assert.Equal(t, tc.expectedUser.Email, user.Email)
				assert.Equal(t, tc.expectedUser.Password, user.Password)
				assert.WithinDuration(t, time.Now().Local(), user.CreatedAt, time.Second)
				assert.WithinDuration(t, time.Now().Local(), user.UpdatedAt, time.Second)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("there were unfulfilled expectations: %s", err)
			}
		})
	}
}
