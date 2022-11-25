package bootstrap

import (
	"context"
	"example/go-api/controllers"
	"example/go-api/infrastructure"
	"example/go-api/lib"
	"example/go-api/repository"
	"example/go-api/routes"
	"example/go-api/service"
	"fmt"

	"go.uber.org/fx"
)

var CommonModules = fx.Options(
	infrastructure.Module,
	lib.Module,
	repository.Module,
	controllers.Module,
	routes.Module,
	service.Module,
	fx.Invoke(registerHooks),
)

func registerHooks(lifecycle fx.Lifecycle, h lib.RequestHandler,userRoute routes.UserRoutes) {
	lifecycle.Append(
		fx.Hook{
			OnStart: func(context.Context) error {
				fmt.Println("Starting application in :5000")
				//userRoute.Setup()
				go h.Gin.Run(":5000")
				return nil
			},
			OnStop: func(context.Context) error {
				fmt.Println("Stopping application")
				return nil
			},
		},
	)
}