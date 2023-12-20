package inspection

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.Use(middleware.Authenticate())

	g.GET("/", h.NeedCheckupShip)
	g.PUT("/update-checkup/:log_id", h.UpdateShipCheckup)
}
