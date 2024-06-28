package repository

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"oosa_rewild/internal/config"
	"oosa_rewild/internal/helpers"

	"github.com/gin-gonic/gin"
)

type CloudflareResponse struct {
	Errors   []string         `json:"errors"`
	Messages []string         `json:"messages"`
	Result   CloudflareResult `json:"result"`
	Success  bool             `json:"success"`
}

type CloudflareResult struct {
	Filename          string   `json:"filename"`
	Id                string   `json:"id"`
	RequireSignedUrLs bool     `json:"requireSignedURLs"`
	Uploaded          string   `json:"uploaded"`
	Variants          []string `json:"variants"`
}

type CloudflareRepository struct{}

func (r CloudflareRepository) ImageDelivery(imageId string, variantName string) string {
	endpoint := "https://imagedelivery.net/" + config.APP.ClourdlareImageAccountHash + "/" + imageId + "/" + variantName
	return endpoint
}

func (r CloudflareRepository) Read(c *gin.Context) {
	imageId := c.Param("imageId")
	url := r.ImageDelivery(imageId, "public")
	c.JSON(200, url)
}

func (r CloudflareRepository) Retrieve(c *gin.Context) {
	endpoint := "https://api.cloudflare.com/client/v4/user/tokens/verify"
	client := &http.Client{}
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+config.APP.CloudflareImageAuthToken)

	fmt.Println("Endpoint: " + endpoint)
	fmt.Println("Bearer " + config.APP.CloudflareImageAuthToken)

	// Create Response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var responseRaw map[string]interface{}
	json.Unmarshal(body, &responseRaw)

	c.JSON(200, responseRaw)
}

func (r CloudflareRepository) Upload(c *gin.Context) {
	file, fileErr := c.FormFile("photo")
	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}
	clourflareResponse, postErr := r.Post(c, file)

	// The file cannot be received.
	if postErr != nil {
		helpers.ResponseBadRequestError(c, postErr.Error())
		return
	}
	c.JSON(200, clourflareResponse.Result)
}

func (r CloudflareRepository) Post(c *gin.Context, file *multipart.FileHeader) (CloudflareResponse, error) {
	var cloudflareResponse CloudflareResponse
	endpoint := "https://api.cloudflare.com/client/v4/accounts/" + config.APP.ClourdlareImageAccountId + "/images/v1"
	var (
		buffer = new(bytes.Buffer)
		writer = multipart.NewWriter(buffer)
	)

	part, err := writer.CreateFormFile("file", file.Filename)

	if err != nil {
		log.Println(err)
		return cloudflareResponse, errors.New("can't create file")
	}

	uploadedFile, err := file.Open()

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "Unable to open file",
		})
		return cloudflareResponse, errors.New("unable to open file")
	}
	b, _ := io.ReadAll(uploadedFile)
	part.Write(b)

	writer.Close()

	// Create Request
	client := &http.Client{}
	req, err := http.NewRequest("POST", endpoint, buffer)
	if err != nil {
		log.Fatal(err)
		return cloudflareResponse, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+config.APP.CloudflareImageAuthToken)

	fmt.Println("Endpoint: " + endpoint)
	fmt.Println("Bearer " + config.APP.CloudflareImageAuthToken)

	// Create Response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
		return cloudflareResponse, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	json.Unmarshal(body, &cloudflareResponse)

	if cloudflareResponse.Success {
		fmt.Println("SUCCESS!", cloudflareResponse.Result)
	}

	return cloudflareResponse, nil
}
