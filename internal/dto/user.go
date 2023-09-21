package dto

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
)
