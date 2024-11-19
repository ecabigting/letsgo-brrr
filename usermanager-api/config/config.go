package config

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ecabigting/letsgo-brrr/usermanager-api/models"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/services"
	"github.com/jaswdr/faker/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type Config struct {
	ConfigService *services.ConfigService
	Database      string
	MongoURI      string
	JWTSecret     string
	Port          string
	DefP          string
	DefE          string
}

var AppConfig Config

func LoadConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error Loading .env file")
		os.Exit(1)
	}
	AppConfig = Config{
		MongoURI:  getEnv("MONGO_URI", "mongodb://localhost:27017"),
		Database:  getEnv("DATABASE", "usermanager-api"),
		JWTSecret: getEnv("JWT_SECRET", ""),
		Port:      getEnv("PORT", "9090"),
		DefP:      getEnv("DEFP", "d3f@ult0415"),
		DefE:      getEnv("DEFE", "ericcabigting@outlook.com"),
	}
}

/*
* your env file
* MONGO_URI=mongodb://localhost:27017
* Database=usermanager-api
* JWT_SECRET=your_jwt_secret
* PORT=8080
* DEFP=d3f@ult0415
* DEFE=ericcabigting@outlook.com
* */

func getEnv(key, fallback string) string {
	if value, exist := os.LookupEnv(key); exist {
		return value
	}
	return fallback
}

func SeedDB(db *mongo.Database) error {
	userCollection := db.Collection("users")
	count, err := userCollection.CountDocuments(context.Background(), bson.M{"role": "Admin"})
	if err != nil {
		log.Println(err)
		return err
	}

	if count < 1 {
		fmt.Println("Seeding Database...")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(AppConfig.DefP), bcrypt.DefaultCost)
		if err != nil {
			return nil
		}
		fake := faker.New()
		adminUser := models.User{
			Email:             AppConfig.DefE,
			Password:          string(hashedPassword),
			Role:              "Admin",
			FirstName:         fake.Person().FirstName(),
			MiddleName:        fake.Person().NameMale(),
			LastName:          fake.Person().LastName(),
			Gender:            fake.Person().Gender(),
			DateOfBirth:       fake.Person().Faker.Time().Time(time.Now().AddDate(-20, 0, 0)),
			IsEnabled:         true,
			IsEnabledByDate:   time.Now(),
			VerificationToken: "",
			VerifiedDate:      time.Now(),
			CreatedDate:       time.Now(),
		}
		_, err = userCollection.InsertOne(context.Background(), adminUser)
		if err != nil {
			return err
		}
		fmt.Println("Added new user: ", AppConfig.DefE)
	}
	return nil
}
