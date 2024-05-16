package repository

import (
	"bytes"
	"context"
	"encoding/xml"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"
	"oosa_rewild/internal/models"
	"oosa_rewild/pkg/flickr"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	flickrEndpoint            = "https://www.flickr.com/services/rest/"
	flickrTokenEndpoint       = "https://www.flickr.com/services/oauth/request_token"
	flickrAccessTokenEndpoint = "https://www.flickr.com/services/oauth/access_token"
	flickrUploadEndpoint      = "https://up.flickr.com/services/upload/"
)

type FlickrRepository struct{}

func (r FlickrRepository) Retrieve(c *gin.Context) {
	c.JSON(200, "")
}

func (r FlickrRepository) Read(c *gin.Context) {
	photoId := c.Param("id")

	flickrReq := &flickr.FlickrRequest{
		ApiKey:    config.APP.FlickrApiKey,
		SecretKey: config.APP.FlickrSecret,
		Method:    "GET",
		Url:       flickrEndpoint,
		Args: map[string]string{
			"api_key":  config.APP.FlickrApiKey,
			"method":   "flickr.photos.getInfo",
			"photo_id": photoId,
			"format":   "json",
		},
	}

	flickrReq.Do()

	req, _ := http.NewRequest("GET", flickrReq.RequestUrl, nil)
	res, error := http.DefaultClient.Do(req)

	if error != nil {
		return
	}

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	//var response ZwApiResponse
	//json.Unmarshal(body, &response)

	c.JSON(200, string(body))
}

func (r FlickrRepository) Oauth(c *gin.Context) {
	currentTimestamp := time.Now()
	flickrReq := &flickr.FlickrRequest{
		ApiKey:    config.APP.FlickrApiKey,
		SecretKey: config.APP.FlickrSecret,
		Method:    "GET",
		Url:       flickrTokenEndpoint,
		Args: map[string]string{
			"oauth_nonce":            currentTimestamp.Format("20060102150405"),
			"oauth_timestamp":        strconv.Itoa(int(currentTimestamp.Unix())),
			"oauth_signature_method": "HMAC-SHA1",
			"oauth_version":          "1.0",
			"oauth_consumer_key":     config.APP.FlickrApiKey,
			"oauth_callback":         url.QueryEscape("http://127.0.0.1:6722/flickr/oauth/callback"),
		},
	}

	flickrReq.Sign()
	flickrReq.Do()

	req, _ := http.NewRequest(flickrReq.Method, flickrReq.RequestUrl, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	values, _ := url.ParseQuery(string(body))
	if res.StatusCode == http.StatusOK {
		insert := models.Flickr{
			FlickrOauthToken:       values["oauth_token"][0],
			FlickrOauthTokenSecret: values["oauth_token_secret"][0],
		}
		config.DB.Collection("Flickr").InsertOne(context.TODO(), insert)
		c.JSON(res.StatusCode, "https://www.flickr.com/services/oauth/authorize?oauth_token="+insert.FlickrOauthToken)
		return
	}

	c.JSON(res.StatusCode, values["oauth_problem"][0])
}

func (r FlickrRepository) OauthCb(c *gin.Context) {
	currentTimestamp := time.Now()
	oauthToken := c.Query("oauth_token")
	oauthVerifier := c.Query("oauth_verifier")

	var Flickr models.Flickr
	filter := bson.D{{Key: "flickr_oauth_token", Value: oauthToken}}
	err := config.DB.Collection("Flickr").FindOne(context.TODO(), filter).Decode(&Flickr)
	if err != nil {
		helpers.ResultEmpty(c, err)
	}

	flickrReq := &flickr.FlickrRequest{
		ApiKey:      config.APP.FlickrApiKey,
		TokenSecret: Flickr.FlickrOauthTokenSecret,
		SecretKey:   config.APP.FlickrSecret,
		Method:      "GET",
		Url:         flickrAccessTokenEndpoint,
		Args: map[string]string{
			"oauth_nonce":            currentTimestamp.Format("20060102150405"),
			"oauth_timestamp":        strconv.Itoa(int(currentTimestamp.Unix())),
			"oauth_verifier":         oauthVerifier,
			"oauth_consumer_key":     config.APP.FlickrApiKey,
			"oauth_signature_method": "HMAC-SHA1",
			"oauth_version":          "1.0",
			"oauth_token":            oauthToken,
		},
	}

	flickrReq.Sign()
	flickrReq.Do()

	req, _ := http.NewRequest(flickrReq.Method, flickrReq.RequestUrl, nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	values, _ := url.ParseQuery(string(body))
	if res.StatusCode == http.StatusOK {
		filter := bson.D{{Key: "_id", Value: Flickr.FlickrId}}
		Flickr.FlickrFullName = values["fullname"][0]
		Flickr.FlickrOauthToken = values["oauth_token"][0]
		Flickr.FlickrOauthTokenSecret = values["oauth_token_secret"][0]
		Flickr.FlickrUserNsID = values["user_nsid"][0]
		Flickr.FlickrUserName = values["username"][0]
		upd := bson.D{{Key: "$set", Value: Flickr}}
		config.DB.Collection("Flickr").UpdateOne(context.TODO(), filter, upd)
		c.JSON(res.StatusCode, Flickr)
		return
	}

	c.JSON(res.StatusCode, values["oauth_problem"][0])
}

type FlickrUploadDataResponse struct {
	Type    string `xml:"stat,attr"`
	PhotoId string `xml:"photoid"`
}

func (r FlickrRepository) Upload(c *gin.Context) {
	var (
		currentTimestamp = time.Now()
		buffer           = new(bytes.Buffer)
		writer           = multipart.NewWriter(buffer)
	)

	var Flickr models.Flickr
	filter := bson.D{{Key: "flickr_user_name", Value: config.APP.FlickrUserName}}
	err := config.DB.Collection("Flickr").FindOne(context.TODO(), filter).Decode(&Flickr)
	if err != nil {
		helpers.ResultEmpty(c, err)
		return
	}

	photoDescription := c.Request.FormValue("photo_description")
	flickrReq := &flickr.FlickrRequest{
		ApiKey:      config.APP.FlickrApiKey,
		TokenSecret: Flickr.FlickrOauthTokenSecret,
		SecretKey:   config.APP.FlickrSecret,
		Method:      "POST",
		Url:         flickrUploadEndpoint,
		Args: map[string]string{
			"oauth_nonce":            currentTimestamp.Format("20060102150405"),
			"oauth_timestamp":        strconv.Itoa(int(currentTimestamp.Unix())),
			"oauth_signature_method": "HMAC-SHA1",
			"oauth_consumer_key":     config.APP.FlickrApiKey,
			"oauth_token":            Flickr.FlickrOauthToken,
			"title":                  photoDescription,
			"hidden":                 "2",
			"is_public":              "1",
			"is_friend":              "1",
			"is_family":              "1",
		},
	}

	signature := flickrReq.Sign()

	// The file cannot be received.
	file, fileErr := c.FormFile("photo")
	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	for key, val := range flickrReq.Args {
		writer.WriteField(key, val)
	}

	writer.WriteField("oauth_signature", signature)

	part, err := writer.CreateFormFile("photo", file.Filename)

	if err != nil {
		log.Println(err)
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
	part.Write(b)

	writer.Close()

	// Create Request
	client := &http.Client{}
	req, err := http.NewRequest("POST", flickrUploadEndpoint, buffer)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Create Response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	xmlData := []byte(string(body))
	var data FlickrUploadDataResponse
	flickerXmlErr := xml.Unmarshal(xmlData, &data)
	if flickerXmlErr != nil {
		return
	}
	c.JSON(200, data)
}
