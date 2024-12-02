package middleware

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	customContext "github.com/federicodosantos/socialize/pkg/context"
	"github.com/federicodosantos/socialize/pkg/jwt"
	response "github.com/federicodosantos/socialize/pkg/response"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

type MiddlewareItf interface {
	JwtAuthMiddleware(next http.Handler) http.Handler
	LoggingMiddleware(next http.Handler) http.Handler
	ValidateMiddleware(next http.HandlerFunc, structToValidate interface{}) http.HandlerFunc
}

type Middleware struct {
	jwt    jwt.JWTItf
	logger *zap.SugaredLogger
	validate *validator.Validate
}

func NewMiddleware(jwt jwt.JWTItf, logger *zap.SugaredLogger) MiddlewareItf {
	return &Middleware{
		jwt: jwt,
		logger: logger,
		validate: validator.New(),}
}

func (m *Middleware) JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		bearerToken := r.Header.Get("Authorization")

		if bearerToken == "" {
			response.FailedResponse(w, http.StatusUnauthorized, "Authorization token is required", nil)
			return
		}

		token := strings.Split(bearerToken, " ")[1]

		userID, err := m.jwt.VerifyToken(token)
		if err != nil {
			response.FailedResponse(w, http.StatusUnauthorized, err.Error(), nil)
			return
		}

		// Set userID in context
		ctx := context.WithValue(r.Context(), customContext.UserIDKey, userID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *Middleware) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log informasi request
		m.logger.Infow("Incoming request",
			"method", r.Method,
			"url", r.URL.String(),
			"remote_addr", r.RemoteAddr,
		)

		// Create a response writer to capture the status code
		rr := &responseRecorder{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(rr, r)

		// Log informasi response
		duration := time.Since(start)
		m.logger.Infow("Request processed",
			"method", r.Method,
			"url", r.URL.String(),
			"duration", duration,
			"status", rr.statusCode,
		)
	})
}

func (m *Middleware) ValidateMiddleware(next http.HandlerFunc, structToValidate interface{}) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" || r.Method == "PUT" || r.Method == "PATCH" {
			if r.Body == nil {
				log.Println("Request body is nil")
				response.FailedResponse(w, http.StatusBadRequest, "Request body cannot be empty", nil)
				return
			}

			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil || len(bodyBytes) == 0 {
				m.logger.Info("Request body is empty or could not be read")
				response.FailedResponse(w, http.StatusBadRequest, "Request body cannot be empty", nil)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

			if err := json.NewDecoder(r.Body).Decode(structToValidate); err != nil {
				log.Printf("Failed to decode JSON: %v", err)
				response.FailedResponse(w, http.StatusBadRequest, err.Error(), nil)
				return
			}

			if err := m.validate.Struct(structToValidate); err != nil {
				log.Printf("Validation error: %v", err)
				validationErrors := make(map[string]string)
				for _, err := range err.(validator.ValidationErrors) {
					validationErrors[err.Field()] = formatValidationError(err)
				}
				response.FailedResponse(w, http.StatusBadRequest, "Validation error", validationErrors)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		next(w, r)
	}
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

func formatValidationError(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return fmt.Sprintf("The %s field is required.", err.Field())
	case "email":
		return fmt.Sprintf("The %s field must be a valid email address.", err.Field())
	case "min":
		return fmt.Sprintf("The %s field must have at least %s characters.", err.Field(), err.Param())
	case "max":
		return fmt.Sprintf("The %s field must have at most %s characters.", err.Field(), err.Param())
	case "eqfield":
		return fmt.Sprintf("The %s field must be equal to the %s field.", err.Field(), err.Param())
	case "url":
		return fmt.Sprintf("The %s field must be a valid URL.", err.Field())
	case "gt":
		return fmt.Sprintf("The %s field must be greater than %s.", err.Field(), err.Param())
	default:
		return fmt.Sprintf("The %s field is invalid.", err.Field())
	}
}

