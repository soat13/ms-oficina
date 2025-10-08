package http

import (
	"errors"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	app "github.com/soat13/fase-1-oficina/internal/product/application"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Handler struct {
	create        *app.CreateProduct
	update        *app.UpdateProduct
	delete        *app.DeleteProduct
	get           *app.GetProduct
	list          *app.ListProducts
	bodyValidator *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
	now           func() time.Time
}

func NewHandler(
	create *app.CreateProduct,
	update *app.UpdateProduct,
	del *app.DeleteProduct,
	get *app.GetProduct,
	list *app.ListProducts,
	errorHandler *fiberHelper.ErrorHandler,
) *Handler {
	bodyValidator := validator.New(validator.WithRequiredStructEnabled())

	return &Handler{
		create:        create,
		update:        update,
		delete:        del,
		get:           get,
		list:          list,
		bodyValidator: bodyValidator,
		errorHandler:  errorHandler,
		now:           time.Now,
	}
}

func Register(app *fiber.App, h *Handler) {
	group := app.Group("/admin/products")
	group.Post("/", h.Create)
	group.Put("/:id", h.Update)
	group.Delete("/:id", h.Delete)
	group.Get("/:id", h.GetByID)
	group.Get("/", h.List)
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

type listJSON struct {
	Products []productJSON `json:"products"`
}

type getJSON struct {
	Product productJSON `json:"product"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var body createBody

	if err := c.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(c, sharedErrors.ErrInvalidJSON)
	}

	if err := h.bodyValidator.Struct(body); err != nil {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	if len(strings.TrimSpace(body.Name)) < 3 || body.Price <= 0 || body.Stock < 0 {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", "invalid name, price or stock")
	}

	price, err := money.New(body.Price)
	if err != nil {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", "invalid product price")
	}

	out, err := h.create.Execute(c.Context(), app.CreateInput{
		Name:  body.Name,
		Price: price,
		Stock: body.Stock,
		Now:   h.now(),
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return h.json(c, fiber.StatusCreated, getJSON{Product: toJSON(out.Product)})
}

func (h *Handler) Update(c *fiber.Ctx) error {
	id, err := h.parseIDParam(c, "id")
	if err != nil {
		return err
	}

	var body updateBody
	if err := h.bindAndValidate(c, &body); err != nil {
		return err
	}

	pricePtr, err := moneyPtr(body.Price)
	if err != nil {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", "invalid product price")
	}

	out, err := h.update.Execute(c.Context(), app.UpdateInput{
		ID:    id,
		Name:  body.Name,
		Price: pricePtr,
		Stock: body.Stock,
		Now:   h.now(),
	})
	if err != nil {
		return h.handleError(c, err)
	}

	return h.json(c, fiber.StatusOK, getJSON{Product: toJSON(out.Product)})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := h.parseIDParam(c, "id")
	if err != nil {
		return err
	}

	if err := h.delete.Execute(c.Context(), app.DeleteInput{ID: id}); err != nil {
		return h.handleError(c, err)
	}

	return c.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetByID(c *fiber.Ctx) error {
	id, err := h.parseIDParam(c, "id")
	if err != nil {
		return err
	}
	out, err := h.get.Execute(c.Context(), app.GetInput{ID: id})
	if err != nil {
		return h.handleError(c, err)
	}
	return h.json(c, fiber.StatusOK, getJSON{Product: toJSON(out.Product)})
}

func (h *Handler) List(c *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(c, 50, 0)

	out, err := h.list.Execute(c.Context(), app.ListInput{Pager: *pager})
	if err != nil {
		return h.handleError(c, err)
	}

	response := listJSON{
		Products: maps.Map(out.Products, toJSON),
	}

	return h.json(c, fiber.StatusOK, response)
}

func (h *Handler) bindAndValidate(c *fiber.Ctx, dst any) error {
	if err := c.BodyParser(dst); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}
	if err := h.bodyValidator.Struct(dst); err != nil {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}
	return nil
}

func (h *Handler) parseIDParam(c *fiber.Ctx, name string) (uuid.UUID, error) {
	id, err := uuid.Parse(c.Params(name))
	if err != nil {
		return uuid.Nil, writeError(c, fiber.StatusBadRequest, "INVALID_ID", "invalid id")
	}
	return id, nil
}

func (h *Handler) json(c *fiber.Ctx, status int, value any) error {
	return c.Status(status).JSON(value)
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, app.ErrProductNotFound):
		return writeError(c, fiber.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
	case errors.Is(err, app.ErrDuplicateProduct):
		return writeError(c, fiber.StatusConflict, "PRODUCT_ALREADY_EXISTS", err.Error())
	default:
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "invalid product name") ||
			strings.Contains(msg, "invalid product price") ||
			strings.Contains(msg, "invalid product stock") {
			return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
	}
}

func toJSON(v app.ProductView) productJSON {
	return productJSON{
		ID:    v.ID,
		Name:  v.Name,
		Price: v.Price.Cents,
		Stock: v.Stock,
	}
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

func writeError(ctx *fiber.Ctx, status int, code string, msg string) error {
	return ctx.Status(status).JSON(fiber.Map{
		"code":    code,
		"message": msg,
	})
}
