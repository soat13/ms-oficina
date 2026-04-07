package main

import (
	"context"
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"github.com/soat13/fase-1-oficina/assets/docs"
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/customer"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/estimate"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/product"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/repairorder"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/service"
	"github.com/soat13/fase-1-oficina/internal/bootstrap/vehicle"
	"github.com/soat13/fase-1-oficina/scripts/db"
	"github.com/soat13/oficina-utils/pkg/observability"
)

func main() {
	loadEnv()

	if err := db.RunMigrations(os.Getenv("PG_DSN")); err != nil {
		log.Warn().Err(err).Msg("failed to run database migrations")
	}

	container := bootstrap.BuildDefault()
	defer container.Close()

	obs := observability.Setup(container.FiberApp, container.DB)
	defer observability.Shutdown(obs)

	container.Metrics = obs.Metrics

	docs.Register(container.FiberApp)
	estimate.SetupDefault(container)
	repairorder.SetupDefault(container)
	product.SetupDefault(container)
	service.SetupDefault(container)
	customer.SetupDefault(container)
	vehicle.SetupDefault(container)

	container.StartConsumers(context.Background())

	// -----------------------------------------------------------------------------
	// HTTP server start
	// -----------------------------------------------------------------------------
	port := os.Getenv("PORT")
	printUsefulLinks(port)

	if err := container.FiberApp.Listen(":" + port); err != nil {
		log.Fatal().Err(err).Msg("failed to start HTTP server")
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
		log.Warn().Msg(".env not found — using system environment variables")
	}
}
