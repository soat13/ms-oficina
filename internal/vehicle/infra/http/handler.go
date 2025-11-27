package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	app "github.com/soat13/fase-1-oficina/internal/vehicle/application"
	"github.com/soat13/fase-1-oficina/internal/vehicle/domain"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/plate"
)

type Handler struct {
	createUseCase         *app.CreateVehicle
	updateUseCase         *app.UpdateVehicle
	deleteUseCase         *app.DeleteVehicle
	getUseCase            *app.GetVehicle
	listUseCase           *app.ListVehicles
	listByCustomerUseCase *app.ListVehiclesByCustomer
	validate              *validator.Validate
	errorHandler          *fiberHelper.ErrorHandler
}

func NewHandler(
	create *app.CreateVehicle,
	update *app.UpdateVehicle,
	delete *app.DeleteVehicle,
	get *app.GetVehicle,
	list *app.ListVehicles,
	listByCustomer *app.ListVehiclesByCustomer,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	validate := validator.New()
	handler := &Handler{
		createUseCase:         create,
		updateUseCase:         update,
		deleteUseCase:         delete,
		getUseCase:            get,
		listUseCase:           list,
		listByCustomerUseCase: listByCustomer,
		validate:              validate,
		errorHandler:          errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrVehicleNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicatePlate)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidVehiclePlate)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidVehicleBrand)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidVehicleModel)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidVehicleYear)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(plate.ErrInvalidPlate)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/vehicles")
	grp.Post("/", h.create)
	grp.Patch("/:id", h.update)
	grp.Delete("/:id", h.delete)
	grp.Get("/:id", h.getByID)
	grp.Get("/", h.list)

	// Customer-specific vehicle routes
	customerGrp := app.Group("/admin/customers/:customerId/vehicles")
	customerGrp.Get("/", h.listByCustomer)
}

// -------- DTOs --------

type createBody struct {
	CustomerID uuid.UUID `json:"customer_id" validate:"required"`
	Plate      string    `json:"plate"       validate:"required,min=7,max=10"`
	Brand      string    `json:"brand"       validate:"required,min=2,max=100"`
	Model      string    `json:"model"       validate:"required,min=2,max=100"`
	Year       int       `json:"year"        validate:"required,min=1900"`
}

type updateBody struct {
	Plate *string `json:"plate" validate:"omitempty,min=7,max=10"`
	Brand *string `json:"brand" validate:"omitempty,min=2,max=100"`
	Model *string `json:"model" validate:"omitempty,min=2,max=100"`
	Year  *int    `json:"year"  validate:"omitempty,min=1900"`
}

type vehicleJSON struct {
	ID         uuid.UUID `json:"id"`
	CustomerID uuid.UUID `json:"customer_id"`
	Plate      string    `json:"plate"`
	Brand      string    `json:"brand"`
	Model      string    `json:"model"`
	Year       int       `json:"year"`
}

// -------- Helpers JSON --------

func toJSON(v app.VehicleView) vehicleJSON {
	return vehicleJSON{
		ID:         v.ID,
		CustomerID: v.CustomerID,
		Plate:      v.Plate.String(),
		Brand:      v.Brand,
		Model:      v.Model,
		Year:       v.Year,
	}
}

// -------- Handlers --------

func (h *Handler) create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	if err := h.validate.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	output, err := h.createUseCase.Execute(ctx.Context(), app.CreateInput{
		CustomerID: body.CustomerID,
		Plate:      body.Plate,
		Brand:      body.Brand,
		Model:      body.Model,
		Year:       body.Year,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": output.VehicleID})
}

func (h *Handler) update(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	var body updateBody

	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}

	var plateVo plate.Plate

	if body.Plate != nil {
		plateVo, err = plate.New(*body.Plate)
		if err != nil {
			return h.errorHandler.Handle(ctx, err)
		}
	}

	if err := h.validate.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.updateUseCase.Execute(ctx.Context(), app.UpdateInput{
		ID:    id,
		Plate: &plateVo,
		Brand: body.Brand,
		Model: body.Model,
		Year:  body.Year,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) delete(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	if err := h.deleteUseCase.Execute(ctx.Context(), app.DeleteInput{ID: id}); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) getByID(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	out, err := h.getUseCase.Execute(ctx.Context(), app.GetInput{ID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.Vehicle))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listUseCase.Execute(ctx.Context(), app.ListInput{Pager: *pager})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	resp := maps.Map(out.Vehicles, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}

func (h *Handler) listByCustomer(ctx *fiber.Ctx) error {
	customerID, err := fiberHelper.GetUuidParam(ctx, "customerId")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listByCustomerUseCase.Execute(ctx.Context(), app.ListByCustomerInput{
		CustomerID: customerID,
		Pager:      *pager,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	resp := maps.Map(out.Vehicles, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}
