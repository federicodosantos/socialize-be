package http

import (
	"net/http"

	"github.com/federicodosantos/socialize/internal/middleware"
	"github.com/federicodosantos/socialize/internal/usecase"
	response "github.com/federicodosantos/socialize/pkg/response"
	"github.com/go-chi/chi/v5"
)

type FileHandler struct {
	fileUsecase usecase.FileUsecaseItf
}

func NewFileHandler(fileUsecase usecase.FileUsecaseItf) *FileHandler {
	return &FileHandler{fileUsecase: fileUsecase}
}

func FileRoutes(router *chi.Mux, fileHandler *FileHandler, middleware middleware.MiddlewareItf) {
	// private routes
	router.Group(func(r chi.Router) {
		r.Use(middleware.JwtAuthMiddleware)
		r.Post("/file/upload", fileHandler.UploadFile)
	})
}

func (h *FileHandler) UploadFile(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		response.FailedResponse(w, http.StatusBadRequest, "File is too big. 2MB maximum.")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		response.FailedResponse(w, http.StatusBadRequest, "Failed to get file.")
		return
	}
	defer file.Close()

	if header.Size > maxUploadSize {
		response.FailedResponse(w, http.StatusBadRequest, "The file size exceeds the 2MB limit.")
		return
	}

	url, err := h.fileUsecase.UploadFile(r.Context(), header)
	if err != nil {
		response.FailedResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	response.SuccessResponse(w, http.StatusOK, "File uploaded successfully", url)

}
