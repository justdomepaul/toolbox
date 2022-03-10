package config

// Spanner type
type Spanner struct {
	ProjectID             string `split_words:"true" default:"test-project"`
	Instance              string `split_words:"true" default:"test-instance"`
	Database              string `split_words:"true" default:"test-database"`
	EndPoint              string `split_words:"true" default:"localhost:9010"`
	WithoutAuthentication bool   `split_words:"true" default:"true"`
	GRPCInsecure          bool   `split_words:"true" default:"true"`
}
