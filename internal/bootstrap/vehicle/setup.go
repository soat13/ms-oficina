package vehicle

import (
	"github.com/soat13/ms-oficina/internal/bootstrap"
	vehicleApp "github.com/soat13/ms-oficina/internal/vehicle/application"
	vehicleHTTP "github.com/soat13/ms-oficina/internal/vehicle/infra/in/http"
	vehicleDB "github.com/soat13/ms-oficina/internal/vehicle/infra/out/db"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository vehicleApp.VehicleRepository) {
	if repository == nil {
		repository = vehicleDB.NewBunVehicleRepository(container.DB)
	}

	createVehicle := vehicleApp.NewCreateVehicle(repository)
	updateVehicle := vehicleApp.NewUpdateVehicle(repository)
	deleteVehicle := vehicleApp.NewDeleteVehicle(repository)
	getVehicle := vehicleApp.NewGetVehicle(repository)
	listVehicles := vehicleApp.NewListVehicles(repository)
	listVehiclesByCustomer := vehicleApp.NewListVehiclesByCustomer(repository)

	vehicleHandler := vehicleHTTP.NewHandler(
		createVehicle,
		updateVehicle,
		deleteVehicle,
		getVehicle,
		listVehicles,
		listVehiclesByCustomer,
		container.FiberErrorHandler,
	)

	vehicleHTTP.Register(container.FiberApp, vehicleHandler)
}
