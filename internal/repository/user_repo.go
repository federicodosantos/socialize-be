package repository

import (
	"context"
	"database/sql"
	"errors"
	"log"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository/query"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/util"
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
func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	exist, err := r.CheckEmailExist(ctx, user.Email)
	if err != nil {
		return err
	}

	if exist {
		return customError.ErrEmailExist
	}

	res, err := r.db.ExecContext(ctx, query.InsertUserQuery,
		 user.ID, user.Name, user.Email, user.Password, user.CreatedAt, user.UpdatedAt)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customError.ErrRowsAffected
	}

	return util.ErrRowsAffected(rows)
}

// GetUserById implements UserRepoItf.
func (r *UserRepo) GetUserById(ctx context.Context, userId string) (*model.User, error) {
	var user model.User

	err := r.db.QueryRowxContext(ctx, query.GetUserByIdQuery, userId).StructScan(&user)
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
	var user model.User

	err := r.db.QueryRowxContext(ctx, query.GetUserByEmailQuery, email).StructScan(&user)
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
	tx, err := r.db.Beginx()
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

	res, err := tx.ExecContext(ctx, query.UpdateUserDataQuery,
		 user.Name, user.Email, user.Password, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	util.ErrRowsAffected(rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// UpdateUserData implements UserRepoItf.
func (u *UserRepo) UpdateUserPhoto(ctx context.Context, user *model.User) error {
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

	res, err := tx.ExecContext(ctx, query.UpdateUserPhotoQuery,
	user.Photo, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}

	util.ErrRowsAffected(rows)

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

// CheckEmailExist implements UserRepoItf.
func (u *UserRepo) CheckEmailExist(ctx context.Context, email string) (bool, error) {
	var count int

	err := u.db.QueryRowxContext(ctx, query.CheckEmailExistQuery, email).Scan(&count)
	if err != nil {
		return false, err
	}

	return count > 0, nil
}