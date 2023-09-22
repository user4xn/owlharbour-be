package model

type AppSetting struct {
	HarbourCode     int      `gorm:"integer"`
	HarbourName     string   `gorm:"varchar"`
	Mode            ModeType `gorm:"enum:interval,range"`
	ApkMinVersion   string   `gorm:"varchar"`
	Interval        int      `gorm:"integer"`
	Range           int      `gorm:"integer"`
	ApkDownloadLink string   `gorm:"varchar"`
}

func (AppSetting) TableName() string {
	return "app_setting"
}
