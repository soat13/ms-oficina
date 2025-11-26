package http

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"

	app "github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
)

type (
	addItemBody struct {
		ItemID   string `json:"item_id"       validate:"required,uuid"`
		Quantity int    `json:"quantity" validate:"required,min=1"`
	}

	Handler struct {
		approveUseCase    *app.Approve
		rejectUseCase     *app.Reject
		addItemUseCase    *app.AddItem
		removeItemUseCase *app.RemoveItem
		errorHandler      *fiberHelper.ErrorHandler
	}
)

func NewHandler(
	approveUseCase *app.Approve,
	rejectUseCase *app.Reject,
	addItemUseCase *app.AddItem,
	removeItemUseCase *app.RemoveItem,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrEstimateNotFound)
	errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrProductNotAvailable)
	errorHandler.ErrorResolver.RegisterHTTPConflictError(domain.ErrCannotChangeItemsAfterApproval)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(app.ErrProductOrServiceNotFound)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrRepairIDInvalid)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrQuantityInvalid)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrItemNotFound)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrEstimateMustHaveAtLeastOneItem)

	return &Handler{
		approveUseCase:    approveUseCase,
		rejectUseCase:     rejectUseCase,
		addItemUseCase:    addItemUseCase,
		removeItemUseCase: removeItemUseCase,
		errorHandler:      errorHandler,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Post("repair-orders/:id/estimate/approve", h.approve)
	group.Post("repair-orders/:id/estimate/reject", h.reject)
	group.Post("estimates/:id/items", h.addItem)
	group.Delete("estimates/:id/items/:itemId", h.removeItem)
}

func (h *Handler) approve(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.approveUseCase.Execute(ctx.Context(), app.ApproveInput{RepairOrderID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) reject(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.rejectUseCase.Execute(ctx.Context(), app.RejectInput{RepairOrderID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) addItem(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	var body addItemBody
	if ctx.BodyParser(&body) != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	itemID, err := uuid.Parse(body.ItemID)
	if err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidID)
	}

	err = h.addItemUseCase.Execute(ctx.Context(), app.AddItemInput{
		EstimateID: id,
		ItemID:     itemID,
		Quantity:   body.Quantity,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) removeItem(ctx *fiber.Ctx) error {
	estimateID, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	itemID, err := fiberHelper.GetUuidParam(ctx, "itemId")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.removeItemUseCase.Execute(ctx.Context(), app.RemoveItemInput{
		EstimateID: estimateID,
		ItemID:     itemID,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
