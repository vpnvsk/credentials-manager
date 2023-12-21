package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/vpnvsk/p_s/pkg/service"

	"github.com/swaggo/files"
	"github.com/swaggo/gin-swagger"

	_ "github.com/vpnvsk/p_s/docs"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api", h.userIdentity)
	{
		ps := api.Group("/ps")
		{
			ps.POST("/", h.createCredentials)
			ps.GET("/", h.getAllCredentials)
			ps.GET("/:id", h.getCredentialsById)
			ps.PUT("/:id", h.updateCredentials)
			ps.DELETE("/:id", h.deleteCredentials)

		}

	}

	return router
}
