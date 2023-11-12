package http

import (
	"github.com/ericmarcelinotju/gram/dto"
	"github.com/gin-gonic/gin"
)

func Bind[K interface{}](c *gin.Context) (*K, error) {
	var payload K
	if err := c.ShouldBind(&payload); err != nil {
		return nil, err
	}

	return &payload, nil
}

func BindId(c *gin.Context) (string, error) {
	var payload dto.IdDto
	if err := c.ShouldBindUri(&payload); err != nil {
		return "", err
	}

	return payload.Id, nil
}
