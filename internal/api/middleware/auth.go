package middleware

import (
	"dashboard-ac-backend/internal/api/response"
	"dashboard-ac-backend/internal/domain"
	"dashboard-ac-backend/pkg/jwt"

	"github.com/gofiber/fiber/v2"
)

func JWTAuth(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Get authorization header
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return response.Unauthorized(c, "Authorization header is required")
		}

		// Extract token from header
		tokenString, err := jwt.ExtractTokenFromHeader(authHeader)
		if err != nil {
			return response.Unauthorized(c, "Invalid authorization header format")
		}

		// Validate token
		claims, err := jwt.ValidateToken(tokenString, jwtSecret)
		if err != nil {
			return response.Unauthorized(c, "Invalid or expired token")
		}

		// Check if it's an access token
		if claims.Subject != "access" {
			return response.Unauthorized(c, "Invalid token type")
		}

		// Store user info in context
		c.Locals("user_id", claims.UserID)
		c.Locals("user_email", claims.Email)
		c.Locals("user_role", claims.Role)

		return c.Next()
	}
}

func RequireRole(roles ...domain.Role) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userRole, ok := c.Locals("user_role").(domain.Role)
		if !ok {
			return response.Unauthorized(c, "User role not found in context")
		}

		// Check if user has required role
		for _, role := range roles {
			if userRole == role {
				return c.Next()
			}
		}

		return response.Forbidden(c, "Insufficient permissions")
	}
}

func RequireAdmin() fiber.Handler {
	return RequireRole(domain.RoleAdmin)
}

func RequireAdminOrTechnician() fiber.Handler {
	return RequireRole(domain.RoleAdmin, domain.RoleTechnician)
}

func GetUserFromContext(c *fiber.Ctx) (uint, string, domain.Role, error) {
	userID, ok := c.Locals("user_id").(uint)
	if !ok {
		return 0, "", "", fiber.NewError(fiber.StatusUnauthorized, "User ID not found in context")
	}

	userEmail, ok := c.Locals("user_email").(string)
	if !ok {
		return 0, "", "", fiber.NewError(fiber.StatusUnauthorized, "User email not found in context")
	}

	userRole, ok := c.Locals("user_role").(domain.Role)
	if !ok {
		return 0, "", "", fiber.NewError(fiber.StatusUnauthorized, "User role not found in context")
	}

	return userID, userEmail, userRole, nil
}