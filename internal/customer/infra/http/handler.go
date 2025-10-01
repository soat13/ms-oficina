package http

import (
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	app "github.com/soat13/fase-1-oficina/internal/customer/application"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type Handler struct {
	createUseCase *app.CreateCustomer
	updateUseCase *app.UpdateCustomer
	deleteUseCase *app.DeleteCustomer
	getUseCase    *app.GetCustomer
	listUseCase   *app.ListCustomers
	validate      *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
}

func NewHandler(create *app.CreateCustomer, update *app.UpdateCustomer, delete *app.DeleteCustomer, get *app.GetCustomer, list *app.ListCustomers, errorHandler *fiberHelper.ErrorHandler) *Handler {
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

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrCustomerNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateDocument)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateEmail)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(document.ErrInvalidDocument)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(phone.ErrInvalidPhoneNumber)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(email.ErrInvalidEmail)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/customers") // TODO: PROTECT WITH JWT
	grp.Post("/", h.create)
	grp.Patch("/:id", h.update)
	grp.Delete("/:id", h.delete)
	grp.Get("/:id", h.getByID)
	grp.Get("/", h.list)
}

// -------- DTOs --------

type createBody struct {
	Name        string `json:"name"           validate:"required,min=3"`
	Document    string `json:"document"       validate:"required,min=11"`
	Email       string `json:"email"          validate:"required,email"`
	PhoneNumber string `json:"phone_number"   validate:"required,min=11"`
}

type updateBody struct {
	Name        *string `json:"name"           validate:"omitempty,min=3"`
	Email       *string `json:"email"          validate:"omitempty"`
	PhoneNumber *string `json:"phone_number"   validate:"omitempty,min=11"`
}

type customerJSON struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
}

// -------- Helpers JSON --------

func toJSON(v app.CustomerView) customerJSON {
	return customerJSON{
		ID:           v.ID,
		Name:         v.Name,
		Document:     v.Document.Value,
		DocumentType: v.Document.TypeString(),
		PhoneNumber:  v.PhoneNumber.String(),
		Email:        v.Email.String(),
	}
}

// -------- Handlers --------

func (h *Handler) create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_JSON",
			"message": "invalid JSON body",
		})
	}
	if err := h.validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": err.Error(),
		})
	}

	documentVO, err := document.New(body.Document)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	emailVO, err := email.New(body.Email)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	phoneVO, err := phone.New(body.PhoneNumber)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.createUseCase.Execute(ctx.Context(), app.CreateInput{
		Name:        body.Name,
		Document:    documentVO,
		PhoneNumber: phoneVO,
		Email:       emailVO,
		Now:         time.Now(),
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
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"code":    "INVALID_JSON",
			"message": "invalid JSON body",
		})
	}
	if err := h.validate.Struct(body); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"code":    "INVALID_BODY",
			"message": err.Error(),
		})
	}

	var emailVO *email.Email
	if body.Email != nil {
		vo, err := email.New(*body.Email)
		if err != nil {
			return h.errorHandler.Handle(ctx, err)
		}
		emailVO = &vo
	}

	var phoneVO *phone.PhoneNumber
	if body.PhoneNumber != nil {
		vo, err := phone.New(*body.PhoneNumber)
		if err != nil {
			return h.errorHandler.Handle(ctx, err)
		}
		phoneVO = &vo
	}

	err = h.updateUseCase.Execute(ctx.Context(), app.UpdateInput{
		ID:          id,
		Name:        body.Name,
		PhoneNumber: phoneVO,
		Email:       emailVO,
		Now:         time.Now(),
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
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.Customer))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	limit := atoiDefault(ctx.Query("limit"), 50)
	offset := atoiDefault(ctx.Query("offset"), 0)

	out, err := h.listUseCase.Execute(ctx.Context(), app.ListInput{Limit: limit, Offset: offset})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	resp := make([]customerJSON, 0, len(out.Customers))
	for _, sv := range out.Customers {
		resp = append(resp, toJSON(sv))
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}

// -------- Helpers --------

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	if n, err := strconv.Atoi(s); err == nil {
		return n
	}
	return def
}
