package config

import "time"

// Postgres type
type Postgres struct {
	PostgresUsername  string        `split_words:"true" default:"postgres"`
	PostgresPassword  string        `split_words:"true" default:"postgres"`
	PostgresHost      string        `split_words:"true" default:"localhost"`
	PostgresPort      string        `split_words:"true" default:"5432"`
	PostgresDatabase  string        `split_words:"true" default:"postgres"`
	PostgresSSLMode   string        `split_words:"true" default:"disable"`
	PostgresURL       string        `split_words:"true" default:"postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable"`
	PostgresTxTimeout time.Duration `split_words:"true" default:"5s"`
}
