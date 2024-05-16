package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventSchedules struct {
	EventSchedulesId          primitive.ObjectID `bson:"_id,omitempty" json:"event_schedules_id"`
	EventSchedulesEvent       primitive.ObjectID `bson:"event_schedules_event,omitempty" json:"event_schedules_event,omitempty"`
	EventSchedulesDatetime    primitive.DateTime `bson:"event_schedules_datetime,omitempty" json:"event_schedules_datetime,omitempty"`
	EventSchedulesDescription string             `bson:"event_schedules_description,omitempty" json:"event_schedules_description,omitempty"`
	EventSchedulesCreatedAt   primitive.DateTime `bson:"event_schedules_created_at,omitempty" json:"event_schedules_created_at,omitempty"`
	EventSchedulesCreatedBy   primitive.ObjectID `bson:"event_schedules_created_by,omitempty" json:"event_schedules_created_by,omitempty"`
}
