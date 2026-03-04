package http

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	estimateApp "github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	sharedRepairOrder "github.com/soat13/fase-1-oficina/internal/shared/repairorder"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
)

type (
	Handler struct {
		listUseCase                    *application.ListRepairOrders
		getUseCase                     *application.GetRepairOrder
		startExecutionUseCase          *application.StartExecution
		startDiagnosticsUseCase        *application.StartDiagnostics
		finishDiagnosticsUseCase       *application.FinishDiagnostics
		finishExecutionUseCase         *application.FinishExecution
		releaseVehicleUseCase          *application.ReleaseVehicle
		createUseCase                  *application.Create
		getAverageExecutionTimeUseCase *application.GetAverageExecutionTime
		cancelUseCase                  *application.Cancel
		validator                      *validator.Validate
		errorHandler                   *fiberHelper.ErrorHandler
	}
)

type ItemLinePayload struct {
	ID       string `json:"id"       validate:"required,uuid"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

func NewHandler(
	listUseCase *application.ListRepairOrders,
	getUseCase *application.GetRepairOrder,
	createUseCase *application.Create,
	startDiagnosticsUseCase *application.StartDiagnostics,
	finishDiagnosticsUseCase *application.FinishDiagnostics,
	startExecutionUseCase *application.StartExecution,
	finishExecutionUseCase *application.FinishExecution,
	releaseVehicleUseCase *application.ReleaseVehicle,
	getAverageExecutionTimeUseCase *application.GetAverageExecutionTime,
	cancelUseCase *application.Cancel,
	validator *validator.Validate,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(sharedRepairOrder.ErrRepairOrderNotFound)
	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(application.ErrProductNotFound)
	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(application.ErrServiceNotFound)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(application.ErrVehicleOrCustomerNotFound)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(application.ErrInsufficientStock)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(application.ErrInvalidQuantity)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(estimateApp.ErrProductNotAvailable)

	return &Handler{
		listUseCase:                    listUseCase,
		getUseCase:                     getUseCase,
		startDiagnosticsUseCase:        startDiagnosticsUseCase,
		finishDiagnosticsUseCase:       finishDiagnosticsUseCase,
		startExecutionUseCase:          startExecutionUseCase,
		finishExecutionUseCase:         finishExecutionUseCase,
		releaseVehicleUseCase:          releaseVehicleUseCase,
		createUseCase:                  createUseCase,
		getAverageExecutionTimeUseCase: getAverageExecutionTimeUseCase,
		cancelUseCase:                  cancelUseCase,
		validator:                      validator,
		errorHandler:                   errorHandler,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Get("repair-orders", h.list)
	group.Post("repair-orders", h.create)
	group.Get("repair-orders/average-execution-time", h.getAverageExecutionTime)
	group.Get("repair-orders/:id", h.getByID)
	group.Post("repair-orders/:id/cancel", h.cancel)
	group.Post("repair-orders/:id/start-diagnostics", h.startDiagnostics)
	group.Post("repair-orders/:id/finish-diagnostics", h.finishDiagnostics)
	group.Post("repair-orders/:id/start-execution", h.startExecution)
	group.Post("repair-orders/:id/finish-execution", h.finishExecution)
	group.Post("repair-orders/:id/release-vehicle", h.releaseVehicle)
}

type (
	repairOrderJSON struct {
		ID         uuid.UUID `json:"id"`
		CustomerID uuid.UUID `json:"customer_id"`
		VehicleID  uuid.UUID `json:"vehicle_id"`
		Status     string    `json:"status"`
		CreatedAt  time.Time `json:"created_at"`
		UpdatedAt  time.Time `json:"updated_at"`
	}

	finishDiagnosticsBody struct {
		Products []ItemLinePayload `json:"products" validate:"required,min=1,dive"`
		Services []ItemLinePayload `json:"services" validate:"required,min=1,dive"`
	}
)

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

func (h *Handler) create(ctx *fiber.Ctx) error {
	var req struct {
		CustomerID uuid.UUID `json:"customer_id"`
		VehicleID  uuid.UUID `json:"vehicle_id"`
	}

	if err := ctx.BodyParser(&req); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.CreateInput{
		CustomerID: req.CustomerID,
		VehicleID:  req.VehicleID,
	}

	output, err := h.createUseCase.Execute(ctx.UserContext(), input)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": output.RepairOrderID})
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	output, err := h.listUseCase.Execute(ctx.UserContext(), application.ListInput{Pager: *pager})
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

	output, err := h.getUseCase.Execute(ctx.UserContext(), application.GetInput{ID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(toJSON(output.RepairOrder))
}

func (h *Handler) startDiagnostics(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.StartDiagnosticsInput{
		RepairOrderID: id,
	}

	if err := h.startDiagnosticsUseCase.Execute(ctx.UserContext(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) finishDiagnostics(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	var body finishDiagnosticsBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input, err := getFinishDiagnosticsInput(id, body)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	if err := h.finishDiagnosticsUseCase.Execute(ctx.UserContext(), *input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) startExecution(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")

	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.StartExecutionInput{
		RepairOrderID: id,
	}

	if err := h.startExecutionUseCase.Execute(ctx.UserContext(), input); err != nil {
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

	if err := h.finishExecutionUseCase.Execute(ctx.UserContext(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) releaseVehicle(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")

	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.ReleaseVehicleInput{
		RepairOrderID: id,
	}

	if err := h.releaseVehicleUseCase.Execute(ctx.UserContext(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) cancel(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")

	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	input := application.CancelInput{
		RepairOrderID: id,
	}

	if err := h.cancelUseCase.Execute(ctx.UserContext(), input); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) getAverageExecutionTime(ctx *fiber.Ctx) error {
	output, err := h.getAverageExecutionTimeUseCase.Execute(ctx.UserContext())
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(output)
}

func getFinishDiagnosticsInput(repairOrderID uuid.UUID, body finishDiagnosticsBody) (*application.FinishDiagnosticsInput, error) {
	productsQty, err := parseItemQuantities(body.Products)
	if err != nil {
		return nil, err
	}

	servicesQty, err := parseItemQuantities(body.Services)
	if err != nil {
		return nil, err
	}

	return &application.FinishDiagnosticsInput{
		RepairOrderID: repairOrderID,
		Products:      productsQty,
		Services:      servicesQty,
	}, nil
}

func parseItemQuantities(raw []ItemLinePayload) (map[uuid.UUID]int, error) {
	out := make(map[uuid.UUID]int, len(raw))
	for _, it := range raw {
		id, err := uuid.Parse(it.ID)
		if err != nil {
			return nil, sharedErrors.ErrInvalidID
		}

		out[id] = it.Quantity
	}
	return out, nil
}
