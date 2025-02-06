package middleware

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func NotificationHeaderMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if notif, exists := c.Get("notification"); exists {
			var notifSlice []interface{}
			switch n := notif.(type) {
			case []interface{}:
				notifSlice = n
			default:
				notifSlice = []interface{}{n}
			}
			jsonData, err := json.Marshal(notifSlice)
			if err != nil {
				fmt.Println("ERROR", err)
				return
			}
			encoded := base64.StdEncoding.EncodeToString(jsonData)
			c.Writer.Header().Set("x-notification", encoded)
		}
	}
}
