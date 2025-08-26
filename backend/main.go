package main

import (
	"errors"
	stdlog "log"

	"reonify/clocked/database"
	"reonify/clocked/router"

	"github.com/gofiber/fiber/v2"
	fiberlog "github.com/gofiber/fiber/v2/log"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	app := fiber.New(fiber.Config{
		// Central error handler: log details and return JSON
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			// Log error with request context
			fiberlog.Errorf("%d %s %s - %v", code, c.Method(), c.OriginalURL(), err)

			// Consistent JSON error response
			return c.Status(code).JSON(fiber.Map{
				"status":  "error",
				"message": err.Error(),
				"data":    nil,
			})
		},
	})

	// Recover from panics and turn them into errors handled above
	app.Use(recover.New())
	app.Use(cors.New())

	database.ConnectDB()

	router.SetupRoutes(app)
	stdlog.Fatal(app.Listen(":8000"))
}
