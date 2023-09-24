package repository

import (
	"context"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"

	"gorm.io/gorm"
)

type Ship interface {
	StoreNewShip(ctx context.Context, request dto.PairingRequestResponse) error
	ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error)
	ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipDetailResponse, error)
	GetLastDockedLog(ctx context.Context, ShipID int) (*dto.ShipDockedLog, error)
	StoreDockedLog(ctx context.Context, request dto.ShipDockedLogStore) error
	StoreLocationLog(ctx context.Context, request dto.ShipLocationLogStore) error
	UpdateShip(ctx context.Context, request model.Ship) error
}

type ship struct {
	Db *gorm.DB
}

func NewShipRepository(db *gorm.DB) Ship {
	return &ship{
		Db: db,
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

	return nil
}

func (r *ship) ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.Ship{})

	if request.Status != nil && request.Status[0] != "" && len(request.Status) > 0 {
		query = query.Where("status IN (?)", request.Status)
	}

	if request.Search != "" {
		query = query.Where("ship_name LIKE ? OR device_id LIKE ? OR phone LIKE ?", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%")
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

	return shipList, nil
}

func (r *ship) ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipDetailResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var ship model.Ship
	err := tx.Where("device_id = ?", DeviceID).First(&ship).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	shipDetail := dto.ShipDetailResponse{
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

	return nil
}
