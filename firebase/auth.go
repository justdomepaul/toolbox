package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
)

func NewAuthApp(ctx context.Context, client *firebase.App) (*auth.Client, error) {
	return client.Auth(ctx)
}

func VerifyIDToken(ctx context.Context, app *auth.Client, token string) (*auth.Token, error) {
	return app.VerifyIDToken(ctx, token)
}
