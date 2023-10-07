package repository

import (
	"context"
	"simpel-api/internal/model"

	"gorm.io/gorm"
)

type AppGeofence interface {
	Store(ctx context.Context, data model.AppGeofence) error
	DeleteAll(ctx context.Context) error
}

type appgeofence struct {
	Db *gorm.DB
}

func NewAppGeofenceRepository(db *gorm.DB) AppGeofence {
	return &appgeofence{
		Db: db,
	}
}

func (r *appgeofence) Store(ctx context.Context, data model.AppGeofence) error {
	tx := r.Db.WithContext(ctx)
	if err := tx.Model(model.AppGeofence{}).Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *appgeofence) DeleteAll(ctx context.Context) error {
	db := r.Db.WithContext(ctx)
	if err := db.Exec("DELETE FROM app_geofences").Error; err != nil {
		return err
	}
	return nil
}
