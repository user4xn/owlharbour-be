package ship

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/pairing", h.PairingShip)
	g.POST("/record-log", h.RecordLog)

	g.Use(middleware.Authenticate())
	g.GET("/pairing-request", h.PairingRequestList)
	g.PUT("/pairing/action", h.PairingAction)

	g.GET("/list", h.ShipList)
	g.GET("/by-device/:device_id", h.ShipByDevice)
	g.GET("/detail/:ship_id", h.ShipDetail)
	g.PUT("/update-detail", h.UpdateShipDetail)
}
