package middleware

import (
	"fmt"
	"net/http"
	"regexp"
	"simpel-api/internal/factory"
	"simpel-api/pkg/util"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.Request.Header["Authorization"]

		if len(header) == 0 {
			response := util.APIResponse("Sorry, you didn't enter a valid bearer token", http.StatusUnauthorized, "failed", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		rep := regexp.MustCompile(`(Bearer)\s?`)
		bearerStr := rep.ReplaceAllString(header[0], "")
		parsedToken, err := parseToken(bearerStr)
		if err != nil || !parsedToken.Valid {
			response := util.APIResponse("Invalid bearer token", http.StatusUnauthorized, "failed", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		claims := parsedToken.Claims.(jwt.MapClaims)

		f := factory.NewFactory()
		userId, _ := strconv.Atoi(claims["user_id"].(string))
		user, err := f.UserRepository.FindOne(c, "id,email,name,role,jwt_token", "id = ?", userId)
		if user.JwtToken != bearerStr {
			response := util.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}
		if err != nil {
			response := util.APIResponse("Unauthorized", http.StatusUnauthorized, "failed", nil)
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Set("user", user)
		c.Set("bearer", bearerStr)

		c.Next()
	}
}

func parseToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(util.GetEnv("SECRET_KEY", "fallback"))
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return secretKey, nil
	})

	return token, err
}
