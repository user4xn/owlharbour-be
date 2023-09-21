package user

import (
	"context"
	"errors"
	"fmt"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/internal/model"
	"simpel-api/internal/repository"
	"simpel-api/pkg/util"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type service struct {
	UserRepository repository.UserInterface
}

type Service interface {
	LoginService(ctx context.Context, payload dto.PayloadLogin) (dto.ReturnJwt, error)
	GetProfile(ctx context.Context, userSess any) dto.ProfileUser
}

func NewService(f *factory.Factory) Service {
	return &service{
		UserRepository: f.UserRepository,
	}
}

func (s *service) LoginService(ctx context.Context, payload dto.PayloadLogin) (dto.ReturnJwt, error) {
	loc, err := time.LoadLocation("Asia/Jakarta")
	user, err := s.UserRepository.FindOne(ctx, "id,email,name,password", "email = ?", payload.Email)
	if err != nil && errors.Is(err, gorm.ErrRecordNotFound) == false {
		fmt.Println(err)
		return dto.ReturnJwt{}, err
	}
	err = ComparePasswords(user.Password, payload.Password)
	if err != nil {
		fmt.Println(err)
		return dto.ReturnJwt{}, err
	}
	secretKey := []byte(util.GetEnv("SECRET_KEY", "fallback"))

	jwt, err := generateToken(secretKey, strconv.Itoa(user.ID), user.Email)
	if err != nil {
		fmt.Println(err)
		return dto.ReturnJwt{}, err
	}

	if jwt == "" {
		fmt.Println(err)
		return dto.ReturnJwt{}, err
	}

	dataUser := dto.DataUserLogin{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}
	jwtMode := util.GetEnv("JWT_MODE", "fallback")
	if jwtMode == "realse" {
		expiredTime := time.Now().In(loc).Add(time.Hour * 2191)
		formatExpiredTime := expiredTime.Format("2006-01-02 15:04:05")
		return dto.ReturnJwt{
			TokenJwt:  jwt,
			ExpiredAt: formatExpiredTime,
			DataUser:  &dataUser,
		}, nil
	}
	expiredTime := time.Now().In(loc).Add(time.Hour * 730)
	formatExpiredTime := expiredTime.Format("2006-01-02 15:04:05")
	return dto.ReturnJwt{
		TokenJwt:  jwt,
		ExpiredAt: formatExpiredTime,
		DataUser:  &dataUser,
	}, nil
}

func (s *service) GetProfile(ctx context.Context, userSess any) dto.ProfileUser {
	fmt.Println(userSess)
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

func generateToken(secretKey []byte, userID string, email string) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["email"] = email
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	tokenString, err := token.SignedString(secretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
