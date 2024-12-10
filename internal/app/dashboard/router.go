package dashboard

import (
	"owlharbour-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	g.GET("/ship-monitor/websocket", h.ShipMonitorWebsocket)
	g.Use(middleware.Authenticate())
	g.GET("/statistic", h.HarbourStatistic)
	g.GET("/terrain-chart", h.TerrainChart)
	g.GET("/logs-chart", h.LogsChart)
	g.GET("/lastest-dock-ship", h.LastestDockedShip)
}
