package repairorder

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application"
	eventIn "github.com/soat13/fase-1-oficina/internal/repairorder/infra/in/event"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/in/http"
	"github.com/soat13/fase-1-oficina/internal/repairorder/infra/out/db"
	eventOut "github.com/soat13/fase-1-oficina/internal/repairorder/infra/out/event"
	metricsOut "github.com/soat13/fase-1-oficina/internal/repairorder/infra/out/metrics"
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
	eventPublisher := eventOut.NewEventPublisher(container.Publisher())
	metricsPublisher := metricsOut.NewPublisher(container.Metrics)

	list := application.NewListRepairOrders(repository)
	get := application.NewGetRepairOrder(repository)
	create := application.NewCreate(repository, customerReader, vehicleReader, metricsPublisher)
	startDiagnostics := application.NewStartDiagnostics(repository, metricsPublisher)
	startExecution := application.NewStartExecution(repository, metricsPublisher)
	finishExecution := application.NewFinishExecution(repository, eventPublisher, metricsPublisher)
	releaseVehicle := application.NewReleaseVehicle(repository, metricsPublisher)
	getAverageExecutionTime := application.NewGetAverageExecutionTime(repository)
	cancel := application.NewCancel(repository, eventPublisher, metricsPublisher)
	finishDiagnostics := application.NewFinishDiagnostics(
		repository,
		eventPublisher,
		productReader,
		serviceReader,
		metricsPublisher,
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

	handleEstimateCreated := application.NewHandleEstimateCreated(repository, metricsPublisher)
	handleEstimateRejected := application.NewHandleEstimateRejected(cancel)
	handleStockInsufficient := application.NewHandleStockInsufficient(cancel)
	handleStockReductionConfirmed := application.NewHandleStockReductionConfirmed(repository, metricsPublisher)
	handlePaymentStatusChanged := application.NewHandlePaymentStatusChanged(repository, metricsPublisher)

	container.Subscribe(events.EstimateCreated{}.Topic(), eventIn.OnEstimateCreated(handleEstimateCreated))
	container.Subscribe(events.EstimateRejected{}.Topic(), eventIn.OnEstimateRejected(handleEstimateRejected))
	container.Subscribe(events.StockInsufficientDetected{}.Topic(), eventIn.OnStockInsufficientDetected(handleStockInsufficient))
	container.Subscribe(events.PaymentStatusChanged{}.Topic(), eventIn.OnPaymentStatusChanged(handlePaymentStatusChanged))
	container.Subscribe(events.TopicRepairOrderStockReductionConfirmed, eventIn.OnStockReduceConfirmed(handleStockReductionConfirmed))
}
