package http

import (
	"errors"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"

	app "github.com/soat13/fase-1-oficina/internal/estimate/application"
)

type (
	Handler struct {
		createUseCase  *app.Create
		approveUseCase *app.Approve
		validate       *validator.Validate
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

var (
	ErrInvalidItemID = errors.New("invalid item id")
)

func NewHandler(createUseCase *app.Create, approveUseCase *app.Approve) *Handler {
	bodyValidator := validator.New()
	bodyValidator.RegisterStructValidation(func(structLevel validator.StructLevel) {
		body := structLevel.Current().Interface().(createEstimateBody)
		if len(body.Products) == 0 && len(body.Services) == 0 {
			structLevel.ReportError(body.Products, "products", "Products", "at_least_one", "")
		}
	}, createEstimateBody{})

	return &Handler{
		approveUseCase: approveUseCase,
		createUseCase:  createUseCase,
		validate:       bodyValidator,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("admin")
	group.Post("repair-orders/:id/estimate", h.create)
	group.Post("estimates/:id/approve", h.approve)
}

func (h *Handler) create(ctx *fiber.Ctx) error {

	inputDTO, err := h.getCreateItems(ctx)

	if err != nil {
		return h.handleError(ctx, err)
	}

	_, err = h.createUseCase.Execute(ctx.Context(), inputDTO)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusCreated)
}

func (h *Handler) approve(ctx *fiber.Ctx) error {

	id, err := h.getIDParam(ctx)
	if err != nil {
		return h.handleError(ctx, err)
	}

	err = h.approveUseCase.Execute(ctx.Context(), app.ApproveInput{ID: id})
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) getIDParam(c *fiber.Ctx) (uuid.UUID, error) {
	idStr := c.Params("id")
	repairOrder, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, ErrInvalidRepairOrderID
	}
	return repairOrder, nil
}

func (h *Handler) getCreateItems(ctx *fiber.Ctx) (app.CreateInput, error) {
	inputDTO := app.CreateInput{}

	repairOrderID, err := h.getIDParam(ctx)
	if err != nil {
		return inputDTO, err
	}

	var body createEstimateBody
	if err := ctx.BodyParser(&body); err != nil {
		return inputDTO, err
	}

	productQty, err := parseItemQuantities(body.Products)
	if err != nil {
		return inputDTO, err
	}
	serviceQty, err := parseItemQuantities(body.Services)
	if err != nil {
		return inputDTO, err
	}

	inputDTO.RepairOrderID = repairOrderID
	inputDTO.Now = time.Now()
	inputDTO.Products = productQty
	inputDTO.Services = serviceQty

	return inputDTO, err
}

func (h *Handler) handleError(ctx *fiber.Ctx, err error) error {
	if info, ok := errorMap[err]; ok {
		return ctx.Status(info.Status).JSON(fiber.Map{
			"code":    info.Code,
			"message": err.Error(),
		})
	}

	log.Error().
		Err(err).
		Str("path", ctx.Path()).
		Str("method", ctx.Method()).
		Msg("unexpected error on request")

	return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"code":    "INTERNAL_ERROR",
		"message": "internal error",
	})
}

func parseItemQuantities(raw []linePayload) (map[uuid.UUID]int, error) {
	out := make(map[uuid.UUID]int, len(raw))
	for _, it := range raw {
		id, err := uuid.Parse(it.ID)
		if err != nil {
			return nil, ErrInvalidItemID
		}

		out[id] = it.Quantity
	}
	return out, nil
}
