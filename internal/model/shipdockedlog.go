package model

type ShipDockedLog struct {
	Common
	ShipID      int
	Long        string     `gorm:"varchar"`
	Lat         string     `gorm:"varchar"`
	Status      ShipStatus `gorm:"enum:checkin,checkout"`
	IsInspected int
	IsReported  int
}

func (ShipDockedLog) TableName() string {
	return "ship_docked_logs"
}
