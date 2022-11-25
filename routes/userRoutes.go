package routes

import (
	"example/go-api/controllers"
	"example/go-api/lib"
	"log"
)

// UserRoutes struct
type UserRoutes struct {
	handler        lib.RequestHandler
	loginController controllers.LoginController
	uploadController controllers.UploadController
}

// Setup user routes
func (s UserRoutes) Setup() {
	log.Println("Setting up routes")
	api:=s.handler.Gin
	api.POST("/login", s.loginController.Login)
		
	
}

// NewUserRoutes creates new user controller
func NewUserRoutes(
	handler lib.RequestHandler,
	loginController controllers.LoginController,
	uploadController controllers.UploadController,
) UserRoutes {
	log.Println("Setting up routes")
	api:=handler.Gin
	api.MaxMultipartMemory=8<<20
	api.POST("/login", loginController.Login)
	api.POST("/upload",uploadController.Upload)
	return UserRoutes{
		handler:        handler,
		loginController: loginController,
	}
}
