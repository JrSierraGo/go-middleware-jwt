package middleware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
	"os"
)

func ValidateAuth(context *gin.Context) {
	mySigningKey := os.Getenv("SECRET_SIGN")
	authorization := context.GetHeader("Authorization")
	if authorization == "" {
		context.AbortWithStatusJSON(http.StatusExpectationFailed, gin.H{
			"messageError": "Header not present",
		})
	}

	token, err := jwt.Parse(authorization, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(mySigningKey), nil
	})

	if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"messageError": err.Error(),
		})
	}

	context.Next()
}
