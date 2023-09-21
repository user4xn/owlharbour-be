package repository

import (
	"context"
	"simpel-api/internal/model"
	"simpel-api/pkg/util"

	"gorm.io/gorm"
)

type UserInterface interface {
	Init() *gorm.DB
	Find(ctx context.Context, queries []string, argsSlice ...[]interface{}) (model.User, error)
	Store(ctx context.Context, data model.User) (int, error)
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

func (u *User) Store(ctx context.Context, data model.User) (int, error) {
	tx := u.Database.WithContext(ctx)
	if err := tx.Model(model.User{}).Create(&data).Error; err != nil {
		return 0, err
	}

	return data.ID, nil
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
