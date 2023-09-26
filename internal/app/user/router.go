package user

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.Use(middleware.Authenticate())
	g.GET("/get-profile", h.GetProfile)
	g.GET("/list", h.GetAllUsers)
	g.GET("/detail/:user_id", h.DetailUser)
	g.POST("/store", h.StoreUser)
	g.PUT("/update", h.UpdateUser)
	g.POST("/logout", h.LogoutHandler)
	g.DELETE("/delete/:user_id", h.DeleteUser)
}
