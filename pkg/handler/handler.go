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

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/log-in", h.logIn)
	}
	api := router.Group("/api", h.userIdentity)
	{
		ps := api.Group("/ps")
		{
			ps.POST("/", h.createPS)
			ps.GET("/", h.getAllPS)
			ps.GET("/:id", h.getPSById)
			ps.PUT("/:id", h.updatePS)
			ps.DELETE("/:id", h.deletePS)

		}

	}

	return router
}
