package delivery

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/federicodosantos/socialize/internal/model"
	"github.com/federicodosantos/socialize/internal/usecase"
	customContext "github.com/federicodosantos/socialize/pkg/context"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	response "github.com/federicodosantos/socialize/pkg/response"
)

const maxUploadSize = 2 * 1024 * 1024

type UserHandler struct {
	userUC usecase.UserUCItf
}

func NewUserHandler(userUC usecase.UserUCItf) *UserHandler {
	return &UserHandler{userUC: userUC}
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

	http.SetCookie(w, &http.Cookie{
		Name: "jwt-token",
		Value: token,
		Expires: time.Now().Add(24 * time.Hour),
		Path: "/",
	})

	response.SuccessResponse(w, http.StatusOK, "successfully login to account", nil)
}

func (uh *UserHandler) GetCurrentUser(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(customContext.UserIDKey)
	if userID == "" {
		response.FailedResponse(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	stringUserID, ok := userID.(string)
	if !ok {
		response.FailedResponse(w, http.StatusBadRequest, "invalid or missing userID in context")
		return
	}

	reqCtx := r.Context()

	user, err := uh.userUC.GetUserById(reqCtx, stringUserID)
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
	var req model.UserUpdatePhoto

	if err :=r.ParseMultipartForm(maxUploadSize); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, "File is too big. 2MB maximum.")
		return
	}
	
	file, header, err := r.FormFile("photo")
	if err != nil {
		response.FailedResponse(w, http.StatusBadRequest, "Failed to get file.")
		return
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		response.FailedResponse(w, http.StatusBadRequest, "The file size exceeds the 2MB limit.")
		return
	}

	req.Photo = header

	userID := r.Context().Value(customContext.UserIDKey)
	if userID == "" {
		response.FailedResponse(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	stringUserID, ok := userID.(string)
	if !ok {
		response.FailedResponse(w, http.StatusBadRequest, "invalid or missing userID in context")
		return
	}

	reqCtx := r.Context()

	updatedUser, err := uh.userUC.UpdateUserPhoto(reqCtx, &req, stringUserID)
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

	userID := r.Context().Value(customContext.UserIDKey)
	if userID == "" {
		response.FailedResponse(w, http.StatusUnauthorized, "user id not found in context")
		return
	}

	stringUserID, ok := userID.(string)
	if !ok {
		response.FailedResponse(w, http.StatusBadRequest, "invalid or missing userID in context")
		return
	}

	reqCtx := r.Context()

	updatedUser, err := uh.userUC.UpdateUserData(reqCtx, req, stringUserID)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return 
	}

	response.SuccessResponse(w, http.StatusOK, "Successfully update user Data", updatedUser)
}