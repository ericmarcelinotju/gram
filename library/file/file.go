package file

import (
	"io"
	"io/fs"
	"mime/multipart"
	"os"
	"path/filepath"
)

// File provides an abstraction on top of file manager system
type File struct {
	path string
}

func NewFileManager(basePath string) (*File, error) {
	filepath.Join(basePath)
	path := basePath
	// Check path exist
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(path, os.ModeDir); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}

	return &File{path: path}, nil
}

func (f *File) Upload(file multipart.File, fileName string) error {

	filePath := filepath.Join(f.path, fileName)

	dir := filepath.Dir(filePath)

	// Check path exist
	_, err := os.Stat(dir)
	if os.IsNotExist(err) {
		if err = os.MkdirAll(dir, os.ModeDir); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, file)
	return err
}

func (f *File) Download(fileName string) ([]byte, error) {
	pwd, _ := os.Getwd()
	filePath := filepath.Join(pwd, f.path, fileName)
	return os.ReadFile(filePath)
}

func (f *File) Remove(fileName string) error {
	pwd, _ := os.Getwd()
	filePath := filepath.Join(pwd, f.path, fileName)
	return os.Remove(filePath)
}

func (f *File) IsExist(fileName string) bool {
	filePath := f.path + fileName

	// Check file exist
	stat, err := os.Stat(filePath)

	return err == nil && !stat.IsDir()
}

func (f *File) Path() string {
	pwd, _ := os.Getwd()
	return filepath.Join(pwd, f.path)
}

func (f *File) List(folder string) ([]fs.FileInfo, error) {
	var fileNames []fs.FileInfo

	pwd, _ := os.Getwd()
	filePath := filepath.Join(pwd, f.path, folder)

	// Check file exist
	dirEntries, err := os.ReadDir(filePath)
	if err != nil {
		return nil, err
	}

	for _, dirEntry := range dirEntries {
		fileInfo, err := os.Stat(dirEntry.Name())
		if err != nil {
			return nil, err
		}
		fileNames = append(fileNames, fileInfo)
	}

	return fileNames, nil
}
