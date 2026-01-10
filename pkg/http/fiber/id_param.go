package fiber

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

var ErrInvalidID = errors.New("invalid ID")

func GetUuidParam(c *fiber.Ctx, name string) (uuid.UUID, error) {
	idStr := c.Params(name)
	id, err := uuid.Parse(idStr)
	if err != nil {
		return uuid.Nil, ErrInvalidID
	}

	return id, nil
}
