package repository

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"
	"simpel-api/pkg/helper"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Ship interface {
	StoreNewShip(ctx context.Context, request dto.PairingRequestResponse) error
	ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error)
	ShipByDevice(ctx context.Context, DeviceID int) (*dto.ShipMobileDetailResponse, error)
	ShipByID(ctx context.Context, ShipID int) (*model.Ship, error)
	GetLastDockedLog(ctx context.Context, ShipID int) (*dto.ShipDockedLog, error)
	StoreDockedLog(ctx context.Context, request dto.ShipDockedLogStore) error
	StoreLocationLog(ctx context.Context, request dto.ShipLocationLogStore) error
	UpdateShip(ctx context.Context, request model.Ship) error
	UpdateShipDetail(ctx context.Context, request dto.ShipAddonDetailRequest) error
	ShipDockedLogs(ctx context.Context, ShipID int) ([]dto.DockLogsShip, error)
	ShipLocationLogs(ctx context.Context, ShipID int) ([]dto.LocationLogsShip, error)
	ShipAddonDetail(ctx context.Context, ShipID int) (*dto.ShipAddonDetailResponse, error)
}

type ship struct {
	Db           *gorm.DB
	RedisClient  *redis.Client
	CacheEnabled bool
}

func NewShipRepository(db *gorm.DB, redisClient *redis.Client) Ship {
	return &ship{
		Db:           db,
		RedisClient:  redisClient,
		CacheEnabled: true,
	}
}

func (r *ship) StoreNewShip(ctx context.Context, request dto.PairingRequestResponse) error {
	shipModel := model.Ship{
		Name:            request.ShipName,
		Phone:           request.Phone,
		ResponsibleName: request.ResponsibleName,
		DeviceID:        request.DeviceID,
		FirebaseToken:   request.FirebaseToken,
		Status:          "out of scope",
	}

	tx := r.Db.WithContext(ctx).Begin()

	if err := tx.Create(&shipModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	cacheKey := "ship_list-"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *ship) ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error) {
	paramJSON, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	hash := sha1.Sum(paramJSON)
	uniqueString := fmt.Sprintf("%x", hash)

	cacheKey := "ship_list-" + uniqueString

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo []dto.ShipResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.Ship{})

	if request.Status != nil && request.Status[0] != "" && len(request.Status) > 0 {
		query = query.Where("status IN (?)", request.Status)
	}

	if request.Search != "" {
		searchLower := strings.ToLower(request.Search)
		query = query.Where("lower(name) LIKE ? OR lower(device_id) LIKE ? OR lower(phone) LIKE ?", "%"+searchLower+"%", "%"+searchLower+"%", "%"+searchLower+"%")
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("created_at DESC")

	var ship []model.Ship

	if err := query.Find(&ship).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var shipList []dto.ShipResponse
	for _, e := range ship {
		shipList = append(shipList, dto.ShipResponse{
			ID:              e.ID,
			ShipName:        e.Name,
			ResponsibleName: e.ResponsibleName,
			DeviceID:        e.DeviceID,
			OnGround:        e.OnGround,
			Status:          string(e.Status),
			CreatedAt:       e.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(shipList)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return shipList, nil
}

func (r *ship) ShipByDevice(ctx context.Context, DeviceID int) (*dto.ShipMobileDetailResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var ship model.Ship
	err := tx.Where("device_id = ?", DeviceID).First(&ship).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	shipDetail := dto.ShipMobileDetailResponse{
		ID:              ship.ID,
		ShipName:        ship.Name,
		ResponsibleName: ship.ResponsibleName,
		DeviceID:        ship.DeviceID,
		CurrentLong:     ship.CurrentLong,
		CurrentLat:      ship.CurrentLat,
		FirebaseToken:   ship.FirebaseToken,
		Status:          string(ship.Status),
		OnGround:        ship.OnGround,
		CreatedAt:       ship.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &shipDetail, nil
}

func (r *ship) GetLastDockedLog(ctx context.Context, ShipID int) (*dto.ShipDockedLog, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var log model.ShipDockedLog
	err := tx.Where("ship_id = ?", ShipID).Order("created_at DESC").First(&log).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	logDock := dto.ShipDockedLog{
		ID:        log.ID,
		Long:      log.Long,
		Lat:       log.Lat,
		Status:    string(log.Status),
		CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &logDock, nil
}

func (r *ship) StoreDockedLog(ctx context.Context, request dto.ShipDockedLogStore) error {
	tx := r.Db.WithContext(ctx).Begin()

	dockedModel := model.ShipDockedLog{
		ShipID: request.ShipID,
		Long:   request.Long,
		Lat:    request.Lat,
		Status: model.ShipStatus(request.Status),
	}

	if err := tx.Create(&dockedModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ship) StoreLocationLog(ctx context.Context, request dto.ShipLocationLogStore) error {
	tx := r.Db.WithContext(ctx).Begin()

	locationModel := model.ShipLocationLog{
		ShipID:   request.ShipID,
		Long:     request.Long,
		Lat:      request.Lat,
		OnGround: request.OnGround,
		IsMocked: request.IsMocked,
	}

	if err := tx.Create(&locationModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ship) UpdateShip(ctx context.Context, request model.Ship) error {
	tx := r.Db.WithContext(ctx).Begin()

	updateFields := map[string]interface{}{
		"status":       model.ShipStatus(request.Status),
		"current_lat":  request.CurrentLat,
		"current_long": request.CurrentLong,
		"on_ground": func() int {
			if request.OnGround == 1 {
				return 1
			}
			return 0
		}(),
	}

	if err := tx.Model(&model.Ship{}).Where("id = ?", request.ID).Updates(updateFields).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	cacheKey := "ship_list-"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *ship) UpdateShipDetail(ctx context.Context, request dto.ShipAddonDetailRequest) error {
	tx := r.Db.WithContext(ctx).Begin()

	var existingShip model.Ship
	if err := tx.Where("id = ?", request.ShipID).First(&existingShip).Error; err != nil {
		tx.Rollback()
		return err
	}

	shipDetailModel := model.ShipDetail{
		ShipID:    request.ShipID,
		Type:      model.ShipType(request.Type),
		Dimension: request.Dimension,
		Harbour:   request.Harbour,
		SIUP:      request.SIUP,
		BKP:       request.BKP,
		SelarMark: request.SelarMark,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ship_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "dimension", "harbour", "siup", "bkp", "selar_mark"}),
	}).Create(&shipDetailModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *ship) ShipByID(ctx context.Context, ShipID int) (*model.Ship, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var ship model.Ship
	err := tx.Where("id = ?", ShipID).First(&ship).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &ship, nil
}

func (r *ship) ShipDockedLogs(ctx context.Context, ShipID int) ([]dto.DockLogsShip, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var logs []model.ShipDockedLog
	err := tx.Where("ship_id = ?", ShipID).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var logDock []dto.DockLogsShip
	for _, log := range logs {
		logDock = append(logDock, dto.DockLogsShip{
			LogID:     log.ID,
			Long:      log.Long,
			Lat:       log.Lat,
			Status:    string(log.Status),
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return logDock, nil
}

func (r *ship) ShipLocationLogs(ctx context.Context, ShipID int) ([]dto.LocationLogsShip, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var logs []model.ShipLocationLog
	err := tx.Where("ship_id = ?", ShipID).Order("created_at DESC").Find(&logs).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var logDock []dto.LocationLogsShip
	for _, log := range logs {
		logDock = append(logDock, dto.LocationLogsShip{
			LogID:     log.ID,
			Long:      log.Long,
			Lat:       log.Lat,
			IsMocked:  log.IsMocked,
			OnGround:  log.OnGround,
			CreatedAt: log.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return logDock, nil
}

func (r *ship) ShipAddonDetail(ctx context.Context, ShipID int) (*dto.ShipAddonDetailResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var detail model.ShipDetail
	err := tx.Where("ship_id = ?", ShipID).First(&detail).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	res := dto.ShipAddonDetailResponse{
		Type:      string(detail.Type),
		Dimension: detail.Dimension,
		Harbour:   detail.Harbour,
		SIUP:      detail.SIUP,
		BKP:       detail.BKP,
		SelarMark: detail.SelarMark,
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &res, nil
}
