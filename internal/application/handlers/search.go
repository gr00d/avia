package handlers

import (
	"aviasales/internal/services"
	"aviasales/pkg/entities"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SearchHandlerTypeCheapest      = "cheapest"
	SearchHandlerTypeMostExpensive = "mostExpensive"
	SearchHandlerTypeLongest       = "longest"
	SearchHandlerTypeShortest      = "shortest"
	SearchHandlerTypeOptimal       = "optimal"
)

type SearchHandler struct{}

// swagger:parameters SearchHandlerQuery
type SearchHandlerQuery struct {
	// Required: true
	Source string `json:"source" form:"source" binding:"required"`
	// Required: true
	Destination string `json:"destination" form:"destination" binding:"required"`
	// Possible type: cheapest mostExpensive longest shortest optimal
	Type string `json:"type" form:"type" binding:"omitempty,oneof=cheapest mostExpensive longest shortest optimal"`
}

func (s *SearchHandler) Process(
	ctx *gin.Context,
	services services.IServiceFactory,
) {
	var query SearchHandlerQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	sourceCity, destinationCity := query.Source, query.Destination

	if query.Type != "" {
		var result *entities.Itinerary
		var err error
		switch query.Type {
		case SearchHandlerTypeCheapest:
			result, err = services.Storage().GetCheapest(sourceCity, destinationCity)
		case SearchHandlerTypeMostExpensive:
			result, err = services.Storage().GetMostExpensive(sourceCity, destinationCity)
		case SearchHandlerTypeLongest:
			result, err = services.Storage().GetLongest(sourceCity, destinationCity)
		case SearchHandlerTypeShortest:
			result, err = services.Storage().GetShortest(sourceCity, destinationCity)
		case SearchHandlerTypeOptimal:
			result, err = services.Storage().GetOptimal(sourceCity, destinationCity)
		}

		if err != nil {
			ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
			return
		}
		ctx.JSON(http.StatusOK, result)
		return
	}

	itineraries, err := services.Storage().GetItineraries(query.Source, query.Destination)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
		return
	}

	ctx.JSON(http.StatusOK, itineraries)
}
