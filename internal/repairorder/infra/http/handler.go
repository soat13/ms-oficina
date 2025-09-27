package http

import (
	"github.com/gofiber/fiber/v2"
	
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"

	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
)

type (
	Handler struct {
		startExecUseCase *application.StartExecution
		errorHandler     *fiberHelper.ErrorHandler
	}
)

func NewHandler(startExecUseCase *application.StartExecution, errorHandler *fiberHelper.ErrorHandler) *Handler {
	handler := &Handler{
		startExecUseCase: startExecUseCase,
		errorHandler:     errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(sharedRepairOrder.ErrRepairOrderNotFound)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Post("repair-orders/:id/start-execution", h.startExecution)
}

func (h *Handler) startExecution(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")

	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.StartExecutionInput{
		RepairOrderID: id,
	}

	if err := h.startExecUseCase.Execute(ctx.Context(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	ctx.Status(fiber.StatusNoContent)
	return nil
}
