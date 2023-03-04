package config

// Mongo type
type Mongo struct {
	MongoUsername string `split_words:"true" default:""`
	MongoPassword string `split_words:"true" default:""`
	MongoHost     string `split_words:"true" default:""`
	MongoDatabase string `split_words:"true" default:""`
}
