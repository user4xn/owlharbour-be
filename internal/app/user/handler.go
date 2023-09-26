package user

import (
	"fmt"
	"io"
	"net/http"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/pkg/constants"
	"simpel-api/pkg/util"
	"strconv"

	"github.com/gin-contrib/sessions"
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

func (h *handler) Login(g *gin.Context) {
	var payload dto.PayloadLogin
	if err := g.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Failed Login", http.StatusUnprocessableEntity, "failed", errorMessage)
		g.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	data, err := h.service.LoginService(g, payload)
	if err == constants.UserNotFound {
		response := util.APIResponse(fmt.Sprintf("%s", constants.UserNotFound), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.InvalidPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.InvalidPassword), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorLoadLocationTime {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorLoadLocationTime), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorGenerateJwt), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.EmptyGenerateJwt {
		response := util.APIResponse(fmt.Sprintf("%s", constants.EmptyGenerateJwt), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	session := sessions.Default(g)
	session.Set("user_id", data.DataUser.ID)
	session.Set("token", data.TokenJwt)
	session.Save()

	response := util.APIResponse("Success Login", http.StatusOK, "success", data)
	g.JSON(http.StatusOK, response)
	return
}

func (h *handler) GetProfile(g *gin.Context) {

	data := h.service.GetProfile(g, g.Value("user"))
	response := util.APIResponse("Success Get Profile", http.StatusOK, "success", data)
	g.JSON(http.StatusOK, response)
	return

}

func (h *handler) GetAllUsers(g *gin.Context) {
	search := g.PostForm("search")
	strLimit := g.PostForm("limit")
	strOffset := g.PostForm("offset")

	limit, err := strconv.Atoi(strLimit)
	if err != nil {
		response := util.APIResponse("Failed Get List Users", http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	offset, err := strconv.Atoi(strOffset)
	if err != nil {
		response := util.APIResponse("Failed Get List Users", http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	data := h.service.GetAllUsers(g, search, limit, offset)

	response := util.APIResponse("Success Get List Users", http.StatusOK, "success", data)
	g.JSON(http.StatusOK, response)
}

func (h *handler) StoreUser(g *gin.Context) {
	var payload dto.PayloadStoreUser
	if err := g.ShouldBind(&payload); err != nil {
		errorMessage := gin.H{"errors": "please fill data"}
		if err != io.EOF {
			errors := util.FormatValidationError(err)
			errorMessage = gin.H{"errors": errors}
		}
		response := util.APIResponse("Error Validation", http.StatusUnprocessableEntity, "failed", errorMessage)
		g.JSON(http.StatusUnprocessableEntity, response)
		return
	}

	err := h.service.StoreUser(g, payload)

	if err == constants.DuplicateStoreUser {
		response := util.APIResponse(fmt.Sprintf("%s", constants.DuplicateStoreUser), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	if err == constants.ErrorHashPassword {
		response := util.APIResponse(fmt.Sprintf("%s", constants.ErrorHashPassword), http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusBadRequest, response)
		return
	}

	response := util.APIResponse("Success Store User", http.StatusOK, "success", nil)
	g.JSON(http.StatusOK, response)
	return
}

func (h *handler) LogoutHandler(g *gin.Context) {
	session := sessions.Default(g)
	tokenString := session.Get("token")
	if tokenString != nil {
		session.Clear()
		session.Save()
		response := util.APIResponse("Success Logout", http.StatusOK, "success", nil)
		g.JSON(http.StatusOK, response)
		return
	}

}
