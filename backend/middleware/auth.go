package middleware

import (
	"reonify/clocked/config"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

// Protected protect routes
func Protected() fiber.Handler {
	return jwtware.New(jwtware.Config{
		SigningKey:   jwtware.SigningKey{Key: []byte(config.Config("SECRET"))},
		ErrorHandler: jwtError,
	})
}

func jwtError(c *fiber.Ctx, err error) error {
    if err.Error() == "Missing or malformed JWT" {
        return fiber.NewError(fiber.StatusBadRequest, "Missing or malformed JWT")
    }
    return fiber.NewError(fiber.StatusUnauthorized, "Invalid or expired JWT")
}
