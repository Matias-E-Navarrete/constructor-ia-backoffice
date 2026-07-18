// Package httpmiddleware holds the cross-cutting HTTP guards shared by every
// bounded context: authentication, pro-plan gating, admin gating, and
// feature-flag kill switches.
package httpmiddleware

import (
	"context"
	"net/http"
	"strings"

	"rimu/backend/internal/platform/auth"
	userdomain "rimu/backend/internal/user/domain"
)

type contextKey string

const userContextKey contextKey = "authenticated_user"

// RequireAuth parses the Bearer token, loads the user, and stores it in the
// request context for downstream handlers and other middlewares.
func RequireAuth(jwtSecret string, users userdomain.Repository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			token, ok := strings.CutPrefix(header, "Bearer ")
			if !ok || token == "" {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized", "missing bearer token")
				return
			}

			claims, err := auth.ParseToken(jwtSecret, token)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized", "invalid token")
				return
			}

			user, err := users.FindByID(r.Context(), claims.UserID)
			if err != nil {
				writeJSONError(w, http.StatusUnauthorized, "unauthorized", "user not found")
				return
			}

			ctx := context.WithValue(r.Context(), userContextKey, user)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func UserFromContext(ctx context.Context) (*userdomain.User, bool) {
	u, ok := ctx.Value(userContextKey).(*userdomain.User)
	return u, ok
}

// RequirePro rejects the request unless the authenticated user is on the pro
// plan. Must run after RequireAuth.
func RequirePro(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || !user.IsPro() {
			writeJSONError(w, http.StatusForbidden, "upgrade_required", "this feature requires the pro plan")
			return
		}
		next.ServeHTTP(w, r)
	})
}

// RequireAdmin rejects the request unless the authenticated user has the
// admin role. Must run after RequireAuth.
func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := UserFromContext(r.Context())
		if !ok || !user.IsAdmin() {
			writeJSONError(w, http.StatusForbidden, "forbidden", "admin access required")
			return
		}
		next.ServeHTTP(w, r)
	})
}
