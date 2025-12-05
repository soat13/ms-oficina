package repairorder

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	repairOrderApp "github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/infra"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"
	"github.com/soat13/fase-1-oficina/internal/repairorder/infra/listeners"
	"github.com/soat13/fase-1-oficina/internal/shared/events"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(container.DB)
	vehicleReader := repairOrderDB.NewVehicleReader(container.DB)
	customerReader := repairOrderDB.NewCustomerReader(container.DB)
	productReader := repairOrderDB.NewProductReader(container.DB)
	serviceReader := repairOrderDB.NewServiceReader(container.DB)
	eventPublisher := infra.NewEventPublisher(container.EventBus)

	list := repairOrderApp.NewListRepairOrders(repairOrderRepository)
	get := repairOrderApp.NewGetRepairOrder(repairOrderRepository)
	create := repairOrderApp.NewCreate(repairOrderRepository, customerReader, vehicleReader)
	startDiagnostics := repairOrderApp.NewStartDiagnostics(repairOrderRepository)
	startExecution := repairOrderApp.NewStartExecution(repairOrderRepository)
	finishExecution := repairOrderApp.NewFinishExecution(repairOrderRepository)
	releaseVehicle := repairOrderApp.NewReleaseVehicle(repairOrderRepository)
	getAverageExecutionTime := repairOrderApp.NewGetAverageExecutionTime(repairOrderRepository)
	cancel := repairOrderApp.NewCancel(repairOrderRepository, eventPublisher)
	finishDiagnostics := repairOrderApp.NewFinishDiagnostics(
		repairOrderRepository,
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

	container.EventBus.Subscribe(events.EstimateCreated{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	container.EventBus.Subscribe(events.EstimateRejected{}.Topic(), listeners.OnEstimateRejected(cancel))
	container.EventBus.Subscribe(
		events.StockInsufficientDetected{}.Topic(),
		listeners.OnStockInsufficientDetected(cancel),
	)
	container.EventBus.Subscribe(
		events.StockReduceConfirmed{}.Topic(),
		listeners.OnStockReduceConfirmed(repairOrderRepository),
	)
}
