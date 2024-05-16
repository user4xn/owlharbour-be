package setting

import (
	"owlharbour-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	g.GET("/mobile", h.GetDataSetting)

	g.Use(middleware.Authenticate())
	g.GET("/web", h.GetDataSettingWeb)
	g.POST("/create-or-update", h.Store)
}
