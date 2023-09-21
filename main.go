package main

import (
	"flag"
	"log"
	"simpel-api/database"
	"simpel-api/database/migration"
	"simpel-api/internal/factory"
	"simpel-api/internal/http"
	"simpel-api/pkg/util"

	"github.com/gin-gonic/gin"
)

func main() {
	var m string
	var s string

	database.CreateConnection()

	flag.StringVar(
		&m,
		"m",
		"none",
		`This flag is used for migration`,
	)

	flag.StringVar(
		&s,
		"s",
		"none",
		`This flag is used for migration`,
	)

	flag.Parse()

	if m == "migrate" {
		migration.Migrate()
	}

	if s == "seeder" {

	}

	f := factory.NewFactory() // Database instance initialization
	g := gin.New()

	http.NewHttp(g, f)

	if err := g.Run(":" + util.GetEnv("APP_PORT", "8080")); err != nil {
		log.Fatal("Can't start server.")
	}
}
