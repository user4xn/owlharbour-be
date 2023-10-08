package dto

type (
	GetDataSetting struct {
		HarbourCode     int    `json:"harbour_code"`
		HarbourName     string `json:"harbour_name"`
		Mode            string `json:"mode"`
		ApkMinVersion   string `json:"apk_min_version"`
		Interval        int    `json:"interval"`
		Range           int    `json:"range"`
		ApkDownloadLink string `json:"apk_min_download"`
	}

	GetDataSettingWeb struct {
		HarbourCode     int           `json:"harbour_code"`
		HarbourName     string        `json:"harbour_name"`
		Mode            string        `json:"mode"`
		ApkMinVersion   string        `json:"apk_min_version"`
		Interval        int           `json:"interval"`
		Range           int           `json:"range"`
		ApkDownloadLink string        `json:"apk_min_download"`
		Geofances       []GetGeofance `json:"geofances"`
	}

	GetGeofance struct {
		Long string `json:"long"`
		Lat  string `json:"lat"`
	}

	PayloadStoreSetting struct {
		HarbourCode     int                  `json:"harbour_code" binding:"required"`
		HarbourName     string               `json:"harbour_name" binding:"required"`
		Mode            string               `json:"mode" binding:"required"`
		ApkMinVersion   string               `json:"apk_min_version" binding:"required"`
		Interval        int                  `json:"interval" binding:"required"`
		Range           int                  `json:"range" binding:"required"`
		ApkDownloadLink string               `json:"apk_download_link" binding:"required"`
		Geofence        []PayloadAppGeofence `json:"geofence"`
	}

	PayloadAppGeofence struct {
		Long string `json:"long"`
		Lat  string `json:"lat"`
	}
)
