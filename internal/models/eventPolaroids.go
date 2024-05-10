package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventPolaroids struct {
	EventPolaroidsId        primitive.ObjectID   `bson:"_id,omitempty" json:"event_polaroids_id"`
	EventPolaroidsEvent     primitive.ObjectID   `bson:"event_polaroids_event,omitempty" json:"event_polaroids_event"`
	EventPolaroidsUrl       string               `bson:"event_polaroids_url,omitempty" json:"event_polaroids_url"`
	EventPolaroidsLat       primitive.Decimal128 `bson:"event_polaroids_lat,omitempty" json:"event_polaroids_lat"`
	EventPolaroidsLng       primitive.Decimal128 `bson:"event_polaroids_lng,omitempty" json:"event_polaroids_lng"`
	EventPolaroidsCreatedBy primitive.ObjectID   `bson:"event_polaroids_created_by,omitempty" json:"event_polaroids_created_by"`
	EventPolaroidsCreatedAt primitive.DateTime   `bson:"event_polaroids_created_at,omitempty" json:"event_polaroids_created_at"`
}
