package service

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	serviceApp "github.com/soat13/fase-1-oficina/internal/service/application"
	serviceDB "github.com/soat13/fase-1-oficina/internal/service/infra/db"
	serviceHTTP "github.com/soat13/fase-1-oficina/internal/service/infra/http"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository serviceApp.Repository) {

	if repository == nil {
		repository = serviceDB.NewBunRepository(container.DB)
	}

	serviceRepository := serviceDB.NewBunRepository(container.DB)
	createService := serviceApp.NewCreateService(serviceRepository)
	updateService := serviceApp.NewUpdateService(serviceRepository)
	deleteService := serviceApp.NewDeleteService(serviceRepository)
	getService := serviceApp.NewGetService(serviceRepository)
	listService := serviceApp.NewListServices(serviceRepository)

	serviceHttpHandler := serviceHTTP.NewHandler(
		createService,
		updateService,
		deleteService,
		getService,
		listService,
		container.FiberErrorHandler,
	)

	serviceHTTP.Register(container.FiberApp, serviceHttpHandler)
}
