package application

import (
	"aviasales/internal/application/handlers"
	"net/http"
)

type route struct {
	path    string
	method  string
	handler handlers.IHandler
}

var routes = []*route{
	// swagger:route GET /v1/search SearchHandlerQuery
	// Responses:
	// 200
	{
		path:    "/v1/search",
		method:  http.MethodGet,
		handler: &handlers.SearchHandler{},
	},
	// swagger:route GET /v1/compare CompareHandlerQuery
	// Responses:
	// 200
	{
		path:    "/v1/compare",
		method:  http.MethodGet,
		handler: &handlers.CompareHandler{},
	},
}
