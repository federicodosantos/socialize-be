package delivery

import (
	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/go-chi/chi/v5"
)

func UserRoutes(router *chi.Mux, userHandle *UserHandler, middleware middleware.MiddlewareItf) {
	// public routes
	router.Post("/auth/register", userHandle.Register)
	router.Post("/auth/login", userHandle.Login)

	// private routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)
		r.Get("/auth/current-user", userHandle.GetCurrentUser)
		r.Patch("/auth/update-photo", userHandle.UpdateUserPhoto)
		r.Patch("/auth/update-data", userHandle.UpdateUserData)
	})
}