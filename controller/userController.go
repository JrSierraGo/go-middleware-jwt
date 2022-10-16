package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"middleware/database"
	"middleware/models"
	"net/http"
	"os"
	"time"
)

var userGlobal models.User

const tableName = "users"

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

	user.Id = uuid.New().String()

	if tx := database.Db.Table(tableName).Create(&user); tx.Error != nil {

		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": tx.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"status": "Successful",
	})

}

func LogIn(ctx *gin.Context) {
	mySigningKey := os.Getenv("SECRET_SIGN")
	var userParam models.User
	if err := ctx.Bind(&userParam); err != nil {
		panic(err)
	}

	var userDB models.User
	database.Db.Table(tableName).Where("email = ?", userParam.Email).First(&userDB)

	passwordMatch := userDB.CheckPasswordHash(userParam.Password)

	if !passwordMatch || userParam.Email == "" || userDB == (models.User{}) {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"messageError": "User or password incorrect",
		})
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"email": userParam.Email,
		"exp":   time.Now().Add(1 * time.Minute).Unix(),
	})

	tokenString, _ := token.SignedString([]byte(mySigningKey))

	ctx.JSON(http.StatusOK, gin.H{
		"userEmail":   userParam.Email,
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
