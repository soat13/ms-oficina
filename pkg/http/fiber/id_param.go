package fiber

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	sharedErrors "github.com/soat13/fase-1-oficina/internal/shared/errors"
)

func GetUuidParam(c *fiber.Ctx, name string) (uuid.UUID, error) {
	idStr := c.Params(name)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, sharedErrors.ErrInvalidID
	}

	return id, nil
}
