package http

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"

	app "github.com/soat13/fase-1-oficina/internal/estimate/application"
)

type (
	Handler struct {
		createUseCase  *app.Create
		approveUseCase *app.Approve
		rejectUseCase  *app.Reject
		validate       *validator.Validate
		errorHandler   *fiberHelper.ErrorHandler
	}

	linePayload struct {
		ID       string `json:"id"       validate:"required,uuid"`
		Quantity int    `json:"quantity" validate:"required,min=1"`
	}

	createEstimateBody struct {
		Products []linePayload `json:"products" validate:"dive"`
		Services []linePayload `json:"services" validate:"dive"`
	}
)

func NewHandler(
	createUseCase *app.Create,
	approveUseCase *app.Approve,
	rejectUseCase *app.Reject,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	bodyValidator := validator.New()
	bodyValidator.RegisterStructValidation(func(structLevel validator.StructLevel) {
		body := structLevel.Current().Interface().(createEstimateBody)
		if len(body.Products) == 0 && len(body.Services) == 0 {
			structLevel.ReportError(body.Products, "products", "Products", "at_least_one", "")
		}
	}, createEstimateBody{})

	errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrEstimateNotFound)

	return &Handler{
		approveUseCase: approveUseCase,
		rejectUseCase:  rejectUseCase,
		createUseCase:  createUseCase,
		validate:       bodyValidator,
		errorHandler:   errorHandler,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Post("repair-orders/:id/estimate", h.create)
	group.Post("estimates/:id/approve", h.approve)
	group.Post("estimates/:id/reject", h.reject)
}

func (h *Handler) create(ctx *fiber.Ctx) error {
	repairOrderID, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	body, err := h.getBody(ctx)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	inputDTO, err := h.getCreateInput(repairOrderID, *body)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.createUseCase.Execute(ctx.Context(), *inputDTO)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *Handler) approve(ctx *fiber.Ctx) error {

	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.approveUseCase.Execute(ctx.Context(), app.ApproveInput{ID: id})
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

	err = h.rejectUseCase.Execute(ctx.Context(), app.RejectInput{ID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) getBody(ctx *fiber.Ctx) (*createEstimateBody, error) {
	var body createEstimateBody
	err := ctx.BodyParser(&body)

	if err != nil {
		return nil, err
	}

	err = h.validate.Struct(body)
	if err != nil {
		return nil, err
	}

	return &body, err
}

func (h *Handler) getCreateInput(repairOrderID uuid.UUID, body createEstimateBody) (*app.CreateInput, error) {
	productQty, err := parseItemQuantities(body.Products)
	if err != nil {
		return nil, err
	}
	serviceQty, err := parseItemQuantities(body.Services)
	if err != nil {
		return nil, err
	}

	return &app.CreateInput{
		Now:           time.Now(),
		RepairOrderID: repairOrderID,
		Products:      productQty,
		Services:      serviceQty,
	}, nil
}

func parseItemQuantities(raw []linePayload) (map[uuid.UUID]int, error) {
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
