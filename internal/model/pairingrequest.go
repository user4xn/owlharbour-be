package model

type PairingRequest struct {
	Common
	Name          string        `gorm:"varchar"`
	Phone         string        `gorm:"varchar"`
	DeviceID      string        `gorm:"varchar"`
	FirebaseToken string        `gorm:"varchar"`
	Status        PairingStatus `gorm:"enum:pending,approved,rejected"`
}

func (PairingRequest) TableName() string {
	return "pairing_requests"
}
