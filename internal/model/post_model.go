package model

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        string         `db:"id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	UserID    string         `db:"user_id"`
	Image     sql.NullString `db:"image"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
}

type PostCreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
	Image   string `json:"image"`
}

type PostResponse struct {
	ID        string    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	UserID    string    `json:"user_id"`
	Image     string    `json:"image"`
	CreatedAt time.Time `json:"created_at"`
}

type PostFilter struct {
	Keyword string `json:"keyword"`
}
