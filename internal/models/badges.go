package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Badges struct {
	BadgesId            primitive.ObjectID `bson:"_id,omitempty" json:"badges_id"`
	BadgesCode          string             `bson:"badges_code,omitempty" json:"badges_code,omitempty"`
	BadgesName          string             `bson:"badges_name,omitempty" json:"badges_name,omitempty"`
	BadgesSpecification string             `bson:"badges_specification,omitempty" json:"badges_specification,omitempty"`
	BadgesCondition     string             `bson:"badges_condition,omitempty" json:"badges_condition,omitempty"`
	BadgesIsOnce        bool               `bson:"badges_is_once,omitempty" json:"badges_is_once,omitempty"`
}
