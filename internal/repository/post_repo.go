package repository

import (
	"context"
	"fmt"

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
	DeletVote(ctx context.Context, postID int64, userID int64) error
}

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepoItf {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *model.Post) error {
	createdAtStr := util.ConvertTimeToString(post.CreatedAt)
	updatedAtStr := util.ConvertTimeToString(post.UpdatedAt)

	insertPostQuery := fmt.Sprintf(`
	INSERT INTO posts (
		user_id, title, content, image, created_at, updated_at
	) VALUES (
		%d, '%s', '%s', '%s', '%s', '%s'
	)`, post.UserID, post.Title, post.Content, post.Image.String, createdAtStr, updatedAtStr)

	res, err := r.db.ExecContext(ctx, insertPostQuery)
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

	getAllPostQuery := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content,
		user_id,
		image,
		created_at,
		updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = posts.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = posts.id) AS down_vote  
	FROM posts
	`)

	if filter.Keyword != "" {
		getAllPostQuery = fmt.Sprintf(`%s WHERE content LIKE '%%%s%%'`, getAllPostQuery, filter.Keyword)
	}

	err := r.db.SelectContext(ctx, &posts, getAllPostQuery)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetAllPostByUserID(ctx context.Context, filter model.PostFilter, userID int64) ([]*model.Post, error) {
	var posts []*model.Post

	getAllPostQuery := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content,
		user_id,
		image,
		created_at,
		updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = posts.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = posts.id) AS down_vote  
	FROM posts 
	WHERE user_id = %d
	`, userID)

	if filter.Keyword != "" {
		getAllPostQuery = fmt.Sprintf(`%s WHERE content LIKE '%%%s%%'`, getAllPostQuery, filter.Keyword)
	}

	err := r.db.SelectContext(ctx, &posts, getAllPostQuery)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	var post model.Post

	query := fmt.Sprintf(`
	SELECT 
		id,
		title,
		content,
		user_id,
		image,
		created_at,
		updated_at, 
		(SELECT count(*) from votes WHERE vote = 1 AND post_id = posts.id) AS up_vote, 
		(SELECT count(*) from votes WHERE vote = -1 AND post_id = posts.id) AS down_vote  
	FROM posts 
	WHERE id = %d`, postID)

	err := r.db.GetContext(ctx, &post, query)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) DeletePost(ctx context.Context, postID int64) error {
	query := fmt.Sprintf("DELETE FROM posts WHERE id = %d", postID)

	res, err := r.db.ExecContext(ctx, query)
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
	query := fmt.Sprintf("INSERT INTO votes(post_id, user_id, vote) VALUES(%d, %d, %d)", postID, userID, vote)

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}

func (r *PostRepo) DeletVote(ctx context.Context, postID int64, userID int64) error {
	query := fmt.Sprintf("DELETE FROM votes WHERE post_id = %d AND user_id = %d", postID, userID)

	_, err := r.db.ExecContext(ctx, query)
	if err != nil {
		return err
	}

	return nil
}
