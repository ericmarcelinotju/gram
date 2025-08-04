package storage

import (
	"context"
	"io"
	"io/ioutil"
	"mime/multipart"
	"time"

	googleStorage "cloud.google.com/go/storage"
	firebase "firebase.google.com/go"
	"github.com/ericmarcelinotju/gram/config"
	"google.golang.org/api/option"
)

func initApp(conf *config.Storage) (*firebase.App, error) {
	config := &firebase.Config{}
	if conf != nil {
		config.StorageBucket = conf.Path
	}

	opt := option.WithCredentialsFile("./firebase.json")
	app, err := firebase.NewApp(context.Background(), config, opt)
	if err != nil {
		return nil, err
	}
	return app, nil
}

type FirebaseStorage struct {
	Bucket *googleStorage.BucketHandle
}

func NewFirebaseStorage(conf *config.Storage) (*FirebaseStorage, error) {
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
	return &FirebaseStorage{
		Bucket: bucket,
	}, err
}

func (s *FirebaseStorage) Download(fileName string) (object []byte, err error) {
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

func (s *FirebaseStorage) Remove(fileName string) (err error) {
	ctx := context.Background()

	ctx, cancel := context.WithTimeout(ctx, time.Second*50)
	defer cancel()
	err = s.Bucket.Object(fileName).Delete(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (s *FirebaseStorage) Upload(file multipart.File, filename string) error {
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

func (s *FirebaseStorage) Path() string {
	return ""
}
