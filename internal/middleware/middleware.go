package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	customContext "github.com/federicodosantos/socialize/pkg/context"
	"github.com/federicodosantos/socialize/pkg/jwt"
	response "github.com/federicodosantos/socialize/pkg/response"
	"go.uber.org/zap"
)

type MiddlewareItf interface {
	JwtAuthMiddleware(next http.Handler) http.Handler
	LoggingMiddleware(next http.Handler) http.Handler
}

type Middleware struct {
	jwt    jwt.JWTItf
	logger *zap.SugaredLogger
}

func NewMiddleware(jwt jwt.JWTItf, logger *zap.SugaredLogger) MiddlewareItf {
	return &Middleware{jwt: jwt, logger: logger}
}

func (m *Middleware) JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("jwt-token")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				response.FailedResponse(w, http.StatusUnauthorized, "no cookie provided")
				return
			}
			response.FailedResponse(w, http.StatusBadRequest, "bad request")
			return
		}

		token := cookie.Value

		userID, err := m.jwt.VerifyToken(token)
		if err != nil {
			response.FailedResponse(w, http.StatusUnauthorized, err.Error())
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

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}
