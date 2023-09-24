package user

import (
	"simpel-api/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))

	g.POST("/login", h.Login)
	g.Use(middleware.Authenticate())
	g.GET("/get-profile", h.GetProfile)
	g.POST("/logout", h.logoutHandler)
}
