package model

import (
	"time"

	"github.com/ericmarcelinotju/gram/constant/enums"
	dto "github.com/ericmarcelinotju/gram/router/http/dto/responses"
	timeUtils "github.com/ericmarcelinotju/gram/utils/time"
)

// Log struct defines the response model for a log APIs.
type Log struct {
	ID      string
	Title   string
	Subject string
	Content string
	Level   enums.LogLevel
	Type    enums.LogType

	CreatedAt time.Time
}

func (entity *Log) ToResponseModel() *dto.LogResponse {
	response := &dto.LogResponse{
		ID:      entity.ID,
		Title:   entity.Title,
		Subject: entity.Subject,
		Content: entity.Content,
		Level:   entity.Level,
		Type:    entity.Type,

		CreatedAt: *timeUtils.FormatResponseTime(&entity.CreatedAt),
	}
	return response
}
