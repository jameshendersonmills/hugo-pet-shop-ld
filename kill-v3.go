package main

import (
	"bytes"
	"fmt"
	"net/http"
)

func disableFeatureFlag(apiKey, projectKey, flagKey, environmentKey string) error {
	url := fmt.Sprintf("https://app.launchdarkly.com/api/v2/flags/%s/%s", projectKey, flagKey)
	body := bytes.NewBuffer([]byte(fmt.Sprintf(`[
		{
			"op": "replace",
			"path": "/environments/%s/on",
			"value": false
		}
	]`, environmentKey)))
	req, err := http.NewRequest("PATCH", url, body)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", fmt.Sprintf("api-%s", apiKey))
	req.Header.Add("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to disable feature flag, status code: %d", resp.StatusCode)
	}

	fmt.Println("Feature flag disabled successfully")
	return nil
}

func main() {
	apiKey := "api-73878f9b-9bc3-4b5c-9d39-a6c8130c5b49"
	projectKey := "default"
	flagKey := "v3-feature"
	environmentKey := "production"

	err := disableFeatureFlag(apiKey, projectKey, flagKey, environmentKey)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}