package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type Flickr struct {
	FlickrId               primitive.ObjectID `bson:"_id,omitempty" json:"flickr_id"`
	FlickrOauthToken       string             `bson:"flickr_oauth_token,omitempty" json:"flickr_oauth_token,omitempty"`
	FlickrOauthTokenSecret string             `bson:"flickr_oauth_token_secret,omitempty" json:"flickr_oauth_token_secret,omitempty"`
	FlickrFullName         string             `bson:"flickr_full_name,omitempty" json:"flickr_full_name,omitempty"`
	FlickrUserNsID         string             `bson:"flickr_user_ns_id,omitempty" json:"flickr_user_ns_id,omitempty"`
	FlickrUserName         string             `bson:"flickr_user_name,omitempty" json:"flickr_user_name,omitempty"`
}
