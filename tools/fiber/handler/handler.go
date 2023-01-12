package handler

import (
	"driver/internal"
	"driver/internal/repository"
	"driver/pkg"
	"driver/tools/zap"
	"github.com/gofiber/fiber/v2"
	"strconv"
)

type IHandler interface {
	ErrorHandler(c *fiber.Ctx, err error) error
	GetNearestDriver(c *fiber.Ctx) error
	BulkCreateDrivers(c *fiber.Ctx) error
}

type Handler struct {
	Repository repository.IRepository
}

func NewHandler(repository repository.IRepository) IHandler {
	return &Handler{
		Repository: repository,
	}
}

func (h *Handler) ErrorHandler(c *fiber.Ctx, err error) error {
	// Status code defaults to 500
	code := fiber.StatusInternalServerError

	// Retrieve the custom status code if it's an fiber.*Error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	// Log error for debugging
	zap.Logger.Error(err.Error())

	// Send error back
	return c.Status(code).JSON(pkg.ErrorResponse{
		Error: pkg.ErrInternalServer.Error(),
	})
}

func (h *Handler) GetNearestDriver(c *fiber.Ctx) error {

	lat := c.Query("lat")
	long := c.Query("long")

	floatLat, err := strconv.ParseFloat(lat, 64)
	if err != nil {
		zap.Logger.Error("Error parsing latitude: " + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(pkg.ErrorResponse{Error: pkg.ErrInvalidRequest.Error()})
	}

	floatLong, err := strconv.ParseFloat(long, 64)
	if err != nil {
		zap.Logger.Error("Error parsing longitude: " + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(pkg.ErrorResponse{Error: pkg.ErrInvalidRequest.Error()})
	}

	driver, err := h.Repository.FindNearestDriver(internal.Coordinates{
		Latitude:  floatLat,
		Longitude: floatLong,
	})
	if err != nil {
		if err == pkg.ErrDriverNotFound {
			return c.Status(fiber.StatusNotFound).JSON(pkg.ErrorResponse{Error: pkg.ErrDriverNotFound.Error()})
		}

		zap.Logger.Error("Error finding nearest driver: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.ErrorResponse{Error: pkg.ErrInternalServer.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"id":          driver.ID,
		"distance":    driver.Distance,
		"coordinates": driver.Location.Coordinates,
	})
}

func (h *Handler) BulkCreateDrivers(c *fiber.Ctx) error {

	var drivers []internal.Model

	err := c.BodyParser(&drivers)
	if err != nil {
		zap.Logger.Error("Error parsing request body: " + err.Error())
		return c.Status(fiber.StatusBadRequest).JSON(pkg.ErrorResponse{Error: pkg.ErrInvalidRequest.Error()})
	}

	if len(drivers) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.ErrorResponse{Error: pkg.ErrInvalidRequest.Error()})
	}

	if len(drivers) > 100 {
		return c.Status(fiber.StatusBadRequest).JSON(pkg.ErrorResponse{Error: pkg.ErrBulkCreateDriversLimit.Error()})
	}

	err = h.Repository.BulkCreateDrivers(drivers)
	if err != nil {
		zap.Logger.Error("Error bulk updating drivers: " + err.Error())
		return c.Status(fiber.StatusInternalServerError).JSON(pkg.ErrorResponse{Error: pkg.ErrInternalServer.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
	})
}
