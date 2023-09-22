package dto

type AppInfo struct {
	HarbourCode     int    `json:"harbour_code"`
	HarbourName     string `json:"harbour_name"`
	Mode            string `json:"mode"`
	Interval        int    `json:"interval"`
	Range           int    `json:"range"`
	ApkDownloadLink string `json:"apk_download_link"`
	Geofence        []AppGeofence
}

type AppGeofence struct {
	Long string `json:"long"`
	Lat  string `json:"lat"`
}
