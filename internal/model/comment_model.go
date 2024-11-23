package model

import (
	"database/sql"
	"time"
)

type Comment struct {
	ID        int64     		`db:"id"`
	PostID    int64     		`db:"post_id"`
	UserID    string     		`db:"user_id"`
	Comment   string    		`db:"comment"`
	UserName  string    		`db:"user_name"`   
	UserPhoto sql.NullString    `db:"user_photo"`  
	CreatedAt time.Time 		`db:"created_at"`
}

type CommentCreate struct {
	PostID  int64  `json:"post_id"`
	Comment string `json:"comment"`
}

type CommentResponse struct {
	ID        int64     `json:"id"`
	PostID    int64     `json:"post_id"`
	UserID    string    `json:"user_id"`
	Comment   string    `json:"comment"`
	UserName  string    `json:"user_name"`   
	UserPhoto string    `json:"user_photo"`  
	CreatedAt time.Time `json:"created_at"`
}