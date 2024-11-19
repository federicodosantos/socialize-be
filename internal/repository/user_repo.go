package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/federicodosantos/socialize/internal/model"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/jmoiron/sqlx"
)

type UserRepoItf interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	GetUserById(ctx context.Context, userId int64) (*model.User, error)
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
func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	insertUserQuery := `
		INSERT INTO users (name, email, password, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?)
	`

	exist, err := r.CheckEmailExist(ctx, user.Email)
	if err != nil {
		return err
	}

	if exist {
		return customError.ErrEmailExist
	}

	res, err := r.db.ExecContext(ctx, insertUserQuery, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	lastInsertID, err := res.LastInsertId()
	if err != nil {
		return customError.ErrLastInsertId
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customError.ErrRowsAffected
	}

	if err := util.ErrRowsAffected(rows); err != nil {
		return err
	}

	user.ID = lastInsertID
	return nil
}

// GetUserById implements UserRepoItf.
func (r *UserRepo) GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	query := `SELECT * FROM users WHERE id = ?`

	var user model.User
	err := r.db.QueryRowxContext(ctx, query, userId).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// GetUserByEmail implements UserRepoItf.
func (r *UserRepo) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT * FROM users WHERE email = ?`

	var user model.User
	err := r.db.QueryRowxContext(ctx, query, email).StructScan(&user)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, customError.ErrEmailNotFound
		}
		return nil, err
	}

	return &user, nil
}

// UpdateUserData implements UserRepoItf.
func (r *UserRepo) UpdateUserData(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET name = ?, email = ?, password = ?, updated_at = ?
		WHERE id = ?
	`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("cannot rollback tx: %s", rollbackErr)
			}
		}
	}()

	res, err := tx.ExecContext(ctx, query, user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if err := util.ErrRowsAffected(rows); err != nil {
		return err
	}

	return tx.Commit()
}

// UpdateUserPhoto implements UserRepoItf.
func (r *UserRepo) UpdateUserPhoto(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET photo = ?, updated_at = ?
		WHERE id = ?
	`

	tx, err := r.db.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				log.Printf("cannot rollback tx: %s", rollbackErr)
			}
		}
	}()

	res, err := tx.ExecContext(ctx, query, user.Photo.String, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if err := util.ErrRowsAffected(rows); err != nil {
		return err
	}

	return tx.Commit()
}

// CheckEmailExist implements UserRepoItf.
func (r *UserRepo) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	query := `SELECT COUNT(*) FROM users WHERE email = ?`

	var count int
	err := r.db.QueryRowxContext(ctx, query, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}
