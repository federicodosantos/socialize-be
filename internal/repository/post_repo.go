package repository

import (
	"context"

	"github.com/federicodosantos/socialize/internal/model"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/jmoiron/sqlx"
)

type PostRepoItf interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetAllPost(ctx context.Context, filter model.PostFilter) ([]*model.Post, error)
	GetPostByID(ctx context.Context, postID int64) (*model.Post, error)
	DeletePost(ctx context.Context, postID int64) error

	CreateVote(ctx context.Context, postID int64, userID int64, vote int64) error
	DeleteVote(ctx context.Context, postID int64, userID int64) error
}

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepoItf {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *model.Post) error {
	query := `
	INSERT INTO posts (
		user_id, title, content, image, created_at, updated_at
	) VALUES (?, ?, ?, ?, ?, ?)`

	res, err := r.db.ExecContext(ctx, query, post.UserID, post.Title, post.Content, post.Image.String, post.CreatedAt, post.UpdatedAt)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customError.ErrRowsAffected
	}

	id, err := res.LastInsertId()
	if err != nil {
		return customError.ErrLastInsertId
	}

	post.ID = id

	return util.ErrRowsAffected(rows)
}

func (r *PostRepo) GetAllPost(ctx context.Context, filter model.PostFilter) ([]*model.Post, error) {
	var posts []*model.Post

	query := `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.user_id,
		p.image,
		u.name AS user_name,
		u.photo AS user_photo,
		p.created_at,
		p.updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = p.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = p.id) AS down_vote  
	FROM posts AS p
	JOIN users AS u ON u.id = p.user_id`

	var args []interface{}
	if filter.Keyword != "" {
		query += ` WHERE p.content LIKE ?`
		args = append(args, "%"+filter.Keyword+"%")
	}

	err := r.db.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetAllPostByUserID(ctx context.Context, filter model.PostFilter, userID int64) ([]*model.Post, error) {
	var posts []*model.Post

	query := `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.user_id,
		p.image,
		u.name AS user_name,
		u.photo AS user_photo,
		p.created_at,
		p.updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = p.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = p.id) AS down_vote   
	FROM posts AS p
	JOIN users AS u ON u.id = p.user_id
	WHERE user_id = ?`

	var args []interface{}
	args = append(args, userID)

	if filter.Keyword != "" {
		query += ` AND p.content LIKE ?`
		args = append(args, "%"+filter.Keyword+"%")
	}

	err := r.db.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	var post model.Post

	query := `
	SELECT 
		p.id,
		p.title,
		p.content,
		p.user_id,
		p.image,
		u.name AS user_name,
		u.photo AS user_photo,
		p.created_at,
		p.updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = p.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = p.id) AS down_vote   
	FROM posts AS p
	JOIN users AS u ON u.id = p.user_id 
	WHERE p.id = ?`

	err := r.db.GetContext(ctx, &post, query, postID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) DeletePost(ctx context.Context, postID int64) error {
	query := `DELETE FROM posts WHERE id = ?`

	res, err := r.db.ExecContext(ctx, query, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customError.ErrRowsAffected
	}

	return util.ErrRowsAffected(rows)
}

func (r *PostRepo) CreateVote(ctx context.Context, postID int64, userID int64, vote int64) error {
	query := `INSERT INTO votes(post_id, user_id, vote) VALUES(?, ?, ?)`

	_, err := r.db.ExecContext(ctx, query, postID, userID, vote)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepo) DeleteVote(ctx context.Context, postID int64, userID int64) error {
	query := `DELETE FROM votes WHERE post_id = ? AND user_id = ?`

	_, err := r.db.ExecContext(ctx, query, postID, userID)
	if err != nil {
		return err
	}

	return nil
}
