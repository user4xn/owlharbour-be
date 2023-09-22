package ship

import (
	"github.com/gin-gonic/gin"
)

// This function accepts gin.Routergroup to define a group route
func (h *handler) Router(g *gin.RouterGroup) {
	// g.Use(middleware.FucntionName())
	g.POST("/pairing", h.PairingShip)
	g.GET("/pairing-request", h.PairingRequestList)
	g.PUT("/pairing/action", h.PairingAction)
}
