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
