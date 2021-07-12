package main

import (
	"fmt"
	"net"
	"os"
)

// add type and make receivers

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
	// allows whispers to be received
	writeToConn(conn, "CAP REQ :twitch.tv/commands")
}

func writeToConn(conn net.Conn, message string) {
	fmt.Fprintf(conn, "%s\r\n", message)
}
