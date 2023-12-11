package dto

type (
	DashboardStatisticResponse struct {
		TotalCheckin  int `json:"total_checkin"`
		TotalCheckout int `json:"total_checkout"`
		TotalShip     int `json:"total_ship"`
		TotalFraud    int `json:"total_fraud"`
	}

	ShipTerrainResponse struct {
		OnGround int64 `json:"on_ground"`
		OnWater  int64 `json:"on_water"`
	}

	LogsStatisticResponse struct {
		CheckIN int64 `json:"checkin"`
		CheckOUT int64 `json:"checkout"`
		Fraud int64 `json:"fraud"`
	}
)
