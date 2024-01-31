package repository

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"
	"simpel-api/pkg/helper"
	"simpel-api/pkg/util"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Ship interface {
	StoreNewShip(ctx context.Context, request dto.PairingToNewShip) error
	ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error)
	ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipMobileDetailResponse, error)
	ShipByAuth(ctx context.Context, authUser model.User) (*dto.ShipMobileDetailResponse, error)
	ShipByID(ctx context.Context, ShipID int) (*model.Ship, error)
	GetLastDockedLog(ctx context.Context, ShipID int) (*dto.ShipDockedLog, error)
	StoreDockedLog(ctx context.Context, request dto.ShipDockedLogStore) error
	StoreLocationLog(ctx context.Context, request dto.ShipLocationLogStore) error
	UpdateShip(ctx context.Context, request model.Ship) error
	UpdateShipDetail(ctx context.Context, request dto.ShipAddonDetailRequest) error
	ShipDockedLogs(ctx context.Context, ShipID int, request *dto.ShipLogParam) ([]dto.DockLogsShip, error)
	ShipLocationLogs(ctx context.Context, ShipID int, request *dto.ShipLogParam) ([]dto.LocationLogsShip, error)
	ShipAddonDetail(ctx context.Context, ShipID int) (dto.ShipAddonDetailResponse, error)
	CountShip(ctx context.Context) (int64, error)
	CountStatistic(ctx context.Context) ([]int64, error)
	LastUpdated(ctx context.Context) (time.Time, error)
	ShipInBatch(ctx context.Context, start int, end int) (*[]model.Ship, bool, error)
	ReportShipDocking(ctx context.Context, request dto.ReportShipDockedParam) ([]dto.ReportShipDockingResponse, error)
	ReportShipFraud(ctx context.Context, request dto.ReportShipLocationParam) ([]dto.ReportShipLocationResponse, error)
	CountShipByTerrain(ctx context.Context, onGround int) (int64, error)
	CountShipByStatus(ctx context.Context, startDate string, endDate string, status string) (int64, error)
	CountShipFraud(ctx context.Context, startDate string, endDate string) (int64, error)
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Ship, error)
	FindOneDockedLog(ctx context.Context, selectedFields string, query string, args ...any) (model.ShipDockedLog, error)
	UpdateShipDeviceID(ctx context.Context, deviceID string, user_id int) error
	UpdateShipCheckup(ctx context.Context, request dto.ShipCheckupRequest, id int, data model.ShipDockedLog) error
	NeedCheckupShip(ctx context.Context, request dto.NeedCheckupShipParam) ([]dto.NeedCheckupShipResponse, error)
	LastestDockedShip(ctx context.Context, limit int) ([]dto.DashboardLastDockedShipResponse, error)
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

func (r *ship) LastestDockedShip(ctx context.Context, limit int) ([]dto.DashboardLastDockedShipResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.ShipDockedLog{}).
		Select("ship_docked_logs.*, ships.name as ship_name, ships.id as ship_id, ships.responsible_name as responsible_name, ships.phone as phone").
		Joins("JOIN ships ON ship_docked_logs.ship_id = ships.id")

	limitParam := 10
	if limit != 0 {
		limitParam = limit
	}

	query = query.Where("ship_docked_logs.status = ?", "checkin").
		Limit(limitParam).
		Order("ship_docked_logs.created_at DESC")

	var result []struct {
		model.ShipDockedLog
		ShipName        string `json:"ship_name"`
		ResponsibleName string `json:"resposible_name"`
		Phone           string `json:"phone"`
	}

	if err := query.Find(&result).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var shipDock []dto.DashboardLastDockedShipResponse
	for _, e := range result {
		shipDock = append(shipDock, dto.DashboardLastDockedShipResponse{
			ShipName:        e.ShipName,
			Phone:           e.Phone,
			ResponsibleName: e.ResponsibleName,
			CheckinDate:     e.CreatedAt.Format("2006-01-02 15:04:05"),
			IsInspected:     e.IsInspected,
			IsReported:      e.IsReported,
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return shipDock, nil
}

func (r *ship) NeedCheckupShip(ctx context.Context, request dto.NeedCheckupShipParam) ([]dto.NeedCheckupShipResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.ShipDockedLog{}).
		Select("ship_docked_logs.*, ships.name as ship_name, ships.id as ship_id").
		Joins("JOIN ships ON ship_docked_logs.ship_id = ships.id")

	if request.Search != "" {
		searchLower := strings.ToLower(request.Search)
		query = query.Where("lower(ships.name) LIKE ?", "%"+searchLower+"%")
	}

	query = query.Where("(is_inspected = ? OR is_reported = ?) and ship_docked_logs.status = ?", 0, 0, "checkin").
		Limit(request.Limit).
		Offset(request.Offset).
		Order("ship_docked_logs.created_at DESC")

	var result []struct {
		model.ShipDockedLog
		ShipName string `json:"ship_name"`
	}

	if err := query.Find(&result).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var shipDock []dto.NeedCheckupShipResponse
	for _, e := range result {
		shipDock = append(shipDock, dto.NeedCheckupShipResponse{
			LogID:       e.ID,
			ShipName:    e.ShipName,
			Lat:         e.Lat,
			Long:        e.Long,
			CheckinDate: e.CreatedAt.Format("2006-01-02 15:04:05"),
			IsInspected: e.IsInspected,
			IsReported:  e.IsReported,
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return shipDock, nil
}

func (r *ship) UpdateShipCheckup(ctx context.Context, request dto.ShipCheckupRequest, id int, data model.ShipDockedLog) error {
	tx := r.Db.WithContext(ctx).Begin()

	isInspectedInt := data.IsInspected
	isReportedInt := data.IsReported

	if request.IsInspected {
		isInspectedInt = 1
	}

	if request.IsReported {
		isReportedInt = 1
	}

	updateFields := map[string]interface{}{
		"is_inspected": isInspectedInt,
		"is_reported":  isReportedInt,
	}

	if err := tx.Debug().Model(&model.ShipDockedLog{}).Where("id = ?", id).Updates(updateFields).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	cacheKey := []string{"ship_list-*", "ship_last_update"}

	for i := range cacheKey {
		if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey[i]); err != nil {
			return nil
		}
	}

	return nil
}

func (r *ship) FindOneDockedLog(ctx context.Context, selectedFields string, query string, args ...any) (model.ShipDockedLog, error) {
	var res model.ShipDockedLog

	db := r.Db.WithContext(ctx).Model(model.ShipDockedLog{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return model.ShipDockedLog{}, err
	}

	return res, nil
}

func (r *ship) UpdateShipDeviceID(ctx context.Context, deviceID string, user_id int) error {
	tx := r.Db.WithContext(ctx).Begin()

	updateFields := map[string]interface{}{
		"device_id": deviceID,
	}

	if err := tx.Model(&model.Ship{}).Where("user_id = ?", user_id).Updates(updateFields).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	cacheKey := []string{"ship_list-*", "ship_last_update"}

	for i := range cacheKey {
		if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey[i]); err != nil {
			return nil
		}
	}

	return nil
}

func (r *ship) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.Ship, error) {
	var res model.Ship

	db := r.Db.WithContext(ctx).Model(model.Ship{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return model.Ship{}, err
	}

	return res, nil
}

func (r *ship) CountShipFraud(ctx context.Context, startDate string, endDate string) (int64, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.ShipLocationLog{})

	var res int64

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Where("is_mocked = ?", 1).Count(&res)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return res, nil
}

func (r *ship) CountShipByStatus(ctx context.Context, startDate string, endDate string, status string) (int64, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.Ship{})

	var res int64

	if startDate != "" && endDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", startDate, endDate)
	}

	query.Where("status = ?", status).Count(&res)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return res, nil
}

func (r *ship) CountShipByTerrain(ctx context.Context, onGround int) (int64, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.Ship{})

	var res int64

	query.Where("on_ground = ?", onGround).Count(&res)

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	return res, nil
}

func (r *ship) StoreNewShip(ctx context.Context, request dto.PairingToNewShip) error {
	shipModel := model.Ship{
		Name:            request.ShipName,
		Phone:           request.Phone,
		ResponsibleName: request.ResponsibleName,
		DeviceID:        request.DeviceID,
		FirebaseToken:   request.FirebaseToken,
		Status:          "out of scope",
		UserID:          request.UserID,
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

	cacheKey := []string{"ship_list-*", "ship_count"}

	for i := range cacheKey {
		if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey[i]); err != nil {
			return nil
		}
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

func (r *ship) ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipMobileDetailResponse, error) {
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

func (r *ship) ShipByAuth(ctx context.Context, authUser model.User) (*dto.ShipMobileDetailResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var ship model.Ship
	err := tx.Where("user_id = ?", authUser.ID).First(&ship).Error
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
		ShipID:      request.ShipID,
		Long:        request.Long,
		Lat:         request.Lat,
		Status:      model.ShipStatus(request.Status),
		IsInspected: 0,
		IsReported:  0,
	}

	if err := tx.Create(&dockedModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	cacheKey := "ship_statistic_count"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
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

	cacheKey := []string{"ship_list-*", "ship_last_update"}

	for i := range cacheKey {
		if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey[i]); err != nil {
			return nil
		}
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
		GT:        request.GT,
		OwnerName: request.OwnerName,
	}

	if err := tx.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "ship_id"}},
		DoUpdates: clause.AssignmentColumns([]string{"type", "dimension", "harbour", "siup", "bkp", "selar_mark", "gt", "owner_name"}),
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

func (r *ship) ShipDockedLogs(ctx context.Context, ShipID int, request *dto.ShipLogParam) ([]dto.DockLogsShip, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var logs []model.ShipDockedLog
	query := tx.Where("ship_id = ?", ShipID).Order("created_at DESC")

	if request.StartDate != "" && request.EndDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", request.StartDate, request.EndDate)
	}

	limit := 10
	if request.Limit != 0 {
		limit = request.Limit
	}

	query.Limit(limit).Offset(request.Offset)

	err := query.Find(&logs).Error
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

func (r *ship) ShipLocationLogs(ctx context.Context, ShipID int, request *dto.ShipLogParam) ([]dto.LocationLogsShip, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var logs []model.ShipLocationLog
	query := tx.Where("ship_id = ?", ShipID).Order("created_at DESC")

	if request.StartDate != "" && request.EndDate != "" {
		query = query.Where("created_at BETWEEN ? AND ?", request.StartDate, request.EndDate)
	}

	limit := 10
	if request.Limit != 0 {
		limit = request.Limit
	}

	query.Limit(limit).Offset(request.Offset)

	err := query.Find(&logs).Error
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

func (r *ship) ShipAddonDetail(ctx context.Context, ShipID int) (dto.ShipAddonDetailResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var detail model.ShipDetail
	if err := tx.Where("ship_id = ?", ShipID).First(&detail).Error; err != nil {
		tx.Rollback()
		return dto.ShipAddonDetailResponse{}, err
	}

	res := dto.ShipAddonDetailResponse{
		Type:      string(detail.Type),
		Dimension: detail.Dimension,
		Harbour:   detail.Harbour,
		SIUP:      detail.SIUP,
		BKP:       detail.BKP,
		SelarMark: detail.SelarMark,
		GT:        detail.GT,
		OwnerName: detail.OwnerName,
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return dto.ShipAddonDetailResponse{}, err
	}

	return res, nil
}

func (r *ship) CountShip(ctx context.Context) (int64, error) {
	cacheKey := "ship_count"

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo int64
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.Ship{})

	var totalShip int64

	if err := query.Count(&totalShip).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(totalShip)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return totalShip, nil
}

func (r *ship) CountStatistic(ctx context.Context) ([]int64, error) {
	cacheKey := "ship_statistic_count"

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo []int64
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	var totalCheckin int64
	if err := tx.Model(&model.ShipDockedLog{}).Where("status = ?", "checkin").Count(&totalCheckin).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var totalCheckout int64
	if err := tx.Model(&model.ShipDockedLog{}).Where("status = ?", "checkout").Count(&totalCheckout).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var totalFraud int64
	if err := tx.Model(&model.ShipLocationLog{}).Where("is_mocked = ?", 1).Count(&totalFraud).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Cache the result in Redis
	if r.CacheEnabled {
		jsonData, err := json.Marshal([]int64{totalCheckin, totalCheckout, totalFraud})
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return []int64{totalCheckin, totalCheckout, totalFraud}, nil
}

func (r *ship) LastUpdated(ctx context.Context) (time.Time, error) {
	cacheKey := "ship_last_update"

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo time.Time
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	var maxUpdatedAt time.Time

	query := tx.Model(&model.Ship{}).Select("MAX(updated_at)").Row()
	if err := query.Scan(&maxUpdatedAt); err != nil {
		tx.Rollback()
		return time.Time{}, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return time.Time{}, err
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(maxUpdatedAt)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return maxUpdatedAt, nil
}

func (r *ship) ShipInBatch(ctx context.Context, start int, end int) (*[]model.Ship, bool, error) {
	lastUpdate, err := r.LastUpdated(ctx)
	if err != nil {
		return nil, false, err
	}

	var result []struct {
		model.Ship
		LogLong       string    `json:"log_long"`
		LogLat        string    `json:"log_lat"`
		CreatedAtLogs time.Time `json:"log_created" gorm:"column:log_created"`
		OnGroundLogs  int       `json:"log_onground" gorm:"column:log_onground"`
	}

	cacheKey := "ship_highest_current_update"

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo time.Time
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				if lastUpdate != cachedInfo {
					tx := r.Db.WithContext(ctx).Begin()

					query := tx.Model(&model.Ship{})

					query = query.Select("ships.*, slr.long as log_long, slr.lat as log_lat, slr.created_at as log_created, slr.on_ground as log_onground")
					query = query.Joins(`
						JOIN ship_location_logs AS slr
						ON slr.ship_id = ships.id
						AND slr.on_ground = 0
						AND slr.created_at = (
							SELECT MAX(created_at)
							FROM ship_location_logs
							WHERE ship_id = ships.id AND on_ground = 0
						)
					`)
					query = query.Where("ships.updated_at BETWEEN ? AND ?", lastUpdate.Add(-10*time.Second), lastUpdate.Add(10*time.Second)).Order("ships.created_at DESC")

					if err := query.Find(&result).Error; err != nil {
						tx.Rollback()
						return nil, false, err
					}

					var updatedShips []model.Ship

					for _, e := range result {
						if e.LogLat != "" && e.LogLong != "" && (e.Ship.CurrentLat != e.LogLat || e.Ship.CurrentLong != e.LogLong) && e.Ship.OnGround == 1 {
							if e.CreatedAtLogs.Before(e.Ship.UpdatedAt) {
								e.Ship.CurrentLat = e.LogLat
								e.Ship.CurrentLong = e.LogLong
								e.Ship.OnGround = e.OnGroundLogs
							}
						}

						updatedShips = append(updatedShips, e.Ship)
					}

					if err := tx.Commit().Error; err != nil {
						tx.Rollback()
						return nil, false, err
					}

					jsonData, err := json.Marshal(lastUpdate)
					if err == nil {
						r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
					} else {
						fmt.Println("Error marshalling data for cache:", err)
					}

					return &updatedShips, true, nil
				}
			}
		} else if err == redis.Nil {
			jsonData, err := json.Marshal(lastUpdate)
			if err == nil {
				r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
			} else {
				fmt.Println("Error marshalling data for cache:", err)
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.Ship{})
	query = query.Select("ships.*, slr.long as log_long, slr.lat as log_lat, slr.created_at as log_created, slr.on_ground as log_onground")
	query = query.Joins(`
		JOIN ship_location_logs AS slr
		ON slr.ship_id = ships.id
		AND slr.on_ground = 0
		AND slr.created_at = (
			SELECT MAX(created_at)
			FROM ship_location_logs
			WHERE ship_id = ships.id AND on_ground = 0
		)
	`)

	query = query.Offset(start).Limit(end - start).Order("ships.created_at ASC")

	if err := query.Find(&result).Error; err != nil {
		tx.Rollback()
		return nil, false, err
	}

	var ships []model.Ship

	for _, e := range result {
		if e.LogLat != "" && e.LogLong != "" && (e.Ship.CurrentLat != e.LogLat || e.Ship.CurrentLong != e.LogLong) && e.Ship.OnGround == 1 {
			if e.CreatedAtLogs.Before(e.Ship.UpdatedAt) {
				e.Ship.CurrentLat = e.LogLat
				e.Ship.CurrentLong = e.LogLong
				e.Ship.OnGround = e.OnGroundLogs
			}
		}

		ships = append(ships, e.Ship)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, false, err
	}

	return &ships, false, nil
}

func (r *ship) ReportShipDocking(ctx context.Context, request dto.ReportShipDockedParam) ([]dto.ReportShipDockingResponse, error) {
	paramJSON, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	hash := sha1.Sum(paramJSON)
	uniqueString := fmt.Sprintf("%x", hash)

	cacheKey := "report_ship_list-" + uniqueString

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo []dto.ReportShipDockingResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.ShipDockedLog{}).
		Select("ship_docked_logs.*, ships.name as ship_name, ships.id as ship_id").
		Joins("JOIN ships ON ship_docked_logs.ship_id = ships.id")

	if request.LogType != nil && request.LogType[0] != "" && len(request.LogType) > 0 {
		query = query.Where("ship_docked_logs.status IN (?)", request.LogType)
	}

	if request.Search != "" {
		searchLower := strings.ToLower(request.Search)
		query = query.Where("lower(ships.name) LIKE ?", "%"+searchLower+"%")
	}

	if request.StartDate != "" && request.EndDate != "" {
		query = query.Where("DATE(ship_docked_logs.created_at) BETWEEN ? AND ?", request.StartDate, request.EndDate)
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("ship_docked_logs.created_at DESC")

	var result []struct {
		model.ShipDockedLog
		ShipName string `json:"ship_name"`
	}

	if err := query.Find(&result).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var shipDock []dto.ReportShipDockingResponse
	for _, e := range result {
		shipDock = append(shipDock, dto.ReportShipDockingResponse{
			LogID:    e.ID,
			ShipID:   e.ShipID,
			ShipName: e.ShipName,
			Lat:      e.Lat,
			Long:     e.Long,
			Status:   string(e.Status),
			LogDate:  e.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(shipDock)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return shipDock, nil
}

func (r *ship) ReportShipFraud(ctx context.Context, request dto.ReportShipLocationParam) ([]dto.ReportShipLocationResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.ShipLocationLog{}).
		Select("ship_location_logs.*, ships.name as ship_name, ships.id as ship_id").
		Joins("JOIN ships ON ship_location_logs.ship_id = ships.id").
		Where("is_mocked = 1")

	if request.Search != "" {
		searchLower := strings.ToLower(request.Search)
		query = query.Where("lower(ships.name) LIKE ?", "%"+searchLower+"%")
	}

	if request.StartDate != "" && request.EndDate != "" {
		query = query.Where("DATE(ship_location_logs.created_at) BETWEEN ? AND ?", request.StartDate, request.EndDate)
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("ship_location_logs.created_at DESC")

	var result []struct {
		model.ShipLocationLog
		ShipName string `json:"ship_name"`
	}

	if err := query.Find(&result).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var shipDock []dto.ReportShipLocationResponse
	for _, e := range result {
		shipDock = append(shipDock, dto.ReportShipLocationResponse{
			LogID:    e.ID,
			ShipID:   e.ShipID,
			ShipName: e.ShipName,
			Lat:      e.Lat,
			Long:     e.Long,
			IsMocked: e.IsMocked,
			OnGround: e.OnGround,
			LogDate:  e.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return shipDock, nil
}
