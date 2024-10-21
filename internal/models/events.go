package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Events struct {
	EventsId                         primitive.ObjectID   `bson:"_id,omitempty" json:"events_id"`
	EventsDate                       primitive.DateTime   `bson:"events_date,omitempty" json:"events_date"`
	EventsDateEnd                    primitive.DateTime   `bson:"events_date_end,omitempty" json:"events_date_end"`
	EventsDeadline                   primitive.DateTime   `bson:"events_deadline,omitempty" json:"events_deadline"`
	EventsName                       string               `bson:"events_name,omitempty" json:"events_name"`
	EventsRewilding                  primitive.ObjectID   `bson:"events_rewilding,omitempty" json:"events_rewilding"`
	EventsRewildingAchievementType   string               `bson:"events_rewilding_achievement_type,omitempty" json:"-"`
	EventsRewildingAchievementTypeID primitive.ObjectID   `bson:"events_rewilding_achievement_type_id,omitempty" json:"-"`
	EventsPlace                      string               `bson:"events_place,omitempty" json:"events_place"`
	EventsCityId                     int                  `bson:"events_city_id,omitempty" json:"events_city_id"`
	EventsType                       primitive.ObjectID   `bson:"events_type,omitempty" json:"events_type,omitempty"`
	EventsInvitationTemplate         string               `bson:"events_invitation_template,omitempty" json:"events_invitation_template"`
	EventsInvitationMessage          string               `bson:"events_invitation_message,omitempty" json:"events_invitation_message"`
	EventsParticipantLimit           *int                 `bson:"events_participant_limit,omitempty" json:"events_participant_limit"`
	EventsPaymentRequired            int                  `bson:"events_payment_required,omitempty" json:"events_payment_required"`
	EventsPaymentFee                 *float64             `bson:"events_payment_fee,omitempty" json:"events_payment_fee"`
	EventsRequiresApproval           *int                 `bson:"events_requires_approval,omitempty" json:"events_requires_approval"`
	EventsQuestionnaireLink          string               `bson:"events_questionnaire_link,omitempty" json:"events_questionnaire_link"`
	EventsLat                        float64              `bson:"events_lat,omitempty" json:"events_lat"`
	EventsLng                        float64              `bson:"events_lng,omitempty" json:"events_lng"`
	EventsCountryCode                string               `bson:"events_country_code,omitempty" json:"events_country_code"`
	EventsMeetingPointLat            float64              `bson:"events_meeting_point_lat,omitempty" json:"events_meeting_point_lat"`
	EventsMeetingPointLng            float64              `bson:"events_meeting_point_lng,omitempty" json:"events_meeting_point_lng"`
	EventsMeetingPointName           string               `bson:"events_meeting_point_name,omitempty" json:"events_meeting_point_name"`
	EventsStatisticTime              float64              `bson:"events_statistic_time,omitempty" json:"events_statistic_time"`
	EventsStatisticDistance          float64              `bson:"events_statistic_distance,omitempty" json:"events_statistic_distance"`
	EventsStatisticMemberCount       int                  `bson:"events_statistic_member_count,omitempty" json:"events_statistic_member_count"`
	EventsPhoto                      string               `bson:"events_photo,omitempty" json:"events_photo"`
	EventsDeleted                    *int                 `bson:"events_deleted,omitempty" json:"events_deleted,omitempty"`
	EventsDeletedAt                  primitive.DateTime   `bson:"events_deleted_at,omitempty" json:"events_deleted_at,omitempty"`
	EventsCreatedBy                  primitive.ObjectID   `bson:"events_created_by,omitempty" json:"events_created_by"`
	EventsCreatedAt                  primitive.DateTime   `bson:"events_created_at,omitempty" json:"events_created_at"`
	EventsUpdatedBy                  primitive.ObjectID   `bson:"events_updated_by,omitempty" json:"events_updated_by"`
	EventsUpdatedAt                  primitive.DateTime   `bson:"events_updated_at,omitempty" json:"events_updated_at"`
	EventsParticipants               *EventParticipantObj `bson:"events_participants,omitempty" json:"events_participants,omitempty"`
	EventsCreatedByUser              *UsersAgg            `bson:"events_created_by_user,omitempty" json:"events_created_by_user,omitempty"`
	EventsRewildingDetail            *RewildingDetail     `bson:"events_rewilding_detail,omitempty" json:"events_rewilding_detail,omitempty"`
}

type EventParticipantObj struct {
	LatestTreeUser *[]UsersAgg `json:"latest_tree_user"`
	RemainNumber   int         `json:"remain_number"`
}

type EventsCountryCount struct {
	EventsCountryCode  string `bson:"_id,omitempty" json:"events_country_code"`
	EventsCountryCount int    `bson:"events_country_count,omitempty" json:"events_country_count"`
}
