package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/vpnvsk/p_s/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}
func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signIn)
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
