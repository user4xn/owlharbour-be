package repository

import (
	"context"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"

	"gorm.io/gorm"
)

type App interface {
	AppInfo(ctx context.Context) (*dto.AppInfo, error)
}

type app struct {
	Db *gorm.DB
}

func NewAppRepository(db *gorm.DB) App {
	return &app{
		Db: db,
	}
}

func (r *app) AppInfo(ctx context.Context) (*dto.AppInfo, error) {
	querySetting := r.Db.Model(&model.AppSetting{})
	queryGeofence := r.Db.Model(&model.AppGeofence{})

	var setting model.AppSetting
	var geofence []model.AppGeofence

	if err := querySetting.Find(&setting).Error; err != nil {
		return nil, err
	}

	if err := queryGeofence.Find(&geofence).Error; err != nil {
		return nil, err
	}

	var geofences []dto.AppGeofence
	for _, e := range geofence {
		geofences = append(geofences, dto.AppGeofence{
			Long: e.Long,
			Lat:  e.Lat,
		})
	}

	res := &dto.AppInfo{
		HarbourCode:     setting.HarbourCode,
		HarbourName:     setting.HarbourName,
		Mode:            setting.Mode.String(),
		Interval:        setting.Interval,
		Range:           setting.Range,
		ApkDownloadLink: setting.ApkDownloadLink,
		Geofence:        geofences,
	}

	return res, nil
}
