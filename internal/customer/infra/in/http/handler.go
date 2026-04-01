package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	app "github.com/soat13/fase-1-oficina/internal/customer/application"
	"github.com/soat13/fase-1-oficina/internal/customer/domain"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	fiberHelper "github.com/soat13/oficina-utils/pkg/http/fiber"
	"github.com/soat13/oficina-utils/pkg/maps"
	"github.com/soat13/oficina-utils/pkg/valueobjects/document"
	"github.com/soat13/oficina-utils/pkg/valueobjects/email"
	"github.com/soat13/oficina-utils/pkg/valueobjects/phone"
)

type Handler struct {
	createUseCase *app.CreateCustomer
	updateUseCase *app.UpdateCustomer
	deleteUseCase *app.DeleteCustomer
	getUseCase    *app.GetCustomer
	listUseCase   *app.ListCustomers
	validator     *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
}

func NewHandler(
	create *app.CreateCustomer,
	update *app.UpdateCustomer,
	delete *app.DeleteCustomer,
	get *app.GetCustomer,
	list *app.ListCustomers,
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

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrCustomerNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateDocument)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateEmail)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidCustomerName)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(document.ErrInvalidDocument)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(phone.ErrInvalidPhoneNumber)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(email.ErrInvalidEmail)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/customers")
	grp.Post("/", h.create)
	grp.Patch("/:id", h.update)
	grp.Delete("/:id", h.delete)
	grp.Get("/:id", h.getByID)
	grp.Get("/", h.list)
}

// -------- DTOs --------

type createBody struct {
	Name        string `json:"name"         validate:"required,min=3"`
	Document    string `json:"document"     validate:"required,min=11"`
	Email       string `json:"email"        validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required,min=11"`
}

type updateBody struct {
	Name        *string `json:"name"         validate:"omitempty,min=3"`
	Email       *string `json:"email"        validate:"omitempty,email"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=11"`
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
		DocumentType: v.Document.Type(),
		PhoneNumber:  v.PhoneNumber.String(),
		Email:        v.Email.String(),
	}
}

// -------- Handlers --------

func (h *Handler) create(ctx *fiber.Ctx) error {
	var body createBody
	if err := ctx.BodyParser(&body); err != nil {
		return h.errorHandler.Handle(ctx, sharedErrors.ErrInvalidJSON)
	}
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
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

	output, err := h.createUseCase.Execute(ctx.UserContext(), app.CreateInput{
		Name:        body.Name,
		Document:    documentVO,
		PhoneNumber: phoneVO,
		Email:       emailVO,
	})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	return ctx.Status(fiber.StatusCreated).JSON(fiber.Map{"id": output.CustomerID})
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
	if err := h.validator.Struct(body); err != nil {
		return h.errorHandler.Handle(ctx, err)
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

	err = h.updateUseCase.Execute(ctx.UserContext(), app.UpdateInput{
		ID:          id,
		Name:        body.Name,
		PhoneNumber: phoneVO,
		Email:       emailVO,
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
	if err := h.deleteUseCase.Execute(ctx.UserContext(), app.DeleteInput{ID: id}); err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.SendStatus(fiber.StatusNoContent)
}

func (h *Handler) getByID(ctx *fiber.Ctx) error {
	id, err := fiberHelper.GetUuidParam(ctx, "id")
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	out, err := h.getUseCase.Execute(ctx.UserContext(), app.GetInput{ID: id})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.Customer))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listUseCase.Execute(ctx.UserContext(), app.ListInput{Pager: *pager})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	resp := maps.Map(out.Customers, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}
