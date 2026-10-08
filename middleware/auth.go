package middleware

import (
	"bamboocians/utils"
	"strings"

	"github.com/gofiber/fiber/v2"
)

// JWTSecret is set at startup from config
var JWTSecret string

// Auth validates the Bearer token and stores claims in Locals
func Auth() fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			return utils.Fail(c, fiber.StatusUnauthorized, "missing or invalid authorization header")
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ParseToken(tokenStr, JWTSecret)
		if err != nil {
			return utils.Fail(c, fiber.StatusUnauthorized, "invalid or expired token")
		}

		c.Locals("userID", claims.UserID)
		c.Locals("email", claims.Email)
		c.Locals("role", claims.Role)
		return c.Next()
	}
}

// RequireRole returns 403 if the caller's role is not in the allowed list
func RequireRole(roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		role, _ := c.Locals("role").(string)
		for _, r := range roles {
			if role == r {
				return c.Next()
			}
		}
		return utils.Fail(c, fiber.StatusForbidden, "access denied: insufficient role")
	}
}
