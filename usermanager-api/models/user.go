package models

import (
	"log"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID primitive.ObjectID `bson:"_id,omitempty"`
	Email  string             `bson:"email" json:"email" binding:"required,email"`
	// TODO : add validation for password specs
	Password          string    `bson:"password" json:"password" binding:"required"`
	Role              string    `bson:"role"`
	FirstName         string    `bson:"firstName"`
	MiddleName        string    `bson:"middleName"`
	LastName          string    `bson:"lastName"`
	Gender            string    `bson:"gender"`
	DateOfBirth       time.Time `bson:"dateOfBirth"`
	CreatedDate       time.Time `bson:"createdDate"`
	CreatedByID       string    `bson:"createdById"`
	UpdatedDate       time.Time `bson:"updatedDate"`
	UpdatedByID       string    `bson:"updatedById"`
	IsEnabled         bool      `bson:"isEnabled"`
	IsEnabledByID     string    `bson:"isEnabledByID"`
	IsEnabledByDate   time.Time `bson:"isEnabledByDate"`
	VerificationToken string    `bson:"verificationToken"`
	VerifiedDate      time.Time `bson:"verifiedDate"`
	VerifiedBy        string    `bson:"verifiedBy"`
}

func (u *User) Validate() (string, bool) {
	log.Println(u.Email)
	return "", true
}