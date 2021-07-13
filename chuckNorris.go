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

func chuckNorrisJokeListener(c client, messageReceived string, currentChannel string) {
	if strings.Contains(messageReceived, "!chucknorris") {
		resp, err := http.Get("https://api.chucknorris.io/jokes/random")
		if err != nil {
			fmt.Println("Error getting joke.")
		}
		joke := chuckNorrisJoke{}
		err = json.NewDecoder(resp.Body).Decode(&joke)

		if err != nil {
			fmt.Println("Error parsing joke.")
		}

		message := "PRIVMSG #" + currentChannel + " :" + joke.Value
		printAndWriteMessage(message, c)
	}
}
