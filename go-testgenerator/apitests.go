package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// TestCase represents individual test cases (positive or negative)
type TestCase struct {
	Name   string                 `json:"name"`
	Values map[string]interface{} `json:"values"`
	Output string                 `json:"output"`
}

// RouteConfig holds the configuration for a single route
type RouteConfig struct {
	Params                 []map[string]map[string]interface{} `json:"params"`
	AllowedHTTPStatusCodes []int                               `json:"allowed_http_status_codes"`
	PositiveTestCases      []TestCase                          `json:"positive_test_cases"`
	NegativeTestCases      []TestCase                          `json:"negative_test_cases"`
}

// Config represents the complete JSON configuration
type Config map[string]RouteConfig

// sendRequest generates and sends an HTTP POST request for a test case
func sendRequest(url string, testCase TestCase) {
	// Serialize the test case values into JSON
	requestBody, err := json.Marshal(testCase.Values)
	if err != nil {
		log.Fatalf("Failed to serialize request body for test case %s: %v", testCase.Name, err)
	}

	// Send the HTTP request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBody))
	if err != nil {
		log.Printf("Error in test case '%s': %v", testCase.Name, err)
		return
	}
	defer resp.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Failed to read response body for test case '%s': %v", testCase.Name, err)
		return
	}

	// Log the result
	fmt.Printf("Test Case: %s\n", testCase.Name)
	fmt.Printf("Request Body: %s\n", string(requestBody))
	fmt.Printf("Response Code: %d\n", resp.StatusCode)
	fmt.Printf("Response Body: %s\n", string(body))
	fmt.Println()
}

// processConfig reads the configuration and executes test cases
func processConfig(config Config) {
	for path, routeConfig := range config {
		url := "http://localhost:8080" + path // Adjust base URL as needed

		fmt.Printf("Testing route: %s\n", url)

		// Execute positive test cases
		fmt.Println("Running Positive Test Cases...")
		for _, testCase := range routeConfig.PositiveTestCases {
			sendRequest(url, testCase)
		}

		// Execute negative test cases
		fmt.Println("Running Negative Test Cases...")
		for _, testCase := range routeConfig.NegativeTestCases {
			sendRequest(url, testCase)
		}
	}
}

func main() {
	configFile := "config.json"

	// Read the JSON file
	fileData, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatalf("Failed to read JSON file: %v", err)
	}

	// Parse JSON into the Config struct
	var config Config
	err = json.Unmarshal(fileData, &config)
	if err != nil {
		log.Fatalf("Failed to parse JSON: %v", err)
	}
	fmt.Println(" \n ##################################################\n json config : \n", config, "\n ###################################\n")

	// Process the configuration and inject HTTP traffic
	processConfig(config)
}
