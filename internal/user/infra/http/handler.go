package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
	app "github.com/soat13/fase-1-oficina/internal/user/application"
	"github.com/soat13/fase-1-oficina/internal/user/domain"
	"github.com/soat13/fase-1-oficina/internal/user/domain/role"
	fiberHelper "github.com/soat13/fase-1-oficina/pkg/http/fiber"
	"github.com/soat13/fase-1-oficina/pkg/maps"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/document"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/email"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/password"
	"github.com/soat13/fase-1-oficina/pkg/valueobjects/phone"
)

type Handler struct {
	createUseCase *app.CreateUser
	updateUseCase *app.UpdateUser
	deleteUseCase *app.DeleteUser
	getUseCase    *app.GetUser
	listUseCase   *app.ListUsers
	validate      *validator.Validate
	errorHandler  *fiberHelper.ErrorHandler
}

func NewHandler(create *app.CreateUser, update *app.UpdateUser, delete *app.DeleteUser, get *app.GetUser, list *app.ListUsers, errorHandler *fiberHelper.ErrorHandler) *Handler {
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

	handler.errorHandler.ErrorResolver.RegisterHTTPNotFoundError(app.ErrUserNotFound)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateEmail)
	handler.errorHandler.ErrorResolver.RegisterHTTPConflictError(app.ErrDuplicateDocument)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(domain.ErrInvalidUserName)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(role.ErrInvalidRole)
	handler.errorHandler.ErrorResolver.RegisterHTTPUnprocessableError(role.ErrRolesRequired)

	return handler
}

func Register(app *fiber.App, h *Handler) {
	grp := app.Group("/admin/users")
	grp.Post("/", h.create)
	grp.Patch("/:id", h.update)
	grp.Delete("/:id", h.delete)
	grp.Get("/:id", h.getByID)
	grp.Get("/", h.list)
}

// -------- DTOs --------

type createBody struct {
	Name        string   `json:"name"         validate:"required,min=3"`
	Document    string   `json:"document"     validate:"required,min=11"`
	Email       string   `json:"email"        validate:"required,email"`
	PhoneNumber string   `json:"phone_number" validate:"required,min=11"`
	Password    string   `json:"password"     validate:"required,min=8"`
	Roles       []string `json:"roles"        validate:"required,min=1,dive,required"`
}

type updateBody struct {
	Name        *string   `json:"name"         validate:"omitempty,min=3"`
	Email       *string   `json:"email"        validate:"omitempty,email"`
	PhoneNumber *string   `json:"phone_number" validate:"omitempty,min=11"`
	Password    *string   `json:"password"     validate:"omitempty,min=8"`
	Roles       *[]string `json:"roles"        validate:"omitempty,min=1,dive,required"`
}

type userJSON struct {
	ID           uuid.UUID `json:"id"`
	Name         string    `json:"name"`
	Document     string    `json:"document"`
	DocumentType string    `json:"document_type"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	Roles        []string  `json:"roles"`
}

// -------- Helpers JSON --------

func toJSON(v app.UserView) userJSON {
	return userJSON{
		ID:           v.ID,
		Name:         v.Name,
		Document:     v.Document.Value,
		DocumentType: v.Document.Type(),
		Email:        v.Email.String(),
		PhoneNumber:  v.PhoneNumber.String(),
		Roles:        v.Roles.Strings(),
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
	passwordVO, err := password.New(body.Password)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	rolesVO, err := role.NewRoles(body.Roles)
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}

	err = h.createUseCase.Execute(ctx.Context(), app.CreateInput{
		Name:        body.Name,
		Document:    documentVO,
		PhoneNumber: phoneVO,
		Email:       emailVO,
		Password:    passwordVO,
		Roles:       rolesVO,
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

	var passwordVO *password.Password
	if body.Password != nil {
		vo, err := password.New(*body.Password)
		if err != nil {
			return h.errorHandler.Handle(ctx, err)
		}
		passwordVO = &vo
	}

	var rolesVO *role.Roles
	if body.Roles != nil {
		vo, err := role.NewRoles(*body.Roles)
		if err != nil {
			return h.errorHandler.Handle(ctx, err)
		}
		rolesVO = &vo
	}

	err = h.updateUseCase.Execute(ctx.Context(), app.UpdateInput{
		ID:          id,
		Name:        body.Name,
		PhoneNumber: phoneVO,
		Email:       emailVO,
		Password:    passwordVO,
		Roles:       rolesVO,
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
	return ctx.Status(fiber.StatusOK).JSON(toJSON(out.User))
}

func (h *Handler) list(ctx *fiber.Ctx) error {
	pager := fiberHelper.NewPagination(ctx, 50, 0)

	out, err := h.listUseCase.Execute(ctx.Context(), app.ListInput{Pager: *pager})
	if err != nil {
		return h.errorHandler.Handle(ctx, err)
	}
	resp := maps.Map(out.Users, toJSON)
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"data": resp})
}
