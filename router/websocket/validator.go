package websocket

import (
	"github.com/gin-gonic/gin"
)

type WebsocketAuth struct {
	Key *string `form:"key" json:"key"`
}

func Bind(c *gin.Context) (*string, error) {
	var json WebsocketAuth
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	return json.Key, nil
}
