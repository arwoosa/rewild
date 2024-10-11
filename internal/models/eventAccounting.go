package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventAccounting struct {
	EventAccountingId            primitive.ObjectID `bson:"_id,omitempty" json:"event_accounting_id"`
	EventAccountingEvent         primitive.ObjectID `bson:"event_accounting_event,omitempty" json:"event_accounting_event,omitempty"`
	EventAccountingMessage       string             `bson:"event_accounting_message,omitempty" json:"event_accounting_message,omitempty"`
	EventAccountingAmount        float64            `bson:"event_accounting_amount,omitempty" json:"event_accounting_amount,omitempty"`
	EventAccountingPaidBy        primitive.ObjectID `bson:"event_accounting_paid_by,omitempty" json:"event_accounting_paid_by,omitempty"`
	EventAccountingCreatedBy     primitive.ObjectID `bson:"event_accounting_created_by,omitempty" json:"event_accounting_created_by,omitempty"`
	EventAccountingCreatedAt     primitive.DateTime `bson:"event_accounting_created_at,omitempty" json:"event_accounting_created_at,omitempty"`
	EventAccountingUpdatedBy     primitive.ObjectID `bson:"event_accounting_updated_by,omitempty" json:"event_accounting_updated_by,omitempty"`
	EventAccountingUpdatedAt     primitive.DateTime `bson:"event_accounting_updated_at,omitempty" json:"event_accounting_updated_at,omitempty"`
	EventAccountingCreatedByUser *UsersAgg          `bson:"event_accounting_created_by_user,omitempty" json:"event_accounting_created_by_user,omitempty"`
}
