package repository

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"oosa_rewild/internal/config"

	"github.com/gin-gonic/gin"
)

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
	endpoint := "https://api.cloudflare.com/client/v4/accounts/" + config.APP.ClourdlareImageAccountId + "/images/v1"
	var (
		buffer = new(bytes.Buffer)
		writer = multipart.NewWriter(buffer)
	)

	// The file cannot be received.
	file, fileErr := c.FormFile("photo")
	if fileErr != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"message": "No file is received",
		})
		return
	}

	part, err := writer.CreateFormFile("file", file.Filename)

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
	req, err := http.NewRequest("POST", endpoint, buffer)
	if err != nil {
		log.Fatal(err)
		return
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
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

	if responseRaw["success"] != nil && responseRaw["success"] == true {
		fmt.Println("SUCCESS!", responseRaw["result"])
	}

	c.JSON(200, responseRaw)
}
