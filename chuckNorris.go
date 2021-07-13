package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type chuckNorrisJoke struct {
	Value string `json:"value"`
}

func chuckNorrisJokeListener(messageReceived string, currentChannel string) string {
	if strings.Contains(messageReceived, "!chucknorris") {
		resp, err := http.Get("https://api.chucknorris.io/jokes/random")
		if err != nil {
			fmt.Println("Error getting joke.")
		}
		joke := chuckNorrisJoke{}
		json.NewDecoder(resp.Body).Decode(&joke)
		return "PRIVMSG #" + currentChannel + " :" + joke.Value
	}
	return ""
}
