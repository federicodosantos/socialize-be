package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID        string         `db:"id"`
	Name      string         `db:"name"`
	Email     string         `db:"email"`
	Password  string         `db:"password"`
	Photo     sql.NullString `db:"photo"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type UserRegister struct {
	Name            string `json:"name" validate:"required"`
	Email           string `json:"email" validate:"required,email"`
	Password        string `json:"password" validate:"required,min=8"`
	ConfirmPassword string `json:"confirm_password" validate:"required,eqfield=Password"`
}

type UserLogin struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}


type UserUpdateData struct {
	Name     string `json:"name,omitempty" validate:"omitempty,min=2"`
	Email    string `json:"email,omitempty" validate:"omitempty,email"`
	Password string `json:"password,omitempty" validate:"omitempty,min=8"`
}

type UserUpdatePhoto struct {
	PhotoUrl string `json:"photo_url" validate:"required,url"`
}

type UserResponse struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Photo     string    `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
