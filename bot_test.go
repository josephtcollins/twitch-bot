package main

import (
	"strings"
	"testing"
)

func TestFormattedOutput(t *testing.T) {
	output := formattedOutput("this is a status")

	if !strings.Contains(output, "this is a status") {
		t.Errorf("Expected output to contain %v but got %v", "this is a status", output)
	}
}

func TestFormattedOutputNameShorten(t *testing.T) {
	output := formattedOutput("status :r1verwater!r1verwater@r1verwater.tmi.twitch.tv")

	if !strings.Contains(output, "status :r1verwater") {
		t.Errorf("Expected output to contain %v but got %v", "status :r1verwater", output)
	}
}

func TestGetUpdatedChannelJOIN(t *testing.T) {
	result := getUpdatedChannel("JOIN", "expected return", "")
	if result != "expected return" {
		t.Errorf("Expected output to be %v but got %v", "expected return", result)
	}
}

func TestGetUpdatedChannelPART(t *testing.T) {
	result := getUpdatedChannel("PART", "some value", "")
	if result != "" {
		t.Errorf("Expected output to be %v but got %v", "", result)
	}
}

func TestGetUpdatedChannelDefault(t *testing.T) {
	result := getUpdatedChannel("WHISPER", "some value", "samechannel")
	if result != "samechannel" {
		t.Errorf("Expected output to be %v but got %v", "samechannel", result)
	}
}
func TestGetFormattedWriteMessageJOIN(t *testing.T) {
	result := getFormattedWriteMessage("JOIN", "newchannel", "user", "oldchannel")
	if result[0] != "PART #oldchannel" || result[1] != "JOIN #newchannel" {
		t.Errorf("Expected output to leave oldchannel and join new channel, but got %v", result[0])
	}
}

func TestGetFormattedWriteMessageLEAVE(t *testing.T) {
	result := getFormattedWriteMessage("LEAVE", "", "user", "channel")
	if result[0] != "PART #channel" {
		t.Errorf("Expected output to be %v, but got %v", "PART #channel", result[0])
	}
}

func TestGetFormattedWriteMessagePRIVMSG(t *testing.T) {
	result := getFormattedWriteMessage("PRIVMSG", "message body", "user", "channel")
	if result[0] != "PRIVMSG #channel :message body" {
		t.Errorf("Expected output to be %v, but got %v", "PRIVMSG #channel :message body", result[0])
	}
}

func TestGetFormattedWriteMessageWHISPER(t *testing.T) {
	result := getFormattedWriteMessage("WHISPER", "touser message body", "fromuser", "channel")
	if result[0] != "PRIVMSG #fromuser :/w touser message body" {
		t.Errorf("Expected output to be %v, but got %v", "PRIVMSG #fromuser :/w touser message body", result[0])
	}
}

func TestGetFormattedWriteMessageQUIT(t *testing.T) {
	result := getFormattedWriteMessage("QUIT", "", "user", "channel")
	if result[0] != "QUIT" {
		t.Errorf("Expected output to be %v, but got %v", "QUIT", result[0])
	}
}

func TestGetFormattedWriteMessageHELP(t *testing.T) {
	result := getFormattedWriteMessage("HELP", "", "user", "channel")
	if len(result) != 0 {
		t.Errorf("Expected output slice to empty, but it was length %v", len(result))
	}
}

func TestGetFormattedWriteMessageDefault(t *testing.T) {
	result := getFormattedWriteMessage("ASDF", "", "user", "channel")
	if len(result) != 0 {
		t.Errorf("Expected output slice to empty, but it was length %v", len(result))
	}
}
