package middleware

import (
	"context"
	"github.com/PavelKilko/NetPulse/database"
	"github.com/go-redis/redis/v8"
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"log"
	"os"
)

func JWTMiddleware() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   []byte(os.Getenv("JWT_SECRET")),
		ErrorHandler: jwtError,
		SuccessHandler: func(c *fiber.Ctx) error {
			// Extract the user from the context
			user := c.Locals("user").(*jwt.Token)
			claims := user.Claims.(jwt.MapClaims)

			// Log all claims for debugging purposes
			log.Printf("JWT Claims: %+v\n", claims)

			// Get JTI from claims and handle if it is missing or of the wrong type
			jtiVal, ok := claims["jti"]
			if !ok {
				log.Println("JWT does not contain jti")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized: missing jti",
				})
			}
			jti, ok := jtiVal.(string)
			if !ok {
				log.Println("Invalid jti format in JWT claims")
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Unauthorized: invalid jti format",
				})
			}

			// Check if the token has been revoked
			_, err := database.RedisClient.Get(context.Background(), jti).Result()
			if err == redis.Nil {
				// Token not found in blacklist, proceed
				return c.Next()
			} else if err != nil {
				// Redis error
				log.Println("Error connecting to Redis:", err)
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
					"error": "Failed to validate token",
				})
			} else {
				// Token found in blacklist, reject the request
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"error": "Token has been revoked",
				})
			}
		},
	})
}

func jwtError(c *fiber.Ctx, err error) error {
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Unauthorized",
		})
	}
	return nil
}
