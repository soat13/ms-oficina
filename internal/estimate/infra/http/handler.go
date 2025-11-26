package http

import (
	"github.com/gofiber/fiber/v2"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"

	app "github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/domain"
)

type (
	Handler struct {
		approveUseCase    *app.Approve
		rejectUseCase     *app.Reject
		removeItemUseCase *app.RemoveItem
		errorHandler      *fiberHelper.ErrorHandler
	}
)

func NewHandler(
	approveUseCase *app.Approve,
	rejectUseCase *app.Reject,
	removeItemUseCase *app.RemoveItem,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrEstimateNotFound)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(app.ErrProductNotAvailable)
	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(domain.ErrItemNotFound)
	errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrCannotChangeItemsAfterApproval)

	return &Handler{
		approveUseCase:    approveUseCase,
		rejectUseCase:     rejectUseCase,
		removeItemUseCase: removeItemUseCase,
		errorHandler:      errorHandler,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Post("repair-orders/:id/estimate/approve", h.approve)
	group.Post("repair-orders/:id/estimate/reject", h.reject)
	group.Delete("repair-orders/:id/estimate/items/:itemId", h.removeItem)
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

func (h *Handler) removeItem(ctx *fiber.Ctx) error {
	repairOrderID, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	itemID, err := fiberHelper.GetUuidParam(ctx, "itemId")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.removeItemUseCase.Execute(ctx.Context(), app.RemoveItemInput{
		RepairOrderID: repairOrderID,
		ItemID:        itemID,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}
