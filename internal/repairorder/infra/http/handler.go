package http

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/soat13/fase-1-oficina/internal/repairorder/application"

	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type (
	Handler struct {
		listUseCase            *application.ListRepairOrders
		getUseCase             *application.GetRepairOrder
		startExecUseCase       *application.StartExecution
		finishExecutionUseCase *application.FinishExecution
		ReleaseVehicleUseCase  *application.ReleaseVehicle
		errorHandler           *fiberHelper.ErrorHandler
	}
)

func NewHandler(
	listUseCase *application.ListRepairOrders,
	getUseCase *application.GetRepairOrder,
	startExecUseCase *application.StartExecution,
	finishExecutionUseCase *application.FinishExecution,
	ReleaseVehicleUseCase *application.ReleaseVehicle,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	return &Handler{
		listUseCase:            listUseCase,
		getUseCase:             getUseCase,
		startExecUseCase:       startExecUseCase,
		finishExecutionUseCase: finishExecutionUseCase,
		ReleaseVehicleUseCase:  ReleaseVehicleUseCase,
		errorHandler:           errorHandler,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Get("repair-orders", h.list)
	group.Get("repair-orders/:id", h.getByID)
	group.Post("repair-orders/:id/start-execution", h.startExecution)
	group.Post("repair-orders/:id/finish-execution", h.finishExecution)
	group.Post("repair-orders/:id/release-vehicle", h.ReleaseVehicle)
}

type repairOrderJSON struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	VehicleID  uuid.UUID `json:"vehicle_id"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func toJSON(view application.RepairOrderView) repairOrderJSON {
	return repairOrderJSON{
		ID:         view.ID,
		CustomerID: view.CustomerID,
		VehicleID:  view.VehicleID,
		Status:     string(view.Status),
		CreatedAt:  view.CreatedAt,
		UpdatedAt:  view.UpdatedAt,
	}
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	output, err := h.listUseCase.Execute(ctx.Context(), application.ListInput{Pager: *pager})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	data := maps.Map(output.RepairOrders, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": data})
}

func (h *Handler) getByID(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	output, err := h.getUseCase.Execute(ctx.Context(), application.GetInput{ID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(toJSON(output.RepairOrder))
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
