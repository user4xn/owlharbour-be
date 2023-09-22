package dto

type PairingRequest struct {
	HarbourCode   int    `json:"harbour_code" binding:"required"`
	ShipName      string `json:"ship_name" binding:"required"`
	Phone         string `json:"phone" binding:"required"`
	DeviceID      string `json:"device_id" binding:"required"`
	FirebaseToken string `json:"firebase_token" binding:"required"`
}

type PairingActionRequest struct {
	PairingID int    `json:"pairing_id" binding:"required"`
	Status    string `json:"status" binding:"required"`
}

type PairingRequestResponse struct {
	ID            int    `json:"id"`
	ShipName      string `json:"ship_name"`
	Phone         string `json:"phone"`
	DeviceID      string `json:"device_id"`
	FirebaseToken string `json:"firebase_token"`
	Status        string `json:"status"`
	CreatedAt     string `json:"created_at"`
}

type PairingListParam struct {
	Offset int      `json:"offset"`
	Limit  int      `json:"limit"`
	Status []string `json:"status"`
	Search string   `json:"search"`
}
