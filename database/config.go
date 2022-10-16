package database

import "fmt"

type Config struct {
	Host     string
	User     string
	Password string
	DBName   string
	Port     string
}

var GetConnectionString = func(config Config) string {

	connectionString := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", config.Host,
		config.User, config.Password,
		config.DBName, config.Port)

	return connectionString
}
