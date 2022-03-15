package config

// Firebase type
type Firebase struct {
	FirebaseConfigJSON string `split_words:"true" default:""`
	FirebaseProjectID  string `split_words:"true" default:""`
}
