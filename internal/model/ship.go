package model

type Ship struct {
	Common
	Name            string     `gorm:"varchar"`
	Phone           string     `gorm:"varchar"`
	ResponsibleName string     `gorm:"varchar"`
	DeviceID        string     `gorm:"varchar"`
	FirebaseToken   string     `gorm:"varchar"`
	Status          ShipStatus `gorm:"enum:checkin,checkout,out of scope"`
	CurrentLat      string     `gorm:"varchar"`
	CurrentLong     string     `gorm:"varchar"`
	OnGround        int
}

func (Ship) TableName() string {
	return "ships"
}
