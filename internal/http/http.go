package http

import (
	Dashboard "simpel-api/internal/app/dashboard"
	Report "simpel-api/internal/app/report"
	Setting "simpel-api/internal/app/setting"
	Ship "simpel-api/internal/app/ship"
	User "simpel-api/internal/app/user"
	"simpel-api/internal/factory"
	"simpel-api/internal/middleware"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

// Here we define route function for user Handlers that accepts gin.Engine and factory parameters
func NewHttp(g *gin.Engine, f *factory.Factory) {

	store := cookie.NewStore([]byte("secret"))
	g.Use(sessions.Sessions("mysession", store))

	Index(g)
	// Here we use logger middleware before the actual API to catch any api call from clients
	g.Use(gin.Logger())
	// Here we use the recovery middleware to catch a panic, if panic occurs recover the application witohut shutting it off
	g.Use(gin.Recovery())

	g.Use(middleware.CORSMiddleware())

	// Here we define a router group
	v1 := g.Group("/api/v1")

	Dashboard.NewHandler(f).Router(v1.Group("/dashboard"))
	Report.NewHandler(f).Router(v1.Group("/report"))
	Setting.NewHandler(f).Router(v1.Group("/setting"))
	Ship.NewHandler(f).Router(v1.Group("/ship"))
	User.NewHandler(f).Router(v1.Group("/user"))
}

func Index(g *gin.Engine) {
	g.GET("/", func(context *gin.Context) {
		context.JSON(200, struct {
			Name string `json:"name"`
		}{
			Name: "Simpel Api",
		})
	})
}
