package firebase

import (
	"context"
	firebase "firebase.google.com/go"
	"github.com/justdomepaul/toolbox/config"
	"google.golang.org/api/option"
)

// env FIREBASE_CONFIG_JSON
// env FIREBASE_PROJECT_ID
// NewClient method
func NewClient(ctx context.Context, firebaseOption config.Firebase) (*firebase.App, error) {
	options := make([]option.ClientOption, 0)
	if firebaseOption.FirebaseConfigJSON != "" {
		options = append(options, option.WithCredentialsJSON([]byte(firebaseOption.FirebaseConfigJSON)))
	}
	var cf *firebase.Config
	if firebaseOption.FirebaseConfigJSON == "" && firebaseOption.FirebaseProjectID != "" {
		cf = &firebase.Config{
			ProjectID: firebaseOption.FirebaseProjectID,
		}
	}
	return firebase.NewApp(ctx, cf, options...)
}
