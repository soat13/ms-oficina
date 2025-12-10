package estimate

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra/http"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra/listeners"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {
	repairOrderReader := db.NewRepairOrderReader(container.DB)
	productCatalogReader := db.NewProductCatalogReader(container.DB)
	serviceCatalogReader := db.NewServiceCatalogReader(container.DB)
	repository := db.NewBunRepository(container.DB)
	eventPublisher := infra.NewEventPublisher(container.EventBus)

	createEstimate := application.NewCreateEstimate(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		repository,
		eventPublisher,
	)

	approve := application.NewApproveEstimate(repository, eventPublisher)
	cancel := application.NewCancelEstimate(repository)
	reject := application.NewRejectEstimate(repository, eventPublisher)
	addItem := application.NewAddItem(productCatalogReader, serviceCatalogReader, repository)
	removeItem := application.NewRemoveItem(repository)

	estimateHttpHandler := http.NewHandler(approve, reject, addItem, removeItem, container.FiberErrorHandler)
	http.Register(container.FiberApp, estimateHttpHandler)

	container.EventBus.Subscribe(
		events.RepairOrderCanceled{}.Topic(),
		listeners.OnRepairOrderCanceled(cancel, repository),
	)
	container.EventBus.Subscribe(
		events.RepairOrderDiagnosticsFinished{}.Topic(),
		listeners.OnDiagnosticsFinished(createEstimate),
	)

}
