package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventReferenceLinks struct {
	EventReferenceLinksId        primitive.ObjectID `bson:"_id,omitempty" json:"event_reference_links_id"`
	EventReferenceLinksEvent     primitive.ObjectID `bson:"event_reference_links_event,omitempty" json:"event_reference_links_event,omitempty"`
	EventReferenceLinksLink      string             `bson:"event_reference_links_link,omitempty" json:"event_reference_links_link,omitempty"`
	EventReferenceLinksCreatedAt primitive.DateTime `bson:"event_reference_links_created_at,omitempty" json:"event_reference_links_created_at,omitempty"`
	EventReferenceLinksCreatedBy primitive.ObjectID `bson:"event_reference_links_created_by,omitempty" json:"event_reference_links_created_by,omitempty"`
}
