package repairorder

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/gofiber/fiber/v2"
	"github.com/soat13/ms-oficina/internal/bootstrap"
	"github.com/soat13/ms-oficina/internal/bootstrap/estimate"
	"github.com/soat13/ms-oficina/internal/bootstrap/repairorder"
	"github.com/soat13/ms-oficina/tests/testsupport"
)

func TestFeatures(t *testing.T) {
	setup := testsupport.SetupHTTP(t, func(app *fiber.App, c *bootstrap.Container) {
		repairorder.SetupDefault(c)
		estimate.SetupDefault(c)
	})

	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			sc.Before(func(ctx context.Context, _ *godog.Scenario) (context.Context, error) {
				return withWorld(ctx, newWorld(setup)), nil
			})

			// Given
			sc.Step(`^que existe um produto "([^"]+)" com estoque (\d+)$`, givenProductExists)
			sc.Step(`^existe um produto "([^"]+)" com estoque (\d+)$`, givenProductExists)
			sc.Step(`^existe um serviço "([^"]+)"$`, givenServiceExists)
			sc.Step(`^que existe uma Ordem de Reparo no status "([^"]+)"$`, givenRepairOrderExistsWithStatus)

			// When
			sc.Step(`^eu inicio o diagnóstico$`, whenStartDiagnostics)
			sc.Step(`^eu finalizo o diagnóstico com:$`, whenFinishDiagnosticsWith)

			// Then
			sc.Step(`^a OS deve estar no status "([^"]+)"$`, thenRepairOrderShouldBeInStatus)
			sc.Step(`^a OS deve permanecer no status "([^"]+)"$`, thenRepairOrderShouldBeInStatus)
			sc.Step(`^a resposta HTTP deve ser (\d+)$`, thenHttpResponseShouldBe)
			sc.Step(`^um orçamento deve ter sido criado para a OS no status "([^"]+)"$`, thenEstimateShouldBeInStatus)
			sc.Step(`^o orçamento deve conter (\d+) itens$`, thenEstimateShouldContainNItems)
		},
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("BDD suite failed")
	}
}
