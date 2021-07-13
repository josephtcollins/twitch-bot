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

// on !chucknorris gets a random joke from the chuck norris API and
// composes IRC message for sending it to the current channel
func chuckNorrisJokeListener(message string, currentChannel string) string {
	if strings.Contains(message, "!chucknorris") {
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
