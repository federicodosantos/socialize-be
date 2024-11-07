package middleware

import (
	"context"
	"errors"
	"net/http"

	"github.com/federicodosantos/socialize/pkg/jwt"
	customContext "github.com/federicodosantos/socialize/pkg/context"
	response "github.com/federicodosantos/socialize/pkg/response"
)

type MiddlewareItf interface {
	JwtAuthMiddleware(next http.Handler) http.Handler
}

type Middleware struct {
	jwt jwt.JWTItf
}

func NewMiddleware(jwt jwt.JWTItf) MiddlewareItf {
	return &Middleware{jwt: jwt}
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
