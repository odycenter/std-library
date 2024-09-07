package websocket_test

import (
	"github.com/odycenter/std-library/websocket"
	"log"
	"testing"
)

func TestConn(t *testing.T) {
	c, _, err := websocket.Dial("https", "127.0.0.1:8765", "/", nil)
	if err != nil {
		log.Panicln(err)
	}
	err = c.SendText([]byte("hello"))
	if err != nil {
		log.Panicln(err)
	}

}
