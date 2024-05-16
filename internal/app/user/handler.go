package user

import (
	"fmt"
	"io"
	"net/http"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/pkg/constants"
	"owlharbour-api/pkg/util"
	"strconv"

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

func (h *handler) Login(c *gin.Context) {
	var payload dto.PayloadLogin
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if payload.Email == "" {
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", "please fill email")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data, err := h.service.LoginService(c, payload, false)
	if err == constants.UserNotFound {
		response := util.APIResponse(fmt.Sprintf("%s", constants.UserNotFound), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.InvalidPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.InvalidPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorLoadLocationTime {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorLoadLocationTime), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorGenerateJwt), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.UserNotVerifyEmail {
		response := util.APIResponse(fmt.Sprintf("%s", constants.UserNotVerifyEmail), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.EmptyGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", constants.EmptyGenerateJwt), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Login", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) LoginMobile(c *gin.Context) {
	var payload dto.PayloadLogin
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if payload.Username == "" {
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", "please fill username")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	if payload.DeviceID == "" {
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", "please fill device_id")
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data, err := h.service.LoginService(c, payload, true)
	if err == constants.UserNotFound {
		response := util.APIResponse(fmt.Sprintf("%s", constants.UserNotFound), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.InvalidPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.InvalidPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorLoadLocationTime {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorLoadLocationTime), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorGenerateJwt), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.EmptyGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", constants.EmptyGenerateJwt), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Login", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) GetProfile(c *gin.Context) {
	data := h.service.GetProfile(c, c.Value("user"))
	response := util.APIResponse("Success Get Profile", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) GetAllUsers(c *gin.Context) {
	ctx := c.Request.Context()

	search := c.Query("search")
	strLimit := c.Query("limit")
	strOffset := c.Query("offset")
	limit, _ := strconv.Atoi(strLimit)
	offset, _ := strconv.Atoi(strOffset)

	request := dto.UserListParam{
		Search: search,
		Limit:  limit,
		Offset: offset,
	}

	data, err := h.service.GetAllUsers(ctx, request)
	if err != nil {
		response := util.APIResponse("Failed to retrieve user list"+err.Error(), http.StatusInternalServerError, "failed", nil)
		c.JSON(http.StatusInternalServerError, response)
		return
	}

	response := util.APIResponse("Success Get List Users", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) StoreUser(c *gin.Context) {
	var payload dto.PayloadStoreUser
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Error Validation", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := h.service.StoreUser(c, payload)

	if err == constants.DuplicateStoreUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.DuplicateStoreUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorHashPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorHashPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Store User", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) LogoutHandler(c *gin.Context) {
	err := h.service.LogoutService(c, c.Value("user"))
	if err != nil {
		response := util.APIResponse("Failed Logout", http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusOK, response)
		return
	}

	response := util.APIResponse("Success Logout", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) DetailUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))

	data, err := h.service.DetailUser(c, userID)

	if err == constants.NotFoundDataUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Get Detail User", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)
}

func (h *handler) UpdateUser(c *gin.Context) {
	var payload dto.PayloadUpdateUser
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("there is an incomplete request", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := h.service.UpdateUser(c, payload)
	if err == constants.NotFoundDataUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorLoadLocationTime {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorLoadLocationTime), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorHashPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorHashPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.FailedUpdateUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.FailedUpdateUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Update User", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) VerifyEmail(c *gin.Context) {
	base64String := c.Param("base_64")

	err := h.service.VerifyEmail(c, base64String)

	if err == constants.NotFoundDataUser {
		urlRedirect := util.GetEnv("FE_URL", "fallback") + "/auth/error-verify"
		c.Redirect(http.StatusSeeOther, urlRedirect)
	}

	if err == constants.FailedDeleteUser {
		urlRedirect := util.GetEnv("FE_URL", "fallback") + "/auth/error-verify"
		c.Redirect(http.StatusSeeOther, urlRedirect)
	}
	urlRedirect := util.GetEnv("FE_URL", "fallback") + "/verification "
	c.Redirect(http.StatusSeeOther, urlRedirect)
}

func (h *handler) DeleteUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))

	err := h.service.DeleteUser(c, userID)

	if err == constants.NotFoundDataUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.FailedDeleteUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.FailedDeleteUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Get Detail User", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ChangePassword(c *gin.Context) {
	user := h.service.GetProfile(c, c.Value("user"))
	var payload dto.PayloadChangePassword
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Error Validation", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := h.service.ChangePassword(c, user.ID, payload)

	if err == constants.NotFoundDataUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.FailedNotSamePassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.FailedNotSamePassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.MinimCharacterPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.MinimCharacterPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.PasswordSameCurrent {
		response := util.APIResponse(fmt.Sprintf("%s", constants.PasswordSameCurrent), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.FailedChangePassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.FailedChangePassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Change Passowrd", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}

func (h *handler) ChangePasswordUser(c *gin.Context) {
	userID, _ := strconv.Atoi(c.Param("user_id"))
	var payload dto.PayloadChangePassword
	if err := c.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Error Validation", http.StatusUnprocessableEntity, "failed", errorMessage)
		c.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := h.service.ChangePassword(c, userID, payload)

	if err == constants.NotFoundDataUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.NotFoundDataUser), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.FailedNotSamePassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.FailedNotSamePassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.MinimCharacterPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.MinimCharacterPassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.PasswordSameCurrent {
		response := util.APIResponse(fmt.Sprintf("%s", constants.PasswordSameCurrent), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.FailedChangePassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.FailedChangePassword), http.StatusBadRequest, "failed", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Change Passowrd", http.StatusOK, "success", nil)
	c.JSON(http.StatusOK, response)
}
