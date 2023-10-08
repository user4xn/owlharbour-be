package repository

import (
	"context"
	"simpel-api/internal/model"
	"simpel-api/pkg/util"

	"gorm.io/gorm"
)

type AppSetting interface {
	Store(ctx context.Context, data model.AppSetting) error
	FindLatest(ctx context.Context, selectedFields string) (model.AppSetting, error)
	Update(ctx context.Context, updatedModels *model.AppSetting, updatedField string, query string, args ...interface{}) error
}

type appsetting struct {
	Db *gorm.DB
}

func NewAppSettingRepository(db *gorm.DB) AppSetting {
	return &appsetting{
		Db: db,
	}
}

func (r *appsetting) Store(ctx context.Context, data model.AppSetting) error {
	tx := r.Db.WithContext(ctx)
	if err := tx.Model(model.AppSetting{}).Create(&data).Error; err != nil {
		return err
	}

	return nil
}

func (r *appsetting) FindLatest(ctx context.Context, selectedFields string) (model.AppSetting, error) {
	var res model.AppSetting

	db := r.Db.WithContext(ctx).Model(model.AppSetting{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Limit(1).Take(&res).Error; err != nil {
		return model.AppSetting{}, err
	}

	return res, nil
}

func (r *appsetting) Update(ctx context.Context, updatedModels *model.AppSetting, updatedField string, query string, args ...interface{}) error {
	modelDb := r.Db.WithContext(ctx).Model(&model.AppSetting{})
	if err := util.SetSelectFields(modelDb, updatedField).Where(query, args...).Updates(updatedModels).Error; err != nil {
		return err
	}
	return nil
}
