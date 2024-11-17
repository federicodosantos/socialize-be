package usecase

import (
	"context"
	"database/sql"
	"time"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository"
)

type PostUsecaseItf interface {
	CreatePost(ctx context.Context, req *model.PostCreate, userID int64) (*model.PostResponse, error)
	GetAllPost(ctx context.Context, filter model.PostFilter) ([]model.PostResponse, error)
	GetPostByID(ctx context.Context, postID int64) (*model.PostResponse, error)
	DeletePost(ctx context.Context, postID int64) error

	CreateUpVote(ctx context.Context, postID int64, userID int64) error
	CreateDownVote(ctx context.Context, postID int64, userID int64) error
}

type PostUsecase struct {
	postRepo repository.PostRepoItf
}

func NewPostUsecase(postRepo repository.PostRepoItf) PostUsecaseItf {
	return &PostUsecase{
		postRepo: postRepo,
	}
}

func (uc *PostUsecase) CreatePost(ctx context.Context, req *model.PostCreate, userID int64) (*model.PostResponse, error) {
	data := &model.Post{
		Title:     req.Title,
		Content:   req.Content,
		Image:     sql.NullString{String: req.Image},
		UserID:    userID,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	err := uc.postRepo.CreatePost(ctx, data)
	if err != nil {
		return nil, err
	}

	res := convertToPostRespone(data)

	return res, nil
}

func (uc *PostUsecase) GetAllPost(ctx context.Context, filter model.PostFilter) ([]model.PostResponse, error) {
	posts, err := uc.postRepo.GetAllPost(ctx, filter)
	if err != nil {
		return nil, err
	}

	var postsResp []model.PostResponse
	for _, post := range posts {
		postsResp = append(postsResp, *convertToPostRespone(post))
	}

	return postsResp, nil
}

func convertToPostRespone(post *model.Post) *model.PostResponse {
	return &model.PostResponse{
		ID:        post.ID,
		UserID:    post.UserID,
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image.String,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
		UpVote:    post.UpVote,
		DownVote:  post.DownVote,
	}
}

func (uc *PostUsecase) GetPostByID(ctx context.Context, postID int64) (*model.PostResponse, error) {
	post, err := uc.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// ToDo : get comment by postID

	return convertToPostRespone(post), nil
}

func (uc *PostUsecase) DeletePost(ctx context.Context, postID int64) error {
	return uc.postRepo.DeletePost(ctx, postID)
}

func (uc *PostUsecase) CreateUpVote(ctx context.Context, postID int64, userID int64) error {
	err := uc.postRepo.DeletVote(ctx, postID, userID)
	if err != nil {
		return err
	}

	return uc.postRepo.CreateVote(ctx, postID, userID, 1)
}

func (uc *PostUsecase) CreateDownVote(ctx context.Context, postID int64, userID int64) error {
	err := uc.postRepo.DeletVote(ctx, postID, userID)
	if err != nil {
		return err
	}

	return uc.postRepo.CreateVote(ctx, postID, userID, -1)
}
