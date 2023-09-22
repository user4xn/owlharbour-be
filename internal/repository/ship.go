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

	return r.Db.Create(&shipModel).Error
}

func (r *ship) ShipList(ctx context.Context, request dto.ShipListParam) ([]dto.ShipResponse, error) {
	query := r.Db.Model(&model.Ship{})

	if request.Status != nil && request.Status[0] != "" && len(request.Status) > 0 {
		query = query.Where("status IN (?)", request.Status)
	}

	if request.Search != "" {
		query = query.Where("ship_name LIKE ? OR device_id LIKE ? OR phone LIKE ?", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%")
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("created_at DESC")

	var ship []model.Ship

	if err := query.Find(&ship).Error; err != nil {
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
			CreatedAt:       e.CreatedAt.Format("2006-01-02"),
		})
	}

	return shipList, nil
}

func (r *ship) ShipByDevice(ctx context.Context, DeviceID string) (*dto.ShipDetailResponse, error) {
	var ship model.Ship
	err := r.Db.Where("device_id = ?", DeviceID).First(&ship).Error
	if err != nil {
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
		CreatedAt:       ship.CreatedAt.Format("2006-01-02"),
	}

	return &shipDetail, nil
}
