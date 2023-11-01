package ship

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	rateLimiter := middleware.NewRateLimiter(10)

	g.POST("/pairing", h.PairingShip)
	g.GET("/pairing/detail/:device_id", h.PairingDetailByDevice)
	g.GET("/by-device/:device_id", h.ShipByDevice)
	g.POST("/record-log", rateLimiter.Limit(), h.RecordLog)

	g.Use(middleware.Authenticate())
	g.GET("/pairing-request", h.PairingRequestList)
	g.PUT("/pairing/action", h.PairingAction)

	g.GET("/list", h.ShipList)
	g.GET("/detail/:ship_id", h.ShipDetail)
	g.PUT("/update-detail", h.UpdateShipDetail)
}
