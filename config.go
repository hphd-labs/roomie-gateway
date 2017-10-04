package main

type Config struct {
	Port         string `default:"8000"`
	AuthUpstream string `default:"localhost:8080" split_words:"true"`
	ShutdownTime int    `default:"30" split_words:"true"`
}
