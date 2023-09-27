package dto

type (
	PairingRequest struct {
		HarbourCode     int    `json:"harbour_code" binding:"required"`
		ShipName        string `json:"ship_name" binding:"required"`
		Phone           string `json:"phone" binding:"required"`
		ResponsibleName string `json:"responsible_name" binding:"required"`
		DeviceID        string `json:"device_id" binding:"required"`
		FirebaseToken   string `json:"firebase_token" binding:"required"`
	}

	PairingActionRequest struct {
		PairingID int    `json:"pairing_id" binding:"required"`
		Status    string `json:"status" binding:"required"`
	}

	PairingRequestResponse struct {
		ID              int    `json:"id"`
		ShipName        string `json:"ship_name"`
		Phone           string `json:"phone"`
		ResponsibleName string `json:"responsible_name"`
		DeviceID        string `json:"device_id"`
		FirebaseToken   string `json:"firebase_token"`
		Status          string `json:"status"`
		CreatedAt       string `json:"created_at"`
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
		Status         string `json:"status"`
		SubmittedAt    string `json:"submitted_at"`
		RespondedAt    string `json:"responded_at"`
	}
)
