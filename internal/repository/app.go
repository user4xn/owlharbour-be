package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App interface {
	AppInfo(ctx context.Context) (*dto.AppInfo, error)
	GetPolygon(ctx context.Context) ([]dto.HarbourGeofences, error)
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
		HarbourCode:     setting.HarbourCode,
		HarbourName:     setting.HarbourName,
		Mode:            setting.Mode.String(),
		Interval:        setting.Interval,
		Range:           setting.Range,
		ApkDownloadLink: setting.ApkDownloadLink,
		Geofence:        geofences,
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
