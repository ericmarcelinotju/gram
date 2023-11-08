package log

import (
	"context"

	"github.com/ericmarcelinotju/gram/domain/model"
)

// Repository provides an abstraction on top of the log data source
type Repository interface {
	InsertLog(context.Context, *model.Log) error
	SelectLog(context.Context, *model.Log) ([]model.Log, error)
	SelectLogByID(context.Context, string) (*model.Log, error)
	DeleteLog(context.Context, *model.Log) error
}
