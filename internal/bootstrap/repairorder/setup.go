package repairorder

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	repairOrderApp "github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	product "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container)
}

func Setup(container *bootstrap.Container) {
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(container.DB)
	vehicleReader := repairOrderDB.NewVehicleReader(container.DB)
	customerReader := repairOrderDB.NewCustomerReader(container.DB)

	list := repairOrderApp.NewListRepairOrders(repairOrderRepository)
	get := repairOrderApp.NewGetRepairOrder(repairOrderRepository)
	create := repairOrderApp.NewCreate(repairOrderRepository, customerReader, vehicleReader)
	startDiagnostics := repairOrderApp.NewStartDiagnostics(repairOrderRepository)
	startExecution := repairOrderApp.NewStartExecution(repairOrderRepository)
	finishExecution := repairOrderApp.NewFinishExecution(repairOrderRepository)
	releaseVehicle := repairOrderApp.NewReleaseVehicle(repairOrderRepository)
	getAverageExecutionTime := repairOrderApp.NewGetAverageExecutionTime(repairOrderRepository)
	cancel := repairOrderApp.NewCancel(repairOrderRepository, container.EventBus)

	repairOrderHTTPHandler := repairOrderHTTP.NewHandler(
		list,
		get,
		create,
		startDiagnostics,
		startExecution,
		finishExecution,
		releaseVehicle,
		getAverageExecutionTime,
		cancel,
		container.FiberErrorHandler,
	)
	repairOrderHTTP.Register(container.FiberApp, repairOrderHTTPHandler)

	container.EventBus.Subscribe(estimate.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	container.EventBus.Subscribe(estimate.Rejected{}.Topic(), listeners.OnEstimateRejected(cancel))
	container.EventBus.Subscribe(
		product.StockInsufficientDetected{}.Topic(),
		listeners.OnStockInsufficientDetected(cancel),
	)
	container.EventBus.Subscribe(
		product.StockReduceConfirmed{}.Topic(),
		listeners.OnStockReduceConfirmed(repairOrderRepository),
	)
}
