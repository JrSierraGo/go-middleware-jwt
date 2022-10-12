package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"middleware/models"
	"net/http"
	"os"
	"time"
)

var userGlobal models.User

func SignIn(ctx *gin.Context) {
	var user models.User

	if err := ctx.Bind(&user); err != nil {
		panic(err)
	}

	if user == (models.User{}) || user.Email == "" || user.Password == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"status":       "Failed",
			"messageError": "Require body not present",
		})
		return
	}

	if err := user.HashPassword(); err != nil {
		panic(err)
	}

	userGlobal = user

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "Successful",
	})

}

func LogIn(ctx *gin.Context) {
	mySigningKey := os.Getenv("SECRET_SIGN")
	var user models.User
	if err := ctx.Bind(&user); err != nil {
		panic(err)
	}

	passwordMatch := userGlobal.CheckPasswordHash(user.Password)

	if !passwordMatch || user.Email == "" {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"messageError": "User or password incorrect",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": user.Email,
		"exp":   time.Now().Add(1 * time.Minute).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(mySigningKey))

	ctx.JSON(http.StatusOK, gin.H{
		"userEmail":   user.Email,
		"accessToken": tokenString,
	})

}

func GetAllUsers(ctx *gin.Context) {
	var users []models.User
	users = append(users, userGlobal)

	ctx.JSON(http.StatusOK, gin.H{
		"users": users,
	})
}
