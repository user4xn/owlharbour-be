package ship

import (
	"owlharbour-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	//100 rps , 1 minute sliding windows
	// rateLimiter := middleware.NewRateLimiter(100, time.Minute)

	g.POST("/pairing", h.PairingShip)
	g.GET("/pairing/detail", h.PairingDetailByUsername)

	g.Use(middleware.Authenticate())
	g.GET("/mobile/profile", h.ShipByAuth)
	g.GET("/mobile/dock-log/:device_id", h.ShipDockLogByDevice)
	g.GET("/mobile/location-log/:device_id", h.ShipLocationLogByDevice)
	g.POST("/record-log", h.RecordRabbitShip)
	g.GET("/pairing-request", h.PairingRequestList)
	g.GET("/pairing-request/count", h.PairingRequestCount)
	g.PUT("/pairing/action", h.PairingAction)

	g.GET("/list", h.ShipList)
	g.GET("/detail/:ship_id", h.ShipDetail)
	g.GET("/dock-log/:ship_id", h.ShipDockLog)
	g.GET("/location-log/:ship_id", h.ShipLocationLog)
	g.PUT("/update-detail", h.UpdateShipDetail)
}
