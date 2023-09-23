package repository

import (
	"context"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"

	"gorm.io/gorm"
)

type PairingRequest interface {
	StorePairingRequests(ctx context.Context, request dto.PairingRequest) error
	PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error)
	UpdatedPairingStatus(ctx context.Context, request dto.PairingActionRequest) (*dto.PairingRequestResponse, error)
}

type pairingRequest struct {
	Db *gorm.DB
}

func NewPairingRequestRepository(db *gorm.DB) PairingRequest {
	return &pairingRequest{
		Db: db,
	}
}

func (r *pairingRequest) StorePairingRequests(ctx context.Context, request dto.PairingRequest) error {
	tx := r.Db.Begin()

	existingDevice := model.PairingRequest{}
	if err := tx.Where("device_id = ?", request.DeviceID).Order("created_at DESC").First(&existingDevice).Error; err == nil {
		if existingDevice.Status == "pending" {
			tx.Rollback()
			return fmt.Errorf("a pending pairing request with the same DeviceID already exists")
		} else if existingDevice.Status == "approved" {
			tx.Rollback()
			return fmt.Errorf("this device already registered at " + existingDevice.CreatedAt.Format("2006-01-02 15:04:05"))
		}
	}

	pairingModel := model.PairingRequest{
		Name:            request.ShipName,
		Phone:           request.Phone,
		ResponsibleName: request.ResponsibleName,
		DeviceID:        request.DeviceID,
		FirebaseToken:   request.FirebaseToken,
		Status:          "pending",
	}

	if err := tx.Create(&pairingModel).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	return nil
}

func (r *pairingRequest) PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error) {
	tx := r.Db.Begin()

	query := tx.Model(&model.PairingRequest{})

	if request.Status != nil && request.Status[0] != "" && len(request.Status) > 0 {
		query = query.Where("status IN (?)", request.Status)
	}

	if request.Search != "" {
		query = query.Where("name LIKE ? OR device_id LIKE ? OR phone LIKE ?", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%")
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("created_at DESC")

	var pairingRequest []model.PairingRequest

	if err := query.Find(&pairingRequest).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var pairingList []dto.PairingRequestResponse
	for _, e := range pairingRequest {
		pairingList = append(pairingList, dto.PairingRequestResponse{
			ID:              e.ID,
			ShipName:        e.Name,
			Phone:           e.Phone,
			ResponsibleName: e.ResponsibleName,
			DeviceID:        e.DeviceID,
			FirebaseToken:   e.FirebaseToken,
			Status:          string(e.Status),
			CreatedAt:       e.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return pairingList, nil
}

func (r *pairingRequest) UpdatedPairingStatus(ctx context.Context, request dto.PairingActionRequest) (*dto.PairingRequestResponse, error) {
	tx := r.Db.Begin()

	var pairing model.PairingRequest
	err := tx.First(&pairing, request.PairingID).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if pairing.Status == "approved" || pairing.Status == "rejected" {
		tx.Rollback()
		return nil, fmt.Errorf("this pairing request already responded")
	}

	pairingModel := model.PairingRequest{
		Status: model.PairingStatus(request.Status),
	}

	err = tx.Model(&model.PairingRequest{}).Where("id = ?", request.PairingID).Updates(&pairingModel).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	pairingData := dto.PairingRequestResponse{
		ID:              pairing.ID,
		ShipName:        pairing.Name,
		Phone:           pairing.Phone,
		ResponsibleName: pairing.ResponsibleName,
		DeviceID:        pairing.DeviceID,
		FirebaseToken:   pairing.FirebaseToken,
		Status:          request.Status,
		CreatedAt:       pairing.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &pairingData, nil
}
