package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type EventMessageBoard struct {
	EventMessageBoardId            primitive.ObjectID `bson:"_id,omitempty" json:"event_message_board_id"`
	EventMessageBoardEvent         primitive.ObjectID `bson:"event_message_board_event,omitempty" json:"event_message_board_event"`
	EventMessageBoardBaseMessage   string             `bson:"event_message_board_base_message,omitempty" json:"event_message_board_base_message,omitempty"`
	EventMessageBoardStatus        int                `bson:"event_message_board_status,omitempty" json:"event_message_board_status,omitempty"`
	EventMessageBoardCategory      string             `bson:"event_message_board_category,omitempty" json:"event_message_board_category,omitempty"`
	EventMessageBoardAnnouncement  string             `bson:"event_message_board_announcement,omitempty" json:"event_message_board_announcement,omitempty"`
	EventMessageBoardMessageId     primitive.ObjectID `bson:"event_message_board_message_id,omitempty" json:"event_message_board_message_id"`
	EventMessageBoardCreatedBy     primitive.ObjectID `bson:"event_message_board_created_by,omitempty" json:"event_message_board_created_by"`
	EventMessageBoardCreatedAt     primitive.DateTime `bson:"event_message_board_created_at,omitempty" json:"event_message_board_created_at"`
	EventMessageBoardIsPinned      *int               `bson:"event_message_board_is_pinned,omitempty" json:"event_message_board_is_pinned"`
	EventMessageBoardCreatedByUser *UsersAgg          `bson:"event_message_board_created_by_user,omitempty" json:"event_message_board_created_by_user,omitempty"`
}
type EventAnnouncement struct {
	EventMessageBoardId            primitive.ObjectID `bson:"_id,omitempty" json:"event_announcement_id"`
	EventMessageBoardEvent         primitive.ObjectID `bson:"event_message_board_event,omitempty" json:"event_announcement_event"`
	EventMessageBoardBaseMessage   string             `bson:"event_message_board_base_message,omitempty" json:"event_announcement_base_message,omitempty"`
	EventMessageBoardStatus        int                `bson:"event_message_board_status,omitempty" json:"event_announcement_status,omitempty"`
	EventMessageBoardCategory      string             `bson:"event_message_board_category,omitempty" json:"event_announcement_category,omitempty"`
	EventMessageBoardAnnouncement  string             `bson:"event_message_board_announcement,omitempty" json:"event_announcement_announcement,omitempty"`
	EventMessageBoardMessageId     primitive.ObjectID `bson:"event_message_board_message_id,omitempty" json:"event_announcement_message_id"`
	EventMessageBoardCreatedBy     primitive.ObjectID `bson:"event_message_board_created_by,omitempty" json:"event_announcement_created_by"`
	EventMessageBoardCreatedAt     primitive.DateTime `bson:"event_message_board_created_at,omitempty" json:"event_announcement_created_at"`
	EventMessageBoardIsPinned      *int               `bson:"event_message_board_is_pinned,omitempty" json:"event_announcement_is_pinned"`
	EventMessageBoardCreatedByUser *UsersAgg          `bson:"event_message_board_created_by_user,omitempty" json:"event_announcement_created_by_user,omitempty"`
}
