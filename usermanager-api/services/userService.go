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
	userCollection       *mongo.Collection
	userDeviceCollection *mongo.Collection
}

func NewUserService(db *mongo.Database) *UserService {
	return &UserService{
		userCollection:       db.Collection("users"),
		userDeviceCollection: db.Collection("userdevices"),
	}
}

// Create New User
func (s *UserService) CreateUser(user *models.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	user.Password = string(hashedPassword)
	user.CreatedDate = time.Now()
	user.VerificationToken, _ = utils.GenerateVerificationToken()

	_, err = s.userCollection.InsertOne(context.Background(), user)
	return err
}

// Verify User with verification Token
func (s *UserService) VerifyUser(userID string, token string) error {
	var user models.User
	// Convert userID to primitive Objectid
	objectId, errOId := primitive.ObjectIDFromHex(userID)
	if errOId != nil {
		return errors.New("Invalid ID:1")
	}
	err := s.userCollection.FindOne(context.Background(), bson.M{"_id": objectId}).Decode(&user)
	if err != nil {
		return errors.New("Invalid ID:2")
	}
	if user.VerificationToken != token {
		return errors.New("Invalid Verification Token")
	}
	user.VerifiedDate = time.Now()
	user.VerificationToken = ""
	user.VerifiedDate = time.Now()
	user.VerifiedBy = objectId
	user.UpdatedByID = objectId
	user.IsEnabled = true
	user.IsEnabledByID = objectId

	_, err = s.userCollection.UpdateOne(context.Background(), bson.M{"_id": objectId}, bson.M{"$set": user})
	return err
}

// GetUser retrieves a user by ID
func (s *UserService) GetUser(userID string) (*models.User, error) {
	var user models.User
	err := s.userCollection.FindOne(context.Background(), bson.M{"_id": userID}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// UpdateUser updates user information
func (s *UserService) UpdateUser(userID string, user *models.User) error {
	_, err := s.userCollection.UpdateOne(context.Background(), bson.M{"_id": userID}, bson.M{"$set": user})
	return err
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(userID string) error {
	_, err := s.userCollection.DeleteOne(context.Background(), bson.M{"_id": userID})
	return err
}

// Check if email address exist
func (s *UserService) CheckIfEmailExist(email string) bool {
	count, err := s.userCollection.CountDocuments(context.Background(), bson.M{"email": email})
	if err != nil {
		log.Println(err)
	}

	if count > 0 {
		return true
	}
	return false
}

// Get user by Email
func (s *UserService) GetUserByEmail(login *models.Login) (*models.User, error) {
	var user models.User

	// Get the User with the associated email
	err := s.userCollection.FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
