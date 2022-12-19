package models

type ApiCallBody struct {
	Type  string      `yaml:"type"`
	Value interface{} `yaml:"value"`
}
