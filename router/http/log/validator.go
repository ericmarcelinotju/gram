package log

import (
	"github.com/gin-gonic/gin"
	"gitlab.com/firelogik/helios/constant/enums"
	"gitlab.com/firelogik/helios/domain/model"
	"gitlab.com/firelogik/helios/router/http/dto"
)

type PostLogPayload struct {
	Title   string         `json:"title" form:"title" binding:"required"`
	Subject string         `json:"subject" form:"subject"`
	Content string         `json:"content" form:"content"`
	Level   enums.LogLevel `json:"level" form:"level" uri:"level" enums:"info,success,warning,danger"`
	Type    enums.LogType  `json:"type" form:"type" uri:"type" enums:"event,system"`
}

type LogFilter struct {
	Title *string         `json:"title" form:"title" uri:"title"`
	Level *enums.LogLevel `json:"level" form:"level" uri:"level" enums:"info,success,warning,danger"`
	Type  *enums.LogType  `json:"type" form:"type" uri:"type" enums:"event,system"`
}

func BindPost(c *gin.Context) (*model.Log, error) {
	var json PostLogPayload
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	// Only allow announcement when created by API
	payload := &model.Log{
		Title:   json.Title,
		Subject: json.Subject,
		Content: json.Content,
		Level:   json.Level,
		Type:    json.Type,
	}

	return payload, nil
}

func BindGet(c *gin.Context) (*model.Log, error) {
	var json LogFilter
	if err := c.ShouldBind(&json); err != nil {
		return nil, err
	}

	var log *model.Log = &model.Log{}

	if json.Title != nil {
		log.Title = *json.Title
	}
	if json.Level != nil {
		log.Level = *json.Level
	}
	if json.Type != nil {
		log.Type = *json.Type
	}
	return log, nil
}

func BindID(c *gin.Context) (string, error) {
	var param dto.IdParam
	if err := c.ShouldBindUri(&param); err != nil {
		return "", err
	}

	return param.Id, nil
}
