package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventAlbumLink struct {
	EventAlbumLinkId            primitive.ObjectID `bson:"_id,omitempty" json:"event_album_link_id"`
	EventAlbumLinkEvent         primitive.ObjectID `bson:"event_album_link_event,omitempty" json:"event_album_link_event"`
	EventAlbumLinkAlbumUrl      string             `bson:"event_album_link_album_url,omitempty" json:"event_album_link_album_url"`
	EventAlbumLinkVisibility    *int64             `bson:"event_album_link_visibility,omitempty" json:"event_album_link_visibility"`
	EventAlbumLinkCreatedBy     primitive.ObjectID `bson:"event_album_link_created_by,omitempty" json:"event_album_link_created_by"`
	EventAlbumLinkCreatedAt     primitive.DateTime `bson:"event_album_link_created_at,omitempty" json:"event_album_link_created_at"`
	EventAlbumLinkCreatedByUser *UsersAgg          `bson:"event_album_link_created_by_user,omitempty" json:"event_album_link_created_by_user,omitempty"`
}
