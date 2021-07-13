package main

import (
	"crypto/tls"
	"fmt"
	"net"
	"os"
)

type client struct {
	conn  net.Conn
	read  chan string
	write chan string
}

func newClient() client {
	// using nil as third param for default tls config
	conn, err := tls.Dial("tcp", "irc.chat.twitch.tv:6697", nil)

	if err != nil {
		// exit program since initial set up failed
		panic(err)
	}

	return client{
		conn:  conn,
		read:  make(chan string),
		write: make(chan string),
	}
}

func (c client) disconnect() {
	c.writeToConn("QUIT Bye")
	os.Exit(0)
}

func (c client) login(defaultUsername string, OAUTHToken string) {
	c.writeToConn("PASS " + OAUTHToken)
	c.writeToConn("NICK " + defaultUsername)
	// allows whispers to be received
	c.writeToConn("CAP REQ :twitch.tv/commands")
}

func (c client) writeToConn(message string) {
	fmt.Fprintf(c.conn, "%s\r\n", message)
}
