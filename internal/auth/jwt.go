package auth

import (
	"context"
	"fmt"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var sampleSecretKey = []byte("YourSampleSecretKey")

func GenerateJWT(user models.Users) (string, error) {
	claims := jwt.MapClaims{
		"iss": "OOSA-AUTH",
		"sub": user.UsersId,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(sampleSecretKey)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func VerifyToken(tokenString string) (models.Users, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return sampleSecretKey, nil
	})

	if err != nil {
		return models.Users{}, err
	}

	if !token.Valid {
		return models.Users{}, fmt.Errorf("invalid token")
	}
	tokenId, _ := token.Claims.GetSubject()
	id, _ := primitive.ObjectIDFromHex(tokenId)

	var User models.Users
	config.DB.Collection("Users").FindOne(context.TODO(), bson.D{{Key: "_id", Value: id}}).Decode(&User)
	return User, nil
}
