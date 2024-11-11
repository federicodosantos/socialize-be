package usecase

import (
	"context"
	"database/sql"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/repository"
)

type PostUsecaseItf interface {
}

type PostUsecase struct {
	postRepo repository.PostRepoItf
}

func NewPostUsecase(postRepo repository.PostRepoItf) PostUsecaseItf {
	return &PostUsecase{
		postRepo: postRepo,
	}
}

func (uc *PostUsecase) CreatePost(ctx context.Context, req *model.PostCreate, userID string) (int64, error) {
	data := model.Post{
		Title:   req.Title,
		Content: req.Content,
		Image:   sql.NullString{String: req.Image},
		UserID:  userID,
	}

	return uc.postRepo.CreatePost(ctx, &data)
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
		Title:     post.Title,
		Content:   post.Content,
		Image:     post.Image.String,
		CreatedAt: post.CreatedAt,
	}
}

func (uc *PostUsecase) GetPostByID(ctx context.Context, postID string) (*model.PostResponse, error) {
	post, err := uc.postRepo.GetPostByID(ctx, postID)
	if err != nil {
		return nil, err
	}

	// get comment by postID

	return convertToPostRespone(post), nil
}

func (uc *PostUsecase) DeletePost(ctx context.Context, postID string) error {
	return uc.postRepo.DeletePost(ctx, postID)
}
