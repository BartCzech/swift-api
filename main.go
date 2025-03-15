package main

import (
	"github.com/BartCzech/swift-api/repository"
	"github.com/BartCzech/swift-api/routes"
)

func main() {
	repository.ConnectDB()

	router := routes.SetupRouter()
	router.Run(":8080")
}
