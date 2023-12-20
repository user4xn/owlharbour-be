package dto

type (
	AppInfo struct {
		HarbourCode     int    `json:"harbour_code"`
		HarbourName     string `json:"harbour_name"`
		Mode            string `json:"mode"`
		Interval        int    `json:"interval"`
		Range           int    `json:"range"`
		ApkDownloadLink string `json:"apk_download_link"`
		Geofence        []AppGeofence
	}

	AppGeofence struct {
		Long string `json:"long"`
		Lat  string `json:"lat"`
	}

	GetDataSetting struct {
		HarbourCode int           `json:"harbour_code"`
		HarbourName string        `json:"harbour_name"`
		Mode        string        `json:"mode"`
		Interval    int           `json:"interval"`
		Range       int           `json:"range"`
		Geofences   []AppGeofence `json:"geofences"`
	}

	GetDataSettingWeb struct {
		HarbourCode int           `json:"harbour_code"`
		HarbourName string        `json:"harbour_name"`
		Mode        string        `json:"mode"`
		Interval    int           `json:"interval"`
		Range       int           `json:"range"`
		Geofences   []AppGeofence `json:"geofences"`
	}

	PayloadStoreSetting struct {
		HarbourCode int                  `json:"harbour_code" binding:"required"`
		HarbourName string               `json:"harbour_name" binding:"required"`
		Mode        string               `json:"mode" binding:"required"`
		Interval    int                  `json:"interval" binding:"required"`
		Range       int                  `json:"range" binding:"required"`
		Geofence    []PayloadAppGeofence `json:"geofence"`
	}

	PayloadAppGeofence struct {
		Long string `json:"long"`
		Lat  string `json:"lat"`
	}
)
