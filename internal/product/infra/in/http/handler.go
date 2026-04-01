package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/product/domain"
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/oficina-utils/pkg/http/fiber"
	"github.com/soat13/oficina-utils/pkg/maps"
	"github.com/soat13/oficina-utils/pkg/money"
)

type Handler struct {
	createUseCase *application.CreateProduct
	updateUseCase *application.UpdateProduct
	deleteUseCase *application.DeleteProduct
	getUseCase    *application.GetProduct
	listUseCase   *application.ListProducts
	validator     *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
}

func NewHandler(
	create *application.CreateProduct,
	update *application.UpdateProduct,
	delete *application.DeleteProduct,
	get *application.GetProduct,
	list *application.ListProducts,
	validator *validator.Validate,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	handler := &Handler{
		createUseCase: create,
		updateUseCase: update,
		deleteUseCase: delete,
		getUseCase:    get,
		listUseCase:   list,
		validator:     validator,
		errorHandler:  errorHandler,
	}

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(application.ErrProductNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(application.ErrDuplicateProduct)
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

func toJSON(v application.ProductView) productJSON {
	return productJSON{
		ID:    v.ID,
		Name:  v.Name,
		Price: v.Price.Cents,
		Stock: v.Stock,
	}
}

func (h *Handler) create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, errors.ErrInvalidJSON)
	}
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	price, err := money.New(body.Price)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	output, err := h.createUseCase.Execute(ctx.UserContext(), application.CreateInput{
		Name:  body.Name,
		Price: price,
		Stock: body.Stock,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": output.ProductID})
}

func (h *Handler) update(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	var body updateBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, errors.ErrInvalidJSON)
	}
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	pricePtr, err := moneyPtr(body.Price)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.updateUseCase.Execute(ctx.UserContext(), application.UpdateInput{
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
	if err := h.deleteUseCase.Execute(ctx.UserContext(), application.DeleteInput{ID: id}); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) getByID(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	out, err := h.getUseCase.Execute(ctx.UserContext(), application.GetInput{ID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.Product))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listUseCase.Execute(ctx.UserContext(), application.ListInput{Pager: *pager})
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
