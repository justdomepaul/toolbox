package config

// Mongo type
type Mongo struct {
	MongoProtocol   string `split_words:"true" default:"mongodb+srv"`
	MongoUsername   string `split_words:"true" default:""`
	MongoPassword   string `split_words:"true" default:""`
	MongoHost       string `split_words:"true" default:""`
	MongoDatabase   string `split_words:"true" default:""`
	MongoAuthSource bool   `split_words:"true" default:""`
}
