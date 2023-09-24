package user

import (
	"context"
	"fmt"
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
	GetAllUsers(ctx context.Context, Search string) []dto.AllUser
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

func (s *service) GetAllUsers(ctx context.Context, Search string) []dto.AllUser {
	var AllUser []dto.AllUser
	selectedFields := "id, name, email, role, created_at, updated_at"
	searchQuery := "Name Like ?"
	args := []interface{}{"%" + Search + "%"}
	users, err := s.UserRepository.GetAll(ctx, selectedFields, searchQuery, args...)
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
