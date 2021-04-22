package handlers

import (
	"aviasales/internal/services"

	"github.com/gin-gonic/gin"
)

type IHandler interface {
	Process(ctx *gin.Context, services services.IServiceFactory)
}
