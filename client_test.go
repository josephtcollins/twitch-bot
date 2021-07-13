package main

import (
	"testing"
)

func TestNewClient(t *testing.T) {
	client := newClient()
	if client.conn == nil {
		t.Errorf("error")
	}
}
