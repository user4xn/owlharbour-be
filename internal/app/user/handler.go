package user

import (
	"io"
	"net/http"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/pkg/util"

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
	if err != nil {
		response := util.APIResponse("Failed Email Or Password", http.StatusBadRequest, "failed", nil)
		g.JSON(http.StatusUnprocessableEntity, response)
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

func (h *handler) logoutHandler(g *gin.Context) {
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
