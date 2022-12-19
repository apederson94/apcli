package utils

import (
	"apcli/models"
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
	"os"
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
