package http

import (
	Dashboard "owlharbour-api/internal/app/dashboard"
	Inspection "owlharbour-api/internal/app/inspection"
	Report "owlharbour-api/internal/app/report"
	Setting "owlharbour-api/internal/app/setting"
	Ship "owlharbour-api/internal/app/ship"
	User "owlharbour-api/internal/app/user"
	"owlharbour-api/internal/factory"
	"owlharbour-api/internal/middleware"
	"owlharbour-api/pkg/util"
	"time"

	ratelimit "github.com/JGLTechnologies/gin-rate-limit"
	"github.com/gin-gonic/gin"
)

func keyFunc(c *gin.Context) string {
	return c.ClientIP()
}

func errorHandler(c *gin.Context, info ratelimit.Info) {
	c.String(429, "Too many requests. Try again in "+time.Until(info.ResetTime).String())
}

// Here we define route function for user Handlers that accepts gin.Engine and factory parameters
func NewHttp(g *gin.Engine, f *factory.Factory) {

	Index(g)
	// Here we use logger middleware before the actual API to catch any api call from clients
	g.Use(gin.Logger())
	// Here we use the recovery middleware to catch a panic, if panic occurs recover the application witohut shutting it off
	g.Use(gin.Recovery())

	g.Use(middleware.CORSMiddleware())

	rate := ratelimit.InMemoryStore(&ratelimit.InMemoryOptions{
		Rate:  time.Minute,
		Limit: 100,
	})

	limiter := ratelimit.RateLimiter(rate, &ratelimit.Options{
		ErrorHandler: errorHandler,
		KeyFunc:      keyFunc,
	})

	// Here we define a router group
	v1 := g.Group("/api/v1")

	v1.Use(limiter)

	Dashboard.NewHandler(f).Router(v1.Group("/dashboard"))
	Report.NewHandler(f).Router(v1.Group("/report"))
	Setting.NewHandler(f).Router(v1.Group("/setting"))
	Ship.NewHandler(f).Router(v1.Group("/ship"))
	User.NewHandler(f).Router(v1.Group("/user"))
	Inspection.NewHandler(f).Router(v1.Group("/inspection"))
}

func Index(g *gin.Engine) {
	g.GET("/", func(context *gin.Context) {
		context.JSON(200, struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		}{
			Name:    "owlharbour-api",
			Version: util.GetEnv("APP_VERSION", "1.0"),
		})
	})
}
