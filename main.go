package main

import (
	"os"
)

func main() {
	runTwitchBot(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH_TOKEN"), func(s1, s2 string) []string {
		return []string{
			chuckNorrisJokeListener(s1, s2),
		}
	})
}
