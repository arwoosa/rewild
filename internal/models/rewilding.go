package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rewilding struct {
	RewildingID             primitive.ObjectID        `bson:"_id,omitempty" json:"rewilding_id"`
	RewildingCity           string                    `bson:"rewilding_city,omitempty" json:"rewilding_city,omitempty"`
	RewildingArea           string                    `bson:"rewilding_area,omitempty" json:"rewilding_area,omitempty"`
	RewildingLocation       []string                  `bson:"rewilding_location,omitempty" json:"rewilding_location,omitempty"`
	RewildingCountryCode    string                    `bson:"rewilding_country_code,omitempty" json:"rewilding_country_code,omitempty"`
	RewildingName           string                    `bson:"rewilding_name,omitempty" json:"rewilding_name"`
	RewildingRating         int                       `bson:"rewilding_rating,omitempty" json:"rewilding_rating"`
	RewildingLat            float64                   `bson:"rewilding_lat,omitempty" json:"rewilding_lat"`
	RewildingLng            float64                   `bson:"rewilding_lng,omitempty" json:"rewilding_lng"`
	RewildingPlaceId        string                    `bson:"rewilding_place_id,omitempty" json:"rewilding_place_id"`
	RewildingElevation      float64                   `bson:"rewilding_elevation,omitempty" json:"rewilding_elevation"`
	RewildingPhotos         []RewildingPhotos         `bson:"rewilding_photos,omitempty" json:"rewilding_photos"`
	RewildingReferenceLinks []RewildingReferenceLinks `bson:"rewilding_reference_links,omitempty" json:"rewilding_reference_links"`
	RewildingApplyOfficial  *bool                     `bson:"rewilding_apply_official,default:false" json:"rewilding_apply_official"`
	RewildingCreatedBy      primitive.ObjectID        `bson:"rewilding_created_by,omitempty" json:"rewilding_created_by"`
	RewildingCreatedAt      primitive.DateTime        `bson:"rewilding_created_at,omitempty" json:"rewilding_created_at"`
	RewildingDeletedBy      *primitive.ObjectID       `bson:"rewilding_deleted_by,omitempty" json:"rewilding_deleted_by,omitempty"`
	RewildingDeletedAt      *primitive.DateTime       `bson:"rewilding_deleted_at,omitempty" json:"rewilding_deleted_at,omitempty"`
	RewildingCreatedByUser  *UsersAgg                 `bson:"rewilding_created_by_user,omitempty" json:"rewilding_created_by_user,omitempty"`
}

type RewildingPhotos struct {
	RewildingPhotosID   primitive.ObjectID `bson:"_id,omitempty" json:"rewilding_photos_id"`
	RewildingPhotosPath string             `bson:"rewilding_photos_path,omitempty" json:"rewilding_photos_path,omitempty"`
	RewildingPhotosData []byte             `bson:"rewilding_photos_data,omitempty" json:"-"`
}

type RewildingDetail struct {
	RewildingID       primitive.ObjectID `bson:"_id,omitempty" json:"rewilding_id"`
	RewildingName     string             `bson:"rewilding_name,omitempty" json:"rewilding_name"`
	RewildingLocation []string           `bson:"rewilding_location,omitempty" json:"rewilding_location,omitempty"`
}
