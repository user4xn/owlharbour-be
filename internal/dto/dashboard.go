package dto

type (
	DashboardStatisticResponse struct {
		TotalCheckin  int `json:"total_checkin"`
		TotalCheckout int `json:"total_checkout"`
		TotalShip     int `json:"total_ship"`
		TotalFraud    int `json:"total_fraud"`
	}
)
