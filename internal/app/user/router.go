package user

import (
	"simpel-api/internal/middleware"

	"github.com/gin-gonic/gin"
)

func (h *handler) Router(g *gin.RouterGroup) {
	g.POST("/login", h.Login)
	g.POST("mobile/login", h.LoginMobile)
	g.GET("/verify/email/:base_64", h.VerifyEmail)

	g.Use(middleware.Authenticate())
	g.GET("/get-profile", h.GetProfile)
	g.POST("/change-password", h.ChangePassword)
	g.POST("/admin/change-password/:user_id", h.ChangePasswordUser)
	g.POST("/logout", h.LogoutHandler)
	g.GET("/list", h.GetAllUsers)
	g.GET("/detail/:user_id", h.DetailUser)
	g.POST("/store", h.StoreUser)
	g.PUT("/update", h.UpdateUser)
	g.DELETE("/delete/:user_id", h.DeleteUser)
}
