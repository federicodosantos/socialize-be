package model

import (
	"database/sql"
	"time"
)

type Post struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	UserID    int64          `db:"user_id"`
	Image     sql.NullString `db:"image"`
	CreatedAt time.Time      `db:"created_at"`
	UpdatedAt time.Time      `db:"updated_at"`
	UpVote    int64          `db:"up_vote"`
	DownVote  int64          `db:"down_vote"`
}

type PostCreate struct {
	Title   string `json:"title" validate:"required"`
	Content string `json:"content" validate:"required"`
	Image   string `json:"image"`
}

type PostResponse struct {
	ID           int64              `json:"id"`
	Title        string             `json:"title"`
	Content      string             `json:"content"`
	UserID       int64              `json:"user_id"`
	UserPhoto    string             `json:"user_photo"`
	UserName     string             `json:"user_name"`
	Image        string             `json:"image"`
	TotalComment int64              `json:"total_comment"`
	Comment      []*CommentResponse `json:"comment,omitempty"`
	UpVote       int64              `json:"up_vote"`
	DownVote     int64              `json:"down_vote"`
	CreatedAt    time.Time          `json:"created_at"`
	UpdatedAt    time.Time          `json:"updated_at"`
}

type PostFilter struct {
	Keyword string `json:"keyword"`
}
