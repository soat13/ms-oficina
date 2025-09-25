package http

import (
	"errors"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	app "github.com/soat13/fase-1-oficina/internal/vehicle/application"
)

type Handler struct {
	create *app.CreateVehicle
	update *app.UpdateVehicle
	delete *app.DeleteVehicle
	get    *app.GetVehicle
	list   *app.ListVehicles
	v      *validator.Validate
}

func NewHandler(
	create *app.CreateVehicle,
	update *app.UpdateVehicle,
	delete *app.DeleteVehicle,
	get *app.GetVehicle,
	list *app.ListVehicles,
) *Handler {
	v := validator.New()
	return &Handler{
		create: create,
		update: update,
		delete: delete,
		get:    get,
		list:   list,
		v:      v,
	}
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/vehicle")
	grp.Post("/", h.Create)
	grp.Put("/:id", h.Update)
	grp.Delete("/:id", h.Delete)
	grp.Get("/:id", h.Get)
	grp.Get("/", h.List)
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
		if strings.Contains(msg, "invalid vehicle") { // Generic validation check
			return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", msg)
		}
		return writeError(ctx, fiber.StatusInternalServerError, "INTERNAL_ERROR", msg)
	}
}

// ------- DTOs ------------

type createBody struct {
	CustomerID string `json:"customer_id" validate:"required,uuid"`
	Plate      string `json:"plate"  validate:"required,min=7"`
	Model      string `json:"model"`
	Brand      string `json:"brand"`
	Year       int    `json:"year" validate:"required,numeric,gte=1886"`
}

type updateBody struct {
	Model string `json:"model"`
	Brand string `json:"brand"`
	Year  int    `json:"year" validate:"required,numeric,gte=1886"`
}

type vehicleJSON struct {
	ID         uuid.UUID `json:"id"`
	CustomerID string    `json:"customer_id"`
	Plate      string    `json:"plate"`
	Model      string    `json:"model"`
	Brand      string    `json:"brand"`
	Year       int       `json:"year"`
}

// --- Handlers ---

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
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"vehicle": toJSON(out.Vehicle)})
}

func (h *Handler) Update(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_ID", "invalid vehicle ID")
	}

	var body updateBody
	if err := ctx.BodyParser(&body); err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}
	if err := h.v.Struct(body); err != nil {
		return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	err = h.update.Execute(ctx.Context(), app.UpdateInput{
		ID:    id,
		Model: body.Model,
		Brand: body.Brand,
		Year:  body.Year,
	})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) Delete(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_ID", "invalid vehicle ID")
	}

	if err := h.delete.Execute(ctx.Context(), id); err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) Get(ctx *fiber.Ctx) error {
	id, err := uuid.Parse(ctx.Params("id"))
	if err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_ID", "invalid vehicle ID")
	}

	vehicle, err := h.get.Execute(ctx.Context(), id)
	if err != nil {
		return h.handleError(ctx, err)
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"vehicle": toJSON(*vehicle)})
}

func (h *Handler) List(ctx *fiber.Ctx) error {
	limit, _ := strconv.Atoi(ctx.Query("limit", "10"))
	offset, _ := strconv.Atoi(ctx.Query("offset", "0"))

	vehicles, err := h.list.Execute(ctx.Context(), app.ListInput{Limit: limit, Offset: offset})
	if err != nil {
		return h.handleError(ctx, err)
	}

	var result []vehicleJSON
	for _, v := range vehicles {
		result = append(result, toJSON(v))
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"vehicles": result})
}
