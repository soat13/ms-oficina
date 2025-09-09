package http

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	app "github.com/soat13/fase-1-oficina/internal/estimate/application"
)

var (
	ErrInvalidRepairOrderID = errors.New("invalid repair order id")
)

type ErrorInfo struct {
	Status int
	Code   string
}

var errorMap = map[error]ErrorInfo{
	ErrInvalidRepairOrderID: {
		Status: fiber.StatusBadRequest,
		Code:   "INVALID_REPAIR_ORDER_ID",
	},
	app.ErrRepairOrderNotFound: {
		Status: fiber.StatusNotFound,
		Code:   "REPAIR_ORDER_NOT_FOUND",
	},
	app.ErrInvalidRepairOrderStatus: {
		Status: fiber.StatusUnprocessableEntity,
		Code:   "INVALID_REPAIR_ORDER_STATUS",
	},
}
