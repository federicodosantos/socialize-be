package http

import (
	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/federicodosantos/socialize/internal/usecase"
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
	})
}
