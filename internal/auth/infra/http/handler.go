package http

import (
    "strings"

    "github.com/go-playground/validator/v10"
    "github.com/gofiber/fiber/v2"

    app "github.com/soat13/fase-1-oficina/internal/auth/application"
)

type Handler struct {
    login *app.LoginService
    v     *validator.Validate
}

func NewHandler(login *app.LoginService) *Handler {
    return &Handler{login: login, v: validator.New()}
}

func Register(app *fiber.App, h *Handler) {
    grp := app.Group("/auth")
    grp.Post("/login", h.Login)
}

type loginBody struct {
    Email    string `json:"email"    validate:"required,email"`
    Password string `json:"password" validate:"required,min=6"`
}

func (h *Handler) Login(ctx *fiber.Ctx) error {
    var body loginBody
    if err := ctx.BodyParser(&body); err != nil {
        return writeError(ctx, fiber.StatusBadRequest, "INVALID_JSON", "invalid JSON body")
    }
    if err := h.v.Struct(body); err != nil {
        return writeError(ctx, fiber.StatusUnprocessableEntity, "INVALID_BODY", err.Error())
    }
    out, err := h.login.Execute(ctx.Context(), app.LoginInput{Email: body.Email, Password: body.Password})
    if err != nil {
        if strings.Contains(err.Error(), "invalid credentials") {
            return writeError(ctx, fiber.StatusUnauthorized, "INVALID_CREDENTIALS", "invalid credentials")
        }
        return writeError(ctx, fiber.StatusInternalServerError, "INTERNAL_ERROR", "internal error")
    }
    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"token": out.Token})
}

