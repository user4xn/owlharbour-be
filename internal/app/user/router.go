package user

import (
	"github.com/gin-gonic/gin"
	"simpel-api/internal/middleware"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.Use(middleware.Authenticate())
	g.GET("/get-profile", h.GetProfile)
	g.POST("/logout", h.logoutHandler)
}
