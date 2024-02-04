package main

import (
	"trip/internals/handler"

	"github.com/gin-gonic/gin"
)

func setupApp() *gin.Engine {
	app := gin.Default()
	app.Use(handler.ErrorHandler)
	app.Use(handler.LimitRequestBodySize(1 << 20))
	app.Use(handler.BindBodyJson)

	return app
}

func main() {
	app := setupApp()

	app.POST("/trip", handler.CreateTripHandler)

	app.GET("/trip/price", handler.GetTripPriceHandler)

	app.GET("/trip/status/:status", handler.ListTripsHandlers)

	app.Run(`:8080`)
}
