package models

type CallOverrides struct {
	Headers         map[string]string `yaml:"headers"`
	Body            map[string]string `yaml:"body"`
	QueryParameters map[string]string `yaml:"queryParameters"`
}
