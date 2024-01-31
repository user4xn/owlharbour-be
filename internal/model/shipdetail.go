package model

type ShipDetail struct {
	ShipID    int      `gorm:"primaryKey"`
	Type      ShipType `gorm:"enum:kapal angkut,kapal tangkap"`
	Dimension string   `gorm:"varchar"`
	Harbour   string   `gorm:"varchar"`
	SIUP      string   `gorm:"varchar"`
	BKP       string   `gorm:"varchar"`
	SelarMark string   `gorm:"varchar"`
	GT        string   `gorm:"varchar"`
	OwnerName string   `gorm:"varchar"`
}

func (ShipDetail) TableName() string {
	return "ship_details"
}
