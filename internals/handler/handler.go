package handler

import (
	"io"
	"net/http"
	"trip/types"

	tripService "trip/internals/service"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/mitchellh/mapstructure"
)

// APP LEVEL HANDLERS //
func ErrorHandler(ctx *gin.Context) {
	ctx.Next()
	if ctx.Errors.Errors() != nil {
		ctx.JSON(-1, ctx.Errors.Errors())
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

// ROUTE LEVEL HANDLERS //

func CreateTripHandler(ctx *gin.Context) {
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
}

func GetTripPriceHandler(ctx *gin.Context) {
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
}

func ListTripsHandlers(ctx *gin.Context) {
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
}
