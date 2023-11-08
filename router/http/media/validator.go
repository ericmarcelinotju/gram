package media

import (
	"errors"
	"mime/multipart"

	"github.com/gin-gonic/gin"
)

type FilePayload struct {
	File *multipart.FileHeader `form:"file" swaggerignore:"true"`
}

func Bind(c *gin.Context) (*multipart.File, error) {
	var payload FilePayload
	if err := c.ShouldBind(&payload); err != nil {
		return nil, err
	}

	if payload.File != nil {
		file, err := payload.File.Open()
		if err != nil {
			return nil, err
		}
		return &file, nil
	}

	return nil, errors.New("no file sent")
}
