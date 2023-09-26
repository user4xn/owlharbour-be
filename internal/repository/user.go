package repository

import (
	"context"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/model"
	"simpel-api/pkg/util"

	"gorm.io/gorm"
)

type UserInterface interface {
	Init() *gorm.DB
	GetAll(ctx context.Context, selectedFields string, searchQuery string, limit, offset int, args ...interface{}) ([]dto.AllUser, error)
	Find(ctx context.Context, queries []string, argsSlice ...[]interface{}) (model.User, error)
	Store(ctx context.Context, data model.User) error
	FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.User, error)
}

type User struct {
	Database *gorm.DB
}

func NewUserRepository(db *gorm.DB) *User {
	return &User{
		Database: db,
	}
}

func (u *User) Init() *gorm.DB {
	return u.Database
}
func (u *User) GetAll(ctx context.Context, selectedFields string, searchQuery string, limit, offset int, args ...interface{}) ([]dto.AllUser, error) {
	var res []model.User

	db := u.Database.WithContext(ctx).Model(&model.User{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(searchQuery, args...).Limit(limit).Offset(offset).Find(&res).Error; err != nil {
		return nil, err
	}

	var AllUser []dto.AllUser
	for _, user := range res {
		tCreatedAt := user.CreatedAt
		tUpdatedAt := user.UpdatedAt
		formatCreatedAt := tCreatedAt.Format("2006-01-02 15:04:05")
		formatUpdatedAt := tUpdatedAt.Format("2006-01-02 15:04:05")
		userDTO := dto.AllUser{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      fmt.Sprintf("%s", user.Role),
			CreatedAt: formatCreatedAt,
			UpdatedAt: formatUpdatedAt,
		}
		AllUser = append(AllUser, userDTO)
	}

	return AllUser, nil
}

func (u *User) Store(ctx context.Context, data model.User) error {
	tx := u.Database.WithContext(ctx)
	if err := tx.Model(model.User{}).Create(&data).Error; err != nil {
		return err
	}
	return nil
}

func (u *User) FindOne(ctx context.Context, selectedFields string, query string, args ...any) (model.User, error) {
	var res model.User

	db := u.Database.WithContext(ctx).Model(model.User{})
	db = util.SetSelectFields(db, selectedFields)

	if err := db.Where(query, args...).Take(&res).Error; err != nil {
		return model.User{}, err
	}

	return res, nil
}

func (u *User) Find(ctx context.Context, queries []string, argsSlice ...[]interface{}) (model.User, error) {
	var res model.User

	db := u.Database.WithContext(ctx).Model(model.User{})

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
