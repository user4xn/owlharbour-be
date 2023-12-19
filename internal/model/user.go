package model

import (
	"time"
)

type User struct {
	Common
	Name            string     `gorm:"varchar"`
	Role            RoleType   `gorm:"enum:superadmin,admin"`
	Username        string     `gorm:"varchar"`
	Email           string     `gorm:"varchar"`
	EmailVerifiedAt *time.Time `gorm:"timestamp"`
	Password        string     `gorm:"varchar"`
	JwtToken        string     `gorm:"text"`
}
