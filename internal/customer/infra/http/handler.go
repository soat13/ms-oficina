package http

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	app "github.com/soat13/fase-1-oficina/internal/customer/application"
)

type Handler struct {
	create   *app.CreateCustomer
	update   *app.UpdateCustomer
	del      *app.DeleteCustomer
	get      *app.GetCustomer
	list     *app.ListCustomers
	validate *validator.Validate
}

func NewHandler(create *app.CreateCustomer, update *app.UpdateCustomer, del *app.DeleteCustomer, get *app.GetCustomer, list *app.ListCustomers) *Handler {
	validate := validator.New()
	return &Handler{
		create:   create,
		update:   update,
		del:      del,
		get:      get,
		list:     list,
		validate: validate,
	}
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/customers") // proteja com JWT no setup/main
	grp.Post("/", h.Create)
	grp.Put("/:id", h.Update)
	grp.Delete("/:id", h.Delete)
	grp.Get("/:id", h.GetByID)
	grp.Get("/", h.List)
}

// -------- DTOs --------

type createBody struct {
	Name      string `json:"name"        validate:"required,min=3"`
	Cellphone string `json:"cellphone"   validate:"required,min=11"`
	Document  string `json:"document"    validate:"required,min=11"`
}

type updateBody struct {
	Name      *string `json:"name"        validate:"omitempty,min=3"`
	Cellphone *string `json:"cellphone"   validate:"omitempty,min=11"`
}

type customerJSON struct {
	ID        uuid.UUID `json:"id"`
	Name      string    `json:"name"`
	Cellphone string    `json:"cellphone"`
	Document  string    `json:"document"`
}

// -------- Helpers JSON --------

func toJSON(v app.CustomerView) customerJSON {
	return customerJSON{
		ID:        v.ID,
		Name:      v.Name,
		Cellphone: v.Cellphone,
		Document:  v.Document,
	}
}

// -------- Handlers --------

func (h *Handler) Create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return writeError(ctx, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}
	if err := h.validate.Struct(body); err != nil {
		return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	out, err := h.create.Execute(ctx.Context(), app.CreateInput{
		Name:      body.Name,
		Cellphone: body.Cellphone,
		Document:  body.Document,
		Now:       time.Now(),
	})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"customer": toJSON(out.Customer)})
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
	if err := h.validate.Struct(body); err != nil {
		return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	out, err := h.update.Execute(ctx.Context(), app.UpdateInput{
		ID:        id,
		Name:      body.Name,
		Cellphone: body.Cellphone,
		Now:       time.Now(),
	})
	if err != nil {
		return h.handleError(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"customer": toJSON(out.Customer)})
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
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"customer": toJSON(out.Customer)})
}

func (h *Handler) List(ctx *fiber.Ctx) error {
	limit := atoiDefault(ctx.Query("limit"), 50)
	offset := atoiDefault(ctx.Query("offset"), 0)

	out, err := h.list.Execute(ctx.Context(), app.ListInput{Limit: limit, Offset: offset})
	if err != nil {
		return h.handleError(ctx, err)
	}
	resp := make([]customerJSON, 0, len(out.Customers))
	for _, sv := range out.Customers {
		resp = append(resp, toJSON(sv))
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"customers": resp})
}

// -------- Error mapping + helpers --------

func (h *Handler) handleError(ctx *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, app.ErrCustomerNotFound):
		return writeError(ctx, fiber.StatusNotFound, "CUSTOMER_NOT_FOUND", err.Error())
	case errors.Is(err, app.ErrDuplicateCustomer):
		return writeError(ctx, fiber.StatusConflict, "CUSTOMER_ALREADY_EXISTS", err.Error())
	default:
		msg := err.Error()
		if strings.Contains(msg, "invalid customer name") || strings.Contains(msg, "invalid customer cellphone") || strings.Contains(msg, "invalid customer document") {
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
