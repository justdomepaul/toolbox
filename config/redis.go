package config

// Redis type
type Redis struct {
	RedisUsername   string `split_words:"true" default:""`
	RedisPassword   string `split_words:"true" default:""`
	RedisClientName string `split_words:"true" default:""`
	RedisHost       string `split_words:"true" default:"localhost"`
	RedisPort       string `split_words:"true" default:"6379"`
	RedisPoolSize   int    `split_words:"true" default:"10"`
}
