package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/cvcio/elections-api/pkg/auth"
	"github.com/gin-gonic/gin"
	"github.com/go-chi/render"
	"github.com/pkg/errors"
	"github.com/plagiari-sm/mediawatch/pkg/web"
)

// Auth is used to authenticate and authorize HTTP requests.
type Auth struct {
	Authenticator *auth.Authenticator
}

// Authenticate validates a JWT from the `Authorization` header.
func (a *Auth) Authenticate() gin.HandlerFunc { // func(next http.Handler) http.Handler {
	return func(c *gin.Context) {
		authHdr := c.Request.Header.Get("Authorization")
		if authHdr == "" {
			c.AbortWithError(http.StatusUnauthorized, errors.Errorf("No Authorization Header Error"))
			return
		}

		tknStr, err := parseAuthHeader(authHdr)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		claims, err := a.Authenticator.ParseClaims(tknStr)
		if err != nil {
			c.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		c.Request = c.Request.WithContext(
			context.WithValue(
				c.Request.Context(),
				auth.Key, claims,
			),
		)
		c.Next()
	}
}

// parseAuthHeader parses an authorization header. Expected header is of
// the format `Bearer <token>`.
func parseAuthHeader(bearerStr string) (string, error) {
	split := strings.Split(bearerStr, " ")
	if len(split) != 2 || strings.ToLower(split[0]) != "bearer" {
		return "", errors.New("Expected Authorization header format: Bearer <token>")
	}

	return split[1], nil
}

// HasRole validates that an authenticated user has at least one role from a
// specified list. This method constructs the actual function that is used.
func (a *Auth) HasRole(roles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value(auth.Key).(auth.Claims)
			if !ok {
				render.Render(w, r, web.ErrUnauthorized)
				return
			}

			if !claims.HasRole(roles...) {
				render.Render(w, r, web.ErrForbidden)
				return
			}
			next.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
}
