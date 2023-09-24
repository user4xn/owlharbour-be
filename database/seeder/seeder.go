package seeder

import (
	"simpel-api/database"

	"gorm.io/gorm"
)

type User struct {
	ID       uint
	Name     string
	Email    string
	Role     string
	Password string
}

func Seed() {

	db := database.GetConnection()

	password := "mysecretpassword" // Note: In real-world scenarios, you should use a library to securely hash passwords.

	seedData := []User{
		{
			Name:     "Super Admin",
			Email:    "superadmin@gmail.com",
			Role:     "superadmin",
			Password: password,
		},
		{
			Name:     "Admin",
			Email:    "admin@gmail.com",
			Role:     "admin",
			Password: password,
		},
	}

	for _, data := range seedData {
		user := User{}
		err := db.Where("email = ?", data.Email).First(&user).Error
		if err != nil && err == gorm.ErrRecordNotFound {
			db.Create(&data)
		}
	}
}
