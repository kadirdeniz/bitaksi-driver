package middleware

import (
	"driver/pkg"
	zap_tools "driver/tools/zap"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func IsApiKeyCorrect(c *fiber.Ctx) error {
	apiKey := c.Query("api_key")
	if apiKey == pkg.AppConfigs.GetApplicationConfig().API_KEY {
		return c.Next()
	}

	return c.Status(fiber.StatusUnauthorized).JSON(pkg.ErrorResponse{Error: pkg.ErrInvalidAPIKey.Error()})
}

func SetContentType(c *fiber.Ctx) error {
	c.Set("Content-Type", "application/json")
	return c.Next()
}

func Logger(c *fiber.Ctx) error {

	request := map[string]string{
		"method": c.Method(),
		"url":    c.OriginalURL(),
		"body":   string(c.Body()),
		"ip":     c.IP(),
		"host":   c.Hostname(),
	}

	zap_tools.Logger.Info("Request", zap.Any("Request Object", request))

	return c.Next()
}
