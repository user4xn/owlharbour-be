package user

import (
	"context"
	"fmt"
	"log"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/model"
	"simpel-api/internal/repository"
	"simpel-api/pkg/constants"
	"simpel-api/pkg/util"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type service struct {
	UserRepository repository.UserInterface
}

type Service interface {
	LoginService(ctx context.Context, payload dto.PayloadLogin) (dto.ReturnJwt, error)
	GetProfile(ctx context.Context, userSess any) dto.ProfileUser
	GetAllUsers(ctx context.Context, Search string, limit int, offset int) []dto.AllUser
	DetailUser(ctx context.Context, userID int) (dto.DetailUser, error)
	StoreUser(ctx context.Context, payload dto.PayloadStoreUser) error
	UpdateUser(ctx context.Context, payload dto.PayloadUpdateUser) error
	DeleteUser(ctx context.Context, userID int) error
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
	}
}

func (s *service) LoginService(ctx context.Context, payload dto.PayloadLogin) (dto.ReturnJwt, error) {

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		return dto.ReturnJwt{}, constants.ErrorLoadLocationTime
	}
	user, err := s.UserRepository.FindOne(ctx, "id,email,name,password", "email = ?", payload.Email)
	if err != nil {
		return dto.ReturnJwt{}, constants.UserNotFound
	}

	err = ComparePasswords(user.Password, payload.Password)
	if err != nil {
		fmt.Println(err)
		return dto.ReturnJwt{}, constants.InvalidPassword
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
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
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

func (s *service) GetAllUsers(ctx context.Context, Search string, limit int, offset int) []dto.AllUser {
	var AllUser []dto.AllUser
	selectedFields := "id, name, email, role, created_at, updated_at"
	searchQuery := "Name Like ?"
	args := []interface{}{"%" + Search + "%"}
	users, err := s.UserRepository.GetAll(ctx, selectedFields, searchQuery, limit, offset, args...)
	if err != nil {
		return AllUser
	}

	for _, user := range users {
		user := dto.AllUser{
			ID:        user.ID,
			Name:      user.Name,
			Email:     user.Email,
			Role:      user.Role,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		AllUser = append(AllUser, user)
	}

	return AllUser
}

func (s *service) GetProfile(ctx context.Context, userSess any) dto.ProfileUser {
	return dto.ProfileUser{
		ID:    userSess.(model.User).ID,
		Name:  userSess.(model.User).Name,
		Email: userSess.(model.User).Email,
		Role:  fmt.Sprintf("%s", userSess.(model.User).Role),
	}
}

func (s *service) DetailUser(ctx context.Context, userID int) (dto.DetailUser, error) {
	user, err := s.UserRepository.FindOne(ctx, "id,email,name,role,created_at,updated_at", "id = ?", userID)
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
		Role:      fmt.Sprintf("%s", user.Role),
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
