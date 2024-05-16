package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func Validate(c *gin.Context, arr interface{}) error {
	var errorList []string
	if err := c.BindJSON(&arr); err != nil {
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
		c.JSON(http.StatusBadRequest, gin.H{"message": "Validation error! Please check your inputs!", "data": errorList})
		return err
	}

	return nil
}

func ValidateError(c *gin.Context, err error) error {
	var errorList []string
	validationErrors := err.(validator.ValidationErrors)
	for _, e := range validationErrors {
		//translatedErr := fmt.Errorf(e.Translate(trans))
		//errs = append(errs, translatedErr)
		errorList = append(errorList, e.Error())
	}
	c.JSON(http.StatusBadRequest, gin.H{"message": "Validation error! Please check your inputs!", "data": errorList})
	return err
}
