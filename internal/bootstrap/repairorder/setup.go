package repairorder

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	"github.com/soat13/fase-1-oficina/internal/repairorder/infra/event"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {
	repository := repairOrderDB.NewBunRepairOrderRepository(container.DB)
	vehicleReader := repairOrderDB.NewVehicleReader(container.DB)
	customerReader := repairOrderDB.NewCustomerReader(container.DB)
	productReader := repairOrderDB.NewProductReader(container.DB)
	serviceReader := repairOrderDB.NewServiceReader(container.DB)
	eventPublisher := event.NewEventPublisher(container.EventBus)

	list := application.NewListRepairOrders(repository)
	get := application.NewGetRepairOrder(repository)
	create := application.NewCreate(repository, customerReader, vehicleReader)
	startDiagnostics := application.NewStartDiagnostics(repository)
	startExecution := application.NewStartExecution(repository)
	finishExecution := application.NewFinishExecution(repository)
	releaseVehicle := application.NewReleaseVehicle(repository)
	getAverageExecutionTime := application.NewGetAverageExecutionTime(repository)
	cancel := application.NewCancel(repository, eventPublisher)
	finishDiagnostics := application.NewFinishDiagnostics(
		repository,
		eventPublisher,
		productReader,
		serviceReader,
	)

	repairOrderHTTPHandler := repairOrderHTTP.NewHandler(
		list,
		get,
		create,
		startDiagnostics,
		finishDiagnostics,
		startExecution,
		finishExecution,
		releaseVehicle,
		getAverageExecutionTime,
		cancel,
		container.Validator,
		container.FiberErrorHandler,
	)
	repairOrderHTTP.Register(container.FiberApp, repairOrderHTTPHandler)

	handleEstimateCreated := application.NewHandleEstimateCreated(repository)
	handleEstimateRejected := application.NewHandleEstimateRejected(cancel)
	handleStockInsufficient := application.NewHandleStockInsufficient(cancel)
	handleStockReductionConfirmed := application.NewHandleStockReductionConfirmed(repository)

	container.EventBus.Subscribe(events.EstimateCreated{}.Topic(), event.OnEstimateCreated(handleEstimateCreated))
	container.EventBus.Subscribe(events.EstimateRejected{}.Topic(), event.OnEstimateRejected(handleEstimateRejected))
	container.EventBus.Subscribe(
		events.StockInsufficientDetected{}.Topic(),
		event.OnStockInsufficientDetected(handleStockInsufficient),
	)
	container.EventBus.Subscribe(
		events.StockReductionConfirmed{}.Topic(),
		event.OnStockReduceConfirmed(handleStockReductionConfirmed),
	)
}
