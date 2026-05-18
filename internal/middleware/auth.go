package middleware

import (
	"log"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/gofiber/fiber/v2"
)

// Auth valida el token JWT de Supabase en cada request.
// Si SUPABASE_JWT_SECRET no está configurado, omite la validación (modo dev).
func Auth() fiber.Handler {
	jwtSecret := os.Getenv("SUPABASE_JWT_SECRET")
	if jwtSecret == "" {
		log.Println("⚠️  SUPABASE_JWT_SECRET no configurado — validación JWT desactivada (modo desarrollo)")
	}

	return func(c *fiber.Ctx) error {
		if jwtSecret == "" {
			return c.Next()
		}

		authHeader := c.Get("Authorization")
		if !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token no proporcionado",
			})
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fiber.ErrUnauthorized
			}
			return []byte(jwtSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"error":   true,
				"message": "Token inválido o expirado",
			})
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Locals("userID", claims["sub"])
		}
		return c.Next()
	}
}
