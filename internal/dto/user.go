package dto

import "time"

type RoleType string

const (
	SuperAdmin RoleType = "superadmin"
	Admin      RoleType = "admin"
)

type (
	PayloadLogin struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	PayloadStoreUser struct {
		Name      string    `json:"name" binding:"required"`
		Email     string    `json:"email" binding:"required"`
		Password  string    `json:"password" binding:"required"`
		Role      string    `json:"role" binding:"required"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	ReturnJwt struct {
		TokenJwt  string         `json:"token_jwt"`
		ExpiredAt string         `json:"expired_at"`
		DataUser  *DataUserLogin `json:"data_user"`
	}

	DataUserLogin struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
	}

	ProfileUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	}

	AllUser struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}
)
