package main

import (
	"bufio"
	"fmt"
	"net/textproto"
	"os"
	"strings"
	"time"
)

const welcomeMessage = `
This is a twitch chat bot.

Available commands:
  join <channel_name> - join specified channel
  leave - leaves current channel
  send <message> - once in a channel
  whisper <user> <message>
  quit - will exit program
	help - for more info
`

func getSupportedActions() map[string]string {
	return map[string]string{
		"join":    "JOIN",
		"leave":   "LEAVE",
		"send":    "PRIVMSG",
		"whisper": "WHISPER",
		"quit":    "QUIT",
		"help":    "HELP",
	}
}

func printAndWriteMessage(message string, c client) {
	fmt.Println(formattedOutput(message))
	c.writeToConn(message)
}

func formattedOutput(str string) string {
	t := time.Now()
	return t.Format("2006-01-02T15:04:05.000Z") + " twitchbot$ " + str
}

func handlePing(c client) {
	c.writeToConn("PONG")
	fmt.Println(formattedOutput("PONG"))
}

func handleReader(c client, twitchChannel string, runCustomListener func(client, string, string)) {
	tp := textproto.NewReader(bufio.NewReader(c.conn))

	for {
		line, err := tp.ReadLine()
		if err != nil {
			fmt.Println("Something went wrong:", err)
		}

		// shorten ouput from channel
		lines := strings.Split(line, " ")
		for i, l := range lines {
			if strings.Contains(l, "tmi.twitch.tv") {
				lines[i] = strings.Split(l, "!")[0]
			}
		}

		c.read <- formattedOutput(strings.Join(lines, " "))

		if strings.HasPrefix(line, "PING") {
			handlePing(c)
		}

		// Run additional listening-based logic specified at implementation level
		runCustomListener(c, line, twitchChannel)
	}
}

func handleWriter(c client, defaultUsername string) {
	reader := bufio.NewReader(os.Stdin)
	twitchChannel := defaultUsername

	for {
		fmt.Print(formattedOutput(""))
		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		trimmedText := strings.Trim(text, "\n")
		splitText := strings.Split(trimmedText, " ")

		action, content := splitText[:1][0], splitText[1:]
		formattedContent := strings.Join(content, " ")

		switch getSupportedActions()[action] {
		case "JOIN":
			message := "JOIN #" + formattedContent
			twitchChannel = formattedContent
			printAndWriteMessage(message, c)
		case "LEAVE":
			message := "PART #" + twitchChannel
			printAndWriteMessage(message, c)
		case "PRIVMSG":
			message := "PRIVMSG #" + twitchChannel + " :" + formattedContent
			printAndWriteMessage(message, c)
		case "WHISPER":
			// The below is the apparent format for a whisper
			// PRIVMSG <channel> :/w <user> testing....
			// I'm likely getting blocked either for being a bot or hitting rate limits
			message := "PRIVMSG " + formattedContent
			printAndWriteMessage(message, c)
		case "QUIT":
			fmt.Println("Exiting program.")
			c.disconnect()
		case "HELP":
			fmt.Println(welcomeMessage)
			fmt.Println("Current channel:", twitchChannel)
		default:
			fmt.Println("Command not recognized.")
		}
	}
}

func runTwitchBot(defaultChannel string, OAUTHToken string, runCustomListener func(client, string, string)) {
	client := newClient()
	client.login(defaultChannel, OAUTHToken)

	printWelcome()
	go handleReader(client, defaultChannel, runCustomListener)
	go handleWriter(client, defaultChannel)

	for {
		msg := <-client.read
		fmt.Fprintln(os.Stdout, msg)
	}
}

func printWelcome() {
	// delay one second for UX
	time.Sleep(time.Second)
	fmt.Println(Banner + welcomeMessage)
}
