package dto

import (
	"github.com/ericmarcelinotju/gram/constant/enums"
)

// ListLogResponse struct defines response fields
type ListLogResponse struct {
	Logs []LogResponse `json:"logs"`
}

// LogResponse struct defines log response fields
type LogResponse struct {
	ID      string         `json:"id"`
	Title   string         `json:"title"`
	Subject string         `json:"subject"`
	Content string         `json:"content"`
	Level   enums.LogLevel `json:"level"`
	Type    enums.LogType  `json:"type"`

	CreatedAt string `json:"created_at"`
}
