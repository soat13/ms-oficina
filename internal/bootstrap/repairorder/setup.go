package repairorder

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	eventIn "github.com/soat13/fase-1-oficina/internal/repairorder/infra/in/event"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/in/http"
	"github.com/soat13/fase-1-oficina/internal/repairorder/infra/out/db"
	eventOut "github.com/soat13/fase-1-oficina/internal/repairorder/infra/out/event"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {
	repository := db.NewBunRepairOrderRepository(container.DB)
	vehicleReader := db.NewVehicleReader(container.DB)
	customerReader := db.NewCustomerReader(container.DB)
	productReader := db.NewProductReader(container.DB)
	serviceReader := db.NewServiceReader(container.DB)
	eventPublisher := eventOut.NewEventPublisher(container.EventBus)

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

	container.EventBus.Subscribe(events.EstimateCreated{}.Topic(), eventIn.OnEstimateCreated(handleEstimateCreated))
	container.EventBus.Subscribe(events.EstimateRejected{}.Topic(), eventIn.OnEstimateRejected(handleEstimateRejected))
	container.EventBus.Subscribe(
		events.StockInsufficientDetected{}.Topic(),
		eventIn.OnStockInsufficientDetected(handleStockInsufficient),
	)
	container.EventBus.Subscribe(
		events.StockReductionConfirmed{}.Topic(),
		eventIn.OnStockReduceConfirmed(handleStockReductionConfirmed),
	)
}
