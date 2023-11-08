package storage

import (
	"io/fs"
	"mime/multipart"
)

// Storage provides an abstraction on top of the file storage logic
type Storage interface {
	Upload(multipart.File, string) error
	Download(string) ([]byte, error)
	Remove(string) error
	List(string) ([]fs.FileInfo, error)

	Path() string
}
