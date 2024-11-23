package repository

import (
	"context"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository/query"
	customerror "github.com/federicodosantos/socialize/pkg/custom-error"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/jmoiron/sqlx"
)

type CommentRepoItf interface {
	CreateComment(ctx context.Context, comment *model.Comment) error
	GetAllCommentsByPostId(ctx context.Context, postId int64) ([]*model.Comment, error)
	DeleteComment(ctx context.Context, id int64) error
}

type CommentRepo struct {
	db *sqlx.DB
}

func NewCommentRepo(db *sqlx.DB) CommentRepoItf {
	return &CommentRepo{db: db}
}
func (r *CommentRepo) CreateComment(ctx context.Context, comment *model.Comment) error {
	res, err := r.db.ExecContext(ctx, query.InsertCommentQuery, 
		comment.UserID, comment.PostID, comment.Comment, comment.CreatedAt)	
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customerror.ErrRowsAffected
	}

	id, err := res.LastInsertId()
	if err != nil {
		return customerror.ErrLastInsertId
	}

	if err := util.ErrRowsAffected(rows); err != nil {
		return err
	}

	comment.ID = id

	return nil
}

func (r *CommentRepo) GetAllCommentsByPostId(ctx context.Context, postId int64) ([]*model.Comment, error) {
	var comments []*model.Comment

	err := r.db.SelectContext(ctx, &comments, query.GetAllCommentsByPostIdQuery, postId)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, id int64) error {
	res, err := r.db.ExecContext(ctx, query.DeleteCommentQuery, id)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customerror.ErrRowsAffected
	}

	return util.ErrRowsAffected(rows)
}
