package models

type WorkflowCall struct {
	Call         string            `yaml:"call"`
	Overrides    CallOverrides     `yaml:"overrides"`
	FieldsToSave map[string]string `yaml:"fieldsToSave"`
}
