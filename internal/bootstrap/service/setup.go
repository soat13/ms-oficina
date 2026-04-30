package service

import (
	"github.com/soat13/ms-oficina/internal/bootstrap"
	serviceApp "github.com/soat13/ms-oficina/internal/service/application"
	serviceHTTP "github.com/soat13/ms-oficina/internal/service/infra/in/http"
	serviceDB "github.com/soat13/ms-oficina/internal/service/infra/out/db"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository serviceApp.Repository) {
	if repository == nil {
		repository = serviceDB.NewBunServiceRepository(container.DB)
	}

	createService := serviceApp.NewCreateService(repository)
	updateService := serviceApp.NewUpdateService(repository)
	deleteService := serviceApp.NewDeleteService(repository)
	getService := serviceApp.NewGetService(repository)
	listService := serviceApp.NewListServices(repository)

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
