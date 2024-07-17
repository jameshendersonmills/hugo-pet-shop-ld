package main

import (
	"fmt"
	"net/http"
)

func triggerWebhook(webhookURL string) error {
	req, err := http.NewRequest("POST", webhookURL, nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("failed to trigger webhook, status code: %d", resp.StatusCode)
	}

	fmt.Println("Webhook triggered successfully")
	return nil
}

func main() {
	webhookURL := "https://app.launchdarkly.com/webhook/triggers/66982d637a37db0fe8a3410a/0c5271df-14fb-428f-b6dc-ceff5d12119b"

	err := triggerWebhook(webhookURL)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}