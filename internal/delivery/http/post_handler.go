package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/usecase"
	response "github.com/federicodosantos/socialize/pkg/response"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/go-chi/chi/v5"
)

type PostHandler struct {
	postUsecase usecase.PostUsecaseItf
}

func NewPostHandler(postUsecase usecase.PostUsecaseItf) *PostHandler {
	return &PostHandler{postUsecase: postUsecase}
}

func PostRoutes(router *chi.Mux, postHandle *PostHandler, middleware middleware.MiddlewareItf) {
	// private routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)
		r.Post("/post", postHandle.CreatePost)
		r.Get("/post", postHandle.GetAllPost)
		r.Get("/post/{postID}", postHandle.GetPostByID)
		r.Delete("/post/{postID}", postHandle.DeletePost)
		r.Post("/post/{postID}/up-vote", postHandle.UpVote)
		r.Post("/post/{postID}/down-vote", postHandle.DownVote)
	})
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
	var model *model.PostCreate

	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	userID, err := util.GetUserIdFromContext(w, r)
	if err != nil {
		return
	}

	id, err := h.postUsecase.CreatePost(reqCtx, model, userID)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Post created successfully", id)
}

func (h *PostHandler) GetAllPost(w http.ResponseWriter, r *http.Request) {
	reqCtx := r.Context()

	var filter model.PostFilter
	if err := util.ParsePostFilter(r, &filter); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	posts, err := h.postUsecase.GetAllPost(reqCtx, filter)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Get all post successfully", posts)
}

func (h *PostHandler) GetPostByID(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	post, err := h.postUsecase.GetPostByID(reqCtx, postID)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Get post by ID successfully", post)
}

func (h *PostHandler) DeletePost(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	err = h.postUsecase.DeletePost(reqCtx, postID)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Post deleted successfully", nil)
}

func (h *PostHandler) UpVote(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	userID, err := util.GetUserIdFromContext(w, r)
	if err != nil {
		return
	}

	err = h.postUsecase.CreateUpVote(reqCtx, postID, userID)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Upvote successfully", nil)
}

func (h *PostHandler) DownVote(w http.ResponseWriter, r *http.Request) {
	postIDStr := chi.URLParam(r, "postID")
	postID, err := strconv.ParseInt(postIDStr, 10, 64)
	if err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	userID, err := util.GetUserIdFromContext(w, r)
	if err != nil {
		return
	}

	err = h.postUsecase.CreateDownVote(reqCtx, postID, userID)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Downvote successfully", nil)
}
