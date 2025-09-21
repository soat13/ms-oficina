package http

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	app "github.com/soat13/fase-1-oficina/internal/service/application"
	"github.com/soat13/fase-1-oficina/pkg/money"
)

type Handler struct {
	create        *app.CreateService
	update        *app.UpdateService
	del           *app.DeleteService
	get           *app.GetService
	list          *app.ListServices
	bodyValidator *validator.Validate
	now           func() time.Time
}

func NewHandler(create *app.CreateService, update *app.UpdateService, del *app.DeleteService, get *app.GetService, list *app.ListServices) *Handler {
	return &Handler{
		create:        create,
		update:        update,
		del:           del,
		get:           get,
		list:          list,
		bodyValidator: validator.New(validator.WithRequiredStructEnabled()),
		now:           time.Now,
	}
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/services")
	grp.Post("/", h.Create)
	grp.Put("/:id", h.Update)
	grp.Delete("/:id", h.Delete)
	grp.Get("/:id", h.GetByID)
	grp.Get("/", h.List)
}

type createBody struct {
	Name  string `json:"name"  validate:"required,min=3"`
	Price int64  `json:"price" validate:"required,gt=0"`
}

type updateBody struct {
	Name  *string `json:"name,omitempty"  validate:"omitempty,min=3"`
	Price *int64  `json:"price,omitempty" validate:"omitempty,gt=0"`
}

type serviceJSON struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	Price int64     `json:"price"`
}

type listJSON struct {
	Services []serviceJSON `json:"services"`
}

type getJSON struct {
	Service serviceJSON `json:"service"`
}

func (h *Handler) Create(c *fiber.Ctx) error {
	var body createBody
	if err := c.BodyParser(&body); err != nil {
		return writeError(c, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
	}

	if err := h.bodyValidator.Struct(body); err != nil {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
	}

	if strings.TrimSpace(body.Name) == "" || len(strings.TrimSpace(body.Name)) < 3 || body.Price <= 0 {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", "invalid name or price")
	}

	pr, err := money.New(body.Price)
	if err != nil {
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", "invalid service price")
	}

	out, err := h.create.Execute(c.Context(), app.CreateInput{
		Name:  body.Name,
		Price: pr,
		Now:   h.now(),
	})
	if err != nil {
		return h.handleError(c, err)
	}
	return h.json(c, fiber.StatusCreated, getJSON{Service: toJSON(out.Service)})
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
		return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", "invalid service price")
	}

	out, err := h.update.Execute(c.Context(), app.UpdateInput{
		ID:    id,
		Name:  body.Name,
		Price: pricePtr,
		Now:   h.now(),
	})
	if err != nil {
		return h.handleError(c, err)
	}
	return h.json(c, fiber.StatusOK, getJSON{Service: toJSON(out.Service)})
}

func (h *Handler) Delete(c *fiber.Ctx) error {
	id, err := h.parseIDParam(c, "id")
	if err != nil {
		return err
	}
	if err := h.del.Execute(c.Context(), app.DeleteInput{ID: id}); err != nil {
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
	return h.json(c, fiber.StatusOK, getJSON{Service: toJSON(out.Service)})
}

func (h *Handler) List(c *fiber.Ctx) error {
	limit := atoiDefault(c.Query("limit"), 50)
	offset := atoiDefault(c.Query("offset"), 0)

	out, err := h.list.Execute(c.Context(), app.ListInput{Limit: limit, Offset: offset})
	if err != nil {
		return h.handleError(c, err)
	}

	resp := make([]serviceJSON, 0, len(out.Services))
	for _, sv := range out.Services {
		resp = append(resp, toJSON(sv))
	}
	return h.json(c, fiber.StatusOK, listJSON{Services: resp})
}

func toJSON(v app.ServiceView) serviceJSON {
	return serviceJSON{
		ID:    v.ID,
		Name:  v.Name,
		Price: v.Price.Cents,
	}
}

func (h *Handler) json(c *fiber.Ctx, status int, v any) error {
	return c.Status(status).JSON(v)
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

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}

func (h *Handler) handleError(c *fiber.Ctx, err error) error {
	switch {
	case errors.Is(err, app.ErrServiceNotFound):
		return writeError(c, fiber.StatusNotFound, "SERVICE_NOT_FOUND", err.Error())
	case errors.Is(err, app.ErrDuplicateService):
		return writeError(c, fiber.StatusConflict, "SERVICE_ALREADY_EXISTS", err.Error())
	default:
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "invalid service name") ||
			strings.Contains(msg, "invalid service price") ||
			strings.Contains(msg, "invalid name") ||
			strings.Contains(msg, "invalid price") {
			return writeError(c, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
		}
		return writeError(c, fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
	}
}
