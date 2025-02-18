package tokens

import (
	"context"
	"net/http"
	"strings"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
)

// RoleMdl validates the JWT token by roles and injects user data into the request context.
func RoleMdl(ctx context.Context, sk string, roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				httpkit.RenderErr(w, problems.Unauthorized("Missing Authorization header"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				httpkit.RenderErr(w, problems.Unauthorized("Invalid Authorization header"))
				return
			}

			tokenString := parts[1]

			tokenData, err := VerifyJWT(ctx, tokenString, sk)
			if err != nil {
				httpkit.RenderErr(w, problems.Unauthorized("Token validation failed"))
				return
			}
			if tokenData.Role == nil {
				httpkit.RenderErr(w, problems.Unauthorized("Token validation failed"))
				return
			}

			roleAllowed := false
			for _, role := range roles {
				if *tokenData.Role == role {
					roleAllowed = true
					break
				}
			}
			if !roleAllowed {
				httpkit.RenderErr(w, problems.Unauthorized("Role not allowed"))
				return
			}

			ctx = context.WithValue(r.Context(), UserIDKey, tokenData.ID)
			ctx = context.WithValue(ctx, RoleKey, tokenData.Role)
			ctx = context.WithValue(ctx, SessionIDKey, tokenData.SessionID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
