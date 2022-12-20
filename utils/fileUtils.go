package utils

import (
	"apcli/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"regexp"
	"strings"
)

func GetApiCall(fileName string) models.ApiCall {
	data, err := os.ReadFile(fileName)

	Check(err, fmt.Sprintf("Failed to read yaml file %s", fileName))

	call := models.ApiCall{}

	err = yaml.Unmarshal(data, &call)

	Check(err, fmt.Sprintf("Failed to unmarshal yaml file %s", fileName))

	return call
}

func GetWorkflow(fileName string) []models.WorkflowCall {
	data, err := os.ReadFile(fileName)

	Check(err, fmt.Sprintf("Failed to read yaml file %s", fileName))

	var calls []models.WorkflowCall

	err = yaml.Unmarshal(data, &calls)

	Check(err, fmt.Sprintf("Failed to unmarshal yaml file %s", fileName))

	return calls
}

func GetConfig(fileName string) models.Config {
	data, err := os.ReadFile(fileName)

	Check(err, fmt.Sprintf("Config file does not exist %s", fileName))

	config := models.Config{}

	err = yaml.Unmarshal(data, &config)

	Check(err, fmt.Sprintf("Failed to unmarshal yaml file %s", fileName))

	return config
}

func Check(e error, message string) {
	if e != nil {
		log.Fatalln(e, message)
	}
}

func WriteResponseToFile(resp []byte, fileName string) {
	err := os.WriteFile(fileName, resp, 0700)
	Check(err, fmt.Sprintf("Failed to write output to %s", fileName))
}

func RetrieveEnvironmentKeysAndTags(s string) map[string]string {
	r, err := regexp.Compile("{{\\s?\\w*\\s?}}")
	Check(err, fmt.Sprintf("Failed to compile regex: %s", "{{\\s?\\w*\\s?}}"))

	keys := r.FindAllString(s, 0)

	var keysAndTags map[string]string

	for _, tag := range keys {
		key := strings.ReplaceAll(tag, "{", "")
		key = strings.ReplaceAll(key, "}", "")
		key = strings.ReplaceAll(key, " ", "")
		keysAndTags[key] = tag
	}

	return keysAndTags
}
