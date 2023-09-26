package dto

type ShipListParam struct {
	Offset int      `json:"offset"`
	Limit  int      `json:"limit"`
	Status []string `json:"status"`
	Search string   `json:"search"`
}

type ShipResponse struct {
	ID              int    `json:"id"`
	ShipName        string `json:"ship_name"`
	ResponsibleName string `json:"responsible_name"`
	DeviceID        string `json:"device_id"`
	Status          string `json:"status"`
	OnGround        int    `json:"on_ground"`
	CreatedAt       string `json:"created_at"`
}

type ShipMobileDetailResponse struct {
	ID              int    `json:"id"`
	ShipName        string `json:"ship_name"`
	ResponsibleName string `json:"responsible_name"`
	DeviceID        string `json:"device_id"`
	CurrentLong     string `json:"current_long"`
	CurrentLat      string `json:"current_lat"`
	FirebaseToken   string `json:"firebase_token"`
	Status          string `json:"status"`
	OnGround        int    `json:"on_ground"`
	CreatedAt       string `json:"created_at"`
	HitMode         string `json:"hit_mode"`
	Range           int    `json:"range"`
	Interval        int    `json:"interval"`
}

type ShipDetailResponse struct {
	ID              int                     `json:"id"`
	ShipName        string                  `json:"ship_name"`
	ResponsibleName string                  `json:"responsible_name"`
	DeviceID        string                  `json:"device_id"`
	DetailShip      ShipAddonDetailResponse `json:"detail"`
	CurrentLong     string                  `json:"current_long"`
	CurrentLat      string                  `json:"current_lat"`
	FirebaseToken   string                  `json:"firebase_token"`
	Status          string                  `json:"status"`
	OnGround        int                     `json:"on_ground"`
	CreatedAt       string                  `json:"created_at"`
	DockLogs        []DockLogsShip          `json:"dock_logs"`
	LocationLogs    []LocationLogsShip      `json:"location_logs"`
}

type ShipAddonDetailResponse struct {
	Type      string `json:"type"`
	Dimension string `json:"dimension"`
	Harbour   string `json:"harbour"`
	SIUP      string `json:"siup"`
	BKP       string `json:"bkp"`
	SelarMark string `json:"selar_mark"`
}

type DockLogsShip struct {
	LogID     int    `json:"log_id"`
	Long      string `json:"long"`
	Lat       string `json:"lat"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type LocationLogsShip struct {
	LogID     int    `json:"log_id"`
	Long      string `json:"long"`
	Lat       string `json:"lat"`
	IsMocked  int    `json:"is_mocked"`
	OnGround  int    `json:"on_ground"`
	CreatedAt string `json:"created_at"`
}

type ShipAddonDetailRequest struct {
	ShipID    int    `json:"ship_id" binding:"required"`
	Type      string `json:"type"`
	Dimension string `json:"dimension"`
	Harbour   string `json:"harbour"`
	SIUP      string `json:"siup"`
	BKP       string `json:"bkp"`
	SelarMark string `json:"selar_mark"`
}

type ShipRecordRequest struct {
	DeviceID int    `json:"device_id" binding:"required"`
	Long     string `json:"long" binding:"required"`
	Lat      string `json:"lat" binding:"required"`
	IsMocked int    `json:"is_mocked"`
}

type ShipDockedLog struct {
	ID        int    `json:"id"`
	Long      string `json:"long"`
	Lat       string `json:"lat"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

type ShipDockedLogStore struct {
	ShipID int    `json:"ship_id"`
	Long   string `json:"long"`
	Lat    string `json:"lat"`
	Status string `json:"status"`
}
type ShipLocationLogStore struct {
	ShipID   int    `json:"ship_id"`
	Long     string `json:"long"`
	Lat      string `json:"lat"`
	IsMocked int    `json:"is_mocked"`
	OnGround int    `json:"on_ground"`
}
