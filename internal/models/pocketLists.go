package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PocketLists struct {
	PocketListsId        primitive.ObjectID `bson:"_id,omitempty" json:"pocket_lists_id"`
	PocketListsUser      primitive.ObjectID `bson:"pocket_lists_user,omitempty" json:"pocket_lists_user"`
	PocketListsName      string             `bson:"pocket_lists_name,omitempty" json:"pocket_lists_name"`
	PocketListsCount     int                `bson:"pocket_lists_count" json:"pocket_lists_count"`
	PocketListsDeletedAt primitive.DateTime `bson:"pocket_lists_deleted_at,omitempty" json:"pocket_lists_deleted_at,omitempty"`
	PocketListsCreatedAt primitive.DateTime `bson:"pocket_lists_created_at,omitempty" json:"pocket_lists_created_at,omitempty"`
}
