package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TransactionLog struct {
	UserID      primitive.ObjectID `bson:"userId"`
	Action      string             `bson:"action"`
	CreatedDate time.Time          `bson:"createdDate"`
	CreatedByID string             `bson:"createdById"`
	UpdatedDate time.Time          `bson:"updatedDate"`
	UpdatedByID string             `bson:"updatedById"`
}
