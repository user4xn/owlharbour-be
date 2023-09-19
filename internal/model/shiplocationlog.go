package model

type ShipLocationLog struct {
	Common
	ShipID   int
	Long     string `gorm:"varchar"`
	Lat      string `gorm:"varchar"`
	IsMocked int
	OnGround int
}

func (ShipLocationLog) TableName() string {
    return "ship_location_logs"
}