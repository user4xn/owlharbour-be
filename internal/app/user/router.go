package user

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.Use(middleware.Authenticate())
	g.GET("/get-profile", h.GetProfile)
	g.GET("/list", h.GetAllUsers)
	g.POST("/store", h.StoreUser)
	g.POST("/logout", h.logoutHandler)
}
