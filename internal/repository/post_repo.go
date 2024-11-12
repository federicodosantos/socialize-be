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
	CreatePost(ctx context.Context, post *model.Post) (int64, error)
	GetAllPost(ctx context.Context, filter model.PostFilter) ([]*model.Post, error)
	GetPostByID(ctx context.Context, postId int) (*model.Post, error)
	DeletePost(ctx context.Context, postId int) error
}

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepoItf {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *model.Post) (int64, error) {
	createdAtStr := util.ConvertTimeToString(post.CreatedAt)
	updatedAtStr := util.ConvertTimeToString(post.UpdatedAt)

	insertPostQuery := fmt.Sprintf(`INSERT INTO posts(id, user_id, content, created_at, updated_at)
				VALUES(%d, '%s', '%s', '%s', '%s')`, post.ID, post.UserID, post.Content, createdAtStr, updatedAtStr)

	res, err := r.db.ExecContext(ctx, insertPostQuery)
	if err != nil {
		return 0, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return 0, customError.ErrRowsAffected
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, customError.ErrLastInsertId
	}

	return id, util.ErrRowsAffected(rows)
}

func (r *PostRepo) GetAllPost(ctx context.Context, filter model.PostFilter) ([]*model.Post, error) {
	var posts []*model.Post

	getAllPostQuery := fmt.Sprintf(`SELECT * FROM posts`)

	if filter.Keyword != "" {
		getAllPostQuery = fmt.Sprintf(`%s WHERE content LIKE '%%%s%%'`, getAllPostQuery, filter.Keyword)
	}

	err := r.db.SelectContext(ctx, &posts, getAllPostQuery)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetPostByID(ctx context.Context, postId int) (*model.Post, error) {
	var post model.Post

	query := fmt.Sprintf("SELECT * FROM posts WHERE id = %d", postId)

	err := r.db.GetContext(ctx, &post, query)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) DeletePost(ctx context.Context, postId int) error {
	query := fmt.Sprintf("DELETE FROM posts WHERE id = %d", postId)

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
