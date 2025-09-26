package http

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/pkg/money"

	app "github.com/soat13/fase-1-oficina/internal/product/application"
)

type Handler struct {
	create        *app.CreateProduct
	update        *app.UpdateProduct
	del           *app.DeleteProduct
	get           *app.GetProduct
	list          *app.ListProducts
	bodyValitator *validator.Validate
}

func NewHandler(create *app.CreateProduct, update *app.UpdateProduct, del *app.DeleteProduct, get *app.GetProduct, list *app.ListProducts) *Handler {
	bodyValitator := validator.New()
	return &Handler{
		create:        create,
		update:        update,
		del:           del,
		get:           get,
		list:          list,
		bodyValitator: bodyValitator,
	}
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/products") // proteja com JWT no setup/main
	grp.Post("/", h.Create)
	grp.Put("/:id", h.Update)
	grp.Delete("/:id", h.Delete)
	grp.Get("/:id", h.GetByID)
	grp.Get("/", h.List)
}

// -------- DTOs --------

type createBody struct {
	Name       string `json:"name"        validate:"required,min=3"`
	PriceCents int64  `json:"price_cents" validate:"required,gt=0"`
	Stock      int    `json:"stock"       validate:"omitempty,gte=0"`
}

type updateBody struct {
	Name       *string `json:"name"        validate:"omitempty,min=3"`
	PriceCents *int64  `json:"price_cents" validate:"omitempty,gt=0"`
	Stock      *int    `json:"stock"       validate:"omitempty,gte=0"`
}

type productJSON struct {
	ID         uuid.UUID `json:"id"`
	Name       string    `json:"name"`
	PriceCents int64     `json:"price_cents"`
	Stock      int       `json:"stock"`
}

// -------- Helpers JSON --------

func toJSON(v app.ProductView) productJSON {
	return productJSON{
		ID:         v.ID,
		Name:       v.Name,
		PriceCents: v.PriceCents,
		Stock:      v.Stock,
	}
}

// -------- Handlers --------

func (h *Handler) Create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}
	if err := h.bodyValitator.Struct(body); err != nil {
		return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	price, err := money.New(body.PriceCents)
	if err != nil {
		return err
	}

	out, err := h.create.Execute(ctx.Context(), app.CreateInput{
		Name:  strings.TrimSpace(body.Name),
		Price: price,
		Stock: body.Stock,
		Now:   time.Now(),
	})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"product": toJSON(out.Product)})
}

func (h *Handler) Update(ctx *fiber.Ctx) error {
	id, err := parseID(ctx.Params("id"))
	if err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_ID", "invalid id")
	}

	var body updateBody
	if err := ctx.BodyParser(&body); err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}
	if err := h.bodyValitator.Struct(body); err != nil {
		return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	out, err := h.update.Execute(ctx.Context(), app.UpdateInput{
		ID:         id,
		Name:       body.Name,
		PriceCents: body.PriceCents,
		Stock:      body.Stock,
		Now:        time.Now(),
	})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"product": toJSON(out.Product)})
}

func (h *Handler) Delete(ctx *fiber.Ctx) error {
	id, err := parseID(ctx.Params("id"))
	if err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_ID", "invalid id")
	}
	if err := h.del.Execute(ctx.Context(), app.DeleteInput{ID: id}); err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) GetByID(ctx *fiber.Ctx) error {
	id, err := parseID(ctx.Params("id"))
	if err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_ID", "invalid id")
	}
	out, err := h.get.Execute(ctx.Context(), app.GetInput{ID: id})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"product": toJSON(out.Product)})
}

func (h *Handler) List(ctx *fiber.Ctx) error {
	limit := atoiDefault(ctx.Query("limit"), 50)
	offset := atoiDefault(ctx.Query("offset"), 0)

	out, err := h.list.Execute(ctx.Context(), app.ListInput{Limit: limit, Offset: offset})
	if err != nil {
		return h.handleError(ctx, err)
	}
	resp := make([]productJSON, 0, len(out.Products))
	for _, pv := range out.Products {
		resp = append(resp, toJSON(pv))
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"products": resp})
}

// -------- Error mapping + helpers --------

func (h *Handler) handleError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, app.ErrProductNotFound):
		return writeError(ctx, fiber.StatusNotFound, "PRODUCT_NOT_FOUND", err.Error())
	case errors.Is(err, app.ErrDuplicateProduct):
		return writeError(ctx, fiber.StatusConflict, "PRODUCT_ALREADY_EXISTS", err.Error())
	default:
		// validações de domínio
		msg := err.Error()
		if strings.Contains(msg, "invalid product name") || strings.Contains(msg, "invalid product price") || strings.Contains(msg, "invalid product stock") {
			return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", msg)
		}
		return writeError(ctx, fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
	}
}

func parseID(s string) (uuid.UUID, error) {
	return uuid.Parse(s)
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
