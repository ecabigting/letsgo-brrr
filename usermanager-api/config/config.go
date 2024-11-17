package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	JWTSecret string
	Port      string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
	}
	AppConfig = Config{
		MongoURI:  getEnv("MONGO_URI", "///"),
		JWTSecret: getEnv("JWTSecret", "///"),
		Port:      getEnv("PORT", "9090"),
	}
}

/*
* your env file
* MONGO_URI=mongodb://localhost:27017
* JWT_SECRET=your_jwt_secret
* PORT=8080
* */

func getEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}
