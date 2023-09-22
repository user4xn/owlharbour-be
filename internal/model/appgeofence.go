package model

type AppGeofence struct {
	Long string `gorm:"varchar"`
	Lat  string `gorm:"varchar"`
}

func (AppGeofence) TableName() string {
    return "app_geofences"
}