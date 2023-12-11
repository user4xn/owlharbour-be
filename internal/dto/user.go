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

	PayloadChangePassword struct {
		Password             string `json:"password" binding:"required"`
		PasswordConfirmation string `json:"password_confirmation" binding:"required"`
	}

	UserListParam struct {
		Search string `json:"search"`
		Limit  int    `json:"limit"`
		Offset int    `json:"offset"`
	}

	PayloadStoreUser struct {
		Name      string    `json:"name" binding:"required"`
		Email     string    `json:"email" binding:"required"`
		Password  string    `json:"password" binding:"required"`
		Role      string    `json:"role" binding:"required"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	PayloadUpdateUser struct {
		ID              int       `json:"id" binding:"required"`
		Name            string    `json:"name" binding:"required"`
		Email           string    `json:"email" binding:"required"`
		Password        string    `json:"password"`
		Role            string    `json:"role" binding:"required"`
		EmailVerifiedAt time.Time `json:"email_verified_at"`
		UpdatedAt       time.Time `json:"updated_at"`
	}

	PayloadUpdateJwtToken struct {
		ID       int    `json:"id"`
		JwtToken string `json:"jwt_token"`
	}

	ReturnJwt struct {
		TokenJwt  string         `json:"token_jwt"`
		ExpiredAt string         `json:"expired_at"`
		DataUser  *DataUserLogin `json:"data_user"`
	}

	DataUserLogin struct {
		ID              int        `json:"id"`
		Email           string     `json:"email"`
		Name            string     `json:"name"`
		EmailVerifiedAt *time.Time `json:"email_verify_at"`
		Role            string     `json:"role"`
	}

	ProfileUser struct {
		ID    int    `json:"id"`
		Email string `json:"email"`
		Name  string `json:"name"`
		Role  string `json:"role"`
	}

	DetailUser struct {
		ID        int    `json:"id"`
		Email     string `json:"email"`
		Name      string `json:"name"`
		Role      string `json:"role"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}

	AllUser struct {
		ID              int    `json:"id"`
		Email           string `json:"email"`
		Name            string `json:"name"`
		Role            string `json:"role"`
		EmailVerifiedAt string `json:"email_verified_at"`
		CreatedAt       string `json:"created_at"`
		UpdatedAt       string `json:"updated_at"`
	}
)
