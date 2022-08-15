package config

// Firestore type
type Firestore struct {
	ProjectID string `split_words:"true" default:""`
	EndPoint  string `split_words:"true" default:""`
}
