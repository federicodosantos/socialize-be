package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/jwt"
	"github.com/federicodosantos/socialize/pkg/md5"
	"github.com/federicodosantos/socialize/pkg/supabase"
)

type UserUsecaseItf interface {
	Register(ctx context.Context, req *model.UserRegister) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.UserLogin) (string, error)
	GetUserById(ctx context.Context, userId int) (*model.UserResponse, error)
	UpdateUserData(ctx context.Context, req *model.UserUpdateData, userId int) (*model.UserResponse, error)
	UpdateUserPhoto(ctx context.Context, req *model.UserUpdatePhoto, userId int) (*model.UserResponse, error)
}

type UserUsecase struct {
	userRepo repository.UserRepoItf
	supabase supabase.SupabaseStorageItf
	jwt      jwt.JWTItf
}

func NewUserUsecase(userRepo repository.UserRepoItf,
	jwt jwt.JWTItf, supabase supabase.SupabaseStorageItf) UserUsecaseItf {
	return &UserUsecase{
		userRepo: userRepo,
		jwt:      jwt,
		supabase: supabase}
}

// Register implements UserUCItf.
func (u *UserUsecase) Register(ctx context.Context, req *model.UserRegister) (*model.UserResponse, error) {
	hashedPassword := md5.HashWithMd5(req.Password)

	createdUser := &model.User{
		Name:      req.Name,
		Email:     req.Email,
		Password:  hashedPassword,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := u.userRepo.CreateUser(ctx, createdUser)
	if err != nil {
		if errors.Is(err, customError.ErrEmailExist) {
			return nil, fmt.Errorf("email : %s already exists: %w", createdUser.Email, err)
		}

		return nil, err
	}

	return convertToUserRespone(createdUser), nil
}

// Login implements UserUCItf.
func (u *UserUsecase) Login(ctx context.Context, req *model.UserLogin) (string, error) {
	user, err := u.userRepo.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return "", err
	}

	reqPassword := md5.HashWithMd5(req.Password)

	if reqPassword != user.Password {
		return "", customError.ErrIncorrectPassword
	}

	token, err := u.jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserById implements UserUCItf.
func (u *UserUsecase) GetUserById(ctx context.Context, userId int) (*model.UserResponse, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return convertToUserRespone(user), nil
}

// UpdateUser implements UserUCItf.
func (u *UserUsecase) UpdateUserData(ctx context.Context, req *model.UserUpdateData, userId int) (*model.UserResponse, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	if req.Name != "" {
		user.Name = req.Name
	}

	if req.Password != "" {
		hashedPassword := md5.HashWithMd5(req.Password)

		user.Password = hashedPassword
	}

	user.UpdatedAt = time.Now()

	err = u.userRepo.UpdateUserData(ctx, user)
	if err != nil {
		return nil, err
	}

	return convertToUserRespone(user), nil
}

func (u *UserUsecase) UpdateUserPhoto(ctx context.Context, req *model.UserUpdatePhoto, userId int) (*model.UserResponse, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if req.Photo != nil {

		photoLink, err := u.supabase.Upload(os.Getenv("SUPABASE_BUCKET_USER"), req.Photo)
		if err != nil {
			return nil, err
		}
		user.Photo = sql.NullString{
			String: photoLink,
			Valid:  true,
		}
	}

	user.UpdatedAt = time.Now()

	err = u.userRepo.UpdateUserPhoto(ctx, user)
	if err != nil {
		return nil, err
	}

	return convertToUserRespone(user), nil
}

func convertToUserRespone(user *model.User) *model.UserResponse {
	return &model.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Photo:     user.Photo,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
