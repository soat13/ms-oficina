package user

import (
	"github.com/soat13/fase-1-oficina/internal/bootstrap"
	userApp "github.com/soat13/fase-1-oficina/internal/user/application"
	userDB "github.com/soat13/fase-1-oficina/internal/user/infra/db"
	userHTTP "github.com/soat13/fase-1-oficina/internal/user/infra/http"
)

func SetupDefault(container *bootstrap.Container) {
	Setup(container, nil)
}

func Setup(container *bootstrap.Container, repository userApp.UserRepository) {
	if repository == nil {
		repository = userDB.NewBunUserRepository(container.DB)
	}

	createUser := userApp.NewCreateUser(repository)
	updateUser := userApp.NewUpdateUser(repository)
	deleteUser := userApp.NewDeleteUser(repository)
	getUser := userApp.NewGetUser(repository)
	listUsers := userApp.NewListUsers(repository)

	userHandler := userHTTP.NewHandler(
		createUser,
		updateUser,
		deleteUser,
		getUser,
		listUsers,
		container.FiberErrorHandler,
	)

	userHTTP.Register(container.FiberApp, userHandler)
}
