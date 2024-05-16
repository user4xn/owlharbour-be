package inspection

import (
	"io"
	"net/http"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/pkg/util"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	service Service
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service: NewService(f),
	}
}

func (h *handler) NeedCheckupShip(c *gin.Context) {
	ctx := c.Request.Context()

	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	searchParam := c.DefaultQuery("search", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	param := dto.NeedCheckupShipParam{
		Offset: offset,
		Limit:  limit,
		Search: searchParam,
	}

	data, err := h.service.NeedCheckupShip(ctx, param)
	if err != nil {
		response := util.APIResponse(err.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data need checkup ship", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateShipCheckup(c *gin.Context) {
	ctx := c.Request.Context()

	id := c.Param("log_id")

	var request dto.ShipCheckupRequest

	if err := c.ShouldBindJSON(&request); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}

		response := util.APIResponse("Invalid request payload", http.StatusBadRequest, "failed", errorMessage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	idInt, _ := strconv.Atoi(id)

	err := h.service.UpdateShipCheckup(ctx, request, idInt)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("invalid ship id, no ship data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Failed to update ship checkup data:"+err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully updated ship checkup data", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
