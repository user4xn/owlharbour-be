package dto

type (
	PairingRequest struct {
		HarbourCode     int    `json:"harbour_code" binding:"required"`
		ShipName        string `json:"ship_name" binding:"required"`
		Phone           string `json:"phone" binding:"required"`
		Username        string `json:"username" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ResponsibleName string `json:"responsible_name" binding:"required"`
		DeviceID        string `json:"device_id" binding:"required"`
		FirebaseToken   string `json:"firebase_token" binding:"required"`
	}

	PairingActionRequest struct {
		PairingID string `json:"pairing_id" binding:"required"`
		Status    string `json:"status" binding:"required"`
	}

	PairingRequestResponseList struct {
		Total int                      `json:"total"`
		Data  []PairingRequestResponse `json:"data"`
	}

	PairingRequestResponse struct {
		ID              int    `json:"id"`
		ShipName        string `json:"ship_name"`
		Phone           string `json:"phone"`
		Username        string `json:"username"`
		Password        string `json:"password"`
		ResponsibleName string `json:"responsible_name"`
		DeviceID        string `json:"device_id"`
		FirebaseToken   string `json:"firebase_token"`
		Status          string `json:"status"`
		CreatedAt       string `json:"created_at"`
	}

	PairingToNewShip struct {
		Phone           string `json:"phone"`
		ResponsibleName string `json:"responsible_name"`
		ShipName        string `json:"ship_name"`
		DeviceID        string `json:"device_id"`
		FirebaseToken   string `json:"firebase_token"`
		UserID          int    `json:"user_id"`
	}

	PairingListParam struct {
		Offset int      `json:"offset"`
		Limit  int      `json:"limit"`
		Status []string `json:"status"`
		Search string   `json:"search"`
	}

	DetailPairingResponse struct {
		ShipName       string           `json:"ship_name"`
		ReponsibleName string           `json:"responsible_name"`
		Phone          string           `json:"phone"`
		Username       string           `json:"username"`
		DeviceID       string           `json:"device_id"`
		Status         string           `json:"status"`
		SubmittedAt    string           `json:"submitted_at"`
		RespondedAt    string           `json:"responded_at"`
		HistoryPairing []HistoryPairing `json:"history_pairing"`
	}

	HistoryPairing struct {
		ShipName       string `json:"ship_name"`
		ReponsibleName string `json:"responsible_name"`
		Phone          string `json:"phone"`
		Username       string `json:"username"`
		Status         string `json:"status"`
		SubmittedAt    string `json:"submitted_at"`
		RespondedAt    string `json:"responded_at"`
	}
)
