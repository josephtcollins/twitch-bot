package main

import (
	"strings"
	"testing"
)

func TestChuckNorrisJokeListener(t *testing.T) {
	result := chuckNorrisJokeListener("!chucknorris", "channelname")
	if !strings.Contains(result, "PRIVMSG") {
		t.Errorf("Expected %v to contain PRIVMSG", result)
	}
}

func TestChuckNorrisJokeListenerNoNorris(t *testing.T) {
	resultNoNorris := chuckNorrisJokeListener("nochucknorris", "channelname")
	if resultNoNorris != "" {
		t.Errorf("Expected %v to be empty", resultNoNorris)
	}
}
