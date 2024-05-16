package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Events struct {
	EventsId               primitive.ObjectID   `bson:"_id,omitempty" json:"events_id"`
	EventsDate             primitive.DateTime   `bson:"events_date,omitempty" json:"events_date"`
	EventsName             string               `bson:"events_name,omitempty" json:"events_name"`
	EventsRewilding        primitive.ObjectID   `bson:"events_rewilding,omitempty" json:"events_rewilding"`
	EventsPlace            string               `bson:"events_place,omitempty" json:"events_place"`
	EventsType             string               `bson:"events_type,omitempty" json:"events_type"`
	EventsPaymentRequired  int                  `bson:"events_payment_required,omitempty" json:"events_payment_required"`
	EventsPaymentFee       float64              `bson:"events_payment_fee,omitempty" json:"events_payment_fee"`
	EventsRequiresApproval int                  `bson:"events_requires_approval,omitempty" json:"events_requires_approval"`
	EventsLat              primitive.Decimal128 `bson:"events_lat,omitempty" json:"events_lat"`
	EventsLng              primitive.Decimal128 `bson:"events_lng,omitempty" json:"events_lng"`
	EventsCreatedBy        primitive.ObjectID   `bson:"events_created_by,omitempty" json:"events_created_by"`
	EventsCreatedAt        primitive.DateTime   `bson:"events_created_at,omitempty" json:"events_created_at"`
	EventsUpdatedBy        primitive.ObjectID   `bson:"events_updated_by,omitempty" json:"events_updated_by"`
	EventsUpdatedAt        primitive.DateTime   `bson:"events_updated_at,omitempty" json:"events_updated_at"`
}
