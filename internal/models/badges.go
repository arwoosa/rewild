package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Badges struct {
	BadgesId            primitive.ObjectID `bson:"_id,omitempty" json:"badges_id"`
	BadgesCode          string             `bson:"badges_code,omitempty" json:"badges_code,omitempty"`
	BadgesName          string             `bson:"badges_name_zh,omitempty" json:"badges_name,omitempty"`
	BadgesSpecification string             `bson:"badges_specification,omitempty" json:"badges_specification,omitempty"`
	BadgesCondition     string             `bson:"badges_condition,omitempty" json:"badges_condition,omitempty"`
	BadgesCategory      string             `bson:"badges_category,omitempty" json:"badges_category,omitempty"`
	BadgesUrl           string             `bson:"badges_url,omitempty" json:"badges_url,omitempty"`
	BadgesIsOnce        bool               `bson:"badges_is_once,omitempty" json:"badges_is_once,omitempty"`
}
