package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"simpel-api/internal/dto"
	"simpel-api/internal/factory"
	"simpel-api/pkg/constants"
	"simpel-api/pkg/util"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func BearerToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header["Authorization"]
		if len(header) == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Common{
				Status:  "failed",
				Code:    401,
				Message: "Unauthenticated",
			})
			return
		}

		rep := regexp.MustCompile(`(Bearer)\s?`)
		bearerStr := rep.ReplaceAllString(header[0], "")
		parsedToken, err := parseToken(bearerStr)
		claims := parsedToken.Claims.(jwt.MapClaims)

		f := factory.NewFactory()
		userId, _ := strconv.Atoi(claims["user_id"].(string))
		user, err := f.UserRepository.FindOne(c, "id,email,name,role", "id = ?", userId)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, dto.Common{
				Status:  "failed",
				Code:    401,
				Message: constants.BearerTokenHasError.Error(),
			})
			return
		}

		c.Set("user", user)
		c.Set("bearer", bearerStr)

		c.Next()
		return
	}
}

func parseToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(util.GetEnv("SECRET_KEY", "fallback"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method")
		}
		return secretKey, nil
	})

	return token, err
}
