package ship

import (
	"io"
	"net/http"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/pkg/util"
	"strconv"
	"strings"

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

func (h *handler) PairingShip(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.PairingRequest

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

	err := h.service.PairingShip(ctx, request)
	if err != nil {
		response := util.APIResponse("failed to sent pairing request: " + err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Pairing request sucessfully sent, please wait for admin approval", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) PairingRequestList(c *gin.Context) {
	ctx := c.Request.Context()

	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	statusParam := c.DefaultQuery("status", "")
	searchParam := c.DefaultQuery("search", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	statusArray := strings.Split(statusParam, ",")
	param := dto.PairingListParam{
		Offset: offset,
		Limit:  limit,
		Status: statusArray,
		Search: searchParam,
	}

	res, err := h.service.PairingRequestList(ctx, param)
	if err != nil {
		response := util.APIResponse("Failed to retrieve pairing request list: " + err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved pairing request list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) PairingAction(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.PairingActionRequest

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

	err := h.service.PairingAction(ctx, request)
	if err != nil {
		response := util.APIResponse("Unable to update pairing data: " + err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Pairing data successfully updated", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipList(c *gin.Context) {
	ctx := c.Request.Context()

	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	statusParam := c.DefaultQuery("status", "")
	searchParam := c.DefaultQuery("search", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}
	statusArray := strings.Split(statusParam, ",")
	param := dto.ShipListParam{
		Offset: offset,
		Limit:  limit,
		Status: statusArray,
		Search: searchParam,
	}

	res, err := h.service.ShipList(ctx, param)
	if err != nil {
		response := util.APIResponse("Failed to retrieve ship list: " + err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved ship list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipByDevice(c *gin.Context) {
	ctx := c.Request.Context()
	deviceIDStr := c.Param("device_id")
	if deviceIDStr == "" {
		response := util.APIResponse("Invalid device_id", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	deviceID, err := strconv.Atoi(deviceIDStr)
	if err != nil {
		response := util.APIResponse("Invalid device_id format", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	shipDetail, err := h.service.ShipByDevice(ctx, deviceID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("invalid device id, no ship data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Failed to retrieve ship data: " + err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully retrieved ship data", http.StatusOK, "success", shipDetail)
	c.JSON(http.StatusOK, response)
}

func (h *handler) RecordLog(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.ShipRecordRequest

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

	err := h.service.RecordLocationShip(ctx, request)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("Invalid device id, no ship data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Unable to record location: " + err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Location successfully recorded", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateShipDetail(c *gin.Context) {
	ctx := c.Request.Context()

	var request dto.ShipAddonDetailRequest

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

	err := h.service.UpdateShipDetail(ctx, request)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("invalid ship id, no ship data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Failed to update ship data:" + err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully updated ship data", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipDetail(c *gin.Context) {
	ctx := c.Request.Context()
	shipIDStr := c.Param("ship_id")
	if shipIDStr == "" {
		response := util.APIResponse("Invalid ship_id", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	shipID, err := strconv.Atoi(shipIDStr)
	if err != nil {
		response := util.APIResponse("Invalid ship_id format", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	shipDetail, err := h.service.ShipDetail(ctx, shipID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("invalid ship id, no ship data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Failed to retrieve detail ship data: " + err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully retrieved ship data", http.StatusOK, "success", shipDetail)
	c.JSON(http.StatusOK, response)
}
