package helper

import (
	"bytes"
	"encoding/json"
	"net/http"
	"owlharbour-api/pkg/util"
)

var apiKey string

func init() {
	apiKey = util.GetEnv("FIREBASE_FCM_KEY", "")
}

func PushNotification(data map[string]interface{}, tokens []string) (map[string]interface{}, error) {
	url := "https://fcm.googleapis.com/fcm/send"

	// Create the request body
	requestBody := map[string]interface{}{
		"registration_ids": tokens,
		"notification":     data,
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	// Create the HTTP request
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+apiKey) // Set your FCM server key here

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Parse the response
	var response map[string]interface{}
	err = json.NewDecoder(resp.Body).Decode(&response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
