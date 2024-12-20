package dashboard

import (
	"net/http"
	"owlharbour-api/internal/factory"
	"owlharbour-api/pkg/log"
	"owlharbour-api/pkg/util"
	"strconv"
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

func (h *handler) LastestDockedShip(c *gin.Context) {
	ctx := c.Request.Context()

	limitParam := c.DefaultQuery("limit", "25")
	limit, _ := strconv.Atoi(limitParam)

	data, err := h.service.LastestDockedShip(ctx, limit)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success retrive data lastest ship", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipMonitorWebsocket(c *gin.Context) {
	ctx := c.Request.Context()

	token := c.Query("token")
	if token == "" || token != util.GetEnv("WEBSOCKET_API_KEY", "fallback") {
		response := util.APIResponse("token is not valid", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Logging("Failed on ws handshake: %s", err.Error()).Error()
		return
	}

	defer conn.Close()

	batchEnv := util.GetEnv("WEBSOCKET_BATCH_SIZE", "30")
	batchSize, err := strconv.Atoi(batchEnv)
	if err != nil {
		response := util.APIResponse("failed to convert integer: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

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
		numBatches := (int(totalShips) + batchSize - 1) / batchSize //100 + 10 - 1 = 109/10 = 10,9

		for batch := 1; batch <= numBatches; batch++ { // 1 2 3 4 5 6 7 8 9 10
			start := (batch - 1) * batchSize // offset ships (1 - 1 = 0 * 10 = 0), first offset 0
			end := start + batchSize         // offset end ships data (0 + 10) = 10
			if end > int(totalShips) {       // if 10 > 100 = end = 110
				end = int(totalShips)
			}
			log.Logging("Fetching ships | %d | %d |", start, end).Info()
			ships, err := h.service.GetShipsInBatch(ctx, start, end)
			if err != nil {
				log.Logging("Error fetching ships %d - %d, Err: %s", start, end, err.Error()).Error()

				conn.Close()
				conn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
				if err != nil {
					log.Logging("Failed to reconnect to WebSocket, Err: %s", err.Error()).Error()
					return
				}
				continue
			}

			if err := conn.WriteJSON(ships); err != nil {
				log.Logging("Error sending ships %d - %d, Err: %s", start, end, err.Error()).Error()

				conn.Close()
				conn, err = upgrader.Upgrade(c.Writer, c.Request, nil)
				if err != nil {
					log.Logging("Failed to reconnect to WebSocket, Err: %s", err.Error()).Error()
					return
				}
				continue
			}

			time.Sleep(rate) // rate sleep interval

			if start+batchSize >= int(totalShips) { // 0+10 >= 110. if the rest all ship are sent
				rate = 10 * time.Second //lower the rate sleep interval
				batch = 0               //batch set to 0
			}
		}
	}
}

func (h *handler) HarbourStatistic(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := h.service.GetStatistic(ctx)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data statistic", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) TerrainChart(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := h.service.TerrainChart(ctx)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data terrain chart", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) LogsChart(c *gin.Context) {
	ctx := c.Request.Context()

	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")

	data, err := h.service.LogsChart(ctx, dateStart, dateEnd)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data logs chart", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
