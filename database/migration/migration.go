package migration

import (
	"simpel-api/database"
	"simpel-api/internal/model"
)

var tables = []interface{}{
	&model.User{},
	&model.AppSetting{},
	&model.AppGeofence{},
	&model.Ship{},
	&model.ShipLocationLog{},
	&model.ShipDockedLog{},
}

func Migrate() {
	conn := database.GetConnection() // Get db connection
	conn.AutoMigrate(tables...)      // migrate the tables
}
