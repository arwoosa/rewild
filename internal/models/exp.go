package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Exp struct {
	ExpId        primitive.ObjectID `bson:"_id,omitempty" json:"user_badges_id"`
	ExpSource    string             `bson:"exp_source,omitempty" json:"exp_source"`
	ExpUser      primitive.ObjectID `bson:"exp_user,omitempty" json:"exp_user"`
	ExpPoints    int                `bson:"exp_points,omitempty" json:"exp_points"`
	ExpRewilding primitive.ObjectID `bson:"exp_rewilding,omitempty" json:"exp_rewilding"`
	ExpCreatedAt primitive.DateTime `bson:"exp_created_at,omitempty" json:"exp_created_at"`
}
