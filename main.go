package main

import (
	//"example/go-api/controllers"
	"example/go-api/bootstrap"

	"go.uber.org/fx"
)

// func init(){
// 	initialize.LoadEnv()
// 	initialize.ConnectDB()
// }

func main() {
	// router.POST("/post", controllers.CreatePost)
	// router.GET("/post", controllers.GetPosts)
	// router.GET("/post/:id", controllers.FetchOnePost)
	// router.PUT("/post/:id", controllers.UpdatePost)
	// routes.PostRoutes(router);
	fx.New(
		bootstrap.CommonModules,
	).Run()
	// listen and serve on 0.0.0.0:3000
}