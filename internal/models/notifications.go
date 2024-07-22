package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Notifications struct {
	NotificationsId         primitive.ObjectID `bson:"_id,omitempty" json:"notifications_id"`
	NotificationsCode       string             `bson:"notifications_code,omitempty" json:"notifications_code,omitempty"`
	NotificationsType       string             `bson:"notifications_type,omitempty" json:"notifications_type,omitempty"`
	NotificationsUser       primitive.ObjectID `bson:"notifications_user,omitempty" json:"notifications_user,omitempty"`
	NotificationsMessage    string             `bson:"notifications_message,omitempty" json:"notifications_message,omitempty"`
	NotificationsIdentifier primitive.ObjectID `bson:"notifications_identifier,omitempty" json:"notifications_identifier,omitempty"`
	NotificationsCreatedAt  primitive.DateTime `bson:"notifications_created_at,omitempty" json:"notifications_created_at,omitempty"`
	NotificationsCreatedBy  primitive.ObjectID `bson:"notifications_created_by,omitempty" json:"notifications_created_by,omitempty"`
}
