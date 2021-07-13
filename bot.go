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

func handlePing(c client) {
	c.writeToConn("PONG")
	fmt.Println(formattedOutput("PONG"))
}

func handleReader(c client, twitchChannel string, customListener func(string, string) []string) {
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
			handlePing(c)
		}

		// Run additional listening-based logic specified at implementation level
		messages := customListener(line, twitchChannel)
		for _, message := range messages {
			if len(message) > 0 {
				printAndWriteMessage(message, c)
			}
		}
	}
}

func handleWriter(c client, defaultUsername string) {
	reader := bufio.NewReader(os.Stdin)
	twitchChannel := defaultUsername

	for {
		// doesn't work reliably
		// attempting to get reliable leading text for write
		time.Sleep(time.Millisecond * 250)
		fmt.Print(formattedOutput(""))

		text, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		}

		action, content := inferredActionAndContent(text)

		switch getSupportedActions()[strings.ToLower(action)] {
		case "JOIN":
			message := "JOIN #" + content
			twitchChannel = content
			printAndWriteMessage(message, c)
		case "LEAVE":
			message := "PART #" + twitchChannel
			twitchChannel = ""
			printAndWriteMessage(message, c)
		case "PRIVMSG":
			message := "PRIVMSG #" + twitchChannel + " :" + content
			printAndWriteMessage(message, c)
		case "WHISPER":
			// The below is the apparent format for a whisper
			// PRIVMSG <channel> :/w <user> testing....
			// I'm likely getting blocked either for being a bot or hitting rate limits
			message := "PRIVMSG " + content
			printAndWriteMessage(message, c)
		case "QUIT":
			fmt.Println("Exiting program.")
			c.writeToConn("QUIT Bye")
			return
		case "HELP":
			fmt.Println(welcomeMessage)
			fmt.Println("Current channel:", twitchChannel)
		default:
			fmt.Println("Command not recognized.")
		}
	}
}

func inferredActionAndContent(text string) (string, string) {
	trimmedText := strings.Trim(text, "\r\n")
	splitText := strings.Split(trimmedText, " ")
	action, content := splitText[:1][0], splitText[1:]

	return action, strings.Join(content, " ")
}

func runTwitchBot(defaultChannel string, OAUTHToken string, customListener func(string, string) []string) {
	client := newClient()
	client.login(defaultChannel, OAUTHToken)

	printWelcome()

	go handleReader(client, defaultChannel, customListener)
	go handleWriter(client, defaultChannel)
	printReads(client)
}

func printReads(c client) {
	for {
		msg := <-c.read
		if msg == "QUIT" {
			return
		}
		fmt.Println(msg)
	}
}

func printWelcome() {
	// delay one second for UX
	time.Sleep(time.Second)
	fmt.Println(Banner + welcomeMessage)
}
