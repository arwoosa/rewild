package helpers

import (
	"context"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/models"

	"go.mongodb.org/mongo-driver/bson"
)

func RefRewildingTypes() []models.RefRewildingTypes {
	var RefRewildingTypes []models.RefRewildingTypes
	cursor, err := config.DB.Collection("RefRewildingTypes").Find(context.TODO(), bson.D{})
	if err != nil {
		return RefRewildingTypes
	}
	cursor.All(context.TODO(), &RefRewildingTypes)
	return RefRewildingTypes
}
