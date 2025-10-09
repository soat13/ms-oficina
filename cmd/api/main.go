package main

import (
	"log"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	authBootstrap "github.com/soat13/fase-1-oficina/internal/bootstrap/auth"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/customer"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/estimate"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/product"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/user"
	serviceDocs "github.com/soat13/fase-1-oficina/internal/service/infra/docs"
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
	product.SetupDefault(container)
	customer.SetupDefault(container)
	user.SetupDefault(container)

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
