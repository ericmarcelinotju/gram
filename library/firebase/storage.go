package firebase

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"time"

	googleStorage "cloud.google.com/go/storage"
	"gitlab.com/firelogik/helios/config"
)

type Storage struct {
	Bucket *googleStorage.BucketHandle
}

func NewStorage(conf *config.Storage) (*Storage, error) {
	app, err := initApp(conf)
	if err != nil {
		return nil, err
	}
	client, err := app.Storage(context.Background())
	if err != nil {
		return nil, err
	}

	bucket, err := client.DefaultBucket()
	if err != nil {
		return nil, err
	}
	return &Storage{
		Bucket: bucket,
	}, err
}

func (s *Storage) postObject(file multipart.File, fileName string) (err error) {
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := s.Bucket.Object(fileName).NewWriter(ctx)
	if _, err = io.Copy(wc, file); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Download(fileName string) (object []byte, err error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	rc, err := s.Bucket.Object(fileName).NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	data, err := ioutil.ReadAll(rc)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (s *Storage) Remove(fileName string) (err error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.Bucket.Object(fileName).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *Storage) Upload(file multipart.File, filename string) error {
	defer file.Close()

	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()

	wc := s.Bucket.Object(filename).NewWriter(ctx)
	if _, err := io.Copy(wc, file); err != nil {
		return err
	}
	if err := wc.Close(); err != nil {
		return err
	}
	return nil
}

func (s *Storage) Path() string {
	return ""
}
