package estimate

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	estimateApp "github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHttp "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra/listeners"
	repairOrderEvent "github.com/soat13/fase-1-oficina/internal/shared/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {
	repairOrderReader := estimateInfraDB.NewRepairOrderReader(container.DB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(container.DB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(container.DB)
	repository := estimateInfraDB.NewBunRepository(container.DB)
	eventPublisher := infra.NewEventPublisher(container.EventBus)

	createEstimate := estimateApp.NewCreateEstimate(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		repository,
		eventPublisher,
	)

	approve := estimateApp.NewApproveEstimate(repository, eventPublisher)
	cancel := estimateApp.NewCancelEstimate(repository)
	reject := estimateApp.NewRejectEstimate(repository, eventPublisher)
	addItem := estimateApp.NewAddItem(productCatalogReader, serviceCatalogReader, repository)
	removeItem := estimateApp.NewRemoveItem(repository)

	estimateHttpHandler := estimateInfraHttp.NewHandler(approve, reject, addItem, removeItem, container.FiberErrorHandler)
	estimateInfraHttp.Register(container.FiberApp, estimateHttpHandler)

	container.EventBus.Subscribe(
		repairOrderEvent.RepairOrderCanceled{}.Topic(),
		listeners.OnRepairOrderCanceled(cancel, repository),
	)
	container.EventBus.Subscribe(
		repairOrderEvent.RepairOrderDiagnosticsFinished{}.Topic(),
		listeners.OnDiagnosticsFinished(createEstimate),
	)

}
