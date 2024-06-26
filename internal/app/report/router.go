package report

import (
	"owlharbour-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	g.Use(middleware.Authenticate())
	g.GET("/ship-docking", h.ShipDocking)
	g.GET("/ship-fraud", h.ShipFraud)
}
