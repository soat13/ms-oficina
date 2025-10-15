package http

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	app "github.com/soat13/fase-1-oficina/internal/service/application"
	"github.com/soat13/fase-1-oficina/internal/service/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Handler struct {
	createUseCase *app.CreateService
	updateUseCase *app.UpdateService
	deleteUseCase *app.DeleteService
	getUseCase    *app.GetService
	listUseCase   *app.ListServices
	validate      *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
}

func NewHandler(create *app.CreateService, update *app.UpdateService, delete *app.DeleteService, get *app.GetService, list *app.ListServices, errorHandler *fiberHelper.ErrorHandler) *Handler {
	validate := validator.New()
	handler := &Handler{
		createUseCase: create,
		updateUseCase: update,
		deleteUseCase: delete,
		getUseCase:    get,
		listUseCase:   list,
		validate:      validate,
		errorHandler:  errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrServiceNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateService)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidServiceName)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidServicePrice)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/services")
	grp.Post("/", h.create)
	grp.Patch("/:id", h.update)
	grp.Delete("/:id", h.delete)
	grp.Get("/:id", h.getByID)
	grp.Get("/", h.list)
}

// -------- DTOs --------

type createBody struct {
	Name  string `json:"name"  validate:"required,min=3"`
	Price int64  `json:"price" validate:"required,gt=0"`
}

type updateBody struct {
	Name  *string `json:"name"  validate:"omitempty,min=3"`
	Price *int64  `json:"price" validate:"omitempty,gt=0"`
}

type serviceJSON struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int64     `json:"price"`
}

// -------- Helpers JSON --------

func toJSON(v app.ServiceView) serviceJSON {
	return serviceJSON{
		ID:    v.ID,
		Name:  v.Name,
		Price: v.Price.Cents,
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

	price, err := money.New(body.Price)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.createUseCase.Execute(ctx.Context(), app.CreateInput{
		Name:  body.Name,
		Price: price,
		Now:   time.Now(),
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.SendStatus(fiber.StatusCreated)
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
	if err := h.validate.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	var priceVO *money.Money
	if body.Price != nil {
		vo, err := money.New(*body.Price)
		if err != nil {
			return h.errorHandler.Handle(ctx, err)
		}
		priceVO = &vo
	}

	err = h.updateUseCase.Execute(ctx.Context(), app.UpdateInput{
		ID:    id,
		Name:  body.Name,
		Price: priceVO,
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
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.Service))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listUseCase.Execute(ctx.Context(), app.ListInput{Pager: *pager})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	resp := maps.Map(out.Services, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}
