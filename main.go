package main

import "os"

func main() {
	runTwitchBot(os.Getenv("TWITCH_USERNAME"), os.Getenv("TWITCH_OAUTH_TOKEN"))
}
