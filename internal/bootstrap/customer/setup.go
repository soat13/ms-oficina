package customer

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	customerApp "github.com/soat13/fase-1-oficina/internal/customer/application"
	customerDB "github.com/soat13/fase-1-oficina/internal/customer/infra/db"
	customerHTTP "github.com/soat13/fase-1-oficina/internal/customer/infra/http"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository customerApp.CustomerRepository) {
	if repository == nil {
		repository = customerDB.NewBunCustomerRepository(container.DB)
	}

	createCustomer := customerApp.NewCreateCustomer(repository)
	updateCustomer := customerApp.NewUpdateCustomer(repository)
	deleteCustomer := customerApp.NewDeleteCustomer(repository)
	getCustomer := customerApp.NewGetCustomer(repository)
	listCustomers := customerApp.NewListCustomers(repository)

	customerHandler := customerHTTP.NewHandler(
		createCustomer,
		updateCustomer,
		deleteCustomer,
		getCustomer,
		listCustomers,
		container.FiberErrorHandler,
	)

	customerHTTP.Register(container.FiberApp, customerHandler)
}
