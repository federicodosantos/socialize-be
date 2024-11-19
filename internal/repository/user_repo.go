package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
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
	UserLogin(ctx context.Context, email string, password string) (*model.User, error)
}

type UserRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepoItf {
	return &UserRepo{db: db}
}

// CreateUser implements UserRepoItf.
func (r *UserRepo) CreateUser(ctx context.Context, user *model.User) error {
	createdAtStr := util.ConvertTimeToString(user.CreatedAt)
	updatedAtStr := util.ConvertTimeToString(user.UpdatedAt)

	insertUserQuery := fmt.Sprintf(`INSERT INTO users(name, email, password, created_at, updated_at) VALUES ('%s', '%s', '%s', '%s', '%s')`,
	 user.Name, user.Email, user.Password, createdAtStr, updatedAtStr)

	exist, err := r.CheckEmailExist(ctx, user.Email)
	if err != nil {
		return err
	}

	if exist {
		return customError.ErrEmailExist
	}

	res, err := r.db.ExecContext(ctx, insertUserQuery)
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

	user.ID = lastInsertID

	return util.ErrRowsAffected(rows)
}

// GetUserById implements UserRepoItf.
func (r *UserRepo) GetUserById(ctx context.Context, userId int64) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE id = %d", userId)

	var user model.User

	err := r.db.QueryRowxContext(ctx, query).StructScan(&user)
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
	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s'", email)

	var user model.User

	err := r.db.QueryRowxContext(ctx, query).StructScan(&user)
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
	updatedAtStr := util.ConvertTimeToString(user.UpdatedAt)

	query := fmt.Sprintf(`UPDATE users 
	SET name = '%s', email = '%s', password = '%s', updated_at = '%s'
	WHERE id = %d'`, user.Name, user.Email, user.Password, updatedAtStr, user.ID)

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

	res, err := tx.ExecContext(ctx, query)
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
	updatedAtStr := util.ConvertTimeToString(user.UpdatedAt)

	query := fmt.Sprintf(`UPDATE users 
	SET photo = '%s', updated_at = '%s'
	WHERE id = %d`, user.Photo.String, updatedAtStr, user.ID)

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

	util.ErrRowsAffected(rows)

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

func (u *UserRepo) UserLogin(ctx context.Context, email string, password string) (*model.User, error) {
	query := fmt.Sprintf("SELECT * FROM users WHERE email = '%s' AND password = '%s'", email, password)

	fmt.Println(query)

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
