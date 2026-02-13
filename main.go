package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

const (
	TOKEN      string = "TOKEN"
	CHANNEL_ID string = "CHANNEL_ID"
)

var ID = fetchID()

type Author struct {
	ID string `json:"id"`
}

type Message struct {
	MessageID string `json:"id"`
	Author    Author `json:"author"`
}

type IDResponse struct {
	ID string `json:"id"`
}

func fetchID() string {
	req, err := http.NewRequest("GET", "https://discord.com/api/v9/users/@me", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var rp IDResponse
	err = json.NewDecoder(resp.Body).Decode(&rp)
	if err != nil {
		panic(err)
	}

	log.Printf("\033[32m[SUCCESS]\033[0m Fetch User ID (%s)\r\n", rp.ID)

	return rp.ID
}

func parseMessageID(messages []Message) []string {
	var result []string
	for _, m := range messages {
		if m.Author.ID == ID {
			result = append(result, m.MessageID)
		}
	}

	return result
}

func fetchMessages(count int) ([]string, error) {
	apiUrl := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages?limit=%d", CHANNEL_ID, count)
	req, err := http.NewRequest("GET", apiUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var mp []Message
	err = json.NewDecoder(resp.Body).Decode(&mp)
	if err != nil {
		return nil, err
	}

	return parseMessageID(mp), nil
}

func deleteMessage(id string) (int, error) {
	apiUrl := fmt.Sprintf("https://discord.com/api/v9/channels/%s/messages/%s", CHANNEL_ID, id)
	req, err := http.NewRequest("DELETE", apiUrl, nil)
	if err != nil {
		return 0, err
	}
	req.Header.Set("Authorization", TOKEN)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	return resp.StatusCode, nil
}

func main() {
	for {
		results, err := fetchMessages(50)
		if err != nil {
			log.Println(err)
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("\033[32m[SUCCESS]\033[0m Fetch Messages (%d)\r\n", len(results))

		for _, id := range results {
			sc, err := deleteMessage(id)
			if sc != 204 || err != nil {
				time.Sleep(10 * time.Second)
				continue
			}
			log.Printf("\033[31m[DELETED]\033[0m (%s)\r\n", id)

			time.Sleep(1 * time.Second)
		}

		time.Sleep(5 * time.Second)
	}
}
