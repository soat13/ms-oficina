package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	authBootstrap "github.com/soat13/fase-1-oficina/internal/bootstrap/auth"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/estimate"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	serviceDocs "github.com/soat13/fase-1-oficina/internal/service/infra/docs"

	// product
	productApp "github.com/soat13/fase-1-oficina/internal/product/application"
	productDB "github.com/soat13/fase-1-oficina/internal/product/infra/db"
	productHTTP "github.com/soat13/fase-1-oficina/internal/product/infra/http"

	// customer
	customerApp "github.com/soat13/fase-1-oficina/internal/customer/application"
	customerDB "github.com/soat13/fase-1-oficina/internal/customer/infra/db"
	customerHTTP "github.com/soat13/fase-1-oficina/internal/customer/infra/http"

	// user
	userApp "github.com/soat13/fase-1-oficina/internal/user/application"
	userDB "github.com/soat13/fase-1-oficina/internal/user/infra/db"
	userHTTP "github.com/soat13/fase-1-oficina/internal/user/infra/http"
)

func main() {
	loadEnv()

	container := bootstrap.BuildDefault()
	defer container.Close()

	// -----------------------------------------------------------------------------
	// HTTP server setup
	// -----------------------------------------------------------------------------
	fiberApp := container.FiberApp

	authBootstrap.SetupDefault(container)
	serviceDocs.Register(fiberApp)
	estimate.SetupDefault(container)
	repairorder.SetupDefault(container)

	// -----------------------------------------------------------------------------
	// Product wiring
	// -----------------------------------------------------------------------------
	productRepo := productDB.NewBunRepository(container.DB)
	createProduct := productApp.NewCreateProduct(productRepo)
	updateProduct := productApp.NewUpdateProduct(productRepo)
	deleteProduct := productApp.NewDeleteProduct(productRepo)
	getProduct := productApp.NewGetProduct(productRepo)
	listProduct := productApp.NewListProducts(productRepo)

	productHandler := productHTTP.NewHandler(createProduct, updateProduct, deleteProduct, getProduct, listProduct, container.FiberErrorHandler)
	productHTTP.Register(fiberApp, productHandler)

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
	// Customers wiring
	// -----------------------------------------------------------------------------
	customerHandler := customerHTTP.NewHandler(createCus, updateCus, deleteCus, getCus, listCus, container.FiberErrorHandler)
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

	userHandler := userHTTP.NewHandler(createUser, updateUser, deleteUser, getUser, listUser, container.FiberErrorHandler)
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
