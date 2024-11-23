package repository

import (
	"context"
	"fmt"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository/query"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/jmoiron/sqlx"
)

type PostRepoItf interface {
	CreatePost(ctx context.Context, post *model.Post) error
	GetAllPost(ctx context.Context, filter model.PostFilter) ([]*model.Post, error)
	GetPostByID(ctx context.Context, postID int64) (*model.Post, error)
	DeletePost(ctx context.Context, postID int64) error

	CreateVote(ctx context.Context, postID int64, userID string, vote int64) error
	DeletVote(ctx context.Context, postID int64, userID string) error
}

type PostRepo struct {
	db *sqlx.DB
}

func NewPostRepo(db *sqlx.DB) PostRepoItf {
	return &PostRepo{db: db}
}

func (r *PostRepo) CreatePost(ctx context.Context, post *model.Post) error {
	res, err := r.db.ExecContext(ctx, query.InsertPostQuery, 
		post.UserID, post.Title, post.Content, post.Image, post.CreatedAt, post.UpdatedAt)
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

	query := query.GetAllPostQuery 
	var args []interface{}

	if filter.Keyword != "" {
		query += " WHERE p.content LIKE ?"
		args = append(args, "%"+filter.Keyword+"%")
	}

	err := r.db.SelectContext(ctx, &posts, query, args...)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetAllPostByUserID(ctx context.Context, filter model.PostFilter, userID string) ([]*model.Post, error) {
	var posts []*model.Post

	query := query.GetAllPostByUserIDQuery
	var args []interface{}

	if filter.Keyword != "" {
		query += " WHERE p.content LIKE ?"
		args = append(args, "%"+filter.Keyword+"%")
	}

	err := r.db.SelectContext(ctx, &posts, query, userID, args)
	if err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *PostRepo) GetPostByID(ctx context.Context, postID int64) (*model.Post, error) {
	var post model.Post

	err := r.db.GetContext(ctx, &post, query.GetPostByIDQuery, postID)
	if err != nil {
		return nil, err
	}

	return &post, nil
}

func (r *PostRepo) DeletePost(ctx context.Context, postID int64) error {
	res, err := r.db.ExecContext(ctx, query.DeletePostQuery, postID)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customError.ErrRowsAffected
	}

	return util.ErrRowsAffected(rows)
}

func (r *PostRepo) CreateVote(ctx context.Context, postID int64, userID string, vote int64) error {
    var count int

    err := r.db.QueryRowContext(ctx, query.CheckPostExistQuery, postID).Scan(&count)
    if err != nil {
        return fmt.Errorf("failed to check post_id: %w", err)
    }

    if count == 0 {
        return fmt.Errorf("postID %d does not exist in posts", postID)
    }

    err = r.db.QueryRowContext(ctx, query.CheckVoteExistQuery, postID, userID).Scan(&count)
    if err != nil {
        return fmt.Errorf("failed to check existing vote: %w", err)
    }

    if count > 0 {
        _, err = r.db.ExecContext(ctx, query.UpdateVoteQuery, vote, postID, userID)
        if err != nil {
            return fmt.Errorf("failed to update vote: %w", err)
        }
    } else {
        _, err = r.db.ExecContext(ctx, query.InsertVoteQuery, postID, userID, vote)
        if err != nil {
            return fmt.Errorf("failed to insert vote: %w", err)
        }
    }

    return nil
}

func (r *PostRepo) DeletVote(ctx context.Context, postID int64, userID string) error {
	_, err := r.db.ExecContext(ctx, query.DeleteVoteQuery, postID, userID)
	if err != nil {
		return err
	}

	return nil
}
