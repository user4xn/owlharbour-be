package util

import (
	"strings"

	"gorm.io/gorm"
)

func SetSelectFields(db *gorm.DB, selectFieldStr string) *gorm.DB {
	if selectFieldStr == "*" {
		return db
	}

	selectFieldArr := strings.Split(selectFieldStr, ",")
	if len(selectFieldArr) == 0 {
		return db
	}

	return db.Select(selectFieldArr)
}
