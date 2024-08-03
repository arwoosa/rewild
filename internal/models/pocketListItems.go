package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type PocketListItems struct {
	PocketListItemsId              primitive.ObjectID `bson:"_id,omitempty" json:"pocket_list_items_id"`
	PocketListItemsMst             primitive.ObjectID `bson:"pocket_list_items_mst" json:"pocket_list_items_mst,omitempty"`
	PocketListItemsName            string             `bson:"pocket_list_items_name" json:"pocket_list_items_name,omitempty"`
	PocketListItemsEvent           primitive.ObjectID `bson:"pocket_list_items_event,omitempty" json:"pocket_list_items_event,omitempty,omitzero"`
	PocketListItemsRewilding       primitive.ObjectID `bson:"pocket_list_items_rewilding,omitempty" json:"pocket_list_items_rewilding,omitempty"`
	PocketListItemsDeletedAt       primitive.DateTime `bson:"pocket_list_items_deleted_at,omitempty" json:"pocket_list_items_deleted_at,omitempty"`
	PocketListItemsCreatedAt       primitive.DateTime `bson:"pocket_list_items_created_at,omitempty" json:"pocket_list_items_created_at,omitempty"`
	PocketListItemsRewildingDetail *Rewilding         `bson:"pocket_list_items_rewilding_detail,omitempty" json:"pocket_list_items_rewilding_detail,omitempty"`
}
