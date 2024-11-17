package http

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/usecase"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	response "github.com/federicodosantos/socialize/pkg/response"
	"github.com/federicodosantos/socialize/pkg/util"
	"github.com/go-chi/chi/v5"
)

const maxUploadSize = 2 * 1024 * 1024

type UserHandler struct {
	userUC usecase.UserUsecaseItf
}

func NewUserHandler(userUC usecase.UserUsecaseItf) *UserHandler {
	return &UserHandler{userUC: userUC}
}

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

func (uh *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req *model.UserRegister

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	user, err := uh.userUC.Register(reqCtx, req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrEmailExist):
			response.FailedResponse(w, http.StatusConflict, err.Error())
			return
		}
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusCreated, "successfully create user", user)
}

func (uh *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req *model.UserLogin

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	reqCtx := r.Context()

	token, err := uh.userUC.Login(reqCtx, req)
	if err != nil {
		switch {
		case errors.Is(err, customError.ErrEmailNotFound):
			response.FailedResponse(w, http.StatusNotFound, err.Error())
			return
		case errors.Is(err, customError.ErrIncorrectPassword):
			response.FailedResponse(w, http.StatusUnauthorized, err.Error())
			return
		default:
			response.FailedResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
	}

	response.SuccessResponse(w, http.StatusOK, "successfully login to account", token)
}

func (uh *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userId, err := util.GetUserIdFromContext(w,r)
	if err != nil {
		return
	}

	reqCtx := r.Context()

	user, err := uh.userUC.GetUserById(reqCtx, userId)
	if err != nil {
		if errors.Is(err, customError.ErrUserNotFound) {
			response.FailedResponse(w, http.StatusNotFound, err.Error())
			return
		}
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "successfully get current user", user)
}

func (uh *UserHandler) UpdateUserPhoto(w http.ResponseWriter, r *http.Request) {
	var req *model.UserUpdatePhoto

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := util.GetUserIdFromContext(w,r)
	if err != nil {
		return
	}

	reqCtx := r.Context()

	updatedUser, err := uh.userUC.UpdateUserPhoto(reqCtx, req, userId)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Successfully update user data", updatedUser)
}

func (uh *UserHandler) UpdateUserData(w http.ResponseWriter, r *http.Request) {
	var req *model.UserUpdateData

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	userId, err := util.GetUserIdFromContext(w,r)
	if err != nil {
		return
	}

	reqCtx := r.Context()

	updatedUser, err := uh.userUC.UpdateUserData(reqCtx, req, userId)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "Successfully update user Data", updatedUser)
}