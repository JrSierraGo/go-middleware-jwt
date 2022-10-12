package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"middleware/controller"
	"middleware/middleware"
	"os"
)

func main() {
	loadEnvFiles()
	startGin()
}

func startGin() {
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Recovery())

	//Public Routes
	router.POST("/sign-in", controller.SignIn)
	router.POST("/log-in", controller.LogIn)

	//Private Routes
	usersRouter := router.Group("/users")
	usersRouter.Use(gin.Recovery())
	usersRouter.Use(middleware.ValidateAuth)
	usersRouter.GET("/list", controller.GetAllUsers)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Server start in port", port)
	if err := router.Run(":" + port); err != nil {
		log.Panicf("error: %s", err)
	}
}

func loadEnvFiles() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}
}
