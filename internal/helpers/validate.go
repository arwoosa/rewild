package helpers

import (
	"errors"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gabriel-vasile/mimetype"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(c *gin.Context, arr interface{}) error {
	errorList := []string{}
	if err := c.Bind(&arr); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "message": "Validation error! JSON does not match", "data": errorList, "validation": "oosa_api"})
		return err
	}

	validate := validator.New()

	err := validate.Struct(arr)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			//translatedErr := fmt.Errorf(e.Translate(trans))
			//errs = append(errs, translatedErr)
			errorList = append(errorList, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "message": "Validation error! Please check your inputs!", "data": errorList, "validation": "oosa_api"})
		return err
	}

	return nil
}

func ValidateForm(c *gin.Context, payload interface{}) error {
	errorList := []string{}
	if err := c.Bind(payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": -2, "message": "Validation error! Form does not match", "data": errorList, "validation": "oosa_api"})
		return err
	}

	validate := validator.New()

	err := validate.Struct(payload)

	if err != nil {
		validationErrors := err.(validator.ValidationErrors)
		for _, e := range validationErrors {
			//translatedErr := fmt.Errorf(e.Translate(trans))
			//errs = append(errs, translatedErr)
			errorList = append(errorList, e.Error())
		}
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation error! Please check your inputs!", "data": errorList})
		return err
	}

	return nil
}

func ValidateError(c *gin.Context, err error) error {
	errorList := []string{}
	validationErrors := err.(validator.ValidationErrors)
	for _, e := range validationErrors {
		//translatedErr := fmt.Errorf(e.Translate(trans))
		//errs = append(errs, translatedErr)
		errorList = append(errorList, e.Error())
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Validation error! Please check your inputs!", "data": errorList})
	return err
}

func ValidatePhotoRequest(c *gin.Context, imageKey string, required bool) (*multipart.FileHeader, error) {
	file, fileErr := c.FormFile(imageKey)
	if fileErr != nil {
		if required {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "No file is received",
			})
		}
		return nil, nil
	}

	fileValidate, fileValidateErr := ValidatePhoto(c, file, required)
	return fileValidate, fileValidateErr
}

func ValidatePhoto(c *gin.Context, file *multipart.FileHeader, required bool) (*multipart.FileHeader, error) {
	uploadedFile, err := file.Open()

	if err != nil {
		if required {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"message": "Unable to open file",
			})
		}
		return nil, nil
	}

	b, _ := io.ReadAll(uploadedFile)
	mimeType := mimetype.Detect(b)

	switch mimeType.String() {
	case "image/heic_":
	case "image/jpeg":
	case "image/png":

	default:
		mimeError := "Mime: " + mimeType.String() + " not supported"
		c.JSON(http.StatusBadRequest, mimeError)
		return nil, errors.New(mimeError)
	}

	return file, nil
}
