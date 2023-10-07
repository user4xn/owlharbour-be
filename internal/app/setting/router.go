package setting

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	g.GET("/", h.GetDataSetting)

	g.Use(middleware.Authenticate())
	g.POST("/create-or-update", h.Store)
}
