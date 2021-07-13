package main

import (
	"os"
)

func main() {
	runTwitchBot(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH_TOKEN"), func(line, currentChannel string) []string {
		return []string{
			chuckNorrisJokeListener(line, currentChannel),
		}
	})
}
