package seeder

import (
	"fmt"
	"simpel-api/database"

	"golang.org/x/crypto/bcrypt"
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

	password := []byte("mysecretpassword")
	hashedPasswordSA, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		fmt.Println(err)
	}

	seedData := []User{
		{
			Name:     "Super Admin",
			Email:    "superadmin@gmail.com",
			Role:     "superadmin",
			Password: string(hashedPasswordSA),
		},
		{
			Name:     "Admin",
			Email:    "admin@gmail.com",
			Role:     "admin",
			Password: string(hashedPasswordSA),
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
