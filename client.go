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

func getSupportedActions() map[string]string {
	return map[string]string{
		"join": "JOIN",
		"send": "PRIVMSG",
		"quit": "QUIT",
	}
}

func connect() net.Conn {
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")

	// exit program since initial set up failed
	if err != nil {
		panic(err)
	}
	return conn
}

func disconnect(conn net.Conn) {
	writeToConn(conn, "QUIT Bye")
	os.Exit(0)
}

func login(conn net.Conn, defaultUsername string, OAUTHToken string) {
	writeToConn(conn, "PASS "+OAUTHToken)
	writeToConn(conn, "NICK "+defaultUsername)
}

func writeToConn(conn net.Conn, message string) {
	fmt.Fprintf(conn, "%s\r\n", message)
}

func formattedOutput(str string) string {
	t := time.Now()
	return t.Format("2006-01-02T15:04:05.000Z") + " twitchbot$ " + str
}

func handleConnectionUpdates(conn net.Conn) {
	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		status, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}

		fmt.Println(formattedOutput(status))

		if strings.HasPrefix(status, "PING") {
			fmt.Println(formattedOutput("PONG"))
			writeToConn(conn, "PONG")
		}

		if strings.Contains(status, "!chucknorris") {
			fmt.Println("Do chuck norris logic ghere")
		}
	}
}

func handleUserInput(conn net.Conn) {
	reader := bufio.NewReader(os.Stdin)

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
			writeToConn(conn, "JOIN "+formattedContent)
		case "PRIVMSG":
			message := "PRIVMSG #" + "banjomanjo1 :" + formattedContent
			fmt.Println(formattedOutput(message))
			writeToConn(conn, message)
		case "QUIT":
			fmt.Println("Exiting program.")
			disconnect(conn)
		default:
			fmt.Println("Command not recognized.")
		}
	}
}

func runTwitchBot(defaultUsername string, OAUTHToken string) {
	conn := connect()
	login(conn, defaultUsername, OAUTHToken)

	fmt.Println(formattedOutput("Joining default channel: " + defaultUsername))
	writeToConn(conn, "JOIN #"+defaultUsername)
	writeToConn(conn, "PRIVMSG #banjomanjo1 :this is an example message")

	go handleConnectionUpdates(conn)
	handleUserInput(conn)
}
