package vehicle

import (
	"sync"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/soat13/fase-1-oficina/internal/vehicle/application"
	"github.com/soat13/fase-1-oficina/internal/vehicle/infra/db"
	"github.com/soat13/fase-1-oficina/internal/vehicle/infra/http"
	"github.com/soat13/fase-1-oficina/tests/testsupport"
	"github.com/uptrace/bun"
)

type TestEnv struct {
	db  *bun.DB
	app *fiber.App
}

var (
	env  *TestEnv
	once sync.Once
)

func ensureSetup(t *testing.T) {
	once.Do(func() {
		testDB := testsupport.NewTestDB(t)
		env = &TestEnv{
			db: testDB.DB,
		}
		setup(t, env.db)
	})
}

func setup(t *testing.T, bunDB *bun.DB) {
	app := fiber.New()

	// wiring
	repo := db.NewVehicleRepository(bunDB)
	createUC := application.NewCreateVehicle(repo)
	updateUC := application.NewUpdateVehicle(repo)
	deleteUC := application.NewDeleteVehicle(repo)
	getUC := application.NewGetVehicle(repo)
	listUC := application.NewListVehicles(repo)

	handler := http.NewHandler(createUC, updateUC, deleteUC, getUC, listUC)
	http.Register(app, handler)

	env.app = app

	// clean up
	// t.Cleanup(func() {
	// 	err := testsupport.TruncateTables(bunDB)
	// 	require.NoError(t, err)
	// })
}
