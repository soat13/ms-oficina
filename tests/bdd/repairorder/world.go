package repairorder

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
)

type (
	world struct {
		setup         *testsupport.SetupConfig
		productIDs    map[string]uuid.UUID
		serviceIDs    map[string]uuid.UUID
		repairOrderID uuid.UUID
		initialStatus string
		lastResp      *http.Response
	}

	worldCtxKey struct{}
)

func newWorld(setup *testsupport.SetupConfig) *world {
	return &world{
		setup:      setup,
		productIDs: make(map[string]uuid.UUID),
		serviceIDs: make(map[string]uuid.UUID),
	}
}

func withWorld(ctx context.Context, w *world) context.Context {
	return context.WithValue(ctx, worldCtxKey{}, w)
}

func worldFromContext(ctx context.Context) *world {
	w, _ := ctx.Value(worldCtxKey{}).(*world)
	return w
}
