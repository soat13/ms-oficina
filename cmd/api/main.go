package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	estimaterListeners "github.com/soat13/fase-1-oficina/internal/estimate/application/listeners"
	"github.com/soat13/fase-1-oficina/internal/shared/kernel/estimate"
	productEvent "github.com/soat13/fase-1-oficina/internal/shared/kernel/product"
	// estimate
	estimateApp "github.com/soat13/fase-1-oficina/internal/estimate/application"
	estimateInfraDB "github.com/soat13/fase-1-oficina/internal/estimate/infra/db"
	estimateInfraHttp "github.com/soat13/fase-1-oficina/internal/estimate/infra/http"

	serviceDocs "github.com/soat13/fase-1-oficina/internal/service/infra/docs"
	// customer
	customerApp "github.com/soat13/fase-1-oficina/internal/customer/application"
	customerDB "github.com/soat13/fase-1-oficina/internal/customer/infra/db"
	customerHTTP "github.com/soat13/fase-1-oficina/internal/customer/infra/http"

	// user
	userApp "github.com/soat13/fase-1-oficina/internal/user/application"
	userDB "github.com/soat13/fase-1-oficina/internal/user/infra/db"
	userHTTP "github.com/soat13/fase-1-oficina/internal/user/infra/http"

	// repair order listeners
	repairOrderApp "github.com/soat13/fase-1-oficina/internal/repairorder/application"
	"github.com/soat13/fase-1-oficina/internal/repairorder/application/listeners"
	repairOrderDB "github.com/soat13/fase-1-oficina/internal/repairorder/infra/db"
	repairOrderHTTP "github.com/soat13/fase-1-oficina/internal/repairorder/infra/http"

	"github.com/soat13/fase-1-oficina/internal/shared/eventbus"
)

func main() {
	loadEnv()

	container := bootstrap.BuildDefault()
	defer container.Close()

	eventBus := eventbus.NewInMemoryBus()

	// -----------------------------------------------------------------------------
	// HTTP server setup
	// -----------------------------------------------------------------------------
	fiberApp := container.FiberApp

	serviceDocs.Register(fiberApp)

	errorHandler := container.FiberErrorHandler

	// -----------------------------------------------------------------------------
	// Estimate wiring
	// -----------------------------------------------------------------------------
	repairOrderReader := estimateInfraDB.NewRepairOrderReader(container.DB)
	productCatalogReader := estimateInfraDB.NewProductCatalogReader(container.DB)
	serviceCatalogReader := estimateInfraDB.NewServiceCatalogReader(container.DB)
	estimateRepository := estimateInfraDB.NewBunRepository(container.DB)

	createEstimate := estimateApp.NewCreateEstimate(
		repairOrderReader,
		productCatalogReader,
		serviceCatalogReader,
		estimateRepository,
		eventBus,
	)

	approveEstimate := estimateApp.NewApproveEstimate(estimateRepository, eventBus)

	estimateHttpHandler := estimateInfraHttp.NewHandler(createEstimate, approveEstimate)
	estimateInfraHttp.Register(fiberApp, estimateHttpHandler)

	eventBus.Subscribe(
		productEvent.StockReduceConfirmed{}.Topic(),
		estimaterListeners.OnStockReduceConfirmed(estimateRepository, eventBus),
	)

	// -----------------------------------------------------------------------------
	// Customer wiring
	// -----------------------------------------------------------------------------
	customerRepo := customerDB.NewBunCustomerRepository(container.DB)
	createCus := customerApp.NewCreateCustomer(customerRepo)
	updateCus := customerApp.NewUpdateCustomer(customerRepo)
	deleteCus := customerApp.NewDeleteCustomer(customerRepo)
	getCus := customerApp.NewGetCustomer(customerRepo)
	listCus := customerApp.NewListCustomers(customerRepo)

	// -----------------------------------------------------------------------------
	// Repairorder wiring
	// -----------------------------------------------------------------------------
	repairOrderRepository := repairOrderDB.NewBunRepairOrderRepository(container.DB)
	startExecution := repairOrderApp.NewStartExecution(repairOrderRepository)
	finishExecution := repairOrderApp.NewFinishExecution(repairOrderRepository)
	releaseVehicle := repairOrderApp.NewReleaseVehicle(repairOrderRepository)

	repairOrderHTTPHandler := repairOrderHTTP.NewHandler(startExecution, finishExecution, releaseVehicle, errorHandler)
	repairOrderHTTP.Register(fiberApp, repairOrderHTTPHandler)

	// Event listeners
	eventBus.Subscribe(estimate.Created{}.Topic(), listeners.OnEstimateCreated(repairOrderRepository))
	eventBus.Subscribe(estimate.Approved{}.Topic(), listeners.OnEstimateApproved(repairOrderRepository))

	// -----------------------------------------------------------------------------
	// Customers wiring
	// -----------------------------------------------------------------------------
	customerHandler := customerHTTP.NewHandler(createCus, updateCus, deleteCus, getCus, listCus, errorHandler)
	customerHTTP.Register(fiberApp, customerHandler)

	// -----------------------------------------------------------------------------
	// User wiring
	// -----------------------------------------------------------------------------
	userRepo := userDB.NewBunUserRepository(container.DB)
	createUser := userApp.NewCreateUser(userRepo)
	updateUser := userApp.NewUpdateUser(userRepo)
	deleteUser := userApp.NewDeleteUser(userRepo)
	getUser := userApp.NewGetUser(userRepo)
	listUser := userApp.NewListUsers(userRepo)

	userHandler := userHTTP.NewHandler(createUser, updateUser, deleteUser, getUser, listUser, errorHandler)
	userHTTP.Register(fiberApp, userHandler)

	// -----------------------------------------------------------------------------
	// HTTP server start
	// -----------------------------------------------------------------------------
	port := os.Getenv("PORT")
	if err := fiberApp.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado, usando variáveis de ambiente do sistema")
	}
}
