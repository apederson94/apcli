package models

type Config struct {
	Headers        map[string]string `yaml:"headers,omitempty"`
	OutputLocation string            `yaml:"outputLocation"`
}
