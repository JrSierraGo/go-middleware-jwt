package main

import (
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"log"
	"middleware/controller"
	"middleware/database"
	"middleware/middleware"
	"os"
)

func main() {
	loadEnvFiles()
	initDB()
	startGin()
}

func initDB() {
	config := database.Config{
		Host:     os.Getenv("DB_HOST"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   os.Getenv("DB_NAME"),
		Port:     os.Getenv("DB_PORT"),
	}
	connectionString := database.GetConnectionString(config)
	err := database.Connect(connectionString)
	if err != nil {
		panic(err.Error())
	}
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

	port := os.Getenv("Port")
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
