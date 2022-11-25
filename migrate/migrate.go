package main

import (
	initialize "example/go-api/Initialize"
	"example/go-api/models"
)

func init() {
	initialize.LoadEnv()
	initialize.ConnectDB()
}
func main() {
initialize.DB.AutoMigrate(&models.Post{})
}
