package estimate

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/estimate/application"
	eventIn "github.com/soat13/fase-1-oficina/internal/estimate/infra/in/event"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra/in/http"
	"github.com/soat13/fase-1-oficina/internal/estimate/infra/out/db"
	eventOut "github.com/soat13/fase-1-oficina/internal/estimate/infra/out/event"
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
	eventPublisher := eventOut.NewEventPublisher(container.Publisher())

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
	handleDiagnosticsFinished := application.NewHandleDiagnosticsFinished(*createEstimate)
	handleRepairOrderCanceled := application.NewHandleRepairOrderCanceled(*cancel, repository)

	estimateHttpHandler := http.NewHandler(approve, reject, addItem, removeItem, container.FiberErrorHandler)
	http.Register(container.FiberApp, estimateHttpHandler)

	container.Subscribe(events.RepairOrderCanceled{}.Topic(), eventIn.OnRepairOrderCanceled(handleRepairOrderCanceled))
	container.Subscribe(events.RepairOrderDiagnosticsFinished{}.Topic(), eventIn.OnDiagnosticsFinished(handleDiagnosticsFinished))

}
