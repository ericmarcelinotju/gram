package firebase

import (
	"context"

	firebase "firebase.google.com/go"
	"gitlab.com/firelogik/helios/config"
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
