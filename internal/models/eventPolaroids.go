package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventPolaroids struct {
	EventPolaroidsId                  primitive.ObjectID `bson:"_id,omitempty" json:"event_polaroids_id"`
	EventPolaroidsEvent               primitive.ObjectID `bson:"event_polaroids_event,omitempty" json:"event_polaroids_event"`
	EventPolaroidsUrl                 string             `bson:"event_polaroids_url,omitempty" json:"event_polaroids_url"`
	EventPolaroidsLat                 float64            `bson:"event_polaroids_lat,omitempty" json:"event_polaroids_lat"`
	EventPolaroidsLng                 float64            `bson:"event_polaroids_lng,omitempty" json:"event_polaroids_lng"`
	EventPolaroidsRadiusFromEvent     float64            `bson:"event_polaroids_radius_from_event,omitempty" json:"event_polaroids_radius_from_event"`
	EventPolaroidsAchievementEligible *bool              `bson:"event_polaroids_achievement_eligible,omitempty" json:"event_polaroids_achievement_eligible"`
	EventPolaroidsMessage             string             `bson:"event_polaroids_message,omitempty" json:"event_polaroids_message"`
	EventPolaroidsTag                 string             `bson:"event_polaroids_tag,omitempty" json:"event_polaroids_tag"`
	EventPolaroidsIsEventPeriod       *bool              `bson:"event_polaroids_is_event_period,omitempty" json:"event_polaroids_is_event_period"`
	EventPolaroidsStarType            int                `bson:"event_polaroids_star_type,omitempty" json:"event_polaroids_star_type"`
	EventPolaroidsCreatedBy           primitive.ObjectID `bson:"event_polaroids_created_by,omitempty" json:"event_polaroids_created_by"`
	EventPolaroidsCreatedAt           primitive.DateTime `bson:"event_polaroids_created_at,omitempty" json:"event_polaroids_created_at"`
	EventPolaroidsCreatedByUser       *UsersAgg          `bson:"event_polaroids_created_by_user,omitempty" json:"event_polaroids_created_by_user,omitempty"`
}
