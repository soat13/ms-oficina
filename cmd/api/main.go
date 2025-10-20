package main

import (
	"fmt"
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/soat13/fase-1-oficina/assets/docs"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	authBootstrap "github.com/soat13/fase-1-oficina/internal/bootstrap/auth"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/customer"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/estimate"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/product"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/service"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/user"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/vehicle"
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
	docs.Register(fiberApp)
	estimate.SetupDefault(container)
	repairorder.SetupDefault(container)
	product.SetupDefault(container)
	service.SetupDefault(container)
	customer.SetupDefault(container)
	user.SetupDefault(container)
	vehicle.SetupDefault(container)

	// -----------------------------------------------------------------------------
	// HTTP server start
	// -----------------------------------------------------------------------------
	port := os.Getenv("PORT")
	printUsefulLinks(port)

	if err := fiberApp.Listen(":" + port); err != nil {
		log.Fatal(err)
	}
}

func printUsefulLinks(port string) {
	baseURL := "http://localhost:" + port

	fmt.Printf("\n🌐 Base URL:        %-45s\n", baseURL)
	fmt.Printf("📚 Swagger/OpenAPI: %-45s\n", baseURL+"/docs")
	fmt.Printf("📊 SonarQube:       %-45s\n", "http://localhost:9000")
}

func loadEnv() {
	if err := godotenv.Load(); err != nil {
		log.Println("⚠️  .env não encontrado, usando variáveis de ambiente do sistema")
	}
}
