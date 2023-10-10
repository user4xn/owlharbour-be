package repository

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"
	"simpel-api/pkg/helper"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type PairingRequest interface {
	StorePairingRequests(ctx context.Context, request dto.PairingRequest) error
	PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error)
	UpdatedPairingStatus(ctx context.Context, request dto.PairingActionRequest) (*dto.PairingRequestResponse, error)
	PairingDetailByDevice(ctx context.Context, DeviceID string) (*dto.DetailPairingResponse, error)
}

type pairingRequest struct {
	Db           *gorm.DB
	RedisClient  *redis.Client
	CacheEnabled bool
}

func NewPairingRequestRepository(db *gorm.DB, redisClient *redis.Client) PairingRequest {
	return &pairingRequest{
		Db:           db,
		RedisClient:  redisClient,
		CacheEnabled: true,
	}
}

func (r *pairingRequest) StorePairingRequests(ctx context.Context, request dto.PairingRequest) error {
	tx := r.Db.WithContext(ctx).Begin()

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

	cacheKey := "pairing_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return err
	}

	return nil
}

func (r *pairingRequest) PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error) {
	paramJSON, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	hash := sha1.Sum(paramJSON)
	uniqueString := fmt.Sprintf("%x", hash)

	cacheKey := "pairing_list-" + uniqueString

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo []dto.PairingRequestResponse
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

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

	if r.CacheEnabled {
		jsonData, err := json.Marshal(pairingList)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return pairingList, nil
}

func (r *pairingRequest) UpdatedPairingStatus(ctx context.Context, request dto.PairingActionRequest) (*dto.PairingRequestResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

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

	cacheKey := "pairing_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil, err
	}

	return &pairingData, nil
}

func (r *pairingRequest) PairingDetailByDevice(ctx context.Context, DeviceID string) (*dto.DetailPairingResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var pairing model.PairingRequest
	err := tx.Where("device_id = ?", DeviceID).First(&pairing).Order("created_at DESC").Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	pairingDetail := dto.DetailPairingResponse{
		ShipName:       pairing.Name,
		ReponsibleName: pairing.ResponsibleName,
		Phone:          pairing.Phone,
		DeviceID:       pairing.DeviceID,
		Status:         string(pairing.Status),
		SubmittedAt:    pairing.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if pairing.Status != model.Pending {
		pairingDetail.RespondedAt = pairing.UpdatedAt.Format("2006-01-02 15:04:05")
	} else {
		pairingDetail.RespondedAt = ""
	}

	var historyPairing []model.PairingRequest
	if err := tx.Where("device_id = ? AND status != 'pending'", DeviceID).Find(&historyPairing).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var resHitory []dto.HistoryPairing
	for _, history := range historyPairing {
		resHitory = append(resHitory, dto.HistoryPairing{
			ShipName:       history.Name,
			ReponsibleName: history.ResponsibleName,
			Phone:          history.Phone,
			Status:         string(history.Status),
			SubmittedAt:    history.CreatedAt.Format("2006-01-02 15:04:05"),
			RespondedAt:    history.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	pairingDetail.HistoryPairing = resHitory

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	return &pairingDetail, nil
}
