package schema

import (
	"github.com/ericmarcelinotju/gram/constant/enums"
	"github.com/ericmarcelinotju/gram/domain/model"
	"github.com/google/uuid"
)

// Log struct defines the database model for an log.
type Log struct {
	Model
	Title   string
	Subject string
	Content string
	Level   string
	Type    string
}

func NewLogSchema(entity *model.Log) *Log {
	id, _ := uuid.Parse(entity.ID)

	log := &Log{
		Model: Model{
			ID:        id,
			CreatedAt: entity.CreatedAt,
		},
		Title:   entity.Title,
		Subject: entity.Subject,
		Content: entity.Content,
		Level:   string(entity.Level),
		Type:    string(entity.Type),
	}

	return log
}

func (entity *Log) ToDomainModel() *model.Log {
	log := &model.Log{
		ID:        entity.ID.String(),
		Title:     entity.Title,
		Subject:   entity.Subject,
		Content:   entity.Content,
		Level:     enums.LogLevel(entity.Level),
		Type:      enums.LogType(entity.Type),
		CreatedAt: entity.CreatedAt,
	}

	return log
}
