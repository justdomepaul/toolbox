package config

import "time"

// Postgresql type
type Postgresql struct {
	PostgresqlUsername  string        `split_words:"true" default:"root"`
	PostgresqlPassword  string        `split_words:"true" default:""`
	PostgresqlHost      string        `split_words:"true" default:"localhost"`
	PostgresqlPort      string        `split_words:"true" default:"26257"`
	PostgresqlDatabase  string        `split_words:"true" default:"defaultdb"`
	PostgresqlSSLMode   string        `split_words:"true" default:"disable"`
	PostgresqlURL       string        `split_words:"true" default:"postgres://root@localhost:26257/defaultdb?sslmode=disable"`
	PostgresqlTxTimeout time.Duration `split_words:"true" default:"5s"`
}
