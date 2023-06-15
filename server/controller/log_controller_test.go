package controller

import (
	"log"
	"testing"

	"github.com/vjeantet/grok"
)

func Test_splitContent(t *testing.T) {
	message := "2023-06-15 23:08:14.963 INFO MacBook-Pro.local [Attach Listener] [com.jrasp.agent.core.server.socket.SocketServer.process] server socket start init..."
	Grok, _ := grok.New()
	maps, err := Grok.Parse(pattern, message)
	if err != nil {
		log.Fatalf("grok message error,%v", err)
	}
	if maps["host"] != "MacBook-Pro.local" {
		log.Fatalf("grok message error")
	}
}
