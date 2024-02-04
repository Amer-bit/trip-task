package types

type Trip struct {
	Uuid        string
	Origin      string
	Destination string
	Status      string
	Price       float64
}

type Location struct {
	Origin      string
	Destination string
}

type CreateTripDto struct {
	LocationForm string `validate:"required,eq=coordinate|eq=long-lat"`
	// WE COULD USE GOOGLE PLACES API TO VALIDATE ORIGIN AND DESTINATION
	Origin      string `validate:"required"`
	Destination string `validate:"required"`
}

type TripPriceDto struct {
	Origin      string `form:"origin" binding:"required"`
	Destination string `form:"destination" binding:"required"`
}

type ListTripsParamsDto struct {
	Status string `uri:"status" binding:"required,eq=pending|eq=completed"`
}

type PriceResponse struct {
	Price float64
}

type MatrixQueryParams struct {
	Destinations string
	Origins      string
	Units        string
	Key          string
}

type MatrixResponse struct {
	DestinationAddresses []string `json:"destination_addresses"`
	OriginAddresses      []string `json:"origin_addresses"`
	Rows                 []struct {
		Elements []MatrixElement `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

type MatrixElement struct {
	Distance struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	} `json:"distance"`
	Duration struct {
		Text  string `json:"text"`
		Value int    `json:"value"`
	} `json:"duration"`
	Status string `json:"status"`
}
