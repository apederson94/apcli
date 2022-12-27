package main

import (
	"apcli/models"
	"apcli/utils"
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var configLocation = ".apcli.yaml"
var config models.Config
var client = http.Client{}
var savedFields = map[string]interface{}{}
var environment models.Environment

func main() {
	var fileName string
	var workflowFile string

	flag.StringVar(&fileName, "f", "", "Specify the file you want to execute.")
	flag.StringVar(&workflowFile, "wf", "", "Specify the workflow you want to execute.")
	flag.Parse()

	config = utils.GetConfig(configLocation)

	if len(fileName) != 0 {
		doSingleCall(fileName)
	} else if len(workflowFile) != 0 {
		workflow := utils.GetWorkflow(workflowFile)

		for _, call := range workflow {
			apiCall := utils.GetApiCall(call.Call)
			overrideApiParams(apiCall, call.Overrides)
			req := generateRequest(apiCall)
			resp := doRequest(req)
			saveFields(resp, call.FieldsToSave)
			utils.WriteResponseToFile(resp, config.OutputLocation)
		}
	}
}

func doSingleCall(file string) {
	call := utils.GetApiCall(file)
	req := generateRequest(call)
	resp := doRequest(req)
	utils.WriteResponseToFile(resp, config.OutputLocation)
}

func doRequest(req *http.Request) []byte {
	resp, err := client.Do(req)
	utils.Check(err, fmt.Sprintf("Failed to do request %s", req.URL))

	body, err := io.ReadAll(resp.Body)
	utils.Check(err, fmt.Sprintf("Failed to read request body for %s", req.URL))

	var bodyJson interface{}
	err = json.Unmarshal(body, &bodyJson)
	utils.Check(err, fmt.Sprintf("Failed to unmarshal body %s", string(body)))

	out, err := json.MarshalIndent(bodyJson, "", "    ")
	utils.Check(err, fmt.Sprintf("Failed to convert response body into json %s", body))

	return out
}

func generateRequest(call models.ApiCall) *http.Request {
	var body []byte
	var err error
	var contentType string

	if call.Body != nil {
		switch call.Body.Type {
		case "json":
			body, err = json.Marshal(call.Body.Value)
			contentType = "application/json"
		case "form-urlencoded":
			contentType = "application/x-www-form-urlencoded"
			fmt.Println(call.Body.Value)
		}

		utils.Check(err, fmt.Sprintf("Failed to read body file %s", call.Body))
	}

	req, err := http.NewRequest(call.Method, call.Url, bytes.NewBuffer(body))

	utils.Check(err, fmt.Sprintf("Failed create new request with body %s", call.Body))

	for key, value := range call.Headers {
		req.Header.Set(key, value)
	}

	for key, value := range config.Headers {
		req.Header.Set(key, value)
	}

	if len(contentType) != 0 {
		req.Header.Set("Content-Type", contentType)
	}

	q := req.URL.Query()
	for key, value := range call.QueryParameters {
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()

	return req
}

func overrideApiParams(call models.ApiCall, overrides models.CallOverrides) {
	for k, v := range overrides.Headers {
		keysAndTags := utils.RetrieveEnvironmentKeysAndTags(v)

		for k, t := range keysAndTags {
			envValue, ok := environment.SavedValues[k]

			if !ok {
				log.Fatalln("Failed to find key", k, "in environment")
			}

			v = strings.ReplaceAll(v, t, envValue)
		}

		call.Headers[k] = v
	}

	for k, v := range overrides.QueryParameters {
		keysAndTags := utils.RetrieveEnvironmentKeysAndTags(v)

		for k, t := range keysAndTags {
			envValue, ok := environment.SavedValues[k]

			if !ok {
				log.Fatalln("Failed to find key", k, "in environment")
			}

			v = strings.ReplaceAll(v, t, envValue)
		}

		call.QueryParameters[k] = v
	}

	var drill = call.Body.Value
	arrayRegex, err := regexp.Compile("(?<=\\[)\\d+")
	utils.Check(err, fmt.Sprintf("Failed to compile regex: %s", "(?<=\\[)\\d+"))

	for _, o := range overrides.Body {
		keysAndTags := utils.RetrieveEnvironmentKeysAndTags(o.Value)

		for k, t := range keysAndTags {
			envValue, ok := environment.SavedValues[k]

			if !ok {
				log.Fatalln("Failed to find key", k, "in environment")
			}

			o.Value = strings.ReplaceAll(o.Value, t, envValue)
		}

		fieldKeys := strings.Split(o.Key, ".")

		for i, fk := range fieldKeys {
			if strings.ContainsAny(fk, "[]") {
				arrayKey := strings.Split(fk, "[")[0]
				arrayIndex, err := strconv.Atoi(arrayRegex.FindString(fk))
				utils.Check(err, fmt.Sprintf("Failed to find index in key: %s", fk))

				if i == len(fieldKeys) {
					drill.(map[string][]interface{})[arrayKey][arrayIndex] = o.Value
				} else {
					drill = drill.(map[string][]interface{})[arrayKey][arrayIndex]
				}
			} else {
				if i == len(fieldKeys) {
					drill.(map[string]interface{})[fk] = o.Value
				}
				drill = drill.(map[string]interface{})[fk]
			}
		}
	}
}

func saveFields(resp interface{}, fields map[string]string) {
	var drill = resp
	arrayRegex, err := regexp.Compile("(?<=\\[)\\d+")
	utils.Check(err, fmt.Sprintf("Failed to compile regex: %s", "(?<=\\[)\\d+"))

	for k, v := range fields {
		fieldKeys := strings.Split(v, ".")

		for _, fk := range fieldKeys {
			if strings.ContainsAny(fk, "[]") {
				newKey := strings.Split(fk, "[")[0]
				i, err := strconv.Atoi(arrayRegex.FindString(fk))
				utils.Check(err, fmt.Sprintf("Failed to find index in key: %s", fk))

				drill = drill.(map[string][]interface{})[newKey][i]
			} else {
				drill = drill.(map[string]interface{})[fk]
			}
		}

		savedFields[k] = drill
	}
}
