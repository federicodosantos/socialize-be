package repository_test

import (
	"context"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository"
	customerror "github.com/federicodosantos/socialize/pkg/custom-error"
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
		ID:        "1",
		Name:      "Jamal",
		Email:     "jamalunyu@gmail.com",
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM users WHERE email = ?")).
					WithArgs(user.Email).
					WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(0))

				mock.ExpectExec(`INSERT INTO users\(id, name, email, password, created_at, updated_at\) VALUES \(\?, \?, \?, \?, \?, \?\)`).
					WithArgs("1", "Jamal", "jamalunyu@gmail.com", "rahasia123", sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))

			},
			user:          createUser(),
			expectedError: nil,
		},
		{
			name: "Error email already exists - CreateUser",
			setupMock: func(mock sqlmock.Sqlmock, user *model.User) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT COUNT(*) FROM users WHERE email = ?")).
					WithArgs(user.Email).
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
				t.Fatalf("Error creating sql mock and db: %s", err)
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
				t.Errorf("Unfulfilled expectations: %s", err)
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
			name:  "Success - GetUserByEmail",
			email: user.Email,
			setupMock: func(mock sqlmock.Sqlmock, email string) {
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = ?")).
					WithArgs(email).
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
				mock.ExpectQuery(regexp.QuoteMeta("SELECT * FROM users WHERE email = ?")).
					WithArgs(email).
					WillReturnRows(sqlmock.NewRows(nil))
			},
			expectedUser:  nil,
			expectedError: customerror.ErrEmailNotFound,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			db, mock, err := setup()
			if err != nil {
				t.Fatalf("Error creating sql mock and db: %s", err)
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
				assert.Equal(t, tc.expectedUser, user)
			}

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("Unfulfilled expectations: %s", err)
			}
		})
	}
}
