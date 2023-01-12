package fiber

import (
	"driver/internal/repository"
	"driver/tools/fiber/handler"
	"driver/tools/fiber/middleware"
	"driver/tools/zap"
	"fmt"
	swagger "github.com/arsmn/fiber-swagger/v2"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func Router(port int) {
	err := StartServer(8000)
	if err != nil {
		zap.Logger.Fatal(err.Error())
	}
}

func StartServer(port int) error {

	// Create repository
	repository := repository.NewRepository()

	if err := repository.Migration(); err != nil {
		return err
	}

	// Create handler
	handler := handler.NewHandler(repository)

	app := fiber.New(
		fiber.Config{
			ErrorHandler: handler.ErrorHandler,
		},
	)

	// Cors
	app.Use(cors.New())

	// Swagger
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Recovery
	app.Use(recover.New())

	api := app.Group("/api/v1", middleware.Logger, middleware.SetContentType)

	// Create routes
	api.Get("/drivers/nearest", middleware.IsApiKeyCorrect, handler.GetNearestDriver)
	// Bulk Update
	api.Post("/drivers", middleware.IsApiKeyCorrect, handler.BulkCreateDrivers)

	// Health check
	api.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).SendString("OK")
	})

	// Monitoring
	api.Get("/metrics", monitor.New(monitor.Config{Title: "Driver Service Metrics Page"}))

	return app.Listen(fmt.Sprintf(":%d", port))
}
