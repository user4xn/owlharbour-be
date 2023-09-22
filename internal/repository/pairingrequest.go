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
	pairingModel := model.PairingRequest{
		Name:          request.ShipName,
		Phone:         request.Phone,
		DeviceID:      request.DeviceID,
		FirebaseToken: request.FirebaseToken,
		Status:        "pending",
	}

	return r.Db.Create(&pairingModel).Error
}

func (r *pairingRequest) PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error) {
	query := r.Db.Model(&model.PairingRequest{})

	if request.Status != nil && request.Status[0] != "" && len(request.Status) > 0 {
		query = query.Where("status IN (?)", request.Status)
	}

	if request.Search != "" {
		query = query.Where("name LIKE ? OR device_id LIKE ? OR phone LIKE ?", "%"+request.Search+"%", "%"+request.Search+"%", "%"+request.Search+"%")
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("created_at DESC")

	var pairingRequest []model.PairingRequest

	if err := query.Find(&pairingRequest).Error; err != nil {
		return nil, err
	}

	var pairingList []dto.PairingRequestResponse
	for _, e := range pairingRequest {
		pairingList = append(pairingList, dto.PairingRequestResponse{
			ID:            e.ID,
			ShipName:      e.Name,
			Phone:         e.Phone,
			DeviceID:      e.DeviceID,
			FirebaseToken: e.FirebaseToken,
			Status:        string(e.Status),
			CreatedAt:     e.CreatedAt.Format("2006-01-02"),
		})
	}

	return pairingList, nil
}

func (r *pairingRequest) UpdatedPairingStatus(ctx context.Context, request dto.PairingActionRequest) (*dto.PairingRequestResponse, error) {
	var pairing model.PairingRequest
	err := r.Db.First(&pairing, request.PairingID).Error
	if err != nil {
		return nil, err
	}

	if pairing.Status == "approved" || pairing.Status == "rejected" {
		return nil, fmt.Errorf("this pairing request already responded")
	}

	pairingModel := model.PairingRequest{
		Status: model.PairingStatus(request.Status),
	}

	err = r.Db.Model(&model.PairingRequest{}).Where("id = ?", request.PairingID).Updates(&pairingModel).Error
	if err != nil {
		return nil, err
	}

	pairingData := dto.PairingRequestResponse{
		ID:            pairing.ID,
		ShipName:      pairing.Name,
		Phone:         pairing.Phone,
		DeviceID:      pairing.DeviceID,
		FirebaseToken: pairing.FirebaseToken,
		Status:        request.Status,
		CreatedAt:     pairing.CreatedAt.Format("2006-01-02"),
	}

	return &pairingData, nil
}
