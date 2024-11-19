package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	UserID            primitive.ObjectID `bson:"_id,omitempty"`
	Email             string             `bson:"email" json:"email" binding:"required,email"`
	Password          string             `bson:"password" json:"password" validate:"password_complexity"`
	Role              string             `bson:"role"`
	FirstName         string             `bson:"firstName" json:"firstName" binding:"required"`
	MiddleName        string             `bson:"middleName" json:"middleName"`
	LastName          string             `bson:"lastName" json:"lastName" binding:"required"`
	Gender            string             `bson:"gender"`
	DateOfBirth       time.Time          `bson:"dateOfBirth" json:"dateOfBirth"`
	CreatedDate       time.Time          `bson:"createdDate"`
	CreatedByID       string             `bson:"createdById"`
	UpdatedDate       time.Time          `bson:"updatedDate"`
	UpdatedByID       primitive.ObjectID `bson:"updatedById"`
	IsEnabled         bool               `bson:"isEnabled"`
	IsEnabledByID     primitive.ObjectID `bson:"isEnabledByID"`
	IsEnabledByDate   time.Time          `bson:"isEnabledByDate"`
	VerificationToken string             `bson:"verificationToken"`
	VerifiedDate      time.Time          `bson:"verifiedDate"`
	VerifiedBy        primitive.ObjectID `bson:"verifiedBy"`
}
