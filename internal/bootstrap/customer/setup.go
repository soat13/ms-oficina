package customer

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/customer/application"
	"github.com/soat13/fase-1-oficina/internal/customer/infra/in/http"
	"github.com/soat13/fase-1-oficina/internal/customer/infra/out/db"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository application.CustomerRepository) {
	if repository == nil {
		repository = db.NewBunCustomerRepository(container.DB)
	}

	createCustomer := application.NewCreateCustomer(repository)
	updateCustomer := application.NewUpdateCustomer(repository)
	deleteCustomer := application.NewDeleteCustomer(repository)
	getCustomer := application.NewGetCustomer(repository)
	listCustomers := application.NewListCustomers(repository)

	customerHandler := http.NewHandler(
		createCustomer,
		updateCustomer,
		deleteCustomer,
		getCustomer,
		listCustomers,
		container.Validator,
		container.FiberErrorHandler,
	)

	http.Register(container.FiberApp, customerHandler)
}
