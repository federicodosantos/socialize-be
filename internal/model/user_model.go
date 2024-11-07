package model

import (
	"database/sql"
	"mime/multipart"
	"time"
)

type User struct {
	ID string `db:"id"`
	Name string `db:"name"`
	Email string `db:"email"`
	Password string `db:"password"`
	Photo sql.NullString `db:"photo"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserRegister struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
	ConfirmPassword string `json:"confirm_password"`
}

type UserLogin struct {
	Name string `json:"name"`
	Email string `json:"email"`
	Password string `json:"password"`
}

type UserUpdateData struct {
	Name string `json:"name,omitempty"`
	Email string `json:"email,omitempty"`
	Password string `json:"password,omitempty"`
}

type UserUpdatePhoto struct {
	Photo *multipart.FileHeader `form:"photo,omitempty"`
}

type UserResponse struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Email string `json:"email"`
	Photo sql.NullString `json:"photo"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}