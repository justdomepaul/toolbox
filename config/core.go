package config

// Core type
type Core struct {
	LoggerMode string `split_words:"true" default:"customized"`
	SystemName string `split_words:"true" default:"system"`
}
