package estimate

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	estimateApp "github.com/soat13/fase-1-oficina/internal/estimate/application"
	estimaterListeners "github.com/soat13/fase-1-oficina/internal/estimate/application/listeners"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHttp "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
	repairOrderEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/repairorder/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {

	repairOrderReader := estimateInfraDB.NewRepairOrderReader(container.DB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(container.DB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(container.DB)
	repository := estimateInfraDB.NewBunRepository(container.DB)

	createEstimate := estimateApp.NewCreateEstimate(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		repository,
		container.EventBus,
	)

	approve := estimateApp.NewApproveEstimate(repository, container.EventBus)
	cancel := estimateApp.NewCancelEstimate(repository, container.EventBus)
	reject := estimateApp.NewRejectEstimate(repository, container.EventBus)

	estimateHttpHandler := estimateInfraHttp.NewHandler(createEstimate, approve, reject, container.FiberErrorHandler)
	estimateInfraHttp.Register(container.FiberApp, estimateHttpHandler)

	container.EventBus.Subscribe(
		productEvent.StockReduceConfirmed{}.Topic(),
		estimaterListeners.OnStockReduceConfirmed(repository, container.EventBus),
	)

	container.EventBus.Subscribe(
		repairOrderEvent.Canceled{}.Topic(),
		estimaterListeners.OnRepairOrderCanceled(cancel, repository),
	)

}
