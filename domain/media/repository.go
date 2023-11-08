package media

import (
	"context"
	"mime/multipart"
)

// Repository provides an abstraction on top of the building data source
type Repository interface {
	Upload(context.Context, *multipart.File) (string, error)
}
