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
  quit - exits program
  help - for more info and current channel
`

var supportedActions = map[string]string{
	"join":    "JOIN",
	"leave":   "LEAVE",
	"send":    "PRIVMSG",
	"whisper": "WHISPER",
	"quit":    "QUIT",
	"help":    "HELP",
}

func runTwitchBot(username string, OAUTHToken string, customListener func(string, string) []string) {
	fmt.Println("Starting Twitch Bot...")
	client := newClient()

	printWelcome()

	go client.login(username, OAUTHToken)
	go handleReader(client, username, customListener)
	go handleWriter(client, username)

	channelsUpdated(client)
}

func printWelcome() {
	// delay one second for UX
	time.Sleep(time.Second)
	fmt.Println(Banner + welcomeMessage)
}

func handleReader(c client, username string, customListener func(string, string) []string) {
	defer func() {
		c.read <- "QUIT"
	}()

	tp := textproto.NewReader(bufio.NewReader(c.conn))

	for {
		line, err := tp.ReadLine()
		if err != nil {
			return
		}

		c.read <- formattedOutput(line)

		if strings.HasPrefix(line, "PING") {
			printAndWriteMessage("PONG", c)
		}

		// Run additional listening-based logic specified at implementation level
		messages := customListener(line, username)
		for _, message := range messages {
			if len(message) > 0 {
				printAndWriteMessage(message, c)
			}
		}
	}
}

func handleWriter(c client, username string) {
	reader := bufio.NewReader(os.Stdin)
	twitchChannel := ""

	for {
		// attempting to get reliable leading text for write
		// doesn't work reliably
		time.Sleep(time.Millisecond * 250)
		fmt.Print(formattedOutput(""))

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		action, content := splitFirstWordFromRest(text)
		switchAction := supportedActions[strings.ToLower(action)]

		messagesToWrite := getFormattedWriteMessage(switchAction, content, username, twitchChannel)
		for _, m := range messagesToWrite {
			printAndWriteMessage(m, c)
		}
		twitchChannel = getUpdatedChannel(switchAction, content, twitchChannel)
	}
}

func getUpdatedChannel(action string, content string, channel string) string {
	switch action {
	case "JOIN":
		return content
	case "PART":
		return ""
	default:
		return channel
	}
}

func getFormattedWriteMessage(action string, content string, username string, twitchChannel string) []string {
	var slice = []string{}
	switch action {
	case "JOIN":
		// leave current channel before joining new
		slice = append(slice, "PART #"+twitchChannel)
		slice = append(slice, "JOIN #"+content)
	case "LEAVE":
		slice = append(slice, "PART #"+twitchChannel)
	case "PRIVMSG":
		slice = append(slice, "PRIVMSG #"+twitchChannel+" :"+content)
	case "WHISPER":
		// This won't work:
		// I'm getting blocked for being a bot/hitting rate limits
		receivingUser, text := splitFirstWordFromRest(content)
		slice = append(slice, fmt.Sprintf("PRIVMSG #%s :/w %s %s", username, receivingUser, text))
	case "QUIT":
		fmt.Println("Exiting program.")
		slice = append(slice, "QUIT")
	case "HELP":
		fmt.Println(welcomeMessage)
		fmt.Println("Current channel:", twitchChannel)
	default:
		fmt.Println("Command not recognized! Type \"help\" for info.")
	}
	return slice
}

func splitFirstWordFromRest(text string) (string, string) {
	trimmedText := strings.Trim(text, "\r\n")
	splitText := strings.Split(trimmedText, " ")
	action, content := splitText[:1][0], splitText[1:]

	return action, strings.Join(content, " ")
}

func channelsUpdated(c client) {
	for {
		select {
		case readMessage := <-c.read:
			if readMessage == "QUIT" {
				return
			}
			fmt.Println(readMessage)
		case writeMessage := <-c.write:
			c.writeToConn(writeMessage)
		}
	}
}

func printAndWriteMessage(message string, c client) {
	fmt.Println(formattedOutput(message))
	c.write <- message
}

func formattedOutput(str string) string {
	lines := removeRedundantText(str)
	return time.Now().Format("2006-01-02 15:04:05") + " $ " + lines
}

// Lines from twitch often have redundant text in the format of:
// :r1verwater!r1verwater@r1verwater.tmi.twitch.tv
// This func distills that to "":r1verwater"
func removeRedundantText(str string) string {
	strArr := strings.Split(str, " ")
	for i, l := range strArr {
		if strings.Contains(l, "tmi.twitch.tv") {
			strArr[i] = strings.Split(l, "!")[0]
		}
	}
	return strings.Join(strArr, " ")
}
