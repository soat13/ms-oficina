package product

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	productApp "github.com/soat13/fase-1-oficina/internal/product/application"
	productDB "github.com/soat13/fase-1-oficina/internal/product/infra/db"
	productHTTP "github.com/soat13/fase-1-oficina/internal/product/infra/http"
	"github.com/soat13/fase-1-oficina/internal/product/infra/listeners"
	estimateEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository productApp.ProductRepository) {
	if repository == nil {
		repository = productDB.NewBunProductRepository(container.DB)
	}

	createProduct := productApp.NewCreateProduct(repository)
	updateProduct := productApp.NewUpdateProduct(repository)
	deleteProduct := productApp.NewDeleteProduct(repository)
	getProduct := productApp.NewGetProduct(repository)
	listProducts := productApp.NewListProducts(repository)
	reduceStock := productApp.NewReduceStock(repository, container.EventBus)

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
		estimateEvent.Approved{}.Topic(),
		listeners.OnEstimateApproved(reduceStock, container.EventBus),
	)
}
