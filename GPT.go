package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/go-resty/resty/v2"
)

func getGPT() {
	apiKey := chatGPTApiKey
	client := resty.New()
	response, err := client.R().
		SetAuthToken(apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]interface{}{
			"model": "gpt-4o-mini",
			"messages": []interface{}{
				map[string]interface{}{"role": "system", "content": "Hi can you tell me what is the factorial of 10?"},
			},
			"max_tokens": 50,
		}).
		Post(apiEndpoint)

	if err != nil {
		log.Fatalf("Error while sending the request: %v", err)
	}

	body := response.Body()
	var data map[string]interface{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println("Error while decoding JSON response:", err)
		return
	}

	for key, value := range data {
		fmt.Printf("key: %s, value: %s\n", key, value)
	}

	//content := data["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
	fmt.Println("Here")
}
