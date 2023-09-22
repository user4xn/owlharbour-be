package repository

import (
	"context"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"

	"gorm.io/gorm"
)

type Ship interface {
	StoreNewShip(ctx context.Context, request dto.PairingRequestResponse) error
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
		Name:          request.ShipName,
		Phone:         request.Phone,
		DeviceID:      request.DeviceID,
		FirebaseToken: request.FirebaseToken,
		Status:        "out of scope",
	}

	return r.Db.Create(&shipModel).Error
}
