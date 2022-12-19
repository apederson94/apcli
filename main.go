package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"gopkg.in/yaml.v3"
	"io"
	"log"
	"net/http"
	"os"
)

var configLocation = ".apcli.yaml"
var config Config
var client = http.Client{}
var outputLocation = "response.json"
var savedFields = map[string]string{}

func main() {
	var fileName string
	var workflowFile string

	flag.StringVar(&fileName, "f", "", "Specify the file you want to execute.")
	flag.StringVar(&workflowFile, "wf", "", "Specify the workflow you want to execute.")
	flag.Parse()

	config = getConfig()

	if len(fileName) != 0 {
		doSingleCall(fileName)
	} else if len(workflowFile) != 0 {
		workflow := getWorkflow(workflowFile)

		for _, call := range workflow {
			apiCall := decodeApiCall(call.Call)
			req := generateRequest(apiCall)
			overrideRequestParams(req, call.Overrides)
			resp := doRequest(req)
			writeResponseToFile(resp)
		}
	}
}

func doSingleCall(file string) {
	call := decodeApiCall(file)
	req := generateRequest(call)
	resp := doRequest(req)
	writeResponseToFile(resp)
}

func getWorkflow(fileName string) []WorkflowCall {
	data, err := os.ReadFile(fileName)

	check(err, fmt.Sprintf("Failed to read yaml file %s", fileName))

	var calls []WorkflowCall

	err = yaml.Unmarshal(data, &calls)

	check(err, fmt.Sprintf("Failed to unmarshal yaml file %s", fileName))

	return calls
}

func getConfig() Config {
	data, err := os.ReadFile(configLocation)

	check(err, fmt.Sprintf("Config file does not exist %s", configLocation))

	config := Config{}

	err = yaml.Unmarshal(data, &config)

	check(err, fmt.Sprintf("Failed to unmarshal yaml file %s", configLocation))

	return config
}

func check(e error, message string) {
	if e != nil {
		log.Fatalln(e, message)
	}
}

func writeResponseToFile(resp []byte) {
	err := os.WriteFile(outputLocation, resp, 0700)
	check(err, fmt.Sprintf("Failed to write output to %s", outputLocation))
}

func doRequest(req *http.Request) []byte {
	resp, err := client.Do(req)
	check(err, fmt.Sprintf("Failed to do request %s", req.URL))

	body, err := io.ReadAll(resp.Body)
	check(err, fmt.Sprintf("Failed to read request body for %s", req.URL))

	var bodyJson interface{}
	err = json.Unmarshal(body, &bodyJson)
	check(err, fmt.Sprintf("Failed to unmarshal body %s", string(body)))

	out, err := json.MarshalIndent(bodyJson, "", "    ")
	check(err, fmt.Sprintf("Failed to convert response body into json %s", body))

	return out
}

func generateRequest(call ApiCall) *http.Request {
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

		check(err, fmt.Sprintf("Failed to read body file %s", call.Body))
	}

	req, err := http.NewRequest(call.Method, call.Url, bytes.NewBuffer(body))

	check(err, fmt.Sprintf("Failed create new request with body %s", call.Body))

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

func overrideRequestParams(req *http.Request, overrides CallOverrides) {
	// do override logic
	for key, value := range overrides.Headers {
		req.Header.Del(key)
		req.Header.Set(key, value)
	}

	q := req.URL.Query()

	for key, value := range overrides.QueryParameters {
		q.Del(key)
		q.Add(key, value)
	}

	req.URL.RawQuery = q.Encode()
}

func decodeApiCall(fileName string) ApiCall {
	data, err := os.ReadFile(fileName)

	check(err, fmt.Sprintf("Failed to read yaml file %s", fileName))

	call := ApiCall{}

	err = yaml.Unmarshal(data, &call)

	check(err, fmt.Sprintf("Failed to unmarshal yaml file %s", fileName))

	return call
}

type Config struct {
	Headers map[string]string `yaml:"headers,omitempty"`
}

type WorkflowCall struct {
	Call         string            `yaml:"call"`
	Overrides    CallOverrides     `yaml:"overrides"`
	FieldsToSave map[string]string `yaml:"fieldsToSave"`
}

type CallOverrides struct {
	Headers         map[string]string `yaml:"headers"`
	Body            map[string]string `yaml:"body"`
	QueryParameters map[string]string `yaml:"queryParameters"`
}

type ApiCall struct {
	Url             string            `yaml:"url"`
	Method          string            `yaml:"method"`
	Body            *ApiCallBody      `yaml:"body,omitempty"`
	QueryParameters map[string]string `yaml:"queryParameters,omitempty"`
	Headers         map[string]string `yaml:"headers,omitempty"`
}

type ApiCallBody struct {
	Type  string      `yaml:"type"`
	Value interface{} `yaml:"value"`
}
