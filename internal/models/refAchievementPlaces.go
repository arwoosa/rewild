package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefAchievementPlaces struct {
	RefAchievementPlacesID       primitive.ObjectID `bson:"_id,omitempty" json:"ref_achievement_places_id"`
	RefAchievementPlacesType     string             `bson:"ref_achievement_places_type,omitempty" json:"ref_achievement_places_type,omitempty"`
	RefAchievementPlacesName     string             `bson:"ref_achievement_places_name,omitempty" json:"ref_achievement_places_name,omitempty"`
	RefAchievementPlacesLat      float64            `bson:"ref_achievement_places_lat,omitempty" json:"ref_achievement_places_lat,omitempty"`
	RefAchievementPlacesLng      float64            `bson:"ref_achievement_places_lng,omitempty" json:"ref_achievement_places_lng,omitempty"`
	RefAchievementPlacesDistance float64            `bson:"ref_achievement_places_distance,omitempty" json:"ref_achievement_places_distance,omitempty"`
}
