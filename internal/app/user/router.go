package user

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.Use(middleware.BearerToken())
	g.GET("/get-profile", h.GetProfile)
}
