package http

import (
	"errors"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	app "github.com/soat13/fase-1-oficina/internal/vehicle/application"
)

type Handler struct {
	create *app.CreateVehicle
	v      *validator.Validate
}

func NewHandler(create *app.CreateVehicle) *Handler {
	v := validator.New()
	return &Handler{
		create: create,
		v:      v,
	}
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/vehicle")
	grp.Post("/", h.Create)
}

// -------------Helpers -----------------------

func toJSON(v app.VehicleView) vehicleJSON {
	return vehicleJSON{
		ID:         v.ID,
		CustomerID: v.CustomerId,
		Plate:      v.Plate,
		Model:      v.Model,
		Brand:      v.Brand,
		Year:       v.Year,
	}
}

// -------- Error mapping + helpers --------

func (h *Handler) handleError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, app.ErrVehicleNotFound):
		return writeError(ctx, fiber.StatusNotFound, "VEHICLE_NOT_FOUND", err.Error())
	case errors.Is(err, app.ErrDuplicateVehicle):
		return writeError(ctx, fiber.StatusConflict, "VEHICLE_ALREADY_EXISTS", err.Error())
	default:
		// validações de domínio
		msg := err.Error()
		if strings.Contains(msg, "invalid vehicle plate") {
			return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", msg)
		}
		return writeError(ctx, fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
	}
}

// ------- DTOs ------------

type createBody struct {
	CustomerID string `json:"customer_id"`
	Plate      string `json:"plate"  validate:"required,min=7"`
	Model      string `json:"model"`
	Brand      string `json:"brand"`
	Year       int    `json:"year" validate:"required,min=4"`
}

type vehicleJSON struct {
	ID         uuid.UUID
	CustomerID string
	Plate      string
	Model      string
	Brand      string
	Year       int
}

func (h *Handler) Create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}
	if err := h.v.Struct(body); err != nil {
		return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	out, err := h.create.Execute(ctx.Context(), app.CreateInput{
		CustomerId: body.CustomerID,
		Plate:      body.Plate,
		Model:      body.Model,
		Brand:      body.Brand,
		Year:       body.Year,
	})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"service": toJSON(out.Vehicle)})
}
