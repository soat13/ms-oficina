package http

import (
	"github.com/gofiber/fiber/v2"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"

	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
)

type (
	Handler struct {
		startExecUseCase       *application.StartExecution
		finishExecutionUseCase *application.FinishExecution
		ReleaseVehicleUseCase  *application.ReleaseVehicle
		errorHandler           *fiberHelper.ErrorHandler
	}
)

func NewHandler(
	startExecUseCase *application.StartExecution,
	finishExecutionUseCase *application.FinishExecution,
	ReleaseVehicleUseCase *application.ReleaseVehicle,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	return &Handler{
		startExecUseCase:       startExecUseCase,
		finishExecutionUseCase: finishExecutionUseCase,
		ReleaseVehicleUseCase:  ReleaseVehicleUseCase,
		errorHandler:           errorHandler,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Post("repair-orders/:id/start-execution", h.startExecution)
	group.Post("repair-orders/:id/finish-execution", h.finishExecution)
	group.Post("repair-orders/:id/release-vehicle", h.ReleaseVehicle)
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

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) finishExecution(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")

	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.FinishExecutionInput{
		RepairOrderID: id,
	}

	if err := h.finishExecutionUseCase.Execute(ctx.Context(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) ReleaseVehicle(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")

	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.ReleaseVehicleInput{
		RepairOrderID: id,
	}

	if err := h.ReleaseVehicleUseCase.Execute(ctx.Context(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}
