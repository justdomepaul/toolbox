package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/messaging"
)

func NewCloudMessagingApp(ctx context.Context, client *firebase.App) (*messaging.Client, error) {
	return client.Messaging(ctx)
}
