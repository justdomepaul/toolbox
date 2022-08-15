package firebase

import (
	"cloud.google.com/go/firestore"
	"context"
	firebase "firebase.google.com/go"
)

func NewFirestoreApp(ctx context.Context, client *firebase.App) (*firestore.Client, error) {
	return client.Firestore(ctx)
}
