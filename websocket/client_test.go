package websocket_test

import (
	"log"
	"testing"

	"github.com/odycenter/std-library/websocket"
)

func TestConn(t *testing.T) {
	c, _, err := websocket.Dial("127.0.0.1:8765", "/", "", nil)
	if err != nil {
		log.Panicln(err)
	}
	err = c.SendText([]byte("hello"))
	if err != nil {
		log.Panicln(err)
	}

}
