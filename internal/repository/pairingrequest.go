package repository

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/model"
	"owlharbour-api/pkg/constants"
	"owlharbour-api/pkg/helper"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type PairingRequest interface {
	StorePairingRequests(ctx context.Context, request dto.PairingRequest) error
	PairingRequestList(ctx context.Context, request dto.PairingListParam) ([]dto.PairingRequestResponse, error)
	UpdatedPairingStatus(ctx context.Context, id int, status string) (*dto.PairingRequestResponse, error)
	PairingDetailByUsername(ctx context.Context, username string) (*dto.DetailPairingResponse, error)
	PairingRequestCount(ctx context.Context, status string) (int64, error)
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

func (r *pairingRequest) PairingRequestCount(ctx context.Context, status string) (int64, error) {
	paramJSON, err := json.Marshal(status)
	if err != nil {
		fmt.Println("Error:", err)
		return 0, err
	}

	hash := sha1.Sum(paramJSON)
	uniqueString := fmt.Sprintf("%x", hash)

	cacheKey := "pairing_pending_count-" + uniqueString

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

	query := tx.Model(&model.PairingRequest{})

	query = query.Where("status = ?", status)

	var res int64

	if err := query.Count(&res).Error; err != nil {
		tx.Rollback()
		return 0, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return 0, err
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

func (r *pairingRequest) StorePairingRequests(ctx context.Context, request dto.PairingRequest) error {
	tx := r.Db.WithContext(ctx).Begin()

	existingDevice := model.PairingRequest{}
	if err := tx.Where("device_id = ?", request.DeviceID).Order("created_at DESC").First(&existingDevice).Error; err == nil {
		if existingDevice.Status == "pending" {
			tx.Rollback()
			return fmt.Errorf("a pending pairing request with the same DeviceID already exists")
		}
	}

	existingDeviceShip := model.Ship{}
	if err := tx.Where("device_id = ?", request.DeviceID).Order("created_at DESC").First(&existingDeviceShip).Error; err == nil {
		if existingDevice.ID != 0 {
			tx.Rollback()
			return fmt.Errorf("a ship with the same DeviceID already exists")
		}
	}

	password := []byte(request.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		return constants.ErrorHashPassword
	}

	pairingModel := model.PairingRequest{
		Name:            request.ShipName,
		Phone:           request.Phone,
		Username:        request.Username,
		Password:        string(hashedPassword),
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

	cacheKey := []string{"pairing_list-*", "pairing_pending_count-*"}

	for _, ck := range cacheKey {
		if err := helper.DeleteRedisKeysByPattern(r.RedisClient, ck); err != nil {
			return err
		}
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
			Username:        e.Username,
			Password:        "********",
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

func (r *pairingRequest) UpdatedPairingStatus(ctx context.Context, id int, status string) (*dto.PairingRequestResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var pairing model.PairingRequest
	err := tx.First(&pairing, id).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if pairing.Status == "approved" || pairing.Status == "rejected" {
		tx.Rollback()
		return nil, fmt.Errorf("this pairing request already responded")
	}

	pairingModel := model.PairingRequest{
		Status: model.PairingStatus(status),
	}

	err = tx.Model(&model.PairingRequest{}).Where("id = ?", id).Updates(&pairingModel).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	pairingData := dto.PairingRequestResponse{
		ID:              pairing.ID,
		ShipName:        pairing.Name,
		Phone:           pairing.Phone,
		Username:        pairing.Username,
		Password:        pairing.Password,
		ResponsibleName: pairing.ResponsibleName,
		DeviceID:        pairing.DeviceID,
		FirebaseToken:   pairing.FirebaseToken,
		Status:          status,
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

func (r *pairingRequest) PairingDetailByUsername(ctx context.Context, username string) (*dto.DetailPairingResponse, error) {
	tx := r.Db.WithContext(ctx).Begin()

	var pairing model.PairingRequest
	err := tx.Where("username = ?", username).Order("created_at DESC").First(&pairing).Error
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	pairingDetail := dto.DetailPairingResponse{
		ShipName:       pairing.Name,
		ReponsibleName: pairing.ResponsibleName,
		Phone:          pairing.Phone,
		Username:       pairing.Username,
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
	if err := tx.Where("username = ? AND status != 'pending'", username).Find(&historyPairing).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var resHitory []dto.HistoryPairing
	for _, history := range historyPairing {
		resHitory = append(resHitory, dto.HistoryPairing{
			ShipName:       history.Name,
			ReponsibleName: history.ResponsibleName,
			Phone:          history.Phone,
			Username:       history.Username,
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
