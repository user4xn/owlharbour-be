package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/model"
	"owlharbour-api/pkg/helper"
	"owlharbour-api/pkg/util"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App interface {
	AppInfo(ctx context.Context) (*dto.AppInfo, error)
	GetPolygon(ctx context.Context) ([]dto.HarbourGeofences, error)
	StoreSetting(ctx context.Context, data model.AppSetting) error
	FindLatestSetting(ctx context.Context, selectedFields string) (model.AppSetting, error)
	UpsertSetting(ctx context.Context, updatedModels *model.AppSetting, updatedField string, query string, args ...interface{}) error
	StoreGeofence(ctx context.Context, data model.AppGeofence) error
	DeleteAllGeofence(ctx context.Context) error
}

type app struct {
	Db           *gorm.DB
	RedisClient  *redis.Client
	CacheEnabled bool
}

func NewAppRepository(db *gorm.DB, redisClient *redis.Client) App {
	return &app{
		Db:           db,
		RedisClient:  redisClient,
		CacheEnabled: true,
	}
}

func (r *app) AppInfo(ctx context.Context) (*dto.AppInfo, error) {
	cacheKey := "app_info"

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo dto.AppInfo
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return &cachedInfo, nil
			}
		}
	}

	querySetting := r.Db.Model(&model.AppSetting{})
	queryGeofence := r.Db.Model(&model.AppGeofence{})

	var setting model.AppSetting
	var geofence []model.AppGeofence

	if err := querySetting.First(&setting).Error; err != nil {
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
		HarbourCode: setting.HarbourCode,
		HarbourName: setting.HarbourName,
		Mode:        setting.Mode.String(),
		Interval:    setting.Interval,
		Range:       setting.Range,
		Geofence:    geofences,
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(res)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return res, nil
}

func (r *app) GetPolygon(ctx context.Context) ([]dto.HarbourGeofences, error) {
	cacheKey := "app_polygon"

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo []dto.HarbourGeofences
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	query := r.Db.Model(&model.AppGeofence{})

	var geofence []model.AppGeofence

	if err := query.Find(&geofence).Error; err != nil {
		return nil, err
	}

	var polygon []dto.HarbourGeofences
	for _, e := range geofence {
		polygon = append(polygon, dto.HarbourGeofences{
			Long: e.Long,
			Lat:  e.Lat,
		})
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(polygon)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return polygon, nil
}

func (r *app) StoreSetting(ctx context.Context, data model.AppSetting) error {
	tx := r.Db.WithContext(ctx)
	if err := tx.Model(model.AppSetting{}).Create(&data).Error; err != nil {
		return err
	}

	cacheKey := "app_info"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *app) FindLatestSetting(ctx context.Context, selectedFields string) (model.AppSetting, error) {
	var res model.AppSetting

	query := r.Db.WithContext(ctx).Model(model.AppSetting{})
	query = util.SetSelectFields(query, selectedFields)

	if err := query.Limit(1).Take(&res).Error; err != nil {
		return model.AppSetting{}, err
	}

	return res, nil
}

func (r *app) UpsertSetting(ctx context.Context, updatedModels *model.AppSetting, updatedField string, query string, args ...interface{}) error {
	setting := r.Db.WithContext(ctx).Model(&model.AppSetting{})
	cacheKey := "app_info"
	var count int64
	if err := r.Db.WithContext(ctx).Model(&model.AppSetting{}).Where(query, args...).Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		if err := util.SetSelectFields(setting, updatedField).Where(query, args...).Updates(updatedModels).Error; err != nil {
			return err
		}
	} else {
		if err := r.Db.WithContext(ctx).Create(updatedModels).Error; err != nil {
			return err
		}
	}

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return err
	}

	return nil
}

func (r *app) StoreGeofence(ctx context.Context, data model.AppGeofence) error {
	tx := r.Db.WithContext(ctx)
	if err := tx.Model(model.AppGeofence{}).Create(&data).Error; err != nil {
		return err
	}

	cacheKey := "app_polygon"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *app) DeleteAllGeofence(ctx context.Context) error {
	db := r.Db.WithContext(ctx)

	if err := db.Exec("DELETE FROM app_geofences").Error; err != nil {
		return err
	}

	return nil
}
