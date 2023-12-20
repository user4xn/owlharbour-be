package dashboard

import (
	"net/http"
	"simpel-api/internal/middleware"
	"simpel-api/pkg/util"

	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	apiKey := util.GetEnv("WEBSOCKET_API_KEY", "fallback")
	g.GET("/ship-monitor/websocket", func(c *gin.Context) {
		// Extract the API key from the WebSocket URL
		providedAPIKey := c.GetHeader("X-Websocket-Key")

		// Check if the provided API key matches the expected value
		if providedAPIKey != apiKey {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		// Upgrade to WebSocket connection and handle it in your handler function
		h.ShipMonitorWebsocket(c)
	})

	g.GET("/ship-monitor/open-websocket", h.ShipMonitorWebsocket)

	g.Use(middleware.Authenticate())
	g.GET("/statistic", h.HarbourStatistic)
	g.GET("/terrain-chart", h.TerrainChart)
	g.GET("/logs-chart", h.LogsChart)
	g.GET("/lastest-dock-ship", h.LastestDockedShip)
}
