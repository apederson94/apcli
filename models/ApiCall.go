package models

type ApiCall struct {
	Url             string            `yaml:"url"`
	Method          string            `yaml:"method"`
	Body            *ApiCallBody      `yaml:"body,omitempty"`
	QueryParameters map[string]string `yaml:"queryParameters,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
}
