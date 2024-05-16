package migration

import (
	"owlharbour-api/database"
	"owlharbour-api/internal/model"
)

var tables = []interface{}{
	&model.User{},
	&model.AppSetting{},
	&model.AppGeofence{},
	&model.PairingRequest{},
	&model.Ship{},
	&model.ShipDetail{},
	&model.ShipLocationLog{},
	&model.ShipDockedLog{},
}

func Migrate() {
	conn := database.GetConnection() // Get db connection
	conn.AutoMigrate(tables...)      // migrate the tables
}
