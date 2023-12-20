package dto

type (
	NeedCheckupShipParam struct {
		Offset int    `json:"offset"`
		Limit  int    `json:"limit"`
		Search string `json:"search"`
	}

	NeedCheckupShipResponse struct {
		LogID       int    `json:"log_id"`
		ShipName    string `json:"ship_name"`
		Lat         string `json:"lat"`
		Long        string `json:"long"`
		IsInspected int    `json:"is_inspected"`
		IsReported  int    `json:"is_reported"`
		CheckinDate string `json:"checkin_date"`
	}

	ShipCheckupRequest struct {
		IsInspected bool `json:"is_inspected"`
		IsReported  bool `json:"is_reported"`
	}
)
