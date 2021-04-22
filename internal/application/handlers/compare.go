package handlers

import (
	"aviasales/internal/services"
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nsf/jsondiff"
)

type CompareHandler struct{}

// swagger:parameters CompareHandlerQuery
type CompareHandlerQuery struct {
	// Required: true
	Ticket1 string `json:"ticket1" form:"ticket1" binding:"required"`
	// Required: true
	Ticket2 string `json:"ticket2" form:"ticket2" binding:"required"`
}

func (s *CompareHandler) Process(
	ctx *gin.Context,
	services services.IServiceFactory,
) {
	var query CompareHandlerQuery
	if err := ctx.ShouldBindQuery(&query); err != nil {
		ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	diff, err := getCompare(services, query.Ticket1, query.Ticket2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))
		return
	}

	ctx.String(http.StatusOK, diff)
}

func getCompare(services services.IServiceFactory, t1, t2 string) (string, error) {
	ticket1, err := services.Storage().GetByUUID(t1)
	if err != nil {
		return "", err
	}
	ticket2, err := services.Storage().GetByUUID(t2)
	if err != nil {
		return "", err
	}

	ticketJSON1, err := json.Marshal(ticket1)
	if err != nil {
		return "", err
	}
	ticketJSON2, err := json.Marshal(ticket2)
	if err != nil {
		return "", err
	}

	opts := jsondiff.DefaultJSONOptions()
	_, diff := jsondiff.Compare(ticketJSON1, ticketJSON2, &opts)

	return diff, nil
}
