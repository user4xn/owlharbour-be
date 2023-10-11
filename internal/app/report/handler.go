package report

import (
	"net/http"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/pkg/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) ShipDocking(c *gin.Context) {
	ctx := c.Request.Context()

	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	logType := c.DefaultQuery("type", "")
	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")
	searchParam := c.DefaultQuery("search", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	typeArray := strings.Split(logType, ",")

	param := dto.ReportShipDockedParam{
		Offset:    offset,
		Limit:     limit,
		LogType:   typeArray,
		Search:    searchParam,
		StartDate: dateStart,
		EndDate:   dateEnd,
	}

	data, err := h.service.ShipDocking(ctx, param)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data docking", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipLocation(c *gin.Context) {
	ctx := c.Request.Context()

	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")
	searchParam := c.DefaultQuery("search", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	param := dto.ReportShipLocationParam{
		Offset:    offset,
		Limit:     limit,
		Search:    searchParam,
		StartDate: dateStart,
		EndDate:   dateEnd,
	}

	data, err := h.service.ShipLocation(ctx, param)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data location", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}
