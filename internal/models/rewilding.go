package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Rewilding struct {
	RewildingID            primitive.ObjectID   `bson:"_id,omitempty" json:"rewilding_id"`
	RewildingType          primitive.ObjectID   `bson:"rewilding_type,omitempty" json:"rewilding_type"`
	RewildingTypeData      RefRewildingTypes    `bson:"rewilding_type_data,omitempty" json:"rewilding_type_data"`
	RewildingCity          string               `bson:"rewilding_city,omitempty" json:"rewilding_city"`
	RewildingArea          string               `bson:"rewilding_area,omitempty" json:"rewilding_area"`
	RewildingName          string               `bson:"rewilding_name,omitempty" json:"rewilding_name"`
	RewildingRating        int                  `bson:"rewilding_rating,omitempty" json:"rewilding_rating"`
	RewildingLat           primitive.Decimal128 `bson:"rewilding_lat,omitempty" json:"rewilding_lat"`
	RewildingLng           primitive.Decimal128 `bson:"rewilding_lng,omitempty" json:"rewilding_lng"`
	RewildingPlaceId       string               `bson:"rewilding_place_id,omitempty" json:"rewilding_place_id"`
	RewildingElevation     float64              `bson:"rewilding_elevation,omitempty" json:"rewilding_elevation"`
	RewildingPhotos        []RewildingPhotos    `bson:"rewilding_photos,omitempty" json:"rewilding_photos"`
	RewildingApplyOfficial *bool                `bson:"rewilding_apply_official,omitempty" json:"rewilding_apply_official"`
	RewildingCreatedBy     primitive.ObjectID   `bson:"rewilding_created_by,omitempty" json:"rewilding_created_by"`
	RewildingCreatedAt     primitive.DateTime   `bson:"rewilding_created_at,omitempty" json:"rewilding_created_at"`
	RewildingCreatedByUser *UsersAgg            `bson:"rewilding_created_by_user,omitempty" json:"rewilding_created_by_user,omitempty"`
}

type RewildingPhotos struct {
	RewildingPhotosID   primitive.ObjectID `bson:"_id,omitempty" json:"rewilding_photos_id"`
	RewildingPhotosPath string             `bson:"rewilding_photos_path,omitempty" json:"rewilding_photos_path,omitempty"`
	RewildingPhotosData []byte             `bson:"rewilding_photos_data,omitempty" json:"-"`
}
