package firebase

import (
	"context"
	"errors"
	"log"

	"firebase.google.com/go/messaging"
	"github.com/ericmarcelinotju/gram/config"
	"github.com/ericmarcelinotju/gram/domain/model"
)

type PushNotifier struct {
	ctx     context.Context
	client  *messaging.Client
	segment uint
}

func NewPushNotifier(conf *config.Storage) (*PushNotifier, error) {
	app, err := initApp(conf)
	if err != nil {
		return nil, err
	}

	ctx := context.Background()
	client, err := app.Messaging(ctx)
	if err != nil {
		return nil, err
	}
	return &PushNotifier{
		client: client,
		ctx:    ctx,
	}, nil
}

func (n *PushNotifier) Notify(topic string, data interface{}, recipient *model.User) error {
	messageMap, ok := data.(map[string]string)
	if !ok {
		return errors.New("data format invalid")
	}

	notification := messaging.Notification{
		Title: messageMap["title"],
		Body:  messageMap["body"],
	}

	message := &messaging.Message{
		Notification: &notification,
		Topic:        topic,
		Android: &messaging.AndroidConfig{
			Priority: "high",
		},
	}

	// Send message to spesific user
	if recipient != nil && recipient.NotificationToken != nil {
		message.Token = *recipient.NotificationToken
	}

	// Send a message to the devices subscribed to the provided topic.
	response, err := n.client.Send(n.ctx, message)
	if err != nil {
		return err
	}
	// Response is a message ID string.
	log.Println("Successfully sent message:", response, topic)

	return nil
}

func (n *PushNotifier) Subscribe(token string, topic string) error {
	_, err := n.client.SubscribeToTopic(n.ctx, []string{token}, topic)
	if err != nil {
		return err
	}
	return nil
}

func (n *PushNotifier) Unsubscribe(token string, topic string) error {
	_, err := n.client.UnsubscribeFromTopic(n.ctx, []string{token}, topic)
	if err != nil {
		return err
	}
	return nil
}
