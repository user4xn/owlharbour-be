package dashboard

import (
	"net/http"
	"simpel-api/internal/factory"
	"simpel-api/pkg/log"
	"simpel-api/pkg/util"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type handler struct {
	service Service
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) shipMonitorWebsocket(c *gin.Context) {
	ctx := c.Request.Context()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection to WebSocket"})
		return
	}
	defer conn.Close()

	batchSize := 10
	totalShips, err := h.service.CountShip(ctx)
	if err != nil {
		response := util.APIResponse("failed to get count ship: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	initialRate := 1 * time.Second
	rate := initialRate
	lastCountFetchTime := time.Now()
	countFetchInterval := 1 * time.Hour

	for {
		if time.Since(lastCountFetchTime) >= countFetchInterval {
			newTotalShips, err := h.service.CountShip(ctx)
			if err != nil {
				log.Logging("Error fetch new total ships Err: %s", err.Error()).Error()
			} else {
				totalShips = newTotalShips
				lastCountFetchTime = time.Now()
			}
		}

		// Calculate the number of batches required
		numBatches := (int(totalShips) + batchSize - 1) / batchSize

		for batch := 1; batch <= numBatches; batch++ {
			start := (batch - 1) * batchSize
			end := start + batchSize
			if end > int(totalShips) {
				end = int(totalShips)
			}

			// log.Logging("Batch: %s | Start-End: %s | Rate: %s", batch, []int{start, end}, rate).Info()

			ships, err := h.service.GetShipsInBatch(ctx, start, end)
			if err != nil {
				break
			}

			if err := conn.WriteJSON(ships); err != nil {
				log.Logging("Error sent ships %s - %s, Err: %s", start, end, err.Error()).Error()
				break
			}

			time.Sleep(rate)

			if start+batchSize >= int(totalShips) {
				rate = 10 * time.Second
				batch = 0
			}
		}
	}
}
