package user

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"owlharbour-api/internal/dto"
	"owlharbour-api/internal/factory"
	"owlharbour-api/internal/model"
	"owlharbour-api/internal/repository"
	"owlharbour-api/pkg/constants"
	"owlharbour-api/pkg/helper"
	"owlharbour-api/pkg/util"
	"strconv"
	"text/template"
	"time"
	"unicode/utf8"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	UserRepository repository.User
	ShipRepository repository.Ship
}

type Service interface {
	LoginService(ctx context.Context, payload dto.PayloadLogin, is_mobile bool) (dto.ReturnJwt, error)
	GetProfile(ctx context.Context, userSess any) dto.ProfileUser
	GetAllUsers(ctx context.Context, request dto.UserListParam) ([]dto.AllUser, error)
	DetailUser(ctx context.Context, userID int) (dto.DetailUser, error)
	StoreUser(ctx context.Context, payload dto.PayloadStoreUser) error
	UpdateUser(ctx context.Context, payload dto.PayloadUpdateUser) error
	DeleteUser(ctx context.Context, userID int) error
	ChangePassword(ctx context.Context, userID int, payload dto.PayloadChangePassword) error
	VerifyEmail(ctx context.Context, base64String string) error
	LogoutService(ctx context.Context, userSess any) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
		ShipRepository: f.ShipRepository,
	}
}

func (s *service) LoginService(ctx context.Context, payload dto.PayloadLogin, is_mobile bool) (dto.ReturnJwt, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return dto.ReturnJwt{}, constants.ErrorLoadLocationTime
	}

	param := "email = ?"
	value := payload.Email

	if is_mobile {
		param = "username = ?"
		value = payload.Username
	}

	user, err := s.UserRepository.FindOne(ctx, "id, email, name, password, email_verified_at, role", param, value)
	if err != nil {
		return dto.ReturnJwt{}, constants.UserNotFound
	}

	err = ComparePasswords(user.Password, payload.Password)
	if err != nil {
		fmt.Println(err)
		return dto.ReturnJwt{}, constants.InvalidPassword
	}

	if user.EmailVerifiedAt == nil && !is_mobile {

		tmpl, err := template.ParseFiles("pkg/resource/email_verify.html")
		if err != nil {
			fmt.Println("Error parsing template:", err)
			return dto.ReturnJwt{}, constants.InvalidPassword
		}

		emailByte := []byte(user.Email)
		encodedString := base64.StdEncoding.EncodeToString(emailByte)
		urlVerify := "/api/v1/user/verify/email/"

		data := struct {
			Name string
			Url  string
		}{
			Name: user.Name,
			Url:  util.GetEnv("URL_EMAIL_VERIF", "fallback") + ":" + util.GetEnv("APP_PORT", "fallback") + urlVerify + encodedString,
		}

		var tplBuffer = new(bytes.Buffer)
		errExecute := tmpl.Execute(tplBuffer, data)
		if errExecute != nil {
			fmt.Println("Error executing template:", err)
			return dto.ReturnJwt{}, constants.InvalidPassword
		}

		go helper.SendMail(user.Email, "Verifikasi Akun Simpel", tplBuffer.String())

		return dto.ReturnJwt{}, constants.UserNotVerifyEmail
	}

	secretKey := []byte(util.GetEnv("SECRET_KEY", "fallback"))

	jwt, err := GenerateToken(secretKey, strconv.Itoa(user.ID), user.Email)
	if err != nil {
		return dto.ReturnJwt{}, constants.ErrorGenerateJwt
	}

	if jwt == "" {
		return dto.ReturnJwt{}, constants.EmptyGenerateJwt
	}

	dataUser := dto.DataUserLogin{
		ID:              user.ID,
		Email:           user.Email,
		Name:            user.Name,
		EmailVerifiedAt: user.EmailVerifiedAt,
		Role:            string(user.Role),
	}

	payloadJwtToken := dto.PayloadUpdateJwtToken{
		ID:       user.ID,
		JwtToken: jwt,
	}
	err = s.UserRepository.UpdateJwtToken(ctx, &payloadJwtToken, "jwt_token", "id = ?", user.ID)
	if err != nil {
		log.Println("Error updating jwt token:", err)
		return dto.ReturnJwt{}, constants.ErrorGenerateJwt
	}

	if is_mobile {
		ship, err := s.ShipRepository.FindOne(ctx, "device_id", "user_id = ?", user.ID)
		if err != nil {
			return dto.ReturnJwt{}, err
		}

		if ship.DeviceID != payload.DeviceID {
			err = s.ShipRepository.UpdateShipDeviceID(ctx, payload.DeviceID, user.ID)
			if err != nil {
				return dto.ReturnJwt{}, err
			}
		}
	}

	jwtMode := util.GetEnv("JWT_MODE", "fallback")
	expiredTime := time.Now().In(loc).Add(time.Hour * 730)
	formatExpiredTime := expiredTime.Format("2006-01-02 15:04:05")

	if jwtMode == "release" {
		expiredTime := time.Now().In(loc).Add(time.Hour * 2191)
		formatExpiredTime := expiredTime.Format("2006-01-02 15:04:05")
		return dto.ReturnJwt{
			TokenJwt:  jwt,
			ExpiredAt: formatExpiredTime,
			DataUser:  &dataUser,
		}, nil
	}

	return dto.ReturnJwt{
		TokenJwt:  jwt,
		ExpiredAt: formatExpiredTime,
		DataUser:  &dataUser,
	}, nil
}

func (s *service) GetAllUsers(ctx context.Context, request dto.UserListParam) ([]dto.AllUser, error) {
	var AllUser []dto.AllUser

	users, err := s.UserRepository.GetAll(ctx, request)
	if err != nil {
		return nil, err
	}
	for _, user := range users {
		user := dto.AllUser{
			ID:              user.ID,
			Name:            user.Name,
			Username:        user.Username,
			Email:           user.Email,
			Role:            user.Role,
			EmailVerifiedAt: user.EmailVerifiedAt,
			CreatedAt:       user.CreatedAt,
			UpdatedAt:       user.UpdatedAt,
		}
		AllUser = append(AllUser, user)
	}

	return AllUser, nil
}

func (s *service) GetProfile(ctx context.Context, userSess any) dto.ProfileUser {
	return dto.ProfileUser{
		ID:    userSess.(model.User).ID,
		Name:  userSess.(model.User).Name,
		Email: userSess.(model.User).Email,
		Role:  string(userSess.(model.User).Role),
	}
}

func (s *service) LogoutService(ctx context.Context, userSess any) error {
	ID := userSess.(model.User).ID
	payloadJwtToken := dto.PayloadUpdateJwtToken{
		ID:       ID,
		JwtToken: "",
	}
	err := s.UserRepository.RemoveJwtToken(ctx, &payloadJwtToken, "jwt_token", "id = ?", ID)
	if err != nil {
		log.Println("Error updating jwt token:", err)
		return err
	}

	return nil
}

func (s *service) DetailUser(ctx context.Context, userID int) (dto.DetailUser, error) {
	user, err := s.UserRepository.FindOne(ctx, "id, email, name, role, created_at, updated_at", "id = ?", userID)
	if err != nil {
		return dto.DetailUser{}, constants.NotFoundDataUser
	}

	tCreatedAt := user.CreatedAt
	tUpdatedAt := user.UpdatedAt
	formatCreatedAt := tCreatedAt.Format("2006-01-02 15:04:05")
	formatUpdateddAt := tUpdatedAt.Format("2006-01-02 15:04:05")

	data := dto.DetailUser{
		ID:        userID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      string(user.Role),
		CreatedAt: formatCreatedAt,
		UpdatedAt: formatUpdateddAt,
	}

	return data, nil
}

func (s *service) StoreUser(ctx context.Context, payload dto.PayloadStoreUser) error {
	_, err := s.UserRepository.FindOne(ctx, "id,email,name,password", "email = ?", payload.Email)

	if err != nil {
		password := []byte(payload.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			return constants.ErrorHashPassword
		}

		tmpl, err := template.ParseFiles("pkg/resource/email_verify.html")

		emailByte := []byte(payload.Email)
		encodedString := base64.StdEncoding.EncodeToString(emailByte)
		urlVerify := "/api/v1/user/verify/email/"

		data := struct {
			Name string
			Url  string
		}{
			Name: payload.Name,
			Url:  util.GetEnv("URL_EMAIL_VERIF", "fallback") + ":" + util.GetEnv("APP_PORT", "fallback") + urlVerify + encodedString,
		}

		var tplBuffer = new(bytes.Buffer)
		errExecute := tmpl.Execute(tplBuffer, data)
		if errExecute != nil {
			fmt.Println("Error executing template:", err)
			return constants.DuplicateStoreUser
		}

		go helper.SendMail(payload.Email, "Verifikasi Akun Simpel", tplBuffer.String())
		dataStore := model.User{
			Name:     payload.Name,
			Email:    payload.Email,
			Password: string(hashedPassword),
			Role:     model.RoleType(payload.Role),
		}
		s.UserRepository.Store(ctx, dataStore)
		return nil
	}

	return constants.DuplicateStoreUser
}

func (s *service) UpdateUser(ctx context.Context, payload dto.PayloadUpdateUser) error {
	user, err := s.UserRepository.FindOne(ctx, "id", "id = ?", payload.ID)
	if err != nil {
		return constants.NotFoundDataUser
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return constants.ErrorLoadLocationTime
	}

	updateUser := dto.PayloadUpdateUser{
		Name:      payload.Name,
		Email:     payload.Email,
		Role:      payload.Role,
		UpdatedAt: time.Now().In(loc),
	}

	if payload.Password != "" {
		password := []byte(payload.Password)
		hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			return constants.ErrorHashPassword
		}
		updateUser.Password = string(hashedPassword)
		err = s.UserRepository.UpdateOne(ctx, &updateUser, "name,email,role,password,updated_at", "id = ?", user.ID)
		if err != nil {
			log.Println("Error updating user:", err)
			return constants.FailedUpdateUser
		}
		return nil
	}

	err = s.UserRepository.UpdateOne(ctx, &updateUser, "name,email,role,updated_at", "id = ?", user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return constants.FailedUpdateUser
	}

	return nil
}

func (s *service) ChangePassword(ctx context.Context, userID int, payload dto.PayloadChangePassword) error {
	user, err := s.UserRepository.FindOne(ctx, "id,password", "id = ?", userID)
	if err != nil {
		return constants.NotFoundDataUser
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return constants.ErrorLoadLocationTime
	}

	if payload.PasswordConfirmation != payload.Password {
		return constants.FailedNotSamePassword
	}

	charCount := utf8.RuneCountInString(payload.Password)

	if charCount < 8 {
		return constants.MinimCharacterPassword
	}

	password := []byte(payload.Password)
	hashedPassword, err := bcrypt.GenerateFromPassword(password, bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return constants.ErrorHashPassword
	}
	checkPassword := ComparePasswords(user.Password, payload.Password)
	if checkPassword == nil {
		return constants.PasswordSameCurrent
	}

	updateUser := dto.
		PayloadUpdateUser{
		Password:  string(hashedPassword),
		UpdatedAt: time.Now().In(loc),
	}
	err = s.UserRepository.UpdateOne(ctx, &updateUser, "password,updated_at", "id = ?", user.ID)
	if err != nil {
		log.Println("Error change password:", err)
		return constants.FailedUpdateUser
	}
	return nil
}

func (s *service) VerifyEmail(ctx context.Context, base64String string) error {
	emailDecode, err := base64.StdEncoding.DecodeString(base64String)
	if err != nil {
		fmt.Println("Error:", err)
		return constants.ErrorDecodeBase64
	}
	user, err := s.UserRepository.FindOne(ctx, "id,email_verified_at", "email = ?", emailDecode)
	if err != nil {
		return constants.NotFoundDataUser
	}

	if user.EmailVerifiedAt != nil {
		return nil
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return constants.ErrorLoadLocationTime
	}

	updateUser := dto.PayloadUpdateUser{
		EmailVerifiedAt: time.Now().In(loc),
	}

	err = s.UserRepository.UpdateOne(ctx, &updateUser, "email_verified_at", "id = ?", user.ID)
	if err != nil {
		log.Println("Error updating user:", err)
		return constants.FailedVerifyEmail
	}

	return nil
}

func (s *service) DeleteUser(ctx context.Context, userID int) error {
	user, err := s.UserRepository.FindOne(ctx, "id", "id = ?", userID)
	if err != nil {
		return constants.NotFoundDataUser
	}

	err = s.UserRepository.DeleteOne(ctx, "id = ?", user.ID)
	if err != nil {
		return constants.FailedDeleteUser
	}

	return nil
}

func ComparePasswords(hashedPassword, inputPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(inputPassword))
}

func GenerateToken(secretKey []byte, userID string, email string) (string, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return "", err
	}

	jwtMode := util.GetEnv("JWT_MODE", "fallback")
	expiredTime := time.Now().In(loc).Add(time.Hour * 730).Unix()

	if jwtMode == "release" {
		expiredTime := time.Now().In(loc).Add(time.Hour * 2191).Unix()
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id": userID,
			"email":   email,
			"exp":     expiredTime,
		})
		if err != nil {
			return "", err
		}
		tokenString, err := token.SignedString(secretKey)
		if err != nil {
			return "", err
		}

		return tokenString, nil
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     expiredTime,
	})

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
