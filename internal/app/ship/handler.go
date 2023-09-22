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

	service := h.service.PairingShip(ctx, request)
	if service != nil {
		response := util.APIResponse(service.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Pairing request sucessfully sent, please wait for admin approval.", http.StatusOK, "success", service)
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
		response := util.APIResponse("Failed to retrieve pairing request list", http.StatusInternalServerError, "failed", nil)
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

	service := h.service.PairingAction(ctx, request)
	if service != nil {
		response := util.APIResponse(service.Error(), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Pairing data successfully updated.", http.StatusOK, "success", service)
	c.JSON(http.StatusOK, response)
}
