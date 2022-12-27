package models

type CallOverrides struct {
	Headers         map[string]string `yaml:"headers"`
	Body            []BodyOverride    `yaml:"body"`
	QueryParameters map[string]string `yaml:"queryParameters"`
}

type BodyOverride struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}
