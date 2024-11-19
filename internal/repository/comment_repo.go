package repository

import (
	"context"
	"fmt"

	"github.com/federicodosantos/socialize/internal/model"
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
	createdAtStr := util.ConvertTimeToString(comment.CreatedAt)
	
	createCommentQuery := fmt.Sprintf(
		`INSERT INTO comments(user_id, post_id, comment, created_at)
		VALUES(%d, %d, '%s', '%s')`, comment.UserID, comment.PostID, comment.Comment, createdAtStr)

	res, err := r.db.ExecContext(ctx, createCommentQuery)	
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

	util.ErrRowsAffected(rows)

	comment.ID = id

	return nil
}

func (r *CommentRepo) GetAllCommentsByPostId(ctx context.Context, postId int64) ([]*model.Comment, error) {
	var comments []*model.Comment

	getAllCommentsByPostIdQUery := fmt.Sprintf(`
	SELECT 
		c.id,
		c.post_id,
		c.user_id,
		c.comment,
		c.created_at,
		u.name AS user_name,      
		u.photo AS user_photo     
	FROM comments AS c
	JOIN users AS u ON u.id = c.user_id
	WHERE c.post_id = %d`, postId)

	err := r.db.SelectContext(ctx, &comments, getAllCommentsByPostIdQUery)
	if err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *CommentRepo) DeleteComment(ctx context.Context, id int64) error {
	deleteCommentQuery := fmt.Sprintf(`DELETE FROM comments where id = %d`, id)

	res, err := r.db.ExecContext(ctx, deleteCommentQuery)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return customerror.ErrRowsAffected
	}

	util.ErrRowsAffected(rows)

	return nil
}


