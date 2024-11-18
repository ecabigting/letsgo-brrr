package services

import (
	"context"
	"errors"
	"log"
	"time"

	"github.com/ecabigting/letsgo-brrr/usermanager-api/models"
	"github.com/ecabigting/letsgo-brrr/usermanager-api/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	collection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		collection: db.Collection("users"),
	}
}

func (s *UserService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedDate = time.Now()
	user.VerificationToken, _ = utils.GenerateVerificationToken()

	_, err = s.collection.InsertOne(context.Background(), user)
	return err
}

func (s *UserService) VerifyUser(userID string, token string) error {
	var user models.User
	// Convert userID to primitive Objectid
	objectId, errOId := primitive.ObjectIDFromHex(userID)
	if errOId != nil {
		log.Println("Invalid id")
	}
	err := s.collection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return err
	}
	if user.VerificationToken != token {
		return errors.New("invalid verification token")
	}
	user.VerifiedDate = time.Now()
	user.VerificationToken = ""
	user.IsEnabled = true

	_, err = s.collection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": user})
	return err
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := s.collection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID string, user *models.User) error {
	_, err := s.collection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": user})
	return err
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(userID string) error {
	_, err := s.collection.DeleteOne(context.Background(), bson.M{"_id": userID})
	return err
}

// Check if email address exist
func (s *UserService) CheckIfEmailExist(email string) bool {
	count, err := s.collection.CountDocuments(context.Background(), bson.M{"email": email})
	if err != nil {
		log.Println(err)
	}

	if count > 0 {
		return true
	}
	return false
}
