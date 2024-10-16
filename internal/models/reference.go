package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RefCity struct {
	RefCityId   primitive.ObjectID `bson:"ref_city_id"`
	RefCityName string             `bson:"ref_city_name"`
}

type RefJournalTypes struct {
	RefJournalTypesId   primitive.ObjectID `bson:"ref_journal_types_id"`
	RefJournalTypesName string             `bson:"ref_journal_types_name"`
}

type RefPrivacySettings struct {
	RefPrivacySettingId          primitive.ObjectID `bson:"ref_privacy_settings_id"`
	RefPrivacySettingName        string             `bson:"ref_privacy_settings_name"`
	RefPrivacySettingDescription string             `bson:"ref_privacy_settings_description"`
}

type RefRewildingAreas struct {
	RefRewildingAreasId   primitive.ObjectID `bson:"ref_rewilding_areas_id"`
	RefRewildingAreasName string             `bson:"ref_rewilding_areas_name"`
}

type RefRewildingTypes struct {
	RefRewildingTypesId           primitive.ObjectID `bson:"_id,omitempty" json:"ref_rewilding_types_id"`
	RefRewildingTypesName         string             `bson:"ref_rewilding_types_name,omitempty" json:"ref_rewilding_types_name"`
	RefRewildingTypesNameEn       string             `bson:"ref_rewilding_types_name_en,omitempty" json:"ref_rewilding_types_name_en"`
	RefRewildingTypesGoogle       string             `bson:"ref_rewilding_types_google,omitempty" json:"ref_rewilding_types_google"`
	RefRewildingTypesKey          string             `bson:"ref_rewilding_types_key,omitempty" json:"ref_rewilding_types_key"`
	RefRewildingTypesDefaultImage string             `bson:"ref_rewilding_types_default_image,omitempty" json:"ref_rewilding_types_default_image"`
}

type RefRewildingWikiTypes struct {
	RefRewildingWikiTypesId   primitive.ObjectID `bson:"ref_rewilding_wiki_types_id"`
	RefRewildingWikiTypesName string             `bson:"ref_rewilding_wiki_types_name"`
}

type RefUsersSource struct {
	RefUsersSourceId   primitive.ObjectID `bson:"ref_users_source_id"`
	RefUsersSourceName string             `bson:"ref_users_source_name"`
}
