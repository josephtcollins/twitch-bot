package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strings"
	"time"
)

var welcomeMessage string = `
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

func printAndWriteMessage(message string, conn net.Conn) {
	fmt.Println(formattedOutput(message))
	writeToConn(conn, message)
}

func formattedOutput(str string) string {
	t := time.Now()
	return t.Format("2006-01-02T15:04:05.000Z") + " twitchbot$ " + str
}

func handleConnectionUpdates(conn net.Conn, twitchChannel string, customHandler func(net.Conn, string, string)) {
	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		line, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}

		fmt.Fprintln(os.Stdout, formattedOutput(line))

		if strings.HasPrefix(line, "PING") {
			writeToConn(conn, "PONG")
			fmt.Println(formattedOutput("PONG"))
		}

		customHandler(conn, line, twitchChannel)
	}
}

func handleUserInput(conn net.Conn, defaultUsername string) {
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
			printAndWriteMessage(message, conn)
		case "LEAVE":
			message := "PART #" + twitchChannel
			printAndWriteMessage(message, conn)
		case "PRIVMSG":
			message := "PRIVMSG #" + twitchChannel + " :" + formattedContent
			printAndWriteMessage(message, conn)
		case "WHISPER":
			// The below is the apparent format for a whisper
			// PRIVMSG <channel> :/w <user> testing....
			// I'm likely getting blocked either for being a bot or hitting rate limits
			message := "PRIVMSG " + formattedContent
			printAndWriteMessage(message, conn)
		case "QUIT":
			fmt.Println("Exiting program.")
			disconnect(conn)
		case "HELP":
			fmt.Println(welcomeMessage)
			fmt.Println("Current channel:", twitchChannel)
		default:
			fmt.Println("Command not recognized.")
		}
	}
}

func runTwitchBot(defaultUsername string, OAUTHToken string, customHandler func(net.Conn, string, string)) {
	conn := connect()
	login(conn, defaultUsername, OAUTHToken)

	go handleConnectionUpdates(conn, defaultUsername, customHandler)
	time.Sleep(time.Second)

	fmt.Println(getBanner() + welcomeMessage)

	fmt.Println(formattedOutput("Joining default channel: " + defaultUsername))
	writeToConn(conn, "JOIN #"+defaultUsername)

	handleUserInput(conn, defaultUsername)
}
