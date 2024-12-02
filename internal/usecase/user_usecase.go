package usecase

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/jwt"
	"github.com/federicodosantos/socialize/pkg/regex"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserUsecaseItf interface {
	Register(ctx context.Context, req *model.UserRegister) (*model.UserResponse, error)
	Login(ctx context.Context, req *model.UserLogin) (string, error)
	GetUserById(ctx context.Context, userId string) (*model.UserResponse, error)
	UpdateUserData(ctx context.Context, req *model.UserUpdateData, userId string) (*model.UserResponse, error)
	UpdateUserPhoto(ctx context.Context, req *model.UserUpdatePhoto, userId string) (*model.UserResponse, error)
}

type UserUsecase struct {
	userRepo repository.UserRepoItf
	jwt      jwt.JWTItf
}

func NewUserUsecase(userRepo repository.UserRepoItf,
	jwt jwt.JWTItf) UserUsecaseItf {
	return &UserUsecase{
		userRepo: userRepo,
		jwt:      jwt,
	}
}

// Register implements UserUCItf.
func (u *UserUsecase) Register(ctx context.Context, req *model.UserRegister) (*model.UserResponse, error) {
	err := regex.Password(req.Password)
	if err != nil {
		return nil, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	createdUser := &model.User{
		ID:        uuid.NewString(),
		Name:      req.Name,
		Email:     req.Email,
		Password:  string(hashedPassword),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err = u.userRepo.CreateUser(ctx, createdUser)
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

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return "", err
	}

	token, err := u.jwt.CreateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

// GetUserById implements UserUCItf.
func (u *UserUsecase) GetUserById(ctx context.Context, userId string) (*model.UserResponse, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	return convertToUserRespone(user), nil
}

// UpdateUser implements UserUCItf.
func (u *UserUsecase) UpdateUserData(ctx context.Context, req *model.UserUpdateData, userId string) (*model.UserResponse, error) {
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
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}

		user.Password = string(hashedPassword)
	}

	user.UpdatedAt = time.Now()

	err = u.userRepo.UpdateUserData(ctx, user)
	if err != nil {
		return nil, err
	}

	return convertToUserRespone(user), nil
}

func (u *UserUsecase) UpdateUserPhoto(ctx context.Context, req *model.UserUpdatePhoto, userId string) (*model.UserResponse, error) {
	user, err := u.userRepo.GetUserById(ctx, userId)
	if err != nil {
		return nil, err
	}

	if req.PhotoUrl != "" {
		user.Photo = sql.NullString{
			String: req.PhotoUrl,
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
		Photo:     user.Photo.String,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
