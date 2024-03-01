package model

type AppSetting struct {
	HarbourCode  int      `gorm:"integer"`
	HarbourName  string   `gorm:"varchar"`
	Mode         ModeType `gorm:"enum:interval,range"`
	Interval     int      `gorm:"integer"`
	Range        int      `gorm:"integer"`
	AdminContact string   `gorm:"varchar"`
}

func (AppSetting) TableName() string {
	return "app_setting"
}
