package product

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/product/application"
	productApp "github.com/soat13/fase-1-oficina/internal/product/application"
	"github.com/soat13/fase-1-oficina/internal/product/infra/in/event"
	productHTTP "github.com/soat13/fase-1-oficina/internal/product/infra/in/http"
	productDB "github.com/soat13/fase-1-oficina/internal/product/infra/out/db"
	eventOut "github.com/soat13/fase-1-oficina/internal/product/infra/out/event"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository productApp.Repository) {
	if repository == nil {
		repository = productDB.NewBunProductRepository(container.DB)
	}

	eventPublisher := eventOut.NewEventPublisher(container.EventBus)
	createProduct := application.NewCreateProduct(repository)
	updateProduct := application.NewUpdateProduct(repository)
	deleteProduct := application.NewDeleteProduct(repository)
	getProduct := application.NewGetProduct(repository)
	listProducts := application.NewListProducts(repository)
	reduceStock := application.NewReduceStock(repository, eventPublisher)
	handleEstimateApproved := application.NewHandleEstimateApproved(reduceStock)

	productHandler := productHTTP.NewHandler(
		createProduct,
		updateProduct,
		deleteProduct,
		getProduct,
		listProducts,
		container.Validator,
		container.FiberErrorHandler,
	)

	productHTTP.Register(container.FiberApp, productHandler)

	container.EventBus.Subscribe(
		estimateEvent.EstimateApproved{}.Topic(),
		event.OnEstimateApproved(handleEstimateApproved),
	)
}
