package config

import "time"

// Cockroach type
type Cockroach struct {
	CockroachUsername  string        `split_words:"true" default:"root"`
	CockroachPassword  string        `split_words:"true" default:""`
	CockroachHost      string        `split_words:"true" default:"localhost"`
	CockroachPort      string        `split_words:"true" default:"26257"`
	CockroachDatabase  string        `split_words:"true" default:"defaultdb"`
	CockroachSSLMode   string        `split_words:"true" default:"disable"`
	CockroachURL       string        `split_words:"true" default:"postgres://root@localhost:26257/defaultdb?sslmode=disable"`
	CockroachTxTimeout time.Duration `split_words:"true" default:"5s"`
}
