package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	IsEnabled         bool               `bson:"isEnabled"`
	UserID            primitive.ObjectID `bson:"_id,omitempty"`
	VerifiedBy        primitive.ObjectID `bson:"verifiedBy"`
	IsEnabledByID     primitive.ObjectID `bson:"isEnabledByID"`
	UpdatedByID       primitive.ObjectID `bson:"updatedById"`
	CreatedDate       time.Time          `bson:"createdDate"`
	DateOfBirth       time.Time          `bson:"dateOfBirth" json:"dateOfBirth"`
	UpdatedDate       time.Time          `bson:"updatedDate"`
	VerifiedDate      time.Time          `bson:"verifiedDate"`
	IsEnabledByDate   time.Time          `bson:"isEnabledByDate"`
	Gender            string             `bson:"gender"`
	LastName          string             `bson:"lastName" json:"lastName" binding:"required"`
	CreatedByID       string             `bson:"createdById"`
	VerificationToken string             `bson:"verificationToken"`
	Email             string             `bson:"email" json:"email" binding:"required,email"`
	Password          string             `bson:"password" json:"password" validate:"password_complexity"`
	Role              string             `bson:"role"`
	FirstName         string             `bson:"firstName" json:"firstName" binding:"required"`
	MiddleName        string             `bson:"middleName" json:"middleName"`
}
