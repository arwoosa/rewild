package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Users struct {
	UsersId                           primitive.ObjectID `bson:"_id,omitempty" json:"users_id"`
	UsersSource                       int                `bson:"users_source,omitempty" json:"users_source"`
	UsersSourceId                     string             `bson:"users_source_id,omitempty" json:"users_source_id"`
	UsersName                         string             `bson:"users_name,omitempty" json:"users_name"`
	UsersEmail                        string             `bson:"users_email,omitempty" json:"users_email"`
	UsersPassword                     string             `bson:"users_password,omitempty" json:"users_password"`
	UsersObject                       string             `bson:"users_object,omitempty" json:"users_object"`
	UsersAvatar                       string             `bson:"users_avatar,omitempty" json:"users_avatar"`
	UsersSettingLanguage              string             `bson:"users_setting_language,omitempty" json:"users_setting_language"`
	UsersSettingVisEvents             int                `bson:"users_setting_vis_events,omitempty" json:"users_setting_vis_events"`
	UsersSettingVisAchievementJournal int                `bson:"users_setting_vis_achievement_journal,omitempty" json:"users_setting_vis_achievement_journal"`
	UsersSettingVisCollabLog          int                `bson:"users_setting_vis_collab_log,omitempty" json:"users_setting_vis_collab_log"`
	UsersSettingVisFollow             int                `bson:"users_setting_vis_follow,omitempty" json:"users_setting_vis_follow"`
	UsersIsSubscribed                 bool               `bson:"users_is_subscribed,omitempty" json:"users_is_subscribed"`
	UsersIsBusiness                   bool               `bson:"users_is_business,omitempty" json:"users_is_business"`
	UsersCreatedAt                    primitive.DateTime `bson:"users_created_at,omitempty" json:"users_created_at"`
}
