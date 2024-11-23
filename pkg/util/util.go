package util

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/federicodosantos/socialize/internal/model"
	customContext "github.com/federicodosantos/socialize/pkg/context"
	customError "github.com/federicodosantos/socialize/pkg/custom-error"
	response "github.com/federicodosantos/socialize/pkg/response"
	"github.com/go-chi/chi/v5"
	"github.com/jmoiron/sqlx"
)

func ErrRowsAffected(rows int64) error {
	if rows != 1 {
		return fmt.Errorf("error : %w, got %d rows affected", customError.ErrRowsAffected, rows)
	}

	return nil
}

func ConvertTimeToString(time time.Time) string {
	return time.Format("2006-01-02 15:04:05")
}

func GetUserIdFromContext(w http.ResponseWriter, r *http.Request) (string, error) {
	userID := r.Context().Value(customContext.UserIDKey)
	if userID == "" {
		response.FailedResponse(w, http.StatusUnauthorized, "User ID tidak ditemukan dalam konteks", nil)
		return "", errors.New("user id not found in context")
	}

	stringUserID, ok := userID.(string)
	if !ok {
		response.FailedResponse(w, http.StatusBadRequest, "User ID tidak valid dalam konteks", nil)
		return "", errors.New("invalid or missing userID in context")
	}

	return stringUserID, nil
}

func HealthCheck(router *chi.Mux, db *sqlx.DB) {
	type HealthStatus struct {
		Status   string `json:"status"`
		Database string `json:"database"`
	}

	router.Get("/health-check", func(w http.ResponseWriter, r *http.Request) {
		status := HealthStatus{
			Status:   "healthy",
			Database: "healthy",
		}

		if err := db.Ping(); err != nil {
			status.Status = "unhealthy"
			status.Database = "unhealthy"
		}

		httpStatus := http.StatusOK
		if status.Status != "healthy" {
			httpStatus = http.StatusServiceUnavailable
		}

		response.SuccessResponse(w, httpStatus, "health check", status)
	})
}

func ParsePostFilter(r *http.Request, filter *model.PostFilter) error {
	if keyword := r.URL.Query().Get("keyword"); keyword != "" {
		filter.Keyword = keyword
	}

	return nil

}
