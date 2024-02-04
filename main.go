package main

import (
	"io"
	"net/http"
	"trip/types"

	tripService "trip/handler"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

func ErrorHandler(ctx *gin.Context) {
	ctx.Next()
	if ctx.Errors != nil {
		ctx.IndentedJSON(-1, ctx.Errors)
	}

}

func LimitRequestBodySize(limit int64) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// Replace the request body with a limited reader
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Copy().Request.Body, limit)
	}
}

func BindBodyJson(ctx *gin.Context) {
	var body any
	if err := ctx.ShouldBindJSON(&body); err != nil {
		if err != io.EOF {
			ctx.AbortWithError(http.StatusBadRequest, err)
			return
		}
	}
	ctx.Set(`body`, body)
}

func setupApp() *gin.Engine {
	app := gin.Default()
	app.Use(ErrorHandler)
	app.Use(LimitRequestBodySize(1 << 20))
	app.Use(BindBodyJson)

	return app
}

func main() {
	app := setupApp()

	app.POST("/trip",
		func(ctx *gin.Context) {
			body, _ := ctx.Get(`body`)
			var createTripDto types.CreateTripDto
			if err := mapstructure.Decode(body, &createTripDto); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
			if validationErr := validator.New(validator.WithRequiredStructEnabled()).Struct(createTripDto); validationErr != nil {
				ctx.AbortWithError(http.StatusBadRequest, validationErr)
				return
			}
			trip, requestErr := tripService.CreateTrip(createTripDto)
			if requestErr != nil {
				ctx.AbortWithError(http.StatusBadRequest, requestErr)
				return
			}
			ctx.JSON(http.StatusCreated, trip)
		},
	)

	app.GET("/trip/price",
		func(ctx *gin.Context) {
			var tripPriceDto types.TripPriceDto
			if err := ctx.ShouldBindQuery(&tripPriceDto); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
			price, requestErr := tripService.CalcTripPrice(tripPriceDto)
			if requestErr != nil {
				ctx.AbortWithError(http.StatusBadRequest, requestErr)
				return
			}
			ctx.JSON(http.StatusOK, types.PriceResponse{Price: price})
		},
	)

	app.GET("/trip/status/:status",
		func(ctx *gin.Context) {
			var listTripsParamsDto types.ListTripsParamsDto
			if err := ctx.BindUri(&listTripsParamsDto); err != nil {
				ctx.AbortWithError(http.StatusBadRequest, err)
				return
			}
			trip, requestErr := tripService.ListTrips(listTripsParamsDto)
			if requestErr != nil {
				ctx.AbortWithError(http.StatusBadRequest, requestErr)
				return
			}
			ctx.JSON(http.StatusOK, trip)
		},
	)

	app.Run(`:8080`)
}
