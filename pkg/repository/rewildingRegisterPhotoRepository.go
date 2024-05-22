package repository

import (
	"context"
	"io"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RewildingRegisterPhotoRepository struct{}

func (r RewildingRegisterPhotoRepository) Retrieve(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	var Rewilding models.Rewilding

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil || Rewilding.RewildingPhotos == nil {
		helpers.ResultEmpty(c, err)
		return
	}

	c.JSON(200, Rewilding.RewildingPhotos)
}

func (r RewildingRegisterPhotoRepository) Create(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	var Rewilding models.Rewilding

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	file, fileErr := c.FormFile("rewilding_photo")
	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	uploadedFile, err := file.Open()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to open file",
		})
		return
	}

	b, _ := io.ReadAll(uploadedFile)
	mimeType := mimetype.Detect(b)

	switch mimeType.String() {
	case "image/heic_":
	case "image/jpeg":
	case "image/png":

	default:
		c.JSON(http.StatusBadRequest, "Mime: "+mimeType.String()+" not supported")
		return
	}

	cloudflare := CloudflareRepository{}
	cloudflareResponse, postErr := cloudflare.Post(c, file)
	if postErr != nil {
		helpers.ResponseBadRequestError(c, postErr.Error())
		return
	}
	fileName := cloudflare.ImageDelivery(cloudflareResponse.Result.Id, "public")

	RwPhoto := models.RewildingPhotos{
		RewildingPhotosID:   primitive.NewObjectID(),
		RewildingPhotosPath: fileName,
		//RewildingPhotosData: b,
	}

	Rewilding.RewildingPhotos = append(Rewilding.RewildingPhotos, RwPhoto)
	upd := bson.D{{Key: "$set", Value: Rewilding}}
	_, updateErr := config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filter, upd)

	if updateErr != nil {
		c.JSON(http.StatusBadRequest, updateErr.Error())
		return
	}

	c.JSON(200, Rewilding)
}

func (r RewildingRegisterPhotoRepository) Read(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	photosId, _ := primitive.ObjectIDFromHex(c.Param("photosId"))
	var Rewilding models.Rewilding

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
		{Key: "rewilding_photos._id", Value: photosId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	idx := 0
	for p, v := range Rewilding.RewildingPhotos {
		if v.RewildingPhotosID == photosId {
			idx = p
		}
	}

	c.Writer.Write(Rewilding.RewildingPhotos[idx].RewildingPhotosData)
}

func (r RewildingRegisterPhotoRepository) Delete(c *gin.Context) {
	id, _ := primitive.ObjectIDFromHex(c.Param("id"))
	userDetail := helpers.GetAuthUser(c)
	photosId, _ := primitive.ObjectIDFromHex(c.Param("photosId"))
	var Rewilding models.Rewilding

	filter := bson.D{
		{Key: "_id", Value: id},
		{Key: "rewilding_created_by", Value: userDetail.UsersId},
		{Key: "rewilding_photos._id", Value: photosId},
	}
	err := config.DB.Collection("Rewilding").FindOne(context.TODO(), filter).Decode(&Rewilding)

	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	idx := 0
	for p, v := range Rewilding.RewildingPhotos {
		if v.RewildingPhotosID == photosId {
			idx = p
		}
	}

	Rewilding.RewildingPhotos = append(Rewilding.RewildingPhotos[:idx], Rewilding.RewildingPhotos[idx+1:]...)
	upd := bson.D{{Key: "$set", Value: Rewilding}}
	_, updateErr := config.DB.Collection("Rewilding").UpdateOne(context.TODO(), filter, upd)

	if updateErr != nil {
		c.JSON(http.StatusBadRequest, updateErr.Error())
		return
	}

	c.JSON(200, "Image deleted")
}
