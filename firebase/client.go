package firebase

import (
	"context"
	"encoding/base64"
	firebase "firebase.google.com/go"
	"github.com/justdomepaul/toolbox/config"
	"google.golang.org/api/option"
)

// NewClient method
// env FIREBASE_CONFIG_JSON
// env FIREBASE_PROJECT_ID
func NewClient(ctx context.Context, firebaseOption config.Firebase) (*firebase.App, error) {
	options := make([]option.ClientOption, 0)
	if firebaseOption.FirebaseConfigJSON != "" {
		options = append(options, option.WithCredentialsJSON([]byte(firebaseOption.FirebaseConfigJSON)))
	}
	if firebaseOption.FirebaseConfigJSONBase64 != "" {
		firebaseConfigJSON, err := base64.StdEncoding.DecodeString(firebaseOption.FirebaseConfigJSONBase64)
		if err != nil {
			return nil, err
		}
		options = append(options, option.WithCredentialsJSON(firebaseConfigJSON))
	}
	var cf *firebase.Config
	if (firebaseOption.FirebaseConfigJSON == "" ||
		firebaseOption.FirebaseConfigJSONBase64 == "") &&
		firebaseOption.FirebaseProjectID != "" {
		cf = &firebase.Config{
			ProjectID: firebaseOption.FirebaseProjectID,
		}
	}
	return firebase.NewApp(ctx, cf, options...)
}
