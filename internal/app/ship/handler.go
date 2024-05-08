package ship

import (
	"io"
	"net/http"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/model"
	"simpel-api/internal/repository"
	"simpel-api/pkg/util"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type handler struct {
	service            Service
	rabbitMqRepository repository.RabbitMq
}

func NewHandler(f *factory.Factory) *handler {
	return &handler{
		service:            NewService(f),
		rabbitMqRepository: f.RabbitMqRepository,
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
		response := util.APIResponse("failed to sent pairing request: "+err.Error(), http.StatusInternalServerError, "failed", nil)
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
		response := util.APIResponse("Failed to retrieve pairing request list: "+err.Error(), http.StatusInternalServerError, "failed", nil)
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
		response := util.APIResponse("Unable to update pairing data: "+err.Error(), http.StatusInternalServerError, "failed", nil)
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
		response := util.APIResponse("Failed to retrieve ship list: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved ship list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipByAuth(c *gin.Context) {
	ctx := c.Request.Context()

	user, ok := c.Get("user")
	if !ok {
		response := util.APIResponse("User information not found", http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	authUser, ok := user.(model.User)
	if !ok {
		response := util.APIResponse("Invalid user type", http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	shipDetail, err := h.service.ShipByAuth(ctx, authUser)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("invalid device id, no ship data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Failed to retrieve ship data: "+err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully retrieved ship data", http.StatusOK, "success", shipDetail)
	c.JSON(http.StatusOK, response)
}

func (h *handler) RecordRabbitShip(c *gin.Context) {
	var request dto.ShipRecordRequest
	err := c.ShouldBindJSON(&request)
	if err != nil {
		response := util.APIResponse("insert rabbit ship record failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	err = h.service.RecordShipRabbit(c.Request.Context(), request)
	if err != nil {
		response := util.APIResponse("insert rabbit ship record failed", http.StatusInternalServerError, "error", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}
	response := util.APIResponse("insert rabbit ship record success", http.StatusOK, "success", nil)
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
			statusCode := http.StatusInternalServerError
			if err.Error() == "HTTP request failed with status code 429" {
				statusCode = http.StatusTooManyRequests
			}

			response := util.APIResponse("Unable to record location: "+err.Error(), statusCode, "failed", nil)
			c.JSON(statusCode, response)
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
			response := util.APIResponse("Failed to update ship data:"+err.Error(), http.StatusInternalServerError, "failed", nil)
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
			response := util.APIResponse("Failed to retrieve detail ship data: "+err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully retrieved ship data", http.StatusOK, "success", shipDetail)
	c.JSON(http.StatusOK, response)
}

func (h *handler) PairingDetailByUsername(c *gin.Context) {
	ctx := c.Request.Context()
	username := c.Query("username")
	if username == "" {
		response := util.APIResponse("Invalid username", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	shipDetail, err := h.service.PairingDetailByUsername(ctx, username)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			response := util.APIResponse("invalid username, no pairing data", http.StatusBadRequest, "failed", nil)
			c.JSON(http.StatusBadRequest, response)
		} else {
			response := util.APIResponse("Failed to retrieve detail pairing data: "+err.Error(), http.StatusInternalServerError, "failed", nil)
			c.JSON(http.StatusInternalServerError, response)
		}
		return
	}

	response := util.APIResponse("Successfully retrieved pairing data", http.StatusOK, "success", shipDetail)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipDockLog(c *gin.Context) {
	ctx := c.Request.Context()

	shipIDStr := c.Param("ship_id")
	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)
	shipID, _ := strconv.Atoi(shipIDStr)

	if limit == 0 {
		limit = 10
	}

	param := dto.ShipLogParam{
		Offset:    offset,
		Limit:     limit,
		StartDate: dateStart,
		EndDate:   dateEnd,
	}

	res, err := h.service.ShipDockLog(ctx, param, shipID)
	if err != nil {
		response := util.APIResponse("Failed to retrieve ship list: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved ship list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipLocationLog(c *gin.Context) {
	ctx := c.Request.Context()

	shipIDStr := c.Param("ship_id")
	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)
	shipID, _ := strconv.Atoi(shipIDStr)

	if limit == 0 {
		limit = 10
	}

	param := dto.ShipLogParam{
		Offset:    offset,
		Limit:     limit,
		StartDate: dateStart,
		EndDate:   dateEnd,
	}

	res, err := h.service.ShipLocationLog(ctx, param, shipID)
	if err != nil {
		response := util.APIResponse("Failed to retrieve ship list: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved ship list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipDockLogByDevice(c *gin.Context) {
	ctx := c.Request.Context()

	deviceID := c.Param("device_id")
	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	param := dto.ShipLogParam{
		Offset:    offset,
		Limit:     limit,
		StartDate: dateStart,
		EndDate:   dateEnd,
	}

	res, err := h.service.ShipDockLog(ctx, param, deviceID)
	if err != nil {
		response := util.APIResponse("Failed to retrieve ship list: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved ship list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ShipLocationLogByDevice(c *gin.Context) {
	ctx := c.Request.Context()

	deviceID := c.Param("device_id")
	offsetParam := c.DefaultQuery("offset", "0")
	limitParam := c.DefaultQuery("limit", "25")
	dateStart := c.DefaultQuery("start_date", "")
	dateEnd := c.DefaultQuery("end_date", "")

	offset, _ := strconv.Atoi(offsetParam)
	limit, _ := strconv.Atoi(limitParam)

	if limit == 0 {
		limit = 10
	}

	param := dto.ShipLogParam{
		Offset:    offset,
		Limit:     limit,
		StartDate: dateStart,
		EndDate:   dateEnd,
	}

	res, err := h.service.ShipLocationLog(ctx, param, deviceID)
	if err != nil {
		response := util.APIResponse("Failed to retrieve ship list: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved ship list", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}

func (h *handler) PairingRequestCount(c *gin.Context) {
	ctx := c.Request.Context()

	res, err := h.service.PairingRequestCount(ctx)
	if err != nil {
		response := util.APIResponse("Failed to retrieve pairing ship count: "+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Successfully retrieved pairing ship count", http.StatusOK, "success", res)
	c.JSON(http.StatusOK, response)
}
