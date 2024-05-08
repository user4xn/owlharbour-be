package main

import (
	"context"
	"flag"
	"log"
	"simpel-api/database"
	"simpel-api/database/migration"
	"simpel-api/database/seeder"
	"simpel-api/internal/app/ship"
	"simpel-api/internal/factory"
	"simpel-api/internal/http"
	"simpel-api/pkg/util"

	"github.com/gin-gonic/gin"
)

func main() {
	var m string
	var s string
	var c string

	database.CreateConnection()

	flag.StringVar(
		&m,
		"m",
		"none",
		`This flag is used for migration`,
	)

	flag.StringVar(
		&c,
		"c",
		"",
		`This flag is used for consumer`,
	)

	flag.StringVar(
		&s,
		"s",
		"none",
		`This flag is used for seeder`,
	)

	flag.Parse()

	if m == "all" {
		migration.Migrate()
		return
	}

	if s == "seed" {
		seeder.Seed()
		return
	}

	f := factory.NewFactory() // Database instance initialization

	if c != "" {
		if c == "ship" {
			ctx := context.Background()
			ship.NewHandler(f).Init()
			ship.NewHandler(f).WorkerRecordLog(ctx)
		}

		return
	}

	g := gin.New()

	http.NewHttp(g, f)

	if err := g.Run(":" + util.GetEnv("APP_PORT", "8080")); err != nil {
		log.Fatal("Can't start server.")
	}
}
