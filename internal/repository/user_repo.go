package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/federicodosantos/socialize/internal/model"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/jmoiron/sqlx"
)

type UserRepoItf interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserById(ctx context.Context, userId string) (*model.User, error)
	CheckEmailExist(ctx context.Context, email string) (bool, error)
	UpdateUserData(ctx context.Context, user *model.User) error
	UpdateUserPhoto(ctx context.Context, user *model.User) error
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepoItf {
	return &UserRepo{db: db}
}

// CreateUser implements UserRepoItf.
func (u *UserRepo) CreateUser(ctx context.Context,user *model.User) error {
	createdAtStr := convertTimeToString(user.CreatedAt)
	updatedAtStr := convertTimeToString(user.UpdatedAt)

	insertUserQuery := fmt.Sprintf(`INSERT INTO users(id, name, email, password, created_at, updated_at)
        VALUES('%s', '%s', '%s', '%s', '%s', '%s')`, user.ID, user.Name, user.Email, user.Password, createdAtStr, updatedAtStr)
	
	exist, err := u.CheckEmailExist(ctx, user.Email)
	if err != nil {
		return err
	}

	if exist {
		return customError.ErrEmailExist
	}

	res, err := u.db.ExecContext(ctx, insertUserQuery)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customError.ErrRowsAffected
	}

	errRowsAffected(rows)

	return nil
}

// GetUserById implements UserRepoItf.
func (u *UserRepo) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	query :=  fmt.Sprintf("SELECT * FROM users WHERE id = '%s'", userId)

	var user model.User

	err := u.db.QueryRowxContext(ctx, query).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail implements UserRepoItf.
func (u *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)
	
	var user model.User

	err := u.db.QueryRowxContext(ctx, query).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrEmailNotFound
		}
		return nil, err
	}

	return &user, nil
}
// UpdateUserData implements UserRepoItf.
func (u *UserRepo) UpdateUserData(ctx context.Context, user *model.User) error {
	updatedAtStr := convertTimeToString(user.UpdatedAt)

	query := fmt.Sprintf(`UPDATE users 
	SET name = '%s', email = '%s', password = '%s', updated_at = '%s'
	WHERE id = '%s'`, user.Name, user.Email, user.Password, updatedAtStr, user.ID)

	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Printf("cannot rollback tx: %s", err)
				return
			}
		}
	}()

	res, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	errRowsAffected(rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserData implements UserRepoItf.
func (u *UserRepo) UpdateUserPhoto(ctx context.Context, user *model.User) error {
	updatedAtStr := convertTimeToString(user.UpdatedAt)

	query := fmt.Sprintf(`UPDATE users 
	SET photo = '%s', updated_at = '%s'
	WHERE id = '%s'`, user.Photo.String, updatedAtStr, user.ID)

	tx, err := u.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if err = tx.Rollback(); err != nil {
				log.Printf("cannot rollback tx: %s", err)
				return
			}
		}
	}()

	res, err := tx.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	errRowsAffected(rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CheckEmailExist implements UserRepoItf.
func (u *UserRepo) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	query := fmt.Sprintf(`SELECT COUNT(*) FROM users WHERE email = '%s'`, email)

	var count int

	err := u.db.QueryRowxContext(ctx, query).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func errRowsAffected(rows int64) error {
	if rows != 1 {
		return fmt.Errorf("error : %w, got %d rows affected",customError.ErrRowsAffected, rows)
	}

	return nil
}

func convertTimeToString(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}