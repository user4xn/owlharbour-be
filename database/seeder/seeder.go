package seeder

import (
	"fmt"
	"owlharbour-api/database"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type User struct {
	ID              uint
	Name            string
	Email           string
	Role            string
	Password        string
	EmailVerifiedAt time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

func Seed() {
	db := database.GetConnection()

	passwordSA := []byte("superadminpassword")
	hashedPasswordSA, err := bcrypt.GenerateFromPassword(passwordSA, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}

	passwordAdmin := []byte("adminpassword")
	hashedPasswordAdmin, err := bcrypt.GenerateFromPassword(passwordAdmin, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
		return
	}

	seedData := []User{
		{
			Name:            "Super Admin",
			Email:           "superadmin@gmail.com",
			Role:            "superadmin",
			Password:        string(hashedPasswordSA),
			EmailVerifiedAt: time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		{
			Name:            "Admin",
			Email:           "admin@gmail.com",
			Role:            "admin",
			Password:        string(hashedPasswordAdmin),
			EmailVerifiedAt: time.Now(),
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
	}

	for _, data := range seedData {
		user := User{}
		err := db.Where("email = ?", data.Email).First(&user).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			if err := db.Create(&data).Error; err != nil {
				fmt.Println(err)
				return
			}
		}
	}
}
