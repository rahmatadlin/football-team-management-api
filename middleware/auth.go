package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"football-team-management-api/config"
	jwtutil "football-team-management-api/utils/jwt"
	"football-team-management-api/utils/response"
)

const ctxAdminID = "admin_id"

func AdminAuth(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		h := c.GetHeader("Authorization")
		if h == "" || !strings.HasPrefix(strings.ToLower(h), "bearer ") {
			response.ErrorMessage(c, http.StatusUnauthorized, "missing or invalid authorization header")
			c.Abort()
			return
		}
		raw := strings.TrimSpace(h[7:])
		claims, err := jwtutil.Parse([]byte(cfg.JWTSecret), raw)
		if err != nil {
			response.ErrorMessage(c, http.StatusUnauthorized, "invalid or expired token")
			c.Abort()
			return
		}
		c.Set(ctxAdminID, claims.AdminID)
		c.Next()
	}
}

func AdminID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get(ctxAdminID)
	if !ok {
		return uuid.Nil, false
	}
	id, ok := v.(uuid.UUID)
	return id, ok
}
