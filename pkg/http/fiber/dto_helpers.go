package fiber

import (
	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/internal/shared/errors"
)

type ItemLinePayload struct {
	ID       string `json:"id"       validate:"required,uuid"`
	Quantity int    `json:"quantity" validate:"required,min=1"`
}

func ParseItemQuantities(raw []ItemLinePayload) (map[uuid.UUID]int, error) {
	out := make(map[uuid.UUID]int, len(raw))
	for _, it := range raw {
		id, err := uuid.Parse(it.ID)
		if err != nil {
			return nil, errors.ErrInvalidID
		}

		out[id] = it.Quantity
	}
	return out, nil
}
