package setting

import (
	"fmt"
	"io"
	"net/http"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/pkg/constants"
	"simpel-api/pkg/util"

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

func (h *handler) GetDataSetting(c *gin.Context) {
	data, err := h.service.GetSetting(c)

	if err == constants.NotFoundDataAppSetting {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataAppSetting), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data setting", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) GetDataSettingWeb(c *gin.Context) {
	data, err := h.service.GetSettingWeb(c)

	if err == constants.NotFoundDataAppSetting {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataAppSetting), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success get data setting", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) Store(c *gin.Context) {
	var payload dto.PayloadStoreSetting
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "Please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Error validation", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := h.service.CreateOrUpdate(c, payload)

	if err == constants.ErrorUpdateAppSetting {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorUpdateAppSetting), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success create or update setting", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
