package helpers

import (
	"net/http"
	"oosa_rewild/internal/structs"

	"github.com/gin-gonic/gin"
)

/*
200 -> Success
204 -> No data
400 -> User input error
404 -> Not found (wrong id)
500 -> Server errors
*/

func ResponseSuccessMessage(c *gin.Context, message string) {
	// 200
	c.JSON(http.StatusOK, structs.Message{Message: message})
}

func ResponseNoData(c *gin.Context, message string) {
	// 204
	c.JSON(http.StatusNoContent, structs.Message{Message: message})
}

func ResponseBadRequestError(c *gin.Context, message string) {
	// 400
	c.JSON(http.StatusBadRequest, structs.Message{Message: message})
}

func ResponseNotFound(c *gin.Context, message string) {
	// 404
	c.JSON(http.StatusNotFound, structs.Message{Message: message})
}

func ResponseError(c *gin.Context, message string) {
	// 500
	c.JSON(http.StatusInternalServerError, structs.Message{Message: message})
}
