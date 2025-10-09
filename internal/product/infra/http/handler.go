package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	app "github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Handler struct {
	createUseCase *app.CreateProduct
	updateUseCase *app.UpdateProduct
	deleteUseCase *app.DeleteProduct
	getUseCase    *app.GetProduct
	listUseCase   *app.ListProducts
	validate      *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
}

func NewHandler(create *app.CreateProduct, update *app.UpdateProduct, del *app.DeleteProduct, get *app.GetProduct, list *app.ListProducts, errorHandler *fiberHelper.ErrorHandler) *Handler {
	validate := validator.New()
	handler := &Handler{
		createUseCase: create,
		updateUseCase: update,
		deleteUseCase: del,
		getUseCase:    get,
		listUseCase:   list,
		validate:      validate,
		errorHandler:  errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrProductNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateProduct)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidProductName)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidProductPrice)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidProductStock)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/products")
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
	Stock int    `json:"stock" validate:"omitempty,gte=0"`
}

type updateBody struct {
	Name  *string `json:"name,omitempty"  validate:"omitempty,min=3"`
	Price *int64  `json:"price,omitempty" validate:"omitempty,gt=0"`
	Stock *int    `json:"stock,omitempty" validate:"omitempty,gte=0"`
}

type productJSON struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int64     `json:"price"`
	Stock int       `json:"stock"`
}

// -------- Helpers JSON --------

func toJSON(v app.ProductView) productJSON {
	return productJSON{
		ID:    v.ID,
		Name:  v.Name,
		Price: v.Price.Cents,
		Stock: v.Stock,
	}
}

// -------- Handlers --------

func (h *Handler) create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	if err := h.validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": err.Error(),
		})
	}

	price, err := money.New(body.Price)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.createUseCase.Execute(ctx.Context(), app.CreateInput{
		Name:  body.Name,
		Price: price,
		Stock: body.Stock,
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
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": err.Error(),
		})
	}

	pricePtr, err := moneyPtr(body.Price)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.updateUseCase.Execute(ctx.Context(), app.UpdateInput{
		ID:    id,
		Name:  body.Name,
		Price: pricePtr,
		Stock: body.Stock,
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
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.Product))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listUseCase.Execute(ctx.Context(), app.ListInput{Pager: *pager})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	resp := maps.Map(out.Products, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}

func moneyPtr(cents *int64) (*money.Money, error) {
	if cents == nil {
		return nil, nil
	}
	m, err := money.New(*cents)
	if err != nil {
		return nil, err
	}
	return &m, nil
}
