package config

import "time"

// Cassandra type
type Cassandra struct {
	CassandraKeySpace          string        `split_words:"true" default:"test"`
	CassandraPort              string        `split_words:"true" default:"9042"`
	CassandraHosts             []string      `split_words:"true" default:"127.0.0.1"`
	CassandraForcePassword     bool          `split_words:"true" default:"false"`
	CassandraUsername          string        `split_words:"true" default:"cassandra"`
	CassandraPassword          string        `split_words:"true" default:"cassandra"`
	CassandraTimeout           time.Duration `split_words:"true" default:"600s"`
	CassandraConnectionTimeout time.Duration `split_words:"true" default:"600s"`
	ReconnectInterval          time.Duration `split_words:"true" default:"60s"`
	CassandraNumConnection     int           `split_words:"true" default:"5"`
	DisableInitialHostLookup   bool          `split_words:"true" default:"false"`
}
