package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	UsersId                               primitive.ObjectID `bson:"_id,omitempty" json:"users_id"`
	UsersSource                           int                `bson:"users_source,omitempty" json:"users_source"`
	UsersSourceId                         string             `bson:"users_source_id,omitempty" json:"users_source_id"`
	UsersBindingFacebook                  string             `bson:"users_binding_facebook,omitempty" json:"users_binding_facebook"`
	UsersName                             string             `bson:"users_name,omitempty" json:"users_name"`
	UsersUsername                         string             `bson:"users_username,omitempty" json:"users_username"`
	UsersEmail                            string             `bson:"users_email,omitempty" json:"users_email"`
	UsersPassword                         string             `bson:"users_password,omitempty" json:"-"`
	UsersObject                           string             `bson:"users_object,omitempty" json:"users_object"`
	UsersAvatar                           string             `bson:"users_avatar,omitempty" json:"users_avatar"`
	UsersSettingLanguage                  string             `bson:"users_setting_language,omitempty" json:"users_setting_language"`
	UsersSettingIsVisibleFriends          int                `bson:"users_setting_is_visible_friends,omitempty" json:"users_setting_is_visible_friends"`
	UsersSettingIsVisibleStatistics       int                `bson:"users_setting_is_visible_statistics,omitempty" json:"users_setting_is_visible_statistics"`
	UsersSettingVisibilityActivitySummary int                `bson:"users_setting_visibility_activity_summary,omitempty" json:"users_setting_visibility_activity_summary"`
	UsersSettingFriendAutoAdd             int                `bson:"users_setting_friend_auto_add,omitempty" json:"users_setting_friend_auto_add"`
	UsersIsSubscribed                     bool               `bson:"users_is_subscribed,omitempty" json:"users_is_subscribed"`
	UsersIsBusiness                       bool               `bson:"users_is_business,omitempty" json:"users_is_business"`
	UsersCreatedAt                        primitive.DateTime `bson:"users_created_at,omitempty" json:"users_created_at"`
	UsersEventCompleted                   int                `bson:"users_event_completed,omitempty" json:"-"`
	UsersEventScheduled                   int                `bson:"users_event_scheduled,omitempty" json:"-"`
	UsersBreathingPoints                  int                `bson:"users_breathing_points,omitempty" json:"users_breathing_points"`
	UsersFollowingCount                   int                `bson:"users_following_count,omitempty" json:"users_following_count"`
	UsersFollowerCount                    int                `bson:"users_follower_count,omitempty" json:"users_follower_count"`
}

type UsersAgg struct {
	UsersId     primitive.ObjectID `bson:"_id,omitempty" json:"user_id"`
	UsersName   string             `bson:"users_name,omitempty" json:"user_name"`
	UsersEmail  string             `bson:"users_email,omitempty" json:"user_email"`
	UsersAvatar string             `bson:"users_avatar,omitempty" json:"user_avatar"`
}
