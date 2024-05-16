package repository

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/model"
	"owlharbour-api/pkg/helper"
	"owlharbour-api/pkg/util"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type User interface {
	GetAll(ctx context.Context, request dto.UserListParam) ([]dto.AllUser, error)
	Find(ctx context.Context, queries []string, argsSlice ...[]interface{}) (model.User, error)
	Store(ctx context.Context, data model.User) error
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.User, error)
	UpdateOne(ctx context.Context, updatedModels *dto.PayloadUpdateUser, updatedField string, query string, args ...interface{}) error
	UpdateJwtToken(ctx context.Context, updatedModels *dto.PayloadUpdateJwtToken, updatedField string, query string, args ...interface{}) error
	DeleteOne(ctx context.Context, query string, args ...interface{}) error
	RemoveJwtToken(ctx context.Context, updatedModels *dto.PayloadUpdateJwtToken, updatedField string, query string, args ...interface{}) error
}

type user struct {
	Db           *gorm.DB
	RedisClient  *redis.Client
	CacheEnabled bool
}

func NewUserRepository(db *gorm.DB, redisClient *redis.Client) User {
	return &user{
		Db:           db,
		RedisClient:  redisClient,
		CacheEnabled: true,
	}
}

func (r *user) GetAll(ctx context.Context, request dto.UserListParam) ([]dto.AllUser, error) {
	paramJSON, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error:", err)
		return nil, err
	}

	hash := sha1.Sum(paramJSON)
	uniqueString := fmt.Sprintf("%x", hash)

	cacheKey := "user_list-" + uniqueString

	if r.CacheEnabled {
		cachedData, err := r.RedisClient.Get(ctx, cacheKey).Result()
		if err == nil {
			var cachedInfo []dto.AllUser
			if err := json.Unmarshal([]byte(cachedData), &cachedInfo); err == nil {
				return cachedInfo, nil
			}
		}
	}

	tx := r.Db.WithContext(ctx).Begin()

	query := tx.Model(&model.User{})

	if request.Search != "" {
		searchLower := strings.ToLower(request.Search)
		query = query.Where("lower(name) LIKE ? OR lower(email) LIKE ?", "%"+searchLower+"%", "%"+searchLower+"%")
	}

	query = query.Limit(request.Limit).Offset(request.Offset).Order("created_at DESC")

	var res []model.User
	if err := query.Find(&res).Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	var AllUser []dto.AllUser // Initialize the slice

	for _, user := range res {
		tCreatedAt := user.CreatedAt
		tUpdatedAt := user.UpdatedAt
		tEmailVerifiedAt := user.EmailVerifiedAt

		formatCreatedAt := tCreatedAt.Format("2006-01-02 15:04:05")
		formatUpdatedAt := tUpdatedAt.Format("2006-01-02 15:04:05")
		formatEmailVerifiedAt := ""

		if tEmailVerifiedAt != nil { // Check if EmailVerifiedAt is not zero time
			formatEmailVerifiedAt = tEmailVerifiedAt.Format("2006-01-02 15:04:05")
		}

		userDTO := dto.AllUser{
			ID:              user.ID,
			Name:            user.Name,
			Username:        user.Username,
			Email:           user.Email,
			Role:            string(user.Role),
			EmailVerifiedAt: formatEmailVerifiedAt,
			CreatedAt:       formatCreatedAt,
			UpdatedAt:       formatUpdatedAt,
		}

		AllUser = append(AllUser, userDTO)
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	if r.CacheEnabled {
		jsonData, err := json.Marshal(AllUser)
		if err == nil {
			r.RedisClient.Set(ctx, cacheKey, jsonData, time.Hour)
		} else {
			fmt.Println("Error marshalling data for cache:", err)
		}
	}

	return AllUser, nil
}

func (r *user) Store(ctx context.Context, data model.User) error {
	tx := r.Db.WithContext(ctx)
	if err := tx.Model(model.User{}).Create(&data).Error; err != nil {
		return err
	}

	cacheKey := "user_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *user) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.User, error) {
	var res model.User

	db := r.Db.WithContext(ctx).Model(model.User{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return model.User{}, err
	}

	return res, nil
}

func (r *user) Find(ctx context.Context, queries []string, argsSlice ...[]interface{}) (model.User, error) {
	var res model.User

	db := r.Db.WithContext(ctx).Model(model.User{})

	for idx, query := range queries {
		if idx < len(argsSlice) {
			db = db.Or(query, argsSlice[idx]...)
		}
	}

	if err := db.Find(&res).Error; err != nil {
		return model.User{}, err
	}

	return res, nil
}

func (r *user) UpdateOne(ctx context.Context, updatedModels *dto.PayloadUpdateUser, updatedField string, query string, args ...interface{}) error {
	fmt.Println(updatedModels)
	modelDb := r.Db.WithContext(ctx).Model(&model.User{})
	if err := util.SetSelectFields(modelDb, updatedField).Where(query, args...).Updates(updatedModels).Error; err != nil {
		return err
	}

	cacheKey := "user_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *user) UpdateJwtToken(ctx context.Context, updatedModels *dto.PayloadUpdateJwtToken, updatedField string, query string, args ...interface{}) error {
	fmt.Println(updatedModels)
	modelDb := r.Db.WithContext(ctx).Model(&model.User{})
	if err := util.SetSelectFields(modelDb, updatedField).Where(query, args...).Updates(updatedModels).Error; err != nil {
		return err
	}

	cacheKey := "user_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *user) RemoveJwtToken(ctx context.Context, updatedModels *dto.PayloadUpdateJwtToken, updatedField string, query string, args ...interface{}) error {
	fmt.Println(updatedModels)
	modelDb := r.Db.WithContext(ctx).Model(&model.User{})
	if err := util.SetSelectFields(modelDb, updatedField).Where(query, args...).Updates(updatedModels).Error; err != nil {
		return err
	}

	cacheKey := "user_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}

func (r *user) DeleteOne(ctx context.Context, query string, args ...interface{}) error {
	db := r.Db.WithContext(ctx).Model(&model.User{})

	if err := db.Where(query, args...).Delete(&model.User{}).Error; err != nil {
		return err
	}

	cacheKey := "user_list-*"

	if err := helper.DeleteRedisKeysByPattern(r.RedisClient, cacheKey); err != nil {
		return nil
	}

	return nil
}
