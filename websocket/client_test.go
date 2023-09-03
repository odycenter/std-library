package websocket_test

import (
	"log"
	"std-library/websocket"
	"testing"
)

func TestConn(t *testing.T) {
	c, _, err := websocket.Dial("127.0.0.1:8765", "/", nil)
	if err != nil {
		log.Panicln(err)
	}
	err = c.SendText([]byte("hello"))
	if err != nil {
		log.Panicln(err)
	}

}
