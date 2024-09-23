package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type News struct {
	NewsId        primitive.ObjectID `bson:"_id,omitempty" json:"news_id"`
	NewsDate      primitive.DateTime `bson:"news_date,omitempty" json:"news_date"`
	NewsTitle     string             `bson:"news_title,omitempty" json:"news_title"`
	NewsContent   string             `bson:"news_content" json:"news_content"`
	NewsDeletedAt primitive.DateTime `bson:"news_deleted_at,omitempty" json:"news_deleted_at,omitempty"`
	NewsCreatedAt primitive.DateTime `bson:"news_created_at,omitempty" json:"news_created_at,omitempty"`
	NewsCreatedBy primitive.ObjectID `bson:"news_created_by,omitempty" json:"news_created_by"`
}
