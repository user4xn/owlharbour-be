package ship

import (
	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	// g.Use(middleware.FucntionName())
	g.POST("/pairing", h.PairingShip)
	g.POST("/record-log", h.RecordLog)
	
	g.GET("/pairing-request", h.PairingRequestList)
	g.PUT("/pairing/action", h.PairingAction)

	g.GET("/list", h.ShipList)
	g.GET("/by-device/:device_id", h.ShipByDevice)
}
