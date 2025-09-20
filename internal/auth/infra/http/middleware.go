package http

import (
    "strings"

    "github.com/gofiber/fiber/v2"
    app "github.com/soat13/fase-1-oficina/internal/auth/application"
)

func JWTMiddleware(tokens app.TokenService) fiber.Handler {
    return func(c *fiber.Ctx) error {
        auth := c.Get("Authorization")
        if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
            return writeError(c, fiber.StatusUnauthorized, "MISSING_TOKEN", "missing bearer token")
        }
        tok := strings.TrimPrefix(auth, "Bearer ")
        claims, err := tokens.ParseToken(strings.TrimSpace(tok))
        if err != nil {
            return writeError(c, fiber.StatusUnauthorized, "INVALID_TOKEN", "invalid token")
        }
        // store user info for downstream handlers if needed
        c.Locals("user_id", claims.UserID.String())
        c.Locals("user_email", claims.Email)
        return c.Next()
    }
}

