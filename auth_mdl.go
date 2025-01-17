package tokens

import (
	"context"
	"net/http"
	"strings"

	"github.com/recovery-flow/comtools/httpkit"
	"github.com/recovery-flow/comtools/httpkit/problems"
)

type contextKey string

const (
	UserIDKey   contextKey = "userID"
	RoleKey     contextKey = "Role"
	DeviceIDKey contextKey = "deviceID"
)

// AuthMdl validates the JWT token and injects user data into the request context.
func (m *TokenManager) AuthMdl(secretKey string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				m.log.Debugf("Missing Authorization header")
				httpkit.RenderErr(w, problems.Unauthorized("Missing Authorization header"))
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
				m.log.Debugf("Invalid Authorization header format")
				httpkit.RenderErr(w, problems.Unauthorized("Invalid Authorization header"))
				return
			}

			tokenString := parts[1]

			m.log.Debugf("Token received: %s", tokenString)

			userData, err := m.VerifyJWTAndExtractClaims(tokenString, secretKey)
			if err != nil {
				m.log.Debugf("Token validation failed: %v", err)
				httpkit.RenderErr(w, problems.Unauthorized("Token validation failed"))
				return
			}
			if userData == nil {
				m.log.Debugf("Token validation failed")
				httpkit.RenderErr(w, problems.Unauthorized("Token validation failed"))
				return
			}

			ctx := context.WithValue(r.Context(), UserIDKey, userData.ID)
			ctx = context.WithValue(ctx, RoleKey, userData.Role)
			ctx = context.WithValue(ctx, DeviceIDKey, userData.DevID)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
